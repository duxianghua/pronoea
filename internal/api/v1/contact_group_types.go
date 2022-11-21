package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// ContactGroupSpec defines the desired state of ContactGroup
type ContactGroupSpec struct {
	Project string   `json:"projects"`
	Members []string `json:"members"`
}

// ContactGroupStatus defines the observed state of ContactGroup
type ContactGroupStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ContactGroup is the Schema for the ContactGroups API
type ContactGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ContactGroupSpec   `json:"spec,omitempty"`
	Status ContactGroupStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ContactGroupList contains a list of ContactGroup
type ContactGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ContactGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ContactGroup{}, &ContactGroupList{})
}
