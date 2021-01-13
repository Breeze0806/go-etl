package mysql

import (
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/storage/database"
)

func testJsonFromString(s string) *config.Json {
	json, err := config.NewJsonFromString(s)
	if err != nil {
		panic(err)
	}
	return json
}

func TestDialect_Source(t *testing.T) {
	type args struct {
		bs *database.BaseSource
	}
	tests := []struct {
		name    string
		d       Dialect
		args    args
		want    database.Source
		wantErr bool
	}{
		{
			name: "1",
			d:    Dialect{},
			args: args{
				bs: database.NewBaseSource(testJsonFromString(`{
					"url" : "tcp(192.168.1.1:3306)/db?parseTime=false",
					"username" : "user",
					"password": "passwd"
				}`)),
			},
			want: &Source{
				BaseSource: database.NewBaseSource(testJsonFromString(`{
					"url" : "tcp(192.168.1.1:3306)/db?parseTime=false",
					"username" : "user",
					"password": "passwd"
				}`)),
				dsn: "user:passwd@tcp(192.168.1.1:3306)/db?parseTime=true",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.Source(tt.args.bs)
			if (err != nil) != tt.wantErr {
				t.Errorf("Dialect.Source() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Dialect.Source() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewSource(t *testing.T) {
	type args struct {
		bs *database.BaseSource
	}
	tests := []struct {
		name    string
		args    args
		wantS   database.Source
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				bs: database.NewBaseSource(testJsonFromString(`{
					"url" : "tcp(192.168.1.1:3306)/db?parseTime=false",
					"username" : "user",
					"password": "passwd"
				}`)),
			},
			wantS: &Source{
				BaseSource: database.NewBaseSource(testJsonFromString(`{
					"url" : "tcp(192.168.1.1:3306)/db?parseTime=false",
					"username" : "user",
					"password": "passwd"
				}`)),
				dsn: "user:passwd@tcp(192.168.1.1:3306)/db?parseTime=true",
			},
		},

		{
			name: "2",
			args: args{
				bs: database.NewBaseSource(testJsonFromString(`{
					"url" : 1,
					"username" : "user",
					"password": "passwd"
				}`)),
			},
			wantErr: true,
		},

		{
			name: "2",
			args: args{
				bs: database.NewBaseSource(testJsonFromString(`{
					"url" : "tcp(192.168.1.1:3306/db?parseTime=false",
					"username" : "user",
					"password": "passwd"
				}`)),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := NewSource(tt.args.bs)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSource() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("NewSource() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func TestNewConfig(t *testing.T) {
	type args struct {
		conf *config.Json
	}
	tests := []struct {
		name    string
		args    args
		wantC   *Config
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				conf: testJsonFromString(`{
					"url" : "1",
					"username" : "1",
					"password": "1"
				}`),
			},
			wantC: &Config{
				URL:      "1",
				Username: "1",
				Password: "1",
			},
		},
		{
			name: "1",
			args: args{
				conf: testJsonFromString(`{
					"url" : 1,
					"username" : "1",
					"password": "1"
				}`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, err := NewConfig(tt.args.conf)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotC, tt.wantC) {
				t.Errorf("NewConfig() = %v, want %v", gotC, tt.wantC)
			}
		})
	}
}

func TestConfig_FormatDSN(t *testing.T) {
	tests := []struct {
		name    string
		c       *Config
		wantDsn string
		wantErr bool
	}{
		{
			name: "1",
			c: &Config{
				URL:      "tcp(192.168.1.1:3306)/db?parseTime=false",
				Username: "user",
				Password: "passwd",
			},
			wantDsn: "user:passwd@tcp(192.168.1.1:3306)/db?parseTime=true",
		},
		{
			name: "1",
			c: &Config{
				URL:      "tcp(192.168.1.1:3306/db?parseTime=false",
				Username: "user",
				Password: "passwd",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDsn, err := tt.c.FormatDSN()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.FormatDSN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDsn != tt.wantDsn {
				t.Errorf("Config.FormatDSN() = %v, want %v", gotDsn, tt.wantDsn)
			}
		})
	}
}

func TestSource_DriverName(t *testing.T) {
	tests := []struct {
		name string
		s    *Source
		want string
	}{
		{
			name: "1",
			s:    &Source{},
			want: "mysql",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.DriverName(); got != tt.want {
				t.Errorf("Source.DriverName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSource_ConnectName(t *testing.T) {
	tests := []struct {
		name string
		s    *Source
		want string
	}{
		{
			name: "1",
			s: &Source{
				dsn: "11111xxx",
			},
			want: "11111xxx",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.ConnectName(); got != tt.want {
				t.Errorf("Source.ConnectName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSource_Table(t *testing.T) {
	type args struct {
		b *database.BaseTable
	}
	tests := []struct {
		name string
		s    *Source
		args args
		want database.Table
	}{
		{
			name: "1",
			s: &Source{
				dsn: "11111xxx",
			},
			args: args{
				b: database.NewBaseTable("db", "schema", "table"),
			},
			want: NewTable(database.NewBaseTable("db", "schema", "table")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Table(tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Source.Table() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDialect_Name(t *testing.T) {
	tests := []struct {
		name string
		d    Dialect
		want string
	}{
		{
			name: "1",
			want: "mysql",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Name(); got != tt.want {
				t.Errorf("Dialect.Name() = %v, want %v", got, tt.want)
			}
		})
	}
}
