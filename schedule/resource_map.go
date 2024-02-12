// Copyright 2020 the go-etl Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package schedule

import (
	"sync"
)

type mappedResouceWrapper struct {
	resource MappedResource
	useCount int
}

// ResourceMap: Resource mapping
type ResourceMap struct {
	mu sync.Mutex

	resources map[string]*mappedResouceWrapper
}

// NewResourceMap: Create a resource mapping
func NewResourceMap() *ResourceMap {
	return &ResourceMap{
		resources: make(map[string]*mappedResouceWrapper),
	}
}

// Get: Retrieve a resource based on the keyword key. If the resource does not exist
// If there is an error creating the resource
func (r *ResourceMap) Get(key string, create func() (MappedResource, error)) (resource MappedResource, err error) {
	var ok bool
	r.mu.Lock()
	if resource, ok = r.loadLocked(key); ok {
		r.mu.Unlock()
		return
	}
	r.mu.Unlock()
	var newResource MappedResource
	if newResource, err = create(); err != nil {
		return nil, err
	}
	r.mu.Lock()
	r.storageLocked(newResource)
	r.mu.Unlock()
	resource = newResource
	return
}

// Release: Release a resource based on the resource itself. If the resource does not exist
// If there is an error creating the resource
func (r *ResourceMap) Release(resource MappedResource) (err error) {
	r.mu.Lock()
	fn := r.releaseLocked(resource)
	r.mu.Unlock()
	return fn()
}

// UseCount: Calculate the number of times a resource has been used based on the resource itself
func (r *ResourceMap) UseCount(resource MappedResource) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.useCountLocked(resource)
}

func (r *ResourceMap) loadLocked(key string) (resource MappedResource, ok bool) {
	var rw *mappedResouceWrapper
	if rw, ok = r.resources[key]; ok {
		resource, rw.useCount = rw.resource, rw.useCount+1
		return
	}
	return
}

func (r *ResourceMap) storageLocked(resource MappedResource) {
	r.resources[resource.Key()] = &mappedResouceWrapper{
		resource: resource,
		useCount: 1,
	}
}

func (r *ResourceMap) releaseLocked(resource MappedResource) func() error {
	if rw, ok := r.resources[resource.Key()]; ok {
		rw.useCount--
		if rw.useCount <= 0 {
			delete(r.resources, resource.Key())
			return func() error {
				return rw.resource.Close()
			}
		}
	}
	return func() error {
		return nil
	}
}

func (r *ResourceMap) useCountLocked(resource MappedResource) (n int) {
	if rw, ok := r.resources[resource.Key()]; ok {
		n = rw.useCount
	}
	return
}
