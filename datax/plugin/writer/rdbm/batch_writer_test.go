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

package rdbm

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Breeze0806/go-etl/datax/common/plugin"
	spiwriter "github.com/Breeze0806/go-etl/datax/common/spi/writer"
	"github.com/Breeze0806/go-etl/datax/core/transport/exchange"
)

func newMockBatchWriter(execer Execer, mode string) *BaseBatchWriter {
	return NewBaseBatchWriter(&Task{
		BaseTask: spiwriter.NewBaseTask(),
		Execer:   execer,
		Config:   &BaseConfig{},
	}, mode, nil)
}

func TestStartWrite(t *testing.T) {
	type args struct {
		ctx      context.Context
		writer   BatchWriter
		receiver plugin.RecordReceiver
	}
	tests := []struct {
		name    string
		args    args
		waitCtx time.Duration
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiver(1000, exchange.ErrTerminate, 1*time.Millisecond),
				writer:   newMockBatchWriter(&MockExecer{}, ""),
			},
		},
		{
			name: "2",
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiverWithoutWait(10000, exchange.ErrTerminate),
				writer:   newMockBatchWriter(&MockExecer{}, "Tx"),
			},
		},
		{
			name: "3",
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiverWithoutWait(10000, errors.New("mock error")),
				writer:   newMockBatchWriter(&MockExecer{}, "Stmt"),
			},
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiverWithoutWait(10000, errors.New("mock error")),
				writer:   newMockBatchWriter(&MockExecer{}, "StmtTx"),
			},
			waitCtx: 5 * time.Microsecond,
			wantErr: false,
		},
		{
			name: "5",
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiverWithoutWait(10000, errors.New("mock error")),
				writer:   newMockBatchWriter(&MockExecer{}, ""),
			},
			waitCtx: 5 * time.Microsecond,
			wantErr: false,
		},
		{
			name: "6",
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiver(1000, exchange.ErrTerminate, 1*time.Millisecond),
				writer: newMockBatchWriter(&MockExecer{
					BatchErr: errors.New("mock error"),
					BatchN:   1,
				}, ""),
			},
			wantErr: true,
		},
		{
			name: "7",
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiverWithoutWait(10000, exchange.ErrTerminate),
				writer: newMockBatchWriter(&MockExecer{
					BatchErr: errors.New("mock error"),
					BatchN:   1,
				}, ""),
			},
			wantErr: true,
		},
		{
			name: "8",
			args: args{
				ctx:      context.TODO(),
				receiver: NewMockReceiver(2, exchange.ErrTerminate, 1*time.Millisecond),
				writer: newMockBatchWriter(&MockExecer{
					BatchErr: errors.New("mock error"),
					BatchN:   0,
				}, ""),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(tt.args.ctx)
			defer cancel()
			if tt.waitCtx != 0 {
				go func() {
					<-time.After(tt.waitCtx)
					cancel()
				}()
			}

			if err := StartWrite(ctx, tt.args.writer, tt.args.receiver); (err != nil) != tt.wantErr {
				t.Errorf("StartWrite() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
