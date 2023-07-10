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

// Code generated by angryjet. DO NOT EDIT.

package v1alpha1

import runtimev1alpha1 "github.com/crossplane/crossplane-runtime/apis/common/v1"

// GetCondition of this CloudantDatabase.
func (mg *CloudantDatabase) GetCondition(ct runtimev1alpha1.ConditionType) runtimev1alpha1.Condition {
	return mg.Status.GetCondition(ct)
}

// GetPublishConnectionDetailsTo of this CloudantDatabase.
func (mg *CloudantDatabase) GetPublishConnectionDetailsTo() *runtimev1alpha1.PublishConnectionDetailsTo {
	return mg.Spec.PublishConnectionDetailsTo
}

// GetDeletionPolicy of this CloudantDatabase.
func (mg *CloudantDatabase) GetDeletionPolicy() runtimev1alpha1.DeletionPolicy {
	return mg.Spec.DeletionPolicy
}

// GetProviderConfigReference of this CloudantDatabase.
func (mg *CloudantDatabase) GetProviderConfigReference() *runtimev1alpha1.Reference {
	return mg.Spec.ProviderConfigReference
}

/*
GetProviderReference of this CloudantDatabase.
Deprecated: Use GetProviderConfigReference.
*/
func (mg *CloudantDatabase) GetProviderReference() *runtimev1alpha1.Reference {
	return mg.Spec.ProviderReference
}

// GetWriteConnectionSecretToReference of this CloudantDatabase.
func (mg *CloudantDatabase) GetWriteConnectionSecretToReference() *runtimev1alpha1.SecretReference {
	return mg.Spec.WriteConnectionSecretToReference
}

// SetConditions of this CloudantDatabase.
func (mg *CloudantDatabase) SetConditions(c ...runtimev1alpha1.Condition) {
	mg.Status.SetConditions(c...)
}

// SetPublishConnectionDetailsTo of this CloudantDatabase.
func (mg *CloudantDatabase) SetPublishConnectionDetailsTo(r *runtimev1alpha1.PublishConnectionDetailsTo) {
	mg.Spec.PublishConnectionDetailsTo = r;
}

// SetDeletionPolicy of this CloudantDatabase.
func (mg *CloudantDatabase) SetDeletionPolicy(r runtimev1alpha1.DeletionPolicy) {
	mg.Spec.DeletionPolicy = r
}

// SetProviderConfigReference of this CloudantDatabase.
func (mg *CloudantDatabase) SetProviderConfigReference(r *runtimev1alpha1.Reference) {
	mg.Spec.ProviderConfigReference = r
}

/*
SetProviderReference of this CloudantDatabase.
Deprecated: Use SetProviderConfigReference.
*/
func (mg *CloudantDatabase) SetProviderReference(r *runtimev1alpha1.Reference) {
	mg.Spec.ProviderReference = r
}

// SetWriteConnectionSecretToReference of this CloudantDatabase.
func (mg *CloudantDatabase) SetWriteConnectionSecretToReference(r *runtimev1alpha1.SecretReference) {
	mg.Spec.WriteConnectionSecretToReference = r
}
