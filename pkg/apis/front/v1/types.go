package v1

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type FrontSpec struct {
	Replicas *int32 `json:"replicas"`
	Image string `json:"image"`
	Ports     []corev1.ServicePort `json:"ports,omitempty"`
}

type FrontStatus struct {
	appsv1.DeploymentStatus `json:",inline"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type Front struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FrontSpec   `json:"spec,omitempty"`
	Status FrontStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type FrontList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Front `json:"items"`
}

