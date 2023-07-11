package v1

import (
	"context"

	"github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
)

const pluralName = "htpasswd-users"

type BasicAuthClient struct {
	restClient rest.Interface
	namespace  string
}

type HtpasswdUserInterface interface {
	Users(ctx context.Context, namespace string) HtpasswdUserQueryInterface
}

type HtpasswdUserQueryInterface interface {
	List(opts metav1.ListOptions) (*HtpasswdUserList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
}

type contextBasicAuthClient struct {
	restClient rest.Interface
	ctx        context.Context
	namespace  string
}

func NewForConfig(c *rest.Config, namespace string) (*BasicAuthClient, error) {
	if err := SchemeBuilder.AddToScheme(scheme.Scheme); err != nil {
		return nil, err
	}

	logrus.Debug("%v", scheme.Scheme)
	crdConfig := *c
	crdConfig.ContentConfig.GroupVersion = &SchemeGroupVersion
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	client, err := rest.RESTClientFor(&crdConfig)
	if err != nil {
		return nil, err
	}

	return &BasicAuthClient{restClient: client, namespace: namespace}, nil
}

func (s *BasicAuthClient) Users(ctx context.Context) HtpasswdUserQueryInterface {
	return &contextBasicAuthClient{restClient: s.restClient, ctx: ctx, namespace: s.namespace}
}

func (c *contextBasicAuthClient) List(opts metav1.ListOptions) (*HtpasswdUserList, error) {
	result := HtpasswdUserList{}
	err := c.restClient.Get().Resource(pluralName).VersionedParams(&opts, scheme.ParameterCodec).Do(c.ctx).Into(&result)

	return &result, err
}

func (c *contextBasicAuthClient) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.restClient.Get().Resource(pluralName).VersionedParams(&opts, scheme.ParameterCodec).Watch(c.ctx)
}
