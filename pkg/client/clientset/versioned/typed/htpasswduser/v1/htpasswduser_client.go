// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"net/http"

	"github.com/mowirth/htpasswd-operator/pkg/apis/htpasswduser/v1"
	"github.com/mowirth/htpasswd-operator/pkg/client/clientset/versioned/scheme"
	"k8s.io/client-go/rest"
)

type FlangaV1Interface interface {
	RESTClient() rest.Interface
	HtpasswdUsersGetter
}

// FlangaV1Client is used to interact with features provided by the flanga.io group.
type FlangaV1Client struct {
	restClient rest.Interface
}

func (c *FlangaV1Client) HtpasswdUsers(namespace string) HtpasswdUserInterface {
	return newHtpasswdUsers(c, namespace)
}

// NewForConfig creates a new FlangaV1Client for the given config.
// NewForConfig is equivalent to NewForConfigAndClient(c, httpClient),
// where httpClient was generated with rest.HTTPClientFor(c).
func NewForConfig(c *rest.Config) (*FlangaV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	httpClient, err := rest.HTTPClientFor(&config)
	if err != nil {
		return nil, err
	}
	return NewForConfigAndClient(&config, httpClient)
}

// NewForConfigAndClient creates a new FlangaV1Client for the given config and http client.
// Note the http client provided takes precedence over the configured transport values.
func NewForConfigAndClient(c *rest.Config, h *http.Client) (*FlangaV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientForConfigAndClient(&config, h)
	if err != nil {
		return nil, err
	}
	return &FlangaV1Client{client}, nil
}

// NewForConfigOrDie creates a new FlangaV1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *FlangaV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new FlangaV1Client for the given RESTClient.
func New(c rest.Interface) *FlangaV1Client {
	return &FlangaV1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FlangaV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
