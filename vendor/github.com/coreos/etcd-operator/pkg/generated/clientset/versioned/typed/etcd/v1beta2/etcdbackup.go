/*
Copyright 2018 The etcd-operator Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package v1beta2

import (
	v1beta2 "github.com/coreos/etcd-operator/pkg/apis/etcd/v1beta2"
	scheme "github.com/coreos/etcd-operator/pkg/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// EtcdBackupsGetter has a method to return a EtcdBackupInterface.
// A group's client should implement this interface.
type EtcdBackupsGetter interface {
	EtcdBackups(namespace string) EtcdBackupInterface
}

// EtcdBackupInterface has methods to work with EtcdBackup resources.
type EtcdBackupInterface interface {
	Create(*v1beta2.EtcdBackup) (*v1beta2.EtcdBackup, error)
	Update(*v1beta2.EtcdBackup) (*v1beta2.EtcdBackup, error)
	UpdateStatus(*v1beta2.EtcdBackup) (*v1beta2.EtcdBackup, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1beta2.EtcdBackup, error)
	List(opts v1.ListOptions) (*v1beta2.EtcdBackupList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta2.EtcdBackup, err error)
	EtcdBackupExpansion
}

// etcdBackups implements EtcdBackupInterface
type etcdBackups struct {
	client rest.Interface
	ns     string
}

// newEtcdBackups returns a EtcdBackups
func newEtcdBackups(c *EtcdV1beta2Client, namespace string) *etcdBackups {
	return &etcdBackups{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the etcdBackup, and returns the corresponding etcdBackup object, and an error if there is any.
func (c *etcdBackups) Get(name string, options v1.GetOptions) (result *v1beta2.EtcdBackup, err error) {
	result = &v1beta2.EtcdBackup{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("etcdbackups").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of EtcdBackups that match those selectors.
func (c *etcdBackups) List(opts v1.ListOptions) (result *v1beta2.EtcdBackupList, err error) {
	result = &v1beta2.EtcdBackupList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("etcdbackups").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested etcdBackups.
func (c *etcdBackups) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("etcdbackups").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a etcdBackup and creates it.  Returns the server's representation of the etcdBackup, and an error, if there is any.
func (c *etcdBackups) Create(etcdBackup *v1beta2.EtcdBackup) (result *v1beta2.EtcdBackup, err error) {
	result = &v1beta2.EtcdBackup{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("etcdbackups").
		Body(etcdBackup).
		Do().
		Into(result)
	return
}

// Update takes the representation of a etcdBackup and updates it. Returns the server's representation of the etcdBackup, and an error, if there is any.
func (c *etcdBackups) Update(etcdBackup *v1beta2.EtcdBackup) (result *v1beta2.EtcdBackup, err error) {
	result = &v1beta2.EtcdBackup{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("etcdbackups").
		Name(etcdBackup.Name).
		Body(etcdBackup).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *etcdBackups) UpdateStatus(etcdBackup *v1beta2.EtcdBackup) (result *v1beta2.EtcdBackup, err error) {
	result = &v1beta2.EtcdBackup{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("etcdbackups").
		Name(etcdBackup.Name).
		SubResource("status").
		Body(etcdBackup).
		Do().
		Into(result)
	return
}

// Delete takes name of the etcdBackup and deletes it. Returns an error if one occurs.
func (c *etcdBackups) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("etcdbackups").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *etcdBackups) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("etcdbackups").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched etcdBackup.
func (c *etcdBackups) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta2.EtcdBackup, err error) {
	result = &v1beta2.EtcdBackup{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("etcdbackups").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
