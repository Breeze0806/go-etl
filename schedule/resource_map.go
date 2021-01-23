package schedule

import (
	"sync"
)

type mappedResouceWrapper struct {
	resource MappedResource
	useCount int
}

type ResourceMap struct {
	mu sync.Mutex

	resources map[string]*mappedResouceWrapper
}

func NewResourceMap() *ResourceMap {
	return &ResourceMap{
		resources: make(map[string]*mappedResouceWrapper),
	}
}

func (r *ResourceMap) Get(loadOrNew func() (MappedResource, error)) (resource MappedResource, err error) {
	var newResource MappedResource
	ok := false
	if newResource, err = loadOrNew(); err != nil {
		return nil, err
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if resource, ok = r.loadLocked(newResource); ok {
		return
	}
	r.storageLocked(newResource)
	return
}

func (r *ResourceMap) Release(resource MappedResource) (err error) {
	r.mu.Lock()
	fn := r.releaseLocked(resource)
	r.mu.Unlock()
	return fn()
}

func (r *ResourceMap) UseCount(resource MappedResource) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.useCountLocked(resource)
}

func (r *ResourceMap) loadLocked(newResource MappedResource) (resource MappedResource, ok bool) {
	var rw *mappedResouceWrapper
	if rw, ok = r.resources[newResource.Key()]; ok {
		resource, rw.useCount = rw.resource, rw.useCount+1
		return
	}

	resource = newResource
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
