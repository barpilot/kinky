/*
Copyright 2017 The Openshift Evangelists

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

// This file was automatically generated by lister-gen

package v1alpha1

import (
	v1alpha1 "github.com/barpilot/kinky/pkg/apis/kinky/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// KinkyLister helps list Kinkies.
type KinkyLister interface {
	// List lists all Kinkies in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.Kinky, err error)
	// Kinkies returns an object that can list and get Kinkies.
	Kinkies(namespace string) KinkyNamespaceLister
	KinkyListerExpansion
}

// kinkyLister implements the KinkyLister interface.
type kinkyLister struct {
	indexer cache.Indexer
}

// NewKinkyLister returns a new KinkyLister.
func NewKinkyLister(indexer cache.Indexer) KinkyLister {
	return &kinkyLister{indexer: indexer}
}

// List lists all Kinkies in the indexer.
func (s *kinkyLister) List(selector labels.Selector) (ret []*v1alpha1.Kinky, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Kinky))
	})
	return ret, err
}

// Kinkies returns an object that can list and get Kinkies.
func (s *kinkyLister) Kinkies(namespace string) KinkyNamespaceLister {
	return kinkyNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// KinkyNamespaceLister helps list and get Kinkies.
type KinkyNamespaceLister interface {
	// List lists all Kinkies in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.Kinky, err error)
	// Get retrieves the Kinky from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.Kinky, error)
	KinkyNamespaceListerExpansion
}

// kinkyNamespaceLister implements the KinkyNamespaceLister
// interface.
type kinkyNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Kinkies in the indexer for a given namespace.
func (s kinkyNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.Kinky, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Kinky))
	})
	return ret, err
}

// Get retrieves the Kinky from the indexer for a given namespace and name.
func (s kinkyNamespaceLister) Get(name string) (*v1alpha1.Kinky, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("kinky"), name)
	}
	return obj.(*v1alpha1.Kinky), nil
}
