package operator

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

// Detects whether the changed secret is referenced by a HtpasswdUser. Runs the supplied onChange function in case an update is required.
func (s *HtpasswdOperator) checkSecretChange(ctx context.Context, secretKey cache.ObjectName, onChange func(ctx context.Context, userKey cache.ObjectName) error) {
	for userKey, userKeyReference := range s.userRefMap {
		s.checkSecretRefChange(ctx, userKey, secretKey, userKeyReference, onChange)
	}
}

// Check if a secret is referenced in either user or password and call updater if it is
func (s *HtpasswdOperator) checkSecretRefChange(ctx context.Context, userKey cache.ObjectName, secretKey cache.ObjectName, ref *CredentialReference, onChange func(ctx context.Context, userKey cache.ObjectName) error) {
	if ref.Username != nil && ref.Username.SecretKeyRef != nil {
		checkRefChange(ctx, userKey, secretKey, ref.Username.SecretKeyRef, onChange)
	}

	if ref.Password != nil && ref.Password.SecretKeyRef != nil {
		checkRefChange(ctx, userKey, secretKey, ref.Password.SecretKeyRef, onChange)
	}
}

/*
Returns the secret from key.
If no secret can be retrieved, returns an error with an empty string
*/
func (s *HtpasswdOperator) retrieveSecretKey(ctx context.Context, ref *KeySource) (string, error) {
	secret, err := s.coreClient.CoreV1().Secrets(ref.Namespace).Get(ctx, ref.SecretKeyRef.Name, metav1.GetOptions{})
	if err != nil {
		return "", err
	}

	return string(secret.Data[ref.SecretKeyRef.Key]), nil
}
