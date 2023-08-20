package v1

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	HtpasswdUserConfigured HtpasswdUserConditionType = "Configured"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.conditions[0].message"
// +kubebuilder:printcolumn:name="Synced",type="string",JSONPath=".status.conditions[0].status"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"
// +genclient
type HtpasswdUser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec HtpasswdUserSpec `json:"spec"`
	// +optional
	Status *HtpasswdUserStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type HtpasswdUserList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []HtpasswdUser `json:"items"`
}

type HtpasswdUserSpec struct {
	Username ReferenceOrValue `json:"username"`
	Password ReferenceOrValue `json:"password"`
}

type ReferenceOrValue struct {
	Value        string     `json:"value"`
	ConfigMapRef *Reference `json:"configMapKeyRef"`
	SecretKeyRef *Reference `json:"secretKeyRef"`
}

type ConfigSecretReference struct {
	ConfigMapRef *Reference `json:"configMapKeyRef"`
	SecretKeyRef *Reference `json:"secretKeyRef"`
}

type Reference struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}

type HtpasswdUserConditionType string

type HtpasswdUserStatus struct {
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,3,opt,name=observedGeneration"`

	// Represents the latest available observations of a operator user's current state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []HtpasswdUserCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,6,rep,name=conditions"`
}

type HtpasswdUserCondition struct {
	Type               HtpasswdUserConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=DeploymentConditionType"`
	Status             v1.ConditionStatus        `json:"status" protobuf:"bytes,2,opt,name=status,casttype=k8s.io/api/core/v1.ConditionStatus"`
	LastUpdateTime     metav1.Time               `json:"lastUpdateTime,omitempty" protobuf:"bytes,6,opt,name=lastUpdateTime"`
	LastTransitionTime metav1.Time               `json:"lastTransitionTime,omitempty" protobuf:"bytes,7,opt,name=lastTransitionTime"`
	Reason             string                    `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`
	Message            string                    `json:"message,omitempty" protobuf:"bytes,5,opt,name=message"`
}

func (s *HtpasswdUser) DeepCopyInto(out *HtpasswdUser) {
	out.TypeMeta = s.TypeMeta
	out.ObjectMeta = s.ObjectMeta
	out.Spec = HtpasswdUserSpec{
		Username: ReferenceOrValue{
			Value: s.Spec.Username.Value,
			ConfigMapRef: &Reference{
				Name: s.Spec.Username.ConfigMapRef.Name,
				Key:  s.Spec.Username.ConfigMapRef.Key,
			},
			SecretKeyRef: &Reference{
				Name: s.Spec.Username.SecretKeyRef.Name,
				Key:  s.Spec.Username.SecretKeyRef.Key,
			},
		},
		Password: ReferenceOrValue{
			Value: s.Spec.Password.Value,
			ConfigMapRef: &Reference{
				Name: s.Spec.Password.ConfigMapRef.Name,
				Key:  s.Spec.Password.ConfigMapRef.Key,
			},
			SecretKeyRef: &Reference{
				Name: s.Spec.Password.SecretKeyRef.Name,
				Key:  s.Spec.Password.SecretKeyRef.Key,
			},
		},
	}
}
