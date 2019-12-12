/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1alpha1 "stash.us.cray.com/dpm/dws-operator/pkg/apis/dws/v1alpha1"
)

// FakeStoragePoolSpecs implements StoragePoolSpecInterface
type FakeStoragePoolSpecs struct {
	Fake *FakeDwsV1alpha1
	ns   string
}

var storagepoolspecsResource = schema.GroupVersionResource{Group: "dws.cray.com", Version: "v1alpha1", Resource: "storagepoolspecs"}

var storagepoolspecsKind = schema.GroupVersionKind{Group: "dws.cray.com", Version: "v1alpha1", Kind: "StoragePoolSpec"}

// Get takes name of the storagePoolSpec, and returns the corresponding storagePoolSpec object, and an error if there is any.
func (c *FakeStoragePoolSpecs) Get(name string, options v1.GetOptions) (result *v1alpha1.StoragePoolSpec, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(storagepoolspecsResource, c.ns, name), &v1alpha1.StoragePoolSpec{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.StoragePoolSpec), err
}

// List takes label and field selectors, and returns the list of StoragePoolSpecs that match those selectors.
func (c *FakeStoragePoolSpecs) List(opts v1.ListOptions) (result *v1alpha1.StoragePoolSpecList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(storagepoolspecsResource, storagepoolspecsKind, c.ns, opts), &v1alpha1.StoragePoolSpecList{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.StoragePoolSpecList), err
}

// Watch returns a watch.Interface that watches the requested storagePoolSpecs.
func (c *FakeStoragePoolSpecs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(storagepoolspecsResource, c.ns, opts))

}

// Create takes the representation of a storagePoolSpec and creates it.  Returns the server's representation of the storagePoolSpec, and an error, if there is any.
func (c *FakeStoragePoolSpecs) Create(storagePoolSpec *v1alpha1.StoragePoolSpec) (result *v1alpha1.StoragePoolSpec, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(storagepoolspecsResource, c.ns, storagePoolSpec), &v1alpha1.StoragePoolSpec{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.StoragePoolSpec), err
}

// Update takes the representation of a storagePoolSpec and updates it. Returns the server's representation of the storagePoolSpec, and an error, if there is any.
func (c *FakeStoragePoolSpecs) Update(storagePoolSpec *v1alpha1.StoragePoolSpec) (result *v1alpha1.StoragePoolSpec, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(storagepoolspecsResource, c.ns, storagePoolSpec), &v1alpha1.StoragePoolSpec{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.StoragePoolSpec), err
}

// Delete takes name of the storagePoolSpec and deletes it. Returns an error if one occurs.
func (c *FakeStoragePoolSpecs) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(storagepoolspecsResource, c.ns, name), &v1alpha1.StoragePoolSpec{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeStoragePoolSpecs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(storagepoolspecsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.StoragePoolSpecList{})
	return err
}

// Patch applies the patch and returns the patched storagePoolSpec.
func (c *FakeStoragePoolSpecs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.StoragePoolSpec, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(storagepoolspecsResource, c.ns, name, pt, data, subresources...), &v1alpha1.StoragePoolSpec{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.StoragePoolSpec), err
}
