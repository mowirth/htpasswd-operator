//go:build tools
// +build tools

package htpasswd_kubernetes_operator

import (
	_ "k8s.io/apimachinery"
	_ "k8s.io/code-generator"
)
