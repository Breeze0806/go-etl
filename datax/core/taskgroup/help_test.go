package taskgroup

import (
	"context"
	"math"
	"math/rand"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
	"github.com/Breeze0806/go-etl/datax/common/plugin/loader"
	"github.com/Breeze0806/go-etl/datax/common/spi/reader"
	"github.com/Breeze0806/go-etl/datax/common/spi/writer"
)

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
	prepareErr error
	postErr    error
}

func (m *mockTask) Prepare(ctx context.Context) error {
	return m.prepareErr
}

func (m *mockTask) Post(ctx context.Context) error {
	return m.postErr
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
}

func (m *mockReaderJob) Split(ctx context.Context, number int) ([]*config.JSON, error) {
	return nil, nil
}

type mockReaderTask struct {
	*mockTask
	*mockPlugin
	startReadErr error
}

func newMockReaderTask(errs []error) *mockReaderTask {
	return &mockReaderTask{
		mockPlugin: &mockPlugin{
			initErr:    errs[0],
			destoryErr: errs[4],
		},
		mockTask: &mockTask{
			prepareErr: errs[1],
			postErr:    errs[3],
		},
		startReadErr: errs[2],
	}
}

func (m *mockReaderTask) StartRead(ctx context.Context, sender plugin.RecordSender) error {
	return m.startReadErr
}

type mockRandReaderTask struct {
	*mockReaderTask
	rand *rand.Rand
}

func newMockRandReaderTask(errs []error) *mockRandReaderTask {
	return &mockRandReaderTask{
		mockReaderTask: newMockReaderTask(errs),
		rand:           rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (m *mockRandReaderTask) StartRead(ctx context.Context, sender plugin.RecordSender) error {
	if x := m.rand.Int31n(math.MaxInt16); x < math.MaxInt16/2 {
		return m.startReadErr
	}
	return nil
}

type mockWriterJob struct {
	*mockJob
}

func (m *mockWriterJob) Split(ctx context.Context, number int) ([]*config.JSON, error) {
	return nil, nil
}

type mockWriterTask struct {
	*mockPlugin
	*mockTask
	startWriteErr error
}

func newMockWriterTask(errs []error) *mockWriterTask {
	return &mockWriterTask{
		mockPlugin: &mockPlugin{
			initErr:    errs[0],
			destoryErr: errs[4],
		},
		mockTask: &mockTask{
			prepareErr: errs[1],
			postErr:    errs[3],
		},
		startWriteErr: errs[2],
	}
}

func (m *mockWriterTask) StartWrite(ctx context.Context, receiver plugin.RecordReceiver) error {
	return m.startWriteErr
}

func (m *mockWriterTask) SupportFailOver() bool {
	return true
}

type mockReader struct {
	errs []error
}

func newMockReader(errs []error) *mockReader {
	return &mockReader{
		errs: errs,
	}
}

func (m *mockReader) Job() reader.Job {
	return &mockReaderJob{}
}

func (m *mockReader) Task() reader.Task {
	return newMockReaderTask(m.errs)
}

type mockRandReader struct {
	errs []error
}

func newMockRandReader(errs []error) *mockRandReader {
	return &mockRandReader{
		errs: errs,
	}
}

func (m *mockRandReader) Job() reader.Job {
	return &mockReaderJob{}
}

func (m *mockRandReader) Task() reader.Task {
	return newMockRandReaderTask(m.errs)
}

type mockWriter struct {
	errs []error
}

func newMockWriter(errs []error) *mockWriter {
	return &mockWriter{
		errs: errs,
	}
}

func (m *mockWriter) Job() writer.Job {
	return &mockWriterJob{}
}

func (m *mockWriter) Task() writer.Task {
	return newMockWriterTask(m.errs)
}

func testJsonFromString(s string) *config.JSON {
	j, err := config.NewJSONFromString(s)
	if err != nil {
		panic(err)
	}
	return j
}

func resetLoader() {
	loader.UnregisterReaders()
	loader.UnregisterWriters()
}
