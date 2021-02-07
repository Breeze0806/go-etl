package mysql

import (
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/element"
)

func Test_tableParam_Query(t *testing.T) {
	type args struct {
		in0 []element.Record
	}
	tests := []struct {
		name    string
		t       *tableParam
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			t:    newTableParam(newParameter(&paramConfig{}, &mockExecer{})),
			args: args{
				in0: nil,
			},
			wantErr: true,
		},
		{
			name: "2",
			t: newTableParam(newParameter(&paramConfig{
				Column: []string{
					"f1", "f2", "f3",
				},
				Connection: connConfig{
					Table: tableConfig{
						Db:   "db",
						Name: "table",
					},
				},
			}, &mockExecer{})),
			args: args{
				in0: nil,
			},
			want: "select f1,f2,f3 from db.table where 1 = 2",
		},
		{
			name: "3",
			t: newTableParam(newParameter(&paramConfig{
				Column: []string{
					"f1",
				},
				Connection: connConfig{
					Table: tableConfig{
						Db:   "db",
						Name: "table",
					},
				},
			}, &mockExecer{})),
			args: args{
				in0: nil,
			},
			want: "select f1 from db.table where 1 = 2",
		},
		{
			name: "4",
			t: newTableParam(newParameter(&paramConfig{
				Column: []string{
					"*",
				},
				Connection: connConfig{
					Table: tableConfig{
						Db:   "db",
						Name: "table",
					},
				},
			}, &mockExecer{})),
			args: args{
				in0: nil,
			},
			want: "select * from db.table where 1 = 2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.t.Query(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("tableParam.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("tableParam.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tableParam_Agrs(t *testing.T) {
	type args struct {
		in0 []element.Record
	}
	tests := []struct {
		name    string
		t       *tableParam
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name: "1",
			t:    newTableParam(newParameter(&paramConfig{}, &mockExecer{})),
			args: args{
				in0: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.t.Agrs(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("tableParam.Agrs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tableParam.Agrs() = %v, want %v", got, tt.want)
			}
		})
	}
}
