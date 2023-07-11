package listener

import (
	"context"

	v12 "go.flangaapis.com/htpasswd-kubernetes-operator/apis/htpasswduser/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (s *HtpasswdListener) retrieveSecretKey(ctx context.Context, namespace string, ref *v12.Reference) (string, error) {
	secret, err := s.coreClient.CoreV1().Secrets(namespace).Get(ctx, ref.Name, v1.GetOptions{})
	if err != nil {
		return "", err
	}

	return string(secret.Data[ref.Key]), nil
}

func (s *HtpasswdListener) retrieveConfigKey(ctx context.Context, namespace string, ref *v12.Reference) (string, error) {
	configMap, err := s.coreClient.CoreV1().ConfigMaps(namespace).Get(ctx, ref.Name, v1.GetOptions{})
	if err != nil {
		return "", err
	}

	return configMap.Data[ref.Key], nil
}
