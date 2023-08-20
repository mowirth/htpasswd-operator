package operator

import (
	"time"

	"github.com/bep/debounce"
	v1 "htpasswd-operator/pkg/apis/htpasswduser/v1"
)

var debounceSync = debounce.New(5 * time.Second)

// CredentialReference stores information about the used value sources (configMaps, secrets, ...)
type CredentialReference struct {
	Username *KeySource
	Password *KeySource
}

type KeySource struct {
	Namespace    string
	ConfigMapRef *v1.Reference
	SecretKeyRef *v1.Reference
	Value        string
}

type BasicAuthCredentials struct {
	Username string
	Password string
}
