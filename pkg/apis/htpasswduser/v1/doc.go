//go:generate bash -c "go mod download && cd ../../../.. && bash $(go list -mod=mod -m -f '{{.Dir}}' k8s.io/code-generator)/generate-groups.sh deepcopy,client,informer,lister htpasswd-operator/pkg/client htpasswd-operator/apis htpasswduser:v1 --go-header-file apis/htpasswduser/v1/boilerplate.go.txt --trim-path-prefix htpasswd-operator"
// +k8s:deepcopy-gen=package,register
// +groupName=flanga.io

package v1
