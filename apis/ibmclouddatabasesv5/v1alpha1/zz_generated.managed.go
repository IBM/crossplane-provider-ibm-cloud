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

// GetCondition of this AutoscalingGroup.
func (mg *AutoscalingGroup) GetCondition(ct runtimev1alpha1.ConditionType) runtimev1alpha1.Condition {
	return mg.Status.GetCondition(ct)
}

// GetDeletionPolicy of this AutoscalingGroup.
func (mg *AutoscalingGroup) GetDeletionPolicy() runtimev1alpha1.DeletionPolicy {
	return mg.Spec.DeletionPolicy
}

// GetProviderConfigReference of this AutoscalingGroup.
func (mg *AutoscalingGroup) GetProviderConfigReference() *runtimev1alpha1.Reference {
	return mg.Spec.ProviderConfigReference
}

/*
GetProviderReference of this AutoscalingGroup.
Deprecated: Use GetProviderConfigReference.
*/
func (mg *AutoscalingGroup) GetProviderReference() *runtimev1alpha1.Reference {
	return mg.Spec.ProviderReference
}

// GetWriteConnectionSecretToReference of this AutoscalingGroup.
func (mg *AutoscalingGroup) GetWriteConnectionSecretToReference() *runtimev1alpha1.SecretReference {
	return mg.Spec.WriteConnectionSecretToReference
}

// SetConditions of this AutoscalingGroup.
func (mg *AutoscalingGroup) SetConditions(c ...runtimev1alpha1.Condition) {
	mg.Status.SetConditions(c...)
}

// SetDeletionPolicy of this AutoscalingGroup.
func (mg *AutoscalingGroup) SetDeletionPolicy(r runtimev1alpha1.DeletionPolicy) {
	mg.Spec.DeletionPolicy = r
}

// SetProviderConfigReference of this AutoscalingGroup.
func (mg *AutoscalingGroup) SetProviderConfigReference(r *runtimev1alpha1.Reference) {
	mg.Spec.ProviderConfigReference = r
}

/*
SetProviderReference of this AutoscalingGroup.
Deprecated: Use SetProviderConfigReference.
*/
func (mg *AutoscalingGroup) SetProviderReference(r *runtimev1alpha1.Reference) {
	mg.Spec.ProviderReference = r
}

// SetWriteConnectionSecretToReference of this AutoscalingGroup.
func (mg *AutoscalingGroup) SetWriteConnectionSecretToReference(r *runtimev1alpha1.SecretReference) {
	mg.Spec.WriteConnectionSecretToReference = r
}

// GetCondition of this ScalingGroup.
func (mg *ScalingGroup) GetCondition(ct runtimev1alpha1.ConditionType) runtimev1alpha1.Condition {
	return mg.Status.GetCondition(ct)
}

// GetDeletionPolicy of this ScalingGroup.
func (mg *ScalingGroup) GetDeletionPolicy() runtimev1alpha1.DeletionPolicy {
	return mg.Spec.DeletionPolicy
}

// GetProviderConfigReference of this ScalingGroup.
func (mg *ScalingGroup) GetProviderConfigReference() *runtimev1alpha1.Reference {
	return mg.Spec.ProviderConfigReference
}

/*
GetProviderReference of this ScalingGroup.
Deprecated: Use GetProviderConfigReference.
*/
func (mg *ScalingGroup) GetProviderReference() *runtimev1alpha1.Reference {
	return mg.Spec.ProviderReference
}

// GetWriteConnectionSecretToReference of this ScalingGroup.
func (mg *ScalingGroup) GetWriteConnectionSecretToReference() *runtimev1alpha1.SecretReference {
	return mg.Spec.WriteConnectionSecretToReference
}

// SetConditions of this ScalingGroup.
func (mg *ScalingGroup) SetConditions(c ...runtimev1alpha1.Condition) {
	mg.Status.SetConditions(c...)
}

// SetDeletionPolicy of this ScalingGroup.
func (mg *ScalingGroup) SetDeletionPolicy(r runtimev1alpha1.DeletionPolicy) {
	mg.Spec.DeletionPolicy = r
}

// SetProviderConfigReference of this ScalingGroup.
func (mg *ScalingGroup) SetProviderConfigReference(r *runtimev1alpha1.Reference) {
	mg.Spec.ProviderConfigReference = r
}

/*
SetProviderReference of this ScalingGroup.
Deprecated: Use SetProviderConfigReference.
*/
func (mg *ScalingGroup) SetProviderReference(r *runtimev1alpha1.Reference) {
	mg.Spec.ProviderReference = r
}

// SetWriteConnectionSecretToReference of this ScalingGroup.
func (mg *ScalingGroup) SetWriteConnectionSecretToReference(r *runtimev1alpha1.SecretReference) {
	mg.Spec.WriteConnectionSecretToReference = r
}

// GetCondition of this Whitelist.
func (mg *Whitelist) GetCondition(ct runtimev1alpha1.ConditionType) runtimev1alpha1.Condition {
	return mg.Status.GetCondition(ct)
}

// GetDeletionPolicy of this Whitelist.
func (mg *Whitelist) GetDeletionPolicy() runtimev1alpha1.DeletionPolicy {
	return mg.Spec.DeletionPolicy
}

// GetProviderConfigReference of this Whitelist.
func (mg *Whitelist) GetProviderConfigReference() *runtimev1alpha1.Reference {
	return mg.Spec.ProviderConfigReference
}

/*
GetProviderReference of this Whitelist.
Deprecated: Use GetProviderConfigReference.
*/
func (mg *Whitelist) GetProviderReference() *runtimev1alpha1.Reference {
	return mg.Spec.ProviderReference
}

// GetWriteConnectionSecretToReference of this Whitelist.
func (mg *Whitelist) GetWriteConnectionSecretToReference() *runtimev1alpha1.SecretReference {
	return mg.Spec.WriteConnectionSecretToReference
}

// SetConditions of this Whitelist.
func (mg *Whitelist) SetConditions(c ...runtimev1alpha1.Condition) {
	mg.Status.SetConditions(c...)
}

// SetDeletionPolicy of this Whitelist.
func (mg *Whitelist) SetDeletionPolicy(r runtimev1alpha1.DeletionPolicy) {
	mg.Spec.DeletionPolicy = r
}

// SetProviderConfigReference of this Whitelist.
func (mg *Whitelist) SetProviderConfigReference(r *runtimev1alpha1.Reference) {
	mg.Spec.ProviderConfigReference = r
}

/*
SetProviderReference of this Whitelist.
Deprecated: Use SetProviderConfigReference.
*/
func (mg *Whitelist) SetProviderReference(r *runtimev1alpha1.Reference) {
	mg.Spec.ProviderReference = r
}

// SetWriteConnectionSecretToReference of this Whitelist.
func (mg *Whitelist) SetWriteConnectionSecretToReference(r *runtimev1alpha1.SecretReference) {
	mg.Spec.WriteConnectionSecretToReference = r
}
