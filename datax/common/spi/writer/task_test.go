package writer

import "testing"

func TestBaseTask_SupportFailOver(t *testing.T) {
	tests := []struct {
		name string
		b    *BaseTask
		want bool
	}{
		{
			name: "1",
			b:    NewBaseTask(),
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.b.SupportFailOver(); got != tt.want {
				t.Errorf("BaseTask.SupportFailOver() = %v, want %v", got, tt.want)
			}
		})
	}
}
