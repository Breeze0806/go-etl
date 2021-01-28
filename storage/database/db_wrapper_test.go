package database

import (
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/schedule"
)

func OpenNoErr(name string, conf *config.JSON) *DBWrapper {
	d, err := Open(name, conf)
	if err != nil {
		panic(err)
	}
	return d
}

func TestOpen(t *testing.T) {
	dbMap = schedule.NewResourceMap()
	registerMock()
	type args struct {
		name string
		conf *config.JSON
	}
	tests := []struct {
		name    string
		args    args
		wantDw  *DBWrapper
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				name: "mock",
				conf: testJSONFromString(`{"pool":{"connMaxIdleTime":"1","connMaxLifetime":"1"}}`),
			},
			wantErr: true,
		},
		{
			name: "2",
			args: args{
				name: "mock",
				conf: testJSONFromString("{}"),
			},
			wantDw: &DBWrapper{
				DB: &DB{
					Source: &mockSource{
						BaseSource: NewBaseSource(testJSONFromString("{}")),
						name:       "mock",
					},
				},
			},
		},
		{
			name: "3",
			args: args{
				name: "mockErr",
				conf: testJSONFromString("{}"),
			},
			wantErr: true,
		},
		{
			name: "4",
			args: args{
				name: "mock",
				conf: testJSONFromString("{}"),
			},
			wantDw: &DBWrapper{
				DB: &DB{
					Source: &mockSource{
						BaseSource: NewBaseSource(testJSONFromString("{}")),
						name:       "mock",
					},
				},
			},
		},
	}
	for _, v := range tests {
		gotDw, err := Open(v.args.name, v.args.conf)
		if (err != nil) != v.wantErr {
			t.Errorf("run %v Open() error = %v, wantErr %v", v.name, err, v.wantErr)
			continue
		}
		if gotDw == nil && v.wantDw == nil {
			t.Logf("run %v Open() = nil, wantDw nil", v.name)
			continue
		}
		if !reflect.DeepEqual(gotDw.Source, v.wantDw.Source) {
			t.Errorf("run %v Open() = %v, want %v", v.name, gotDw.Source, v.wantDw.Source)
		}
	}
}

func TestDBWrapper_Close(t *testing.T) {
	dbMap = schedule.NewResourceMap()
	registerMock()
	tests := []struct {
		name    string
		d       *DBWrapper
		want    int
		wantErr bool
	}{

		{
			name: "1",
			d:    OpenNoErr("mock", testJSONFromString("{}")),
			want: 1,
		},
		{
			name: "2",
			d:    OpenNoErr("mock", testJSONFromString("{}")),
			want: 0,
		},
		{
			name: "3",
			d:    &DBWrapper{},
			want: 0,
		},
	}
	for _, tt := range tests {

		if err := tt.d.Close(); (err != nil) != tt.wantErr {
			t.Errorf("run %v DBWrapper.Close() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			continue
		}

		if tt.d.DB == nil {
			continue
		}

		if got := dbMap.UseCount(tt.d); got != tt.want {
			t.Errorf("run %v UseCount()   = %v, want %v", tt.name, got, tt.want)
		}
	}
}
