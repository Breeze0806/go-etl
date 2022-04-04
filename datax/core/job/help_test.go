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

package job

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
	"github.com/Breeze0806/go-etl/datax/common/spi/reader"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
)

func testJSONFromString(s string) *config.JSON {
	j, err := config.NewJSONFromString(s)
	if err != nil {
		panic(err)
	}
	return j
}

func testContainer(conf *config.JSON) *Container {
	c, err := NewContainer(context.TODO(), conf)
	if err != nil {
		panic(err)
	}
	return c
}

func resetLoader() {
	loader.UnregisterReaders()
	loader.UnregisterWriters()
}

type mockPlugin struct {
	initErr    error
	destoryErr error
}

func (m *mockPlugin) Init(ctx context.Context) error {
	return m.initErr
}

func (m *mockPlugin) Destroy(ctx context.Context) error {
	return m.destoryErr
}

type mockTask struct {
	*plugin.BaseTask
	*mockPlugin
}

func (m *mockTask) Prepare(ctx context.Context) error {
	return nil
}

func (m *mockTask) Post(ctx context.Context) error {
	return nil
}

type mockJob struct {
	*plugin.BaseJob
	*mockPlugin
	prepareErr error
	postErr    error
}

func (m *mockJob) Prepare(ctx context.Context) error {
	return m.prepareErr
}

func (m *mockJob) Post(ctx context.Context) error {
	return m.postErr
}

type mockReaderJob struct {
	*mockJob
	splitErr error
	confs    []*config.JSON
}

func newMockReaderJob(errs []error, confs []*config.JSON) *mockReaderJob {
	return &mockReaderJob{
		mockJob: &mockJob{
			BaseJob: plugin.NewBaseJob(),
			mockPlugin: &mockPlugin{
				initErr:    errs[0],
				destoryErr: errs[4],
			},
			prepareErr: errs[1],
			postErr:    errs[3],
		},
		splitErr: errs[2],
		confs:    confs,
	}
}

func (m *mockReaderJob) Split(ctx context.Context, number int) ([]*config.JSON, error) {
	return m.confs, m.splitErr
}

type mockReaderTask struct {
	*mockTask
	*mockPlugin
}

func newMockReaderTask() *mockReaderTask {
	return &mockReaderTask{
		mockPlugin: &mockPlugin{},
		mockTask: &mockTask{
			BaseTask: plugin.NewBaseTask(),
		},
	}
}

func (m *mockReaderTask) StartRead(ctx context.Context, sender plugin.RecordSender) error {
	return nil
}

type mockWriterJob struct {
	*mockJob
	splitErr error
	confs    []*config.JSON
}

func newMockWriterJob(errs []error, confs []*config.JSON) *mockWriterJob {
	return &mockWriterJob{
		mockJob: &mockJob{
			BaseJob: plugin.NewBaseJob(),
			mockPlugin: &mockPlugin{
				initErr:    errs[0],
				destoryErr: errs[4],
			},
			prepareErr: errs[1],
			postErr:    errs[3],
		},
		splitErr: errs[2],
		confs:    confs,
	}
}

func (m *mockWriterJob) Split(ctx context.Context, number int) ([]*config.JSON, error) {
	return m.confs, m.splitErr
}

type mockWriterTask struct {
	*mockPlugin
	*mockTask
}

func newMockWriterTask() *mockWriterTask {
	return &mockWriterTask{
		mockPlugin: &mockPlugin{},
		mockTask: &mockTask{
			BaseTask: plugin.NewBaseTask(),
		},
	}
}

func (m *mockWriterTask) StartWrite(ctx context.Context, receiver plugin.RecordReceiver) error {
	return nil
}

func (m *mockWriterTask) SupportFailOver() bool {
	return false
}

type mockReader struct {
	errs  []error
	confs []*config.JSON
}

func newMockReader(errs []error, confs []*config.JSON) *mockReader {
	return &mockReader{
		errs:  errs,
		confs: confs,
	}
}

func (m *mockReader) Job() reader.Job {
	return newMockReaderJob(m.errs, m.confs)
}

func (m *mockReader) Task() reader.Task {
	return newMockReaderTask()
}

type mockRandReader struct {
	errs []error
}

type mockWriter struct {
	errs  []error
	confs []*config.JSON
}

func newMockWriter(errs []error, confs []*config.JSON) *mockWriter {
	return &mockWriter{
		errs:  errs,
		confs: confs,
	}
}

func (m *mockWriter) Job() writer.Job {
	return newMockReaderJob(m.errs, m.confs)
}

func (m *mockWriter) Task() writer.Task {
	return newMockWriterTask()
}

func equalConfigJSON(gotConfig, wantConfig *config.JSON) bool {
	var got, want interface{}
	err := json.Unmarshal([]byte(gotConfig.String()), &got)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal([]byte(wantConfig.String()), &want)
	if err != nil {
		panic(err)
	}
	return reflect.DeepEqual(got, want)
}
