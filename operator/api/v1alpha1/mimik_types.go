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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MimikSpec defines the desired state of Mimik
type MimikSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of Mimik. Edit Mimik_types.go to remove/update
	Service   string     `json:"service"`
	Version   string     `json:"version"`
	Endpoints []Endpoint `json:"endpoints"`
}

// Endpoint defines what a Mimik instance listens for
type Endpoint struct {
	Path        string       `json:"path"`
	Method      string       `json:"method"`
	Connections []Connection `json:"connections"`
}

// Connection defines an upstream connection
type Connection struct {
	Service string `json:"service"`
	Port    int    `json:"port"`
	Path    string `json:"path"`
	Method  string `json:"method"`
}

// MimikStatus defines the observed state of Mimik
type MimikStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Mimik is the Schema for the mimiks API
type Mimik struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MimikSpec   `json:"spec,omitempty"`
	Status MimikStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// MimikList contains a list of Mimik
type MimikList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Mimik `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Mimik{}, &MimikList{})
}
