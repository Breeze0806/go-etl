package sqlite3

import (
	"database/sql/driver"
	"github.com/Breeze0806/go-etl/storage/database"
	"github.com/mattn/go-sqlite3"
	"net"
	"testing"
)

func TestTable_Quoted(t *testing.T) {
	tests := []struct {
		name string
		tr   *Table
		want string
	}{
		{
			name: "1",
			tr:   NewTable(database.NewBaseTable("", "", "table")),
			want: "`table`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Quoted(); got != tt.want {
				t.Errorf("Table.Quoted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_String(t *testing.T) {
	tests := []struct {
		name string
		tr   *Table
		want string
	}{
		{
			name: "1",
			tr:   NewTable(database.NewBaseTable("", "", "table")),
			want: "`table`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.String(); got != tt.want {
				t.Errorf("Table.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_ShouldRetry(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		tr   *Table
		args args
		want bool
	}{
		{
			name: "1",
			tr:   NewTable(database.NewBaseTable("", "", "table")),
			args: args{
				err: nil,
			},
		},
		{
			name: "2",
			tr:   NewTable(database.NewBaseTable("", "", "table")),
			args: args{
				err: &net.AddrError{},
			},
			want: false,
		},
		{
			name: "3",
			tr:   NewTable(database.NewBaseTable("", "", "table")),
			args: args{
				err: driver.ErrBadConn,
			},
			want: true,
		},
		{
			name: "4",
			tr:   NewTable(database.NewBaseTable("", "", "table")),
			args: args{
				err: &sqlite3.Error{},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ShouldRetry(tt.args.err); got != tt.want {
				t.Errorf("Table.ShouldRetry() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTable_ShouldOneByOne(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		tr   *Table
		args args
		want bool
	}{
		{
			name: "1",
			tr:   NewTable(database.NewBaseTable("", "", "table")),
			args: args{
				err: nil,
			},
		},
		{
			name: "2",
			tr:   NewTable(database.NewBaseTable("", "", "table")),
			args: args{
				err: &net.AddrError{},
			},
		},
		{
			name: "3",
			tr:   NewTable(database.NewBaseTable("", "", "table")),
			args: args{
				err: driver.ErrBadConn,
			},
		},
		{
			name: "4",
			tr:   NewTable(database.NewBaseTable("", "", "table")),
			args: args{
				err: sqlite3.Error{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.ShouldOneByOne(tt.args.err); got != tt.want {
				t.Errorf("Table.ShouldOneByOne() = %v, want %v", got, tt.want)
			}
		})
	}
}
