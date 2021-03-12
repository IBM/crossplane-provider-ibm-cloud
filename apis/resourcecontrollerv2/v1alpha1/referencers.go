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

package v1alpha1

import (
	"context"

	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/crossplane-runtime/pkg/reference"
	"github.com/crossplane/crossplane-runtime/pkg/resource"
)

// ResolveReferences of this ResourceKey
func (mg *ResourceKey) ResolveReferences(ctx context.Context, c client.Reader) error {
	r := reference.NewAPIResolver(c, mg)

	rsp, err := r.Resolve(ctx, reference.ResolutionRequest{
		CurrentValue: reference.FromPtrValue(mg.Spec.ForProvider.Source),
		Reference:    mg.Spec.ForProvider.SourceRef,
		Selector:     mg.Spec.ForProvider.SourceSelector,
		To:           reference.To{Managed: &ResourceInstance{}, List: &ResourceInstanceList{}},
		Extract:      SourceCRN(),
	})
	if err != nil {
		return errors.Wrap(err, "spec.forProvider.source")
	}
	mg.Spec.ForProvider.Source = reference.ToPtrValue(rsp.ResolvedValue)
	mg.Spec.ForProvider.SourceRef = rsp.ResolvedReference
	return nil
}

// SourceCRN extracts the resolved ResourceInstance's CRN
func SourceCRN() reference.ExtractValueFn {
	return func(mg resource.Managed) string {
		cr, ok := mg.(*ResourceInstance)
		if !ok {
			return ""
		}
		return cr.Status.AtProvider.CRN
	}
}
