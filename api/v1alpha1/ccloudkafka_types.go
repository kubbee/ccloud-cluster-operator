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

// CCloudKafkaSpec defines the desired state of CCloudKafka
type CCloudKafkaSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	ClusterName         string              `json:"clusterName,omitempty"`
	Cloud               string              `json:"cloud,omitempty"`
	Region              string              `json:"region,omitempty"`
	Availability        string              `json:"availability,omitempty"` //--availability string     Availability of the cluster. Allowed Values: single-zone, multi-zone. (default "single-zone")
	ClusterType         string              `json:"clusterType,omitempty"`  // --type string             Type of the Kafka cluster. Allowed values: basic, standard, dedicated. (default "basic")
	Environment         string              `json:"environment,omitempty"`
	ApiKeyName          string              `json:"apiKeyName,omitempty"`
	CCloudKafkaDedicate CCloudKafkaDedicate `json:"ccloudKafkaDedicate,omitempty"`
	CCloudKafkaResource CCloudKafkaResource `json:"kafkaResource,omitempty"`
}

type CCloudKafkaResource struct {
	ResourceExist bool `json:"resourceExist,omitempty"`
}

type CCloudKafkaDedicate struct {
	Dedicated bool  `json:"dedicated,omitempty"`
	CKU       int64 `json:"cku,omitempty"`
}

// CCloudKafkaStatus defines the observed state of CCloudKafka
type CCloudKafkaStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ObservedGeneration int64       `json:"observedGeneration,omitempty"`
	Conditions         []Condition `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// CCloudKafka is the Schema for the ccloudkafkas API
type CCloudKafka struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CCloudKafkaSpec   `json:"spec,omitempty"`
	Status CCloudKafkaStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CCloudKafkaList contains a list of CCloudKafka
type CCloudKafkaList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []CCloudKafka `json:"items"`
}

func init() {
	SchemeBuilder.Register(&CCloudKafka{}, &CCloudKafkaList{})
}
