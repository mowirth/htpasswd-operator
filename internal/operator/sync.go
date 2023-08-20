package operator

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"syscall"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// Reload the config if configured over pid_file.
// If no config is provided, silently return.
// Sends a SIGHUP to the process, which signals it to reload configuration files.
func (s *HtpasswdOperator) signalReload() error {
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

// Signal that the htpasswd file should be updated on disk
// Debounced to avoid rewriting (and restarting) the registry in case many changes are done within a small timeframe
func (s *HtpasswdOperator) signalSyncWithVolume() {
	debounceSync(func() {
		if err := s.syncWithVolume(); err != nil {
			logrus.Errorf("Failed to sync with volume: %v", err)
		}

		logrus.Debugf("Synced with volume")
	})
}

// Synchronizes the current authMap with the htpasswd file on the volume
// Sends a reload signal after writing to file to indicate that the file was updated.
func (s *HtpasswdOperator) syncWithVolume() error {
	logrus.Debugf("Syncing to volume")
	if err := s.writePasswordFile(); err != nil {
		return err
	}

	return s.signalReload()
}

// Write the current configmap synced from kubernetes to passwd file
// Access to userRefMap must be protected by calling functions
func (s *HtpasswdOperator) writePasswordFile() error {
	// Htpasswd is bcrypt hashes.
	f, err := os.Create(s.PasswordFile)
	if err != nil {
		return err
	}

	defer f.Close()

	w := bufio.NewWriter(f)
	s.passwordMapMutex.RLock()
	defer s.passwordMapMutex.RUnlock()
	for key, data := range s.passwordMap {
		if data.Username == "" || data.Password == "" {
			logrus.Tracef("Empty username or password for entry %v, skipping", key)
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
