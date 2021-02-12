package schedule

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
