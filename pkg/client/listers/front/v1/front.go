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

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "Crd-End/pkg/apis/front/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// FrontLister helps list Fronts.
// All objects returned here must be treated as read-only.
type FrontLister interface {
	// List lists all Fronts in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Front, err error)
	// Fronts returns an object that can list and get Fronts.
	Fronts(namespace string) FrontNamespaceLister
	FrontListerExpansion
}

// frontLister implements the FrontLister interface.
type frontLister struct {
	indexer cache.Indexer
}

// NewFrontLister returns a new FrontLister.
func NewFrontLister(indexer cache.Indexer) FrontLister {
	return &frontLister{indexer: indexer}
}

// List lists all Fronts in the indexer.
func (s *frontLister) List(selector labels.Selector) (ret []*v1.Front, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Front))
	})
	return ret, err
}

// Fronts returns an object that can list and get Fronts.
func (s *frontLister) Fronts(namespace string) FrontNamespaceLister {
	return frontNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// FrontNamespaceLister helps list and get Fronts.
// All objects returned here must be treated as read-only.
type FrontNamespaceLister interface {
	// List lists all Fronts in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.Front, err error)
	// Get retrieves the Front from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.Front, error)
	FrontNamespaceListerExpansion
}

// frontNamespaceLister implements the FrontNamespaceLister
// interface.
type frontNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Fronts in the indexer for a given namespace.
func (s frontNamespaceLister) List(selector labels.Selector) (ret []*v1.Front, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Front))
	})
	return ret, err
}

// Get retrieves the Front from the indexer for a given namespace and name.
func (s frontNamespaceLister) Get(name string) (*v1.Front, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("front"), name)
	}
	return obj.(*v1.Front), nil
}
