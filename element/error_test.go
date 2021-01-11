package element

import (
	"fmt"
	"strings"
	"testing"
)

func TestTransformError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *TransformError
		want string
	}{
		{
			name: "1",
			e:    NewTransformErrorFormColumnTypes(TypeString, TypeBigInt, fmt.Errorf("test")),
			want: "test",
		},
		{
			name: "2",
			e:    NewTransformErrorFormColumnTypes(TypeString, TypeBigInt, (NewTransformErrorFormColumnTypes(TypeString, TypeBigInt, fmt.Errorf("test1")))),
			want: "test1",
		},

		{
			name: "3",
			e:    NewTransformErrorFormColumnTypes(TypeString, TypeBigInt, nil),
			want: "transform",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); !strings.Contains(got, tt.want) {
				t.Errorf("TransformError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetError_Error(t *testing.T) {
	tests := []struct {
		name string
		e    *SetError
		want string
	}{
		{
			name: "1",
			e:    NewSetError(TypeString, TypeBigInt, fmt.Errorf("test")),
			want: "test",
		},
		{
			name: "2",
			e:    NewSetError(TypeString, TypeBigInt, NewSetError(TypeString, TypeBigInt, fmt.Errorf("test1"))),
			want: "test1",
		},
		{
			name: "3",
			e:    NewSetError(TypeString, TypeBigInt, NewSetError(TypeString, TypeBigInt, nil)),
			want: "set",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.e.Error(); !strings.Contains(got, tt.want) {
				t.Errorf("SetError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}
