/*
Copyright 2021.

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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TestSpec defines the desired state of Test
type TestSpec struct {
	// SpecField
	// +kubebuilder:validation:Required
	SpecField int `json:"specField,omitempty"`
}

// TestStatus defines the observed state of Test
type TestStatus struct {
	// +kubebuilder:validation:Optional
	StatusField int `json:"statusField"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Test is the Schema for the tests API
type Test struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   TestSpec   `json:"spec,omitempty"`
	Status TestStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// TestList contains a list of Test
type TestList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Test `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Test{}, &TestList{})
}
