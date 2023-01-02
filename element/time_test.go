package element

import (
	"reflect"
	"testing"
	"time"
)

func TestStringTimeEncoder_TimeEncode(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		name    string
		e       *StringTimeEncoder
		args    args
		want    time.Time
		wantErr bool
	}{
		{
			name: "1",
			e:    NewStringTimeEncoder(DefaultTimeFormat).(*StringTimeEncoder),
			args: args{
				i: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.e.TimeEncode(tt.args.i)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringTimeEncoder.TimeEncode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StringTimeEncoder.TimeEncode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringTimeDecoder_Layout(t *testing.T) {
	tests := []struct {
		name string
		d    *StringTimeDecoder
		want string
	}{
		{
			name: "1",
			d:    NewStringTimeDecoder(DefaultTimeFormat).(*StringTimeDecoder),
			want: DefaultTimeFormat,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.d.Layout(); got != tt.want {
				t.Errorf("StringTimeDecoder.Layout() = %v, want %v", got, tt.want)
			}
		})
	}
}
