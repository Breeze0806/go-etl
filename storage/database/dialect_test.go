package database

import (
	"testing"
)

func TestRegisterDialect(t *testing.T) {
	UnregisterAllDialects()
	d1 := &mockDialect{
		name: "nil",
	}
	type args struct {
		name    string
		dialect Dialect
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantOk  bool
		want    Dialect
	}{
		{
			name: "1",
			args: args{
				name:    "nil",
				dialect: nil,
			},
			wantErr: true,
			wantOk:  false,
			want:    nil,
		},
		{
			name: "2",
			args: args{
				name:    "nil",
				dialect: d1,
			},
			wantOk: true,
			want:   d1,
		},
		{
			name: "3",
			args: args{
				name:    "nil",
				dialect: &mockDialect{},
			},
			wantErr: true,
			wantOk:  true,
			want:    d1,
		},
	}

	for _, tt := range tests {
		run := func() (err error) {
			defer func() {
				if perr := recover(); perr != nil {
					err = perr.(error)
				}
			}()
			RegisterDialect(tt.args.name, tt.args.dialect)
			return
		}
		err := run()
		if (err != nil) != tt.wantErr {
			t.Errorf("run %v RegisterDialect() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			return
		}

		got, gotOk := dialects.dialect(tt.args.name)
		if gotOk != tt.wantOk {
			t.Errorf("run %v dialects.dialect() gotOk = %v, wantOk %v", tt.name, gotOk, tt.wantOk)
			return
		}
		if got != tt.want {
			t.Errorf("run %v dialects.dialect() got = %v, want %v", tt.name, got, tt.want)
		}

	}
}

func TestUnregisterAllDialects(t *testing.T) {
	UnregisterAllDialects()
	RegisterDialect("nil", &mockDialect{})
	if len(dialects.dialects) == 0 {
		t.Errorf("dialects is empty")
		return
	}
	UnregisterAllDialects()
	if len(dialects.dialects) != 0 {
		t.Errorf("dialects is not empty")
		return
	}
}
