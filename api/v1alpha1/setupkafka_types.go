/*
Copyright 2022 Kubbee Tech.

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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SetupKafkaSpec defines the desired state of SetupKafka
type SetupKafkaSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of SetupKafka. Edit setupkafka_types.go to remove/update
	Cloud                    string                   `json:"cloud,omitempty"`
	Region                   string                   `json:"region,omitempty"`
	Availability             string                   `json:"availability,omitempty"`
	Type                     string                   `json:"type,omitempty"`
	CKU                      string                   `json:"cku,omitempty"`
	EncryptionKey            string                   `json:"encryptionKey,omitempty"`
	Context                  string                   `json:"context,omitempty"`
	EnvironmentReferenceSpec EnvironmentReferenceSpec `json:"environmentReferenceSpec,omitempty"`
}

// EnvironmentReferenceSpec defines the desired state of SetupKafka
type EnvironmentReferenceSpec struct {
	Environment string `json:"environment,omitempty"`
}

// SetupKafkaStatus defines the observed state of SetupKafka
type SetupKafkaStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Nodes []string `json:"nodes,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// SetupKafka is the Schema for the setupkafkas API
type SetupKafka struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SetupKafkaSpec   `json:"spec,omitempty"`
	Status SetupKafkaStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// SetupKafkaList contains a list of SetupKafka
type SetupKafkaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SetupKafka `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SetupKafka{}, &SetupKafkaList{})
}
