package operator

import (
	"context"

	"github.com/mowirth/htpasswd-operator/pkg/apis/htpasswduser/v1"
	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/cache"
)

// upsertUser adds a user to the config.
// if it already exists, it is overwritten.
func (s *HtpasswdOperator) upsertUser(ctx context.Context, new *v1.HtpasswdUser) error {
	s.userRefMapMutex.Lock()
	defer s.userRefMapMutex.Unlock()

	ref := getCredentialReferenceForUser(new)
	cacheObj := cache.MetaObjectToName(new)
	s.userRefMap[cacheObj] = ref

	username, err := s.retrieveKeySource(ctx, ref.Username)
	if err != nil {
		return err
	}

	password, err := s.retrieveKeySource(ctx, ref.Password)
	if err != nil {
		return err
	}

	s.passwordMapMutex.Lock()
	defer s.passwordMapMutex.Unlock()
	s.passwordMap[cacheObj] = &BasicAuthCredentials{
		Username: username,
		Password: password,
	}

	s.signalSyncWithVolume()
	return nil
}

func (s *HtpasswdOperator) fetchAndUpsertHtpasswdUser(ctx context.Context, user cache.ObjectName) (updateErr error) {
	htpasswdUser, err := s.htpasswdUserClient.HtpasswdUsers(user.Namespace).Get(ctx, user.Name, metav1.GetOptions{})
	if err != nil {
		logrus.Errorf("Failed to get htpasswd user on secret create: %v", err)
		return err
	}

	defer func() {
		s.updateHtpasswdUserStatus(ctx, htpasswdUser, updateErr)
	}()

	if err := s.upsertUser(ctx, htpasswdUser); err != nil {
		logrus.Errorf("Failed to upsert user: %v", err)
		return err
	}

	return nil
}

func (s *HtpasswdOperator) removeUserFromPasswordConfig(_ context.Context, old cache.ObjectName) error {
	s.passwordMapMutex.Lock()
	defer s.passwordMapMutex.Unlock()
	delete(s.passwordMap, old)
	s.signalSyncWithVolume()
	return nil
}

func (s *HtpasswdOperator) removeUserInConfig(ctx context.Context, old cache.ObjectName) error {
	s.userRefMapMutex.Lock()
	defer s.userRefMapMutex.Unlock()
	delete(s.userRefMap, old)
	return s.removeUserFromPasswordConfig(ctx, old)
}
