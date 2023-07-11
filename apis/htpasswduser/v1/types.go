package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HtpasswdUser struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec HtpasswdUserSpec `json:"spec"`
}

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
