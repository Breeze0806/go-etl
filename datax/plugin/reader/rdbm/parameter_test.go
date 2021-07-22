package rdbm

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
		t       *TableParam
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			t:    NewTableParam(&BaseConfig{}, &MockQuerier{}, nil),
			args: args{
				in0: nil,
			},
			wantErr: true,
		},
		{
			name: "2",
			t: NewTableParam(&BaseConfig{
				Column: []string{
					"f1", "f2", "f3",
				},
				Connection: ConnConfig{
					Table: TableConfig{
						Db:   "db",
						Name: "table",
					},
				},
			}, &MockQuerier{}, nil),
			args: args{
				in0: nil,
			},
			want: "select f1,f2,f3 from db.table where 1 = 2",
		},
		{
			name: "3",
			t: NewTableParam(&BaseConfig{
				Column: []string{
					"f1",
				},
				Connection: ConnConfig{
					Table: TableConfig{
						Db:   "db",
						Name: "table",
					},
				},
			}, &MockQuerier{}, nil),
			args: args{
				in0: nil,
			},
			want: "select f1 from db.table where 1 = 2",
		},
		{
			name: "4",
			t: NewTableParam(&BaseConfig{
				Column: []string{
					"*",
				},
				Connection: ConnConfig{
					Table: TableConfig{
						Db:   "db",
						Name: "table",
					},
				},
			}, &MockQuerier{}, nil),
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
		t       *TableParam
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name: "1",
			t:    NewTableParam(&BaseConfig{}, &MockQuerier{}, nil),
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

func Test_queryParam_Query(t *testing.T) {
	type args struct {
		in0 []element.Record
	}
	tests := []struct {
		name    string
		q       *QueryParam
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1",
			q:    NewQueryParam(&BaseConfig{}, &MockQuerier{}, nil),
			args: args{
				in0: nil,
			},
			wantErr: true,
		},
		{
			name: "2",
			q: NewQueryParam(&BaseConfig{
				Column: []string{
					"f1", "f2", "f3",
				},
				Connection: ConnConfig{
					Table: TableConfig{
						Db:   "db",
						Name: "table",
					},
				},
			}, &MockQuerier{}, nil),
			args: args{
				in0: nil,
			},
			want: "select f1,f2,f3 from db.table",
		},
		{
			name: "3",
			q: NewQueryParam(&BaseConfig{
				Column: []string{
					"f1",
				},
				Connection: ConnConfig{
					Table: TableConfig{
						Db:   "db",
						Name: "table",
					},
				},
			}, &MockQuerier{}, nil),
			args: args{
				in0: nil,
			},
			want: "select f1 from db.table",
		},
		{
			name: "3",
			q: NewQueryParam(&BaseConfig{
				Column: []string{
					"f1",
				},
				Connection: ConnConfig{
					Table: TableConfig{
						Db:   "db",
						Name: "table",
					},
				},
				Where: "a <> 1",
			}, &MockQuerier{}, nil),
			args: args{
				in0: nil,
			},
			want: "select f1 from db.table where a <> 1",
		},

		{
			name: "4",
			q: NewQueryParam(&BaseConfig{
				Column: []string{
					"*",
				},
				Connection: ConnConfig{
					Table: TableConfig{
						Db:   "db",
						Name: "table",
					},
				},
				Where: "a <> 1",
			}, &MockQuerier{}, nil),
			args: args{
				in0: nil,
			},
			want: "select * from db.table where a <> 1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.q.Query(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("queryParam.Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("queryParam.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_queryParam_Agrs(t *testing.T) {
	type args struct {
		in0 []element.Record
	}
	tests := []struct {
		name    string
		q       *QueryParam
		args    args
		want    []interface{}
		wantErr bool
	}{
		{
			name: "1",
			q:    NewQueryParam(&BaseConfig{}, &MockQuerier{}, nil),
			args: args{
				in0: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.q.Agrs(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("queryParam.Agrs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryParam.Agrs() = %v, want %v", got, tt.want)
			}
		})
	}
}
