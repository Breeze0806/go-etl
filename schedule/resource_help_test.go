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
