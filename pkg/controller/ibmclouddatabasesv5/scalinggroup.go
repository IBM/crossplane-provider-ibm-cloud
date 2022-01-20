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

package ibmclouddatabasesv5

import (
	"context"

	"github.com/IBM/go-sdk-core/v5/core"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	cpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
	"github.com/crossplane/crossplane-runtime/pkg/event"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"github.com/crossplane/crossplane-runtime/pkg/meta"
	"github.com/crossplane/crossplane-runtime/pkg/reconciler/managed"
	"github.com/crossplane/crossplane-runtime/pkg/reference"
	"github.com/crossplane/crossplane-runtime/pkg/resource"

	icdv5 "github.com/IBM/experimental-go-sdk/ibmclouddatabasesv5"

	"github.com/crossplane-contrib/provider-ibm-cloud/apis/ibmclouddatabasesv5/v1alpha1"
	"github.com/crossplane-contrib/provider-ibm-cloud/apis/v1beta1"
	ibmc "github.com/crossplane-contrib/provider-ibm-cloud/pkg/clients"
	ibmcsg "github.com/crossplane-contrib/provider-ibm-cloud/pkg/clients/scalinggroup"
)

const (
	errNotScalingGroup = "managed resource is not a ScalingGroup custom resource"

	errNewClient         = "cannot create new Client"
	errGetAuth           = "error getting auth info"
	errGetInstanceFailed = "error getting instance"
	errLateInitSpec      = "error on late initilize specs"
	errUpdateCR          = "eror updating CR"
	errGenObservation    = "error generating observation"
	errCheckUpToDate     = "cannot determine if instance is up to date"
	errSetOpts           = "error setting scaling group options"
	errResNotAvailable   = "Resource not found or not available"
)

// SetupScalingGroup adds a controller that reconciles ScalingGroup managed resources.
func SetupScalingGroup(mgr ctrl.Manager, l logging.Logger) error {
	name := managed.ControllerName(v1alpha1.ScalingGroupGroupKind)
	log := l.WithValues("ScalingGroup-controller", name)

	r := managed.NewReconciler(mgr,
		resource.ManagedKind(v1alpha1.ScalingGroupGroupVersionKind),
		managed.WithExternalConnecter(&sgConnector{
			kube:     mgr.GetClient(),
			usage:    resource.NewProviderConfigUsageTracker(mgr.GetClient(), &v1beta1.ProviderConfigUsage{}),
			clientFn: ibmc.NewClient,
			logger:   log}),
		managed.WithInitializers(managed.NewDefaultProviderConfig(mgr.GetClient())),
		managed.WithReferenceResolver(managed.NewAPISimpleReferenceResolver(mgr.GetClient())),
		managed.WithLogger(log),
		managed.WithRecorder(event.NewAPIRecorder(mgr.GetEventRecorderFor(name))))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&v1alpha1.ScalingGroup{}).
		Complete(r)
}

// A sgConnector is expected to produce an ExternalClient when its Connect method
// is called.
type sgConnector struct {
	kube     client.Client
	usage    resource.Tracker
	clientFn func(optd ibmc.ClientOptions) (ibmc.ClientSession, error)
	logger   logging.Logger
}

// Connect produces an ExternalClient for IBM Cloud API
func (c *sgConnector) Connect(ctx context.Context, mg resource.Managed) (managed.ExternalClient, error) {
	opts, err := ibmc.GetAuthInfo(ctx, c.kube, mg)
	if err != nil {
		return nil, errors.Wrap(err, errGetAuth)
	}

	service, err := c.clientFn(opts)
	if err != nil {
		return nil, errors.Wrap(err, errNewClient)
	}

	return &sgExternal{client: service, kube: c.kube, logger: c.logger}, nil
}

// An sgExternal observes, then either creates, updates, or deletes an
// external resource to ensure it reflects the managed resource's desired state.
type sgExternal struct {
	client ibmc.ClientSession
	kube   client.Client
	logger logging.Logger
}

