/*
Copyright 2024.

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

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DeclarativeLabelsSpec defines the desired state of DeclarativeLabels
type DeclarativeLabelsSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Time in seconds. This number of seconds defined how often manager checks
	// nodes and their labels. Default 60 seconds.
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:default=60
	// +kubebuilder:validation:Required
	Period *int32 `json:"period"`

	// Number of nodes which would be labelled. Minimum is 1.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:default=1
	MinNodes *int32 `json:"minNodes"`

	// Defines labels that we want minNodes to have.
	// Labels in kubernetes format `label: value` separated by new line.
	NodeLabels map[string]string `json:"nodeLabels"`
}

// DeclarativeLabelsStatus defines the observed state of DeclarativeLabels
type DeclarativeLabelsStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	LastCheck *metav1.Time `json:"lastClusterCheck,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// DeclarativeLabels is the Schema for the declarativelabels API
type DeclarativeLabels struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DeclarativeLabelsSpec   `json:"spec,omitempty"`
	Status DeclarativeLabelsStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DeclarativeLabelsList contains a list of DeclarativeLabels
type DeclarativeLabelsList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DeclarativeLabels `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DeclarativeLabels{}, &DeclarativeLabelsList{})
}
