package listener

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/sirupsen/logrus"
	"go.flangaapis.com/htpasswd-kubernetes-operator/apis/htpasswduser/v1"
	"golang.org/x/crypto/bcrypt"
	"k8s.io/apimachinery/pkg/types"
)

// retrieve a username from the cluster.
// We allow reading from secrets, configMaps or directly from the value.
func (s *HtpasswdListener) retrieveUsername(ctx context.Context, in *v1.HtpasswdUser) (string, error) {
	if in == nil {
		return "", fmt.Errorf("no user or no specification provided, skipping. (User: %v)", in)
	}

	if in.Spec.Username.SecretKeyRef != nil {
		return s.retrieveSecretKey(ctx, in.Namespace, in.Spec.Username.SecretKeyRef)
	}

	if in.Spec.Username.ConfigMapRef != nil {
		return s.retrieveConfigKey(ctx, in.Namespace, in.Spec.Username.ConfigMapRef)
	}

	return in.Spec.Username.Value, nil
}

// retrieve a password from the cluster.
// We allow reading from secrets, configMaps or directly from the value.
func (s *HtpasswdListener) retrievePassword(ctx context.Context, in *v1.HtpasswdUser) (string, error) {
	if in == nil {
		return "", fmt.Errorf("no user or no specification provided, skipping. (User: %v)", in)
	}

	if in.Spec.Password.SecretKeyRef != nil {
		return s.retrieveSecretKey(ctx, in.Namespace, in.Spec.Password.SecretKeyRef)
	}

	if in.Spec.Password.ConfigMapRef != nil {
		return s.retrieveConfigKey(ctx, in.Namespace, in.Spec.Password.ConfigMapRef)
	}

	return in.Spec.Password.Value, nil
}

func (s *HtpasswdListener) addUsersToConfig(ctx context.Context, users *v1.HtpasswdUserList) error {
	for _, user := range users.Items {
		if err := s.addUserToConfig(ctx, &user); err != nil {
			return err
		}
	}

	return s.syncWithVolume()
}

// addAndSyncUserToConfig adds a user to the config and writes it to file.
func (s *HtpasswdListener) addAndSyncUserToConfig(ctx context.Context, new *v1.HtpasswdUser) error {
	if err := s.addUserToConfig(ctx, new); err != nil {
		return err
	}

	return s.syncWithVolume()
}

// addUserToConfig adds a user to the config.
// if it already exists, it is overwritten.
func (s *HtpasswdListener) addUserToConfig(ctx context.Context, new *v1.HtpasswdUser) error {
	username, err := s.retrieveUsername(ctx, new)
	if err != nil {
		return err
	}

	password, err := s.retrievePassword(ctx, new)
	if err != nil {
		return err
	}

	s.authMap[string(new.UID)] = &BasicAuthCredentials{
		Username: username,
		Password: password,
	}

	return nil
}

// update an user in the config. Does write to file.
func (s *HtpasswdListener) updateUserInConfig(ctx context.Context, oldUid types.UID, new *v1.HtpasswdUser) error {
	if oldUid != new.UID {
		if err := s.removeUserInConfig(ctx, oldUid); err != nil {
			return err
		}

		return s.addAndSyncUserToConfig(ctx, new)
	}

	oldUser, ok := s.authMap[string(oldUid)]
	if !ok {
		return nil
	}

	equal, username, password, err := s.equalUser(ctx, oldUser, new)
	if err != nil {
		return err
	}

	if equal {
		return nil
	}

	s.authMap[string(new.UID)] = &BasicAuthCredentials{
		Username: username,
		Password: password,
	}

	return s.syncWithVolume()
}

func (s *HtpasswdListener) removeUserInConfig(_ context.Context, old types.UID) error {
	s.authMap[string(old)] = nil
	return s.syncWithVolume()
}

func (s *HtpasswdListener) equalUser(ctx context.Context, old *BasicAuthCredentials, new *v1.HtpasswdUser) (bool, string, string, error) {
	username, err := s.retrieveUsername(ctx, new)
	if err != nil {
		return false, "", "", err
	}

	password, err := s.retrievePassword(ctx, new)
	if err != nil {
		return false, "", "", err
	}

	return old.Password == password && old.Username == username, username, password, nil
}

func (s *HtpasswdListener) syncWithVolume() error {
	logrus.Debugf("Syncing to new volume")
	if err := s.writePasswordFile(); err != nil {
		return err
	}

	return s.signalReload()
}

// Write the current configmap synced from kubernetes to passwd file
func (s *HtpasswdListener) writePasswordFile() error {
	// Htpasswd is bcrypt hashes.
	f, err := os.Create(s.PasswordFile)
	if err != nil {
		return err
	}

	defer f.Close()

	w := bufio.NewWriter(f)
	for key, data := range s.authMap {
		if data.Username == "" {
			logrus.Debugf("Empty username for entry %v/%v, skipping", key, data)
			continue
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		userStr := fmt.Sprintf("%v:%v", data.Username, string(hash))
		_, err = fmt.Fprintln(w, userStr)
		if err != nil {
			return err
		}
	}

	return w.Flush()
}

// Reload the config if configured over pid_file.
// If no config is provided, silently return.
// Sends a SIGHUP to the process, which signals it to reload configuration files.
func (s *HtpasswdListener) signalReload() error {
	if len(s.PidFile) <= 0 {
		return nil
	}

	pid, err := os.ReadFile(s.PidFile)
	if err != nil {
		return err
	}

	pidInt, err := strconv.Atoi(string(pid))
	if err != nil {
		return err
	}

	proc, err := os.FindProcess(pidInt)
	if err != nil {
		return err
	}

	return proc.Signal(syscall.SIGHUP)
}
