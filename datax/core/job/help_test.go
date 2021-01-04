package job

import (
	"context"
	"encoding/json"
	"reflect"

	"github.com/Breeze0806/go-etl/datax/common/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
	"github.com/Breeze0806/go-etl/datax/common/spi/reader"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
)

func testJsonFromString(s string) *config.Json {
	j, err := config.NewJsonFromString(s)
	if err != nil {
		panic(err)
	}
	return j
}

func testContainer(conf *config.Json) *Container {
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
	confs    []*config.Json
}

func newMockReaderJob(errs []error, confs []*config.Json) *mockReaderJob {
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

func (m *mockReaderJob) Split(ctx context.Context, number int) ([]*config.Json, error) {
	return m.confs, m.splitErr
}

type mockReaderTask struct {
	*mockTask
	*mockPlugin
}

func newMockReaderTask() *mockReaderTask {
	return &mockReaderTask{
		mockPlugin: &mockPlugin{},
		mockTask:   &mockTask{},
	}
}

func (m *mockReaderTask) StartRead(ctx context.Context, sender plugin.RecordSender) error {
	return nil
}

type mockWriterJob struct {
	*mockJob
	splitErr error
	confs    []*config.Json
}

func newMockWriterJob(errs []error, confs []*config.Json) *mockWriterJob {
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

func (m *mockWriterJob) Split(ctx context.Context, number int) ([]*config.Json, error) {
	return m.confs, m.splitErr
}

type mockWriterTask struct {
	*mockPlugin
	*mockTask
}

func newMockWriterTask() *mockWriterTask {
	return &mockWriterTask{
		mockPlugin: &mockPlugin{},
		mockTask:   &mockTask{},
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
	confs []*config.Json
}

func newMockReader(errs []error, confs []*config.Json) *mockReader {
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
	confs []*config.Json
}

func newMockWriter(errs []error, confs []*config.Json) *mockWriter {
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

func equalConfigJson(gotConfig, wantConfig *config.Json) bool {
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
