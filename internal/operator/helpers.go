package operator

import (
	"context"
	"fmt"
	"reflect"

	"github.com/mowirth/htpasswd-operator/pkg/apis/htpasswduser/v1"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// retrieve a username from the cluster.
// We allow reading from secrets, configMaps or directly from the value.
func (s *HtpasswdOperator) retrieveKeySource(ctx context.Context, ref *KeySource) (string, error) {
	if ref == nil {
		return "", fmt.Errorf("no user or no specification provided, skipping. (User: %v)", ref)
	}

	if ref.SecretKeyRef != nil {
		return s.retrieveSecretKey(ctx, ref)
	}

	if ref.ConfigMapRef != nil {
		return s.retrieveConfigKey(ctx, ref)
	}

	return ref.Value, nil
}

// updates the status of the HtpasswdUser in Kubernetes
func (s *HtpasswdOperator) updateHtpasswdUserStatus(ctx context.Context, htpasswdUser *v1.HtpasswdUser, userErr error) {
	if htpasswdUser.Status == nil {
		htpasswdUser.Status = &v1.HtpasswdUserStatus{}
	}

	htpasswdUser.Status.ObservedGeneration = htpasswdUser.ObjectMeta.Generation
	updateHtpasswdUsersStatusConditions(htpasswdUser.Status, userErr)
	if _, err := s.htpasswdUserClient.HtpasswdUsers(htpasswdUser.Namespace).UpdateStatus(ctx, htpasswdUser, metav1.UpdateOptions{}); err != nil {
		logrus.Errorf("Failed to update status for user %v: %v", htpasswdUser.Name, err)
	}
}

func updateHtpasswdUsersStatusConditions(st *v1.HtpasswdUserStatus, err error) {
	cond := func() *v1.HtpasswdUserCondition {
		for i := range st.Conditions {
			if st.Conditions[i].Type == v1.HtpasswdUserConfigured {
				return &st.Conditions[i]
			}
		}
		st.Conditions = append(st.Conditions, v1.HtpasswdUserCondition{
			Type: v1.HtpasswdUserConfigured,
		})
		return &st.Conditions[len(st.Conditions)-1]
	}()

	var status corev1.ConditionStatus
	if err == nil {
		status = corev1.ConditionTrue
		cond.Message = ""
	} else {
		status = corev1.ConditionFalse
		cond.Message = err.Error()
	}
	cond.LastUpdateTime = metav1.Now()
	if cond.Status != status {
		cond.LastTransitionTime = cond.LastUpdateTime
		cond.Status = status
	}
}

// check if the htpasswduser itself was updated
func htpasswdUserChanged(oldObj, newObj interface{}) bool {
	oldUser, err := convertHtpasswdUser(oldObj)
	if err != nil {
		return true // any conversion error means we assume it might have changed
	}

	newUser, err := convertHtpasswdUser(newObj)
	if err != nil {
		return true
	}

	return !reflect.DeepEqual(oldUser.Spec, newUser.Spec)
}

// convertHtpasswdUser converts a kubernetes object into a typed HtpasswdUser
func convertHtpasswdUser(obj interface{}) (*v1.HtpasswdUser, error) {
	htpasswdUser, ok := (obj).(*v1.HtpasswdUser)
	if !ok {
		return nil, fmt.Errorf("failed to cast %v into HtpasswdUser", obj)
	}

	if htpasswdUser.APIVersion == "" || htpasswdUser.Kind == "" {
		gv := schema.GroupVersion{Group: v1.GroupName, Version: v1.GroupVersion}
		gvk := gv.WithKind("HtpasswdUser")
		htpasswdUser.APIVersion = gvk.GroupVersion().String()
		htpasswdUser.Kind = gvk.Kind
	}

	return htpasswdUser, nil
}

func getKeySourceForReference(namespace string, ref v1.ReferenceOrValue) *KeySource {
	keySource := &KeySource{Namespace: namespace}

	if ref.SecretKeyRef != nil {
		keySource.SecretKeyRef = ref.SecretKeyRef
		return keySource
	}

	if ref.ConfigMapRef != nil {
		keySource.ConfigMapRef = ref.ConfigMapRef
		return keySource
	}

	if ref.Value != "" {
		keySource.Value = ref.Value
	}

	return keySource
}

// Get a CredentialReference from the User
func getCredentialReferenceForUser(user *v1.HtpasswdUser) *CredentialReference {
	userRef := getKeySourceForReference(user.Namespace, user.Spec.Username)
	secretRef := getKeySourceForReference(user.Namespace, user.Spec.Password)

	return &CredentialReference{
		Username: userRef,
		Password: secretRef,
	}
}
