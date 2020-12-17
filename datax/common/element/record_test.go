package element

import (
	"testing"
)

func TestDefaultRecord(t *testing.T) {
	r := NewDefaultRecord()
	type args struct {
		c Column
	}
	tests := []struct {
		r       *DefaultRecord
		args    args
		wantErr bool
	}{
		{
			r: r,
			args: args{
				NewDefaultColumn(NewNilBigIntColumnValue(), "test", 0),
			},
		},
		{
			r: r,
			args: args{
				NewDefaultColumn(NewNilBigIntColumnValue(), "test", 0),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		if err := tt.r.Add(tt.args.c); (err != nil) != tt.wantErr {
			t.Errorf("DefaultRecord.Add() error = %v, wantErr %v", err, tt.wantErr)
		}
	}

	_, err := r.GetByIndex(0)
	if (err != nil) != false {
		t.Errorf("DefaultRecord.GetByIndex() error = %v, wantErr true", err)
		return
	}

	_, err = r.GetByIndex(1)
	if (err != nil) != true {
		t.Errorf("DefaultRecord.GetByIndex() error = %v, wantErr false", err)
		return
	}

	_, err = r.GetByName("test")
	if (err != nil) != false {
		t.Errorf("DefaultRecord.GetByName() error = %v, wantErr true", err)
		return
	}

	_, err = r.GetByName("")
	if (err != nil) != true {
		t.Errorf("DefaultRecord.GetByName() error = %v, wantErr false", err)
		return
	}

	s := r.ByteSize()
	if s != 0 {
		t.Errorf("DefaultRecord.ByteSize() = %v, want 0", s)
		return
	}

	s = r.MemorySize()
	if s != 8 {
		t.Errorf("DefaultRecord.ByteSize() = %v, want 8", s)
		return
	}
}
