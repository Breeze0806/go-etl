package schedule

type Resource interface {
	Close() error
}

type MappedResource interface {
	Resource

	Key() string
}

type LoadMappedResource struct {
	key string
}

func NewLoadMappedResource(key string) *LoadMappedResource {
	return &LoadMappedResource{
		key: key,
	}
}

func (l *LoadMappedResource) Close() (err error) {
	return
}

func (l *LoadMappedResource) Key() string {
	return l.key
}
