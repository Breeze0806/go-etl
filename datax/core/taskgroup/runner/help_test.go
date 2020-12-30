package runner

import (
	"context"

	"github.com/Breeze0806/go-etl/datax/common/element"
	"github.com/Breeze0806/go-etl/datax/common/plugin"
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
