package schedule

import (
	"errors"
	"sync"
	"testing"
	"time"
)

type mockMappedResource struct {
	close func() error
	key   string
}

func newMockMappedResourceNoClose(key string) *mockMappedResource {
	return newMockMappedResource(
		func() error {
			return nil
		}, key)
}

func newMockMappedResource(close func() error, key string) *mockMappedResource {
	return &mockMappedResource{
		close: close,
		key:   key,
	}
}

func (m *mockMappedResource) Close() error {
	return m.close()
}

func (m *mockMappedResource) Key() string {
	return m.key
}

func TestResourceMap_Get(t *testing.T) {
	r1 := newMockMappedResourceNoClose("mock1")
	r2 := newMockMappedResourceNoClose("mock2")
	resourceMap := NewResourceMap()

	type args struct {
		loadOrNew func() (MappedResource, error)
	}
	tests := []struct {
		name         string
		r            *ResourceMap
		args         args
		wantResource MappedResource
		wantErr      bool
	}{
		{
			name: "1",
			args: args{
				loadOrNew: func() (MappedResource, error) {
					return r1, nil
				},
			},
			wantResource: r1,
		},
		{
			name: "2",
			args: args{
				loadOrNew: func() (MappedResource, error) {
					return r2, nil
				},
			},
			wantResource: r2,
		},
		{
			name: "3",
			args: args{
				loadOrNew: func() (MappedResource, error) {
					return NewLoadMappedResource("mock1"), nil
				},
			},
			wantResource: r1,
		},
		{
			name: "4",
			args: args{
				loadOrNew: func() (MappedResource, error) {
					return NewLoadMappedResource("mock2"), nil
				},
			},
			wantResource: r2,
		},
		{
			name: "5",
			args: args{
				loadOrNew: func() (MappedResource, error) {
					return nil, errors.New("mock error")
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResource, err := resourceMap.Get(tt.args.loadOrNew)
			if (err != nil) != tt.wantErr {
				t.Errorf("ResourceMap.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotResource != tt.wantResource {
				t.Errorf("ResourceMap.Get() = %p, want %p", gotResource, tt.wantResource)
			}
		})
	}
}

func TestResourceMap_UseCount(t *testing.T) {
	resourceMap := NewResourceMap()
	type args struct {
		fn func()
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				fn: func() {
					resourceMap.Get(func() (MappedResource, error) {
						return newMockMappedResourceNoClose("mock1"), nil
					})
				},
			},
			want: 1,
		},

		{
			name: "2",
			args: args{
				fn: func() {
					resourceMap.Get(func() (MappedResource, error) {
						return NewLoadMappedResource("mock1"), nil
					})
				},
			},
			want: 2,
		},

		{
			name: "3",
			args: args{
				fn: func() {
					resourceMap.Release(NewLoadMappedResource("mock1"))
				},
			},
			want: 1,
		},

		{
			name: "4",
			args: args{
				fn: func() {
					resourceMap.Get(func() (MappedResource, error) {
						return NewLoadMappedResource("mock1"), nil
					})
				},
			},
			want: 2,
		},
		{
			name: "5",
			args: args{
				fn: func() {
					resourceMap.Get(func() (MappedResource, error) {
						return NewLoadMappedResource("mock1"), nil
					})
				},
			},
			want: 3,
		},
		{
			name: "6",
			args: args{
				fn: func() {
					resourceMap.Release(NewLoadMappedResource("mock1"))
				},
			},
			want: 2,
		},
		{
			name: "7",
			args: args{
				fn: func() {
					resourceMap.Release(NewLoadMappedResource("mock1"))
				},
			},
			want: 1,
		},
		{
			name: "8",
			args: args{
				fn: func() {
					resourceMap.Release(NewLoadMappedResource("mock1"))
				},
			},
			want: 0,
		},
		{
			name: "9",
			args: args{
				fn: func() {
					resourceMap.Release(NewLoadMappedResource("mock1"))
				},
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		tt.args.fn()
		if got := resourceMap.UseCount(NewLoadMappedResource("mock1")); got != tt.want {
			t.Errorf("run %v got: %v want: %v", tt.name, got, tt.want)
		}
	}
}

func TestResourceMap_Release(t *testing.T) {
	resourceMap := NewResourceMap()
	resourceMap.Get(func() (MappedResource, error) {
		return newMockMappedResourceNoClose("mock1"), nil
	})
	resourceMap.Get(func() (MappedResource, error) {
		return newMockMappedResource(func() error {
			return errors.New("mock error")
		}, "mock2"), nil
	})
	type args struct {
		resource MappedResource
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				resource: NewLoadMappedResource("mock1"),
			},
		},
		{
			name: "2",
			args: args{
				resource: NewLoadMappedResource("mock2"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := resourceMap.Release(tt.args.resource); (err != nil) != tt.wantErr {
				t.Errorf("ResourceMap.Release() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestResourceMap_Block(t *testing.T) {
	stop := make(chan struct{})
	var wg2 sync.WaitGroup
	resourceMap := NewResourceMap()
	resourceMap.Get(func() (MappedResource, error) {
		return newMockMappedResource(func() error {
			wg2.Done()
			select {
			case <-time.After(1 * time.Second):
			case <-stop:
			}
			return nil
		}, "mock2"), nil
	})
	var wg1 sync.WaitGroup
	wg1.Add(1)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		resourceMap.Get(func() (MappedResource, error) {
			wg1.Done()
			select {
			case <-time.After(1 * time.Second):
			case <-stop:
			}
			return newMockMappedResourceNoClose("mock1"), nil
		})
	}()
	wg.Add(1)
	wg2.Add(1)
	go func() {
		defer wg.Done()
		wg1.Wait()
		resourceMap.Release(newMockMappedResourceNoClose("mock2"))
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		wg2.Wait()
		resourceMap.Get(func() (MappedResource, error) {
			return newMockMappedResourceNoClose("mock1"), nil
		})
		resourceMap.Release(newMockMappedResourceNoClose("mock2"))
		close(stop)
	}()

	select {
	case <-time.After(1 * time.Second):
		t.Error("Block")
	case <-stop:
	}
	wg.Wait()
}
