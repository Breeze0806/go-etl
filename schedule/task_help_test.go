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
	"time"

	"github.com/Breeze0806/go-etl/element"
)

type mockTask struct {
	d time.Duration
}

func (m *mockTask) Do() error {
	if m.d != 0 {
		time.Sleep(m.d)
	}
	return nil
}

type mockAsyncTask struct {
	d       time.Duration
	doErr   error
	postErr error
}

func newMockAsyncTask(d time.Duration, errs []error) *mockAsyncTask {
	return &mockAsyncTask{
		d:       d,
		doErr:   errs[0],
		postErr: errs[1],
	}
}

func (m *mockAsyncTask) Do() error {
	if m.d != 0 {
		time.Sleep(m.d)
	}
	return m.doErr
}

func (m *mockAsyncTask) Post() error {
	return m.postErr
}

type mockMappedTask struct {
	taskID int64
}

func (m *mockMappedTask) Key() string {
	return element.FormatInt64(m.taskID)
}
