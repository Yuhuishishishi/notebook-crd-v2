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

// DaskSpec defines the desired state of Dask
type DaskSpec struct {
	NumWorkers        *int32         `json:"numWorkers"`
	SchedulerTemplate WorkerTemplate `json:"schedulerTemplate"`
	WorkerTemplate    WorkerTemplate `json:"workerTemplate"`
}

type WorkerTemplate struct {
	Spec corev1.PodSpec `json:"spec,omitempty"`
}

// DaskStatus defines the observed state of Dask
type DaskStatus struct {
	SchedulerReadyReplicas int32 `json:"schedulerReady"`
	WorkerReadyReplicas    int32 `json:"workerReady"`
	DesiredWorkers         int32 `json:"desiredWorkers"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:path=dasks,singular=dask,scope=Namespaced

// Dask is the Schema for the dasks API
type Dask struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DaskSpec   `json:"spec,omitempty"`
	Status DaskStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// DaskList contains a list of Dask
type DaskList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Dask `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Dask{}, &DaskList{})
}
