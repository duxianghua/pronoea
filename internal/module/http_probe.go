package models

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HttpProbeSpec struct {
	Name string `json:"name,omitempty"`
	Url  string `json:"url,omitempty"`
}

type HttpProbe struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              HttpProbeSpec `json:"spec"`
}

type HttpProbeList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []HttpProbe `json:"items"`
}
