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

package v1alpha

import (
	v1alpha "study/pkg/apis/runtime/v1alpha"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// MeshServerLister helps list MeshServers.
// All objects returned here must be treated as read-only.
type MeshServerLister interface {
	// List lists all MeshServers in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha.MeshServer, err error)
	// MeshServers returns an object that can list and get MeshServers.
	MeshServers(namespace string) MeshServerNamespaceLister
	MeshServerListerExpansion
}

// meshServerLister implements the MeshServerLister interface.
type meshServerLister struct {
	indexer cache.Indexer
}

// NewMeshServerLister returns a new MeshServerLister.
func NewMeshServerLister(indexer cache.Indexer) MeshServerLister {
	return &meshServerLister{indexer: indexer}
}

// List lists all MeshServers in the indexer.
func (s *meshServerLister) List(selector labels.Selector) (ret []*v1alpha.MeshServer, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha.MeshServer))
	})
	return ret, err
}

// MeshServers returns an object that can list and get MeshServers.
func (s *meshServerLister) MeshServers(namespace string) MeshServerNamespaceLister {
	return meshServerNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// MeshServerNamespaceLister helps list and get MeshServers.
// All objects returned here must be treated as read-only.
type MeshServerNamespaceLister interface {
	// List lists all MeshServers in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha.MeshServer, err error)
	// Get retrieves the MeshServer from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha.MeshServer, error)
	MeshServerNamespaceListerExpansion
}

// meshServerNamespaceLister implements the MeshServerNamespaceLister
// interface.
type meshServerNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all MeshServers in the indexer for a given namespace.
func (s meshServerNamespaceLister) List(selector labels.Selector) (ret []*v1alpha.MeshServer, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha.MeshServer))
	})
	return ret, err
}

// Get retrieves the MeshServer from the indexer for a given namespace and name.
func (s meshServerNamespaceLister) Get(name string) (*v1alpha.MeshServer, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha.Resource("meshserver"), name)
	}
	return obj.(*v1alpha.MeshServer), nil
}
