/*
Copyright 2022 The Crossplane Authors.

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
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	xpv1 "github.com/crossplane/crossplane-runtime/apis/common/v1"
)

// PipelineParameters are the configurable fields of a Pipeline.
type PipelineParameters struct {
	ConfigurableField string `json:"configurableField"`
}

// PipelineObservation are the observable fields of a Pipeline.
type PipelineObservation struct {
	ObservableField string `json:"observableField,omitempty"`
}

// A PipelineSpec defines the desired state of a Pipeline.
type PipelineSpec struct {
	xpv1.ResourceSpec `json:",inline"`
	ForProvider       PipelineParameters `json:"forProvider"`
}

// A PipelineStatus represents the observed state of a Pipeline.
type PipelineStatus struct {
	xpv1.ResourceStatus `json:",inline"`
	AtProvider          PipelineObservation `json:"atProvider,omitempty"`
}

// +kubebuilder:object:root=true

// A Pipeline is an example API type.
// +kubebuilder:printcolumn:name="READY",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="SYNCED",type="string",JSONPath=".status.conditions[?(@.type=='Synced')].status"
// +kubebuilder:printcolumn:name="EXTERNAL-NAME",type="string",JSONPath=".metadata.annotations.crossplane\\.io/external-name"
// +kubebuilder:printcolumn:name="AGE",type="date",JSONPath=".metadata.creationTimestamp"
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster,categories={crossplane,managed,nativeproviderjenkins}
type Pipeline struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PipelineSpec   `json:"spec"`
	Status PipelineStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// PipelineList contains a list of Pipeline
type PipelineList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Pipeline `json:"items"`
}

// Pipeline type metadata.
var (
	PipelineKind             = reflect.TypeOf(Pipeline{}).Name()
	PipelineGroupKind        = schema.GroupKind{Group: Group, Kind: PipelineKind}.String()
	PipelineKindAPIVersion   = PipelineKind + "." + SchemeGroupVersion.String()
	PipelineGroupVersionKind = SchemeGroupVersion.WithKind(PipelineKind)
)

func init() {
	SchemeBuilder.Register(&Pipeline{}, &PipelineList{})
}
