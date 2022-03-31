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

package runner

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/element"
)

type mockTaskPlugin struct {
	initErr    error
	destoryErr error
}

func (m *mockTaskPlugin) Init(ctx context.Context) error {
	return m.initErr
}

func (m *mockTaskPlugin) Destroy(ctx context.Context) error {
	return m.destoryErr
}

type mockTask struct {
	*plugin.BaseTask
	prepareErr error
	postErr    error
}

func (m *mockTask) Prepare(ctx context.Context) error {
	return m.prepareErr
}

func (m *mockTask) Post(ctx context.Context) error {
	return m.postErr
}

type mockReaderTask struct {
	*mockTask
	*mockTaskPlugin
	startReadErr error
}

func newMockReaderTask(errors []error) *mockReaderTask {
	return &mockReaderTask{
		mockTaskPlugin: &mockTaskPlugin{
			initErr:    errors[0],
			destoryErr: errors[4],
		},
		mockTask: &mockTask{
			prepareErr: errors[1],
			postErr:    errors[3],
		},
		startReadErr: errors[2],
	}
}

func (m *mockReaderTask) StartRead(ctx context.Context, sender plugin.RecordSender) error {
	return m.startReadErr
}

type mockWriterTask struct {
	*mockTaskPlugin
	*mockTask
	startWriteErr error
}

func newMockWriterTask(errors []error) *mockWriterTask {
	return &mockWriterTask{
		mockTaskPlugin: &mockTaskPlugin{
			initErr:    errors[0],
			destoryErr: errors[4],
		},
		mockTask: &mockTask{
			prepareErr: errors[1],
			postErr:    errors[3],
		},
		startWriteErr: errors[2],
	}
}

func (m *mockWriterTask) StartWrite(ctx context.Context, receiver plugin.RecordReceiver) error {
	return m.startWriteErr
}

func (m *mockWriterTask) SupportFailOver() bool {
	return false
}

type mockRecordSender struct{}

func (m *mockRecordSender) CreateRecord() (element.Record, error) {
	return nil, nil
}

func (m *mockRecordSender) SendWriter(record element.Record) error {
	return nil
}

func (m *mockRecordSender) Flush() error {
	return nil
}

func (m *mockRecordSender) Terminate() error {
	return nil
}

func (m *mockRecordSender) Shutdown() error {
	return nil
}

type mockRecordReceiver struct{}

func (m *mockRecordReceiver) GetFromReader() (element.Record, error) {
	return nil, nil
}

func (m *mockRecordReceiver) Shutdown() error {
	return nil
}
