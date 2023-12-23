// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	"github.com/mowirth/htpasswd-operator/pkg/apis/htpasswduser/v1"
	"github.com/mowirth/htpasswd-operator/pkg/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
)

// HtpasswdUsersGetter has a method to return a HtpasswdUserInterface.
// A group's client should implement this interface.
type HtpasswdUsersGetter interface {
	HtpasswdUsers(namespace string) HtpasswdUserInterface
}

// HtpasswdUserInterface has methods to work with HtpasswdUser resources.
type HtpasswdUserInterface interface {
	Create(ctx context.Context, htpasswdUser *v1.HtpasswdUser, opts metav1.CreateOptions) (*v1.HtpasswdUser, error)
	Update(ctx context.Context, htpasswdUser *v1.HtpasswdUser, opts metav1.UpdateOptions) (*v1.HtpasswdUser, error)
	UpdateStatus(ctx context.Context, htpasswdUser *v1.HtpasswdUser, opts metav1.UpdateOptions) (*v1.HtpasswdUser, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.HtpasswdUser, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.HtpasswdUserList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.HtpasswdUser, err error)
	HtpasswdUserExpansion
}

// htpasswdUsers implements HtpasswdUserInterface
type htpasswdUsers struct {
	client rest.Interface
	ns     string
}

// newHtpasswdUsers returns a HtpasswdUsers
func newHtpasswdUsers(c *FlangaV1Client, namespace string) *htpasswdUsers {
	return &htpasswdUsers{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the htpasswdUser, and returns the corresponding htpasswdUser object, and an error if there is any.
func (c *htpasswdUsers) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.HtpasswdUser, err error) {
	result = &v1.HtpasswdUser{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("htpasswdusers").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of HtpasswdUsers that match those selectors.
func (c *htpasswdUsers) List(ctx context.Context, opts metav1.ListOptions) (result *v1.HtpasswdUserList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.HtpasswdUserList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("htpasswdusers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested htpasswdUsers.
func (c *htpasswdUsers) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("htpasswdusers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a htpasswdUser and creates it.  Returns the server's representation of the htpasswdUser, and an error, if there is any.
func (c *htpasswdUsers) Create(ctx context.Context, htpasswdUser *v1.HtpasswdUser, opts metav1.CreateOptions) (result *v1.HtpasswdUser, err error) {
	result = &v1.HtpasswdUser{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("htpasswdusers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(htpasswdUser).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a htpasswdUser and updates it. Returns the server's representation of the htpasswdUser, and an error, if there is any.
func (c *htpasswdUsers) Update(ctx context.Context, htpasswdUser *v1.HtpasswdUser, opts metav1.UpdateOptions) (result *v1.HtpasswdUser, err error) {
	result = &v1.HtpasswdUser{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("htpasswdusers").
		Name(htpasswdUser.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(htpasswdUser).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *htpasswdUsers) UpdateStatus(ctx context.Context, htpasswdUser *v1.HtpasswdUser, opts metav1.UpdateOptions) (result *v1.HtpasswdUser, err error) {
	result = &v1.HtpasswdUser{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("htpasswdusers").
		Name(htpasswdUser.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(htpasswdUser).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the htpasswdUser and deletes it. Returns an error if one occurs.
func (c *htpasswdUsers) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("htpasswdusers").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *htpasswdUsers) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("htpasswdusers").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched htpasswdUser.
func (c *htpasswdUsers) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.HtpasswdUser, err error) {
	result = &v1.HtpasswdUser{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("htpasswdusers").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
