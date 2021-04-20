package schedule

import (
	"strconv"
	"time"
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
	return strconv.FormatInt(m.taskID, 10)
}
