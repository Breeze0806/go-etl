package transform

import (
	"testing"

	"github.com/Breeze0806/go-etl/element"
)

func TestNilTransformer_DoTransform(t *testing.T) {
	r := element.NewDefaultRecord()
	type args struct {
		record element.Record
	}
	tests := []struct {
		name    string
		n       *NilTransformer
		args    args
		want    element.Record
		wantErr bool
	}{
		{
			name: "1",
			n:    &NilTransformer{},
			args: args{
				record: r,
			},
			want:    r,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.n.DoTransform(tt.args.record)
			if (err != nil) != tt.wantErr {
				t.Errorf("NilTransformer.DoTransform() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("NilTransformer.DoTransform() = %v, want %v", got, tt.want)
			}
		})
	}
}
