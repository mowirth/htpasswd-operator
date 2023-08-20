package operator

import (
	"context"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

// Detects whether the changed secret is referenced by a HtpasswdUser. Runs the supplied onChange function in case an update is required.
func (s *HtpasswdOperator) checkConfigmapChange(ctx context.Context, configMapKey cache.ObjectName, onChange func(ctx context.Context, userKey cache.ObjectName) error) {
	for userKey, userKeyReference := range s.userRefMap {
		s.checkConfigMapRefChange(ctx, userKey, configMapKey, userKeyReference, onChange)
	}
}

// Check if a configMap is referenced in either user or password and call updater if it is
func (s *HtpasswdOperator) checkConfigMapRefChange(ctx context.Context, userKey cache.ObjectName, key cache.ObjectName, ref *CredentialReference, onChange func(ctx context.Context, userKey cache.ObjectName) error) {
	checkKeySourceConfigMapChange(ctx, userKey, key, ref.Username, onChange)
	checkKeySourceConfigMapChange(ctx, userKey, key, ref.Password, onChange)
}

func checkKeySourceConfigMapChange(ctx context.Context, userKey cache.ObjectName, key cache.ObjectName, ref *KeySource, onChange func(ctx context.Context, userKey cache.ObjectName) error) {
	if ref != nil && ref.ConfigMapRef != nil && ref.Namespace == key.Namespace {
		checkRefChange(ctx, userKey, key, ref.ConfigMapRef, onChange)
	}
}

/*
Returns the configKey from configMap
If no secret can be retrieved, returns an error with an empty string
*/
func (s *HtpasswdOperator) retrieveConfigKey(ctx context.Context, ref *KeySource) (string, error) {
	configMap, err := s.coreClient.CoreV1().ConfigMaps(ref.Namespace).Get(ctx, ref.ConfigMapRef.Name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return configMap.Data[ref.ConfigMapRef.Key], nil
}

func splitMetaNamespaceData(key string) (string, string) {
	if strings.Contains(key, "/") {
		split := strings.Split(key, "/")
		return split[0], split[1]
	}

	return "", key
}
