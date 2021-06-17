package postgres

import (
	"reflect"
	"testing"

	"github.com/Breeze0806/go-etl/storage/database"
	_ "github.com/Breeze0806/go/database/pqto"
)

func TestDialect_Name(t *testing.T) {
	tests := []struct {
		name string
		d    Dialect
		want string
	}{
		{
			name: "1",
			d:    Dialect{},
			want: "postgres",
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
				bs: database.NewBaseSource(testJSONFromString(`{
					"url":"postgres://hostname:5432/db",
					"username":"user",
					"password":"password"
				}`)),
			},
			want: &Source{
				BaseSource: database.NewBaseSource(testJSONFromString(`{
					"url":"postgres://hostname:5432/db",
					"username":"user",
					"password":"password"
				}`)),
				dsn: "postgres://user:password@hostname:5432/db",
			},
		},

		{
			name: "2",
			d:    Dialect{},
			args: args{
				bs: database.NewBaseSource(testJSONFromString(`{
					"url":"postgres://hostname:5432/db",
					"username":"user",
					"password":1
				}`)),
			},
			wantErr: true,
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

func TestSource_DriverName(t *testing.T) {
	tests := []struct {
		name string
		s    *Source
		want string
	}{
		{
			name: "1",
			s: &Source{
				BaseSource: database.NewBaseSource(testJSONFromString(`{
					"url":"postgres://hostname:5432/db",
					"username":"user",
					"password":"password"
				}`)),
				dsn: "postgres://user:password@hostname:5432/db",
			},
			want: "pgTimeout",
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
				BaseSource: database.NewBaseSource(testJSONFromString(`{
					"url":"postgres://hostname:5432/db",
					"username":"user",
					"password":"password"
				}`)),
				dsn: "postgres://user:password@hostname:5432/db",
			},
			want: "postgres://user:password@hostname:5432/db",
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

func TestSource_Key(t *testing.T) {
	tests := []struct {
		name string
		s    *Source
		want string
	}{
		{
			name: "1",
			s: &Source{
				BaseSource: database.NewBaseSource(testJSONFromString(`{
					"url":"postgres://hostname:5432/db",
					"username":"user",
					"password":"password"
				}`)),
				dsn: "postgres://user:password@hostname:5432/db",
			},
			want: "postgres://user:password@hostname:5432/db",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Key(); got != tt.want {
				t.Errorf("Source.Key() = %v, want %v", got, tt.want)
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
				BaseSource: database.NewBaseSource(testJSONFromString(`{
					"url":"postgres://hostname:5432/db",
					"username":"user",
					"password":"password"
				}`)),
				dsn: "postgres://user:password@hostname:5432/db",
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
