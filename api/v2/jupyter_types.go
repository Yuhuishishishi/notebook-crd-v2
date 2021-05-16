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

package v2

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// JupyterSpec defines the desired state of Jupyter
type JupyterSpec struct {
	Template JupyterTemplate `json:",omitempty"`
}

type JupyterTemplate struct {
	Spec corev1.PodSpec `json:",omitempty"`
}

// JupyterStatus defines the observed state of Jupyter
type JupyterStatus struct {
	ReadyReplicas  int32                 `json:",omitempty"`
	ContainerState corev1.ContainerState `json:"state,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:path=jupyters,singular=jupyter,scope=Namespaced

// Jupyter is the Schema for the jupyters API
type Jupyter struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   JupyterSpec   `json:"spec,omitempty"`
	Status JupyterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// JupyterList contains a list of Jupyter
type JupyterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Jupyter `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Jupyter{}, &JupyterList{})
}
