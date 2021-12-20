/*
Copyright 2020 The Crossplane Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"gopkg.in/alecthomas/kingpin.v2"
	ca "k8s.io/client-go/tools/cache"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/crossplane/crossplane-runtime/pkg/logging"

	"github.com/crossplane-contrib/provider-ibm-cloud/apis"
	"github.com/crossplane-contrib/provider-ibm-cloud/pkg/controller"
)

func main() {
	var (
		app            = kingpin.New(filepath.Base(os.Args[0]), "IBM Cloud support for Crossplane.").DefaultEnvars()
		debug          = app.Flag("debug", "Run with debug logging.").Short('d').Bool()
		syncPeriod     = app.Flag("sync", "Controller manager sync period such as 300ms, 1.5h, or 2h45m").Short('s').Default("1h").Duration()
		leaderElection = app.Flag("leader-election", "Use leader election for the conroller manager.").Short('l').Default("false").OverrideDefaultFromEnvar("LEADER_ELECTION").Bool()
		watchNamespace = app.Flag("namespace", "Namespace containing nss configuration.").Short('n').Default("ibm-common-services").OverrideDefaultFromEnvar("WATCH_NAMESPACE").String()
	)
	kingpin.MustParse(app.Parse(os.Args[1:]))

	zl := zap.New(zap.UseDevMode(*debug))
	log := logging.NewLogrLogger(zl.WithName("provider-ibm-cloud"))
	if *debug {
		// The controller-runtime runs with a no-op logger by default. It is
		// *very* verbose even at info level, so we only provide it a real
		// logger when we're running in debug mode.
		ctrl.SetLogger(zl)
	}

	log.Debug("Starting", "sync-period", syncPeriod.String())

	cfg, err := ctrl.GetConfig()
	kingpin.FatalIfError(err, "Cannot get API server rest config")

	// IBM Patch: reduce cluster permission
	// we want to restrict cache to watch only a given list of namespaces
	// instead of all (cluster scoped). List of namespaces is read
	// from NamespaceScope resource, if it exists. Changes in this resource
	// should restart Provider's pod.

	// Default to watchNamespace
	namespaces := []string{*watchNamespace}
	nssName := "common-service"

	cfn, err := client.New(cfg, client.Options{})
	kingpin.FatalIfError(err, "Cannot create client for reading NamespaceScope")

	nfn, err := listNamespacesFromNss(cfn, *watchNamespace, nssName)
	// Proceed with informer when no error found during NamespaceScope reading
	if err == nil {
		namespaces = append(namespaces, nfn...)

		// Start informer to watch for changes in NamespaceScope resource
		dc, err := dynamic.NewForConfig(cfg)
		kingpin.FatalIfError(err, "Cannot create client for observing NamespaceScope")

		stopper := startNssInformer(log, dc, *watchNamespace, nssName)
		defer close(stopper)
		log.Debug(fmt.Sprintf("Starting watch on namespaceScope %s", nssName))
	}
	log.Debug(fmt.Sprintf("Creating multinamespaced cache with namespaces: %+q", namespaces))
	// IBM Patch end: reduce cluster permission

	mgr, err := ctrl.NewManager(cfg, ctrl.Options{
		LeaderElection:   *leaderElection,
		LeaderElectionID: "crossplane-leader-election-provider-ibm-cloud",
		SyncPeriod:       syncPeriod,
		NewCache:         cache.MultiNamespacedCacheBuilder(namespaces),
	})
	kingpin.FatalIfError(err, "Cannot create controller manager")

	kingpin.FatalIfError(apis.AddToScheme(mgr.GetScheme()), "Cannot add IBM Cloud APIs to scheme")
	kingpin.FatalIfError(controller.Setup(mgr, log), "Cannot setup IBM Cloud controllers")
	kingpin.FatalIfError(mgr.Start(ctrl.SetupSignalHandler()), "Cannot start controller manager")
}

// IBM Patch:  reduce cluster permission
func listNamespacesFromNss(cfn client.Client, watchNamespace string, nssName string) ([]string, error) {
	var namespaces []string
	nss := &unstructured.Unstructured{}
	nss.SetGroupVersionKind(schema.GroupVersionKind{Version: "operator.ibm.com/v1", Kind: "NamespaceScope"})

	if err := cfn.Get(context.Background(), types.NamespacedName{Namespace: watchNamespace, Name: nssName}, nss); err != nil {
		if errors.IsNotFound(err) {
			// If not found return empty array and no error to allow for further steps.
			return namespaces, nil
		}
		// Block further steps because probably there is no such CRD on the cluster.
		return namespaces, err
	}

	spec := nss.Object["spec"].(map[string]interface{})
	members := spec["namespaceMembers"]
	if members != nil {
		for _, m := range members.([]interface{}) {
			if m.(string) == watchNamespace {
				continue
			}
			namespaces = append(namespaces, m.(string))
		}
	}
	return namespaces, nil
}

func startNssInformer(log logging.Logger, dc dynamic.Interface, watchNamespace string, nssName string) chan struct{} {
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(dc, 0, watchNamespace, nil)
	informer := factory.ForResource(schema.GroupVersionResource{
		Group:    "operator.ibm.com",
		Version:  "v1",
		Resource: "namespacescopes",
	})

	// Handle each update by restarting pod.
	handlers := ca.ResourceEventHandlerFuncs{
		UpdateFunc: func(oldObj, newObj interface{}) {
			if newObj.(*unstructured.Unstructured).GetName() == nssName {
				log.Debug("Observed NamespaceScope has been updated, restarting")
				os.Exit(1)
			}
		},
	}
	informer.Informer().AddEventHandler(handlers)

	stopper := make(chan struct{})
	go informer.Informer().Run(stopper)

	return stopper
}

// IBM Patch end
