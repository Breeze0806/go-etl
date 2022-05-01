package db2

import (
	"testing"

	"github.com/Breeze0806/go-etl/datax/plugin/writer/rdbm"
	"github.com/Breeze0806/go-etl/storage/database"
	_ "github.com/Breeze0806/go-etl/storage/database/db2"
)

func Test_execMode(t *testing.T) {
	type args struct {
		writeMode string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{
				writeMode: database.WriteModeInsert,
			},
			want: rdbm.ExecModeNormal,
		},
		{
			name: "2",
			args: args{
				writeMode: "",
			},
			want: rdbm.ExecModeNormal,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := execMode(tt.args.writeMode); got != tt.want {
				t.Errorf("execMode() = %v, want %v", got, tt.want)
			}
		})
	}
}
