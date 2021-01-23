package schedule

import "testing"

func TestLoadMappedResource_Close(t *testing.T) {
	tests := []struct {
		name    string
		l       *LoadMappedResource
		wantErr bool
	}{
		{
			name: "1",
			l:    NewLoadMappedResource("load"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.l.Close(); (err != nil) != tt.wantErr {
				t.Errorf("LoadResource.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoadMappedResource_Key(t *testing.T) {
	tests := []struct {
		name string
		l    *LoadMappedResource
		want string
	}{
		{
			name: "1",
			l:    NewLoadMappedResource("load"),
			want: "load",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.l.Key(); got != tt.want {
				t.Errorf("LoadResource.Key() = %v, wantErr %v", got, tt.want)
			}
		})
	}
}