func (c *sgExternal) Observe(ctx context.Context, mg resource.Managed) (managed.ExternalObservation, error) { // nolint:gocyclo
	cr, ok := mg.(*v1alpha1.ScalingGroup)
	if !ok {
		return managed.ExternalObservation{}, errors.New(errNotScalingGroup)
	}

	// since we do not really delete an external resource but rather we have a configuration on an existing service
	// we need to look at the deletion timestamp to figure out if the scaling group config was deleted.
	if meta.GetExternalName(cr) == "" || cr.DeletionTimestamp != nil {
		return managed.ExternalObservation{
			ResourceExists: false,
		}, nil
	}

	instance, _, err := c.client.IbmCloudDatabasesV5().GetDeploymentScalingGroups(&icdv5.GetDeploymentScalingGroupsOptions{ID: reference.ToPtrValue(meta.GetExternalName(cr))})
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(resource.Ignore(ibmc.IsResourceNotFound, err), errGetInstanceFailed)
	}

	currentSpec := cr.Spec.ForProvider.DeepCopy()
	if err = ibmcsg.LateInitializeSpec(&cr.Spec.ForProvider, instance); err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errLateInitSpec)
	}
	if !cmp.Equal(currentSpec, &cr.Spec.ForProvider) {
		if err := c.kube.Update(ctx, cr); err != nil {
			return managed.ExternalObservation{}, errors.Wrap(err, errUpdateCR)
		}
	}

	cr.Status.AtProvider, err = ibmcsg.GenerateObservation(instance)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errGenObservation)
	}

	if cr.Status.AtProvider.Groups != nil {
		cr.Status.SetConditions(cpv1.Available())
		cr.Status.AtProvider.State = string(cpv1.Available().Reason)
	} else {
		cr.Status.SetConditions(cpv1.Unavailable())
		cr.Status.AtProvider.State = string(cpv1.Unavailable().Reason)
	}

	upToDate, err := ibmcsg.IsUpToDate(meta.GetExternalName(cr), &cr.Spec.ForProvider, instance, c.logger)
	if err != nil {
		return managed.ExternalObservation{}, errors.Wrap(err, errCheckUpToDate)
	}

	return managed.ExternalObservation{
		ResourceExists:    true,
		ResourceUpToDate:  upToDate,
		ConnectionDetails: nil,
	}, nil
}

func (c *sgExternal) Create(ctx context.Context, mg resource.Managed) (managed.ExternalCreation, error) {
	cr, ok := mg.(*v1alpha1.ScalingGroup)
	if !ok {
		return managed.ExternalCreation{}, errors.New(errNotScalingGroup)
	}

	cr.SetConditions(cpv1.Creating())
	if cr.Spec.ForProvider.ID == nil {
		return managed.ExternalCreation{}, errors.New(errResNotAvailable)
	}

	meta.SetExternalName(cr, reference.FromPtrValue(cr.Spec.ForProvider.ID))
	return managed.ExternalCreation{ExternalNameAssigned: true}, nil
}

func (c *sgExternal) Update(ctx context.Context, mg resource.Managed) (managed.ExternalUpdate, error) {
	cr, ok := mg.(*v1alpha1.ScalingGroup)
	if !ok {
		return managed.ExternalUpdate{}, errors.New(errNotScalingGroup)
	}

	opts := &icdv5.SetDeploymentScalingGroupOptions{}
	err := ibmcsg.GenerateSetDeploymentScalingGroupOptions(meta.GetExternalName(cr), *cr, opts)
	if err != nil {
		return managed.ExternalUpdate{}, errors.Wrap(err, errSetOpts)
	}

	_, resp, err := c.client.IbmCloudDatabasesV5().SetDeploymentScalingGroup(opts)
	if err != nil {
		return managed.ExternalUpdate{}, ibmc.ExtractErrorMessage((*core.DetailedResponse)(resp), err)
	}

	return managed.ExternalUpdate{}, nil
}

func (c *sgExternal) Delete(ctx context.Context, mg resource.Managed) error {
	cr, ok := mg.(*v1alpha1.ScalingGroup)
	if !ok {
		return errors.New(errNotScalingGroup)
	}
	cr.SetConditions(cpv1.Deleting())
	return nil
}
