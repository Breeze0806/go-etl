package plugin

import (
	"testing"
)

func TestType_IsValid(t *testing.T) {
	tests := []struct {
		name string
		t    Type
		want bool
	}{
		{
			name: "1",
			t:    Writer,
			want: true,
		},
		{
			name: "2",
			t:    "",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.IsValid(); got != tt.want {
				t.Errorf("Type.IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestType_String(t *testing.T) {
	tests := []struct {
		name string
		t    Type
		want string
	}{
		{
			name: "1",
			t:    Writer,
			want: string(Writer),
		},
		{
			name: "2",
			t:    NewType(""),
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.want {
				t.Errorf("Type.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
