package config

import (
	"reflect"
	"strings"
	"testing"

	"github.com/Breeze0806/go/encoding"
)

var basicJSON = `{
	"age": 100,
	"place": {
		"here": "B\\\"R"
	},
	"noop": {
		"what is a wren?": "a bird"
	},
	"happy": true,
	"immortal": false,
	"items": [1, 2, 3, {
		"tags": [1, 2, 3],
		"points": [
			[1, 2],
			[3, 4]
		]
	}, 4, 5, 6, 7],
	"arr": ["1", 2, "3", {
		"hello": "world"
	}, "4", 5],
	"vals": [1, 2, 3, {
		"sadf": "asdf"
	}],
	"name": {
		"first": "tom",
		"last": null
	},
	"created": "2014-05-16T08:28:06.989Z",
	"loggy": {
		"programmers": [{
				"firstName": "Brett",
				"lastName": "McLaughlin",
				"email": "aaaa",
				"tag": "good"
			},
			{
				"firstName": "Jason",
				"lastName": "Hunter",
				"email": "bbbb",
				"tag": "bad"
			},
			{
				"firstName": "Elliotte",
				"lastName": "Harold",
				"email": "cccc",
				"tag": "good"
			},
			{
				"firstName": 1002.3,
				"age": 101
			}
		]
	},
	"lastly": {
		"yay": "final"
	},
	"float": 1e1000
}`

var invlaidJSON = `{
	"age": 100,
	"place": {
		"here": "B\\\"R"
	},
	"noop": {
		"what is a wren?": ,"a bird"
	},
	"happy": true,
	"immortal": false,
	"items": [1, 2, 3, {
		"tags": [1, 2, 3],
		"points": [
			[1, 2],
			[3, 4]
		]
	}, 4, 5, 6, 7],
	"arr": ["1", 2, "3", {
		"hello": "world"
	}, "4", 5],
	"vals": [1, 2, 3, {
		"sadf": "asdf"
	}],
	"name": {
		"first": "tom",
		"last": null
	},
	"created": "2014-05-16T08:28:06.989Z",
	"loggy": {
		"programmers": [{
				"firstName": "Brett",
				"lastName": "McLaughlin",
				"email": "aaaa",
				"tag": "good"
			},
			{
				"firstName": "Jason",
				"lastName": "Hunter",
				"email": "bbbb",
				"tag": "bad"
			},
			{
				"firstName": "Elliotte",
				"lastName": "Harold",
				"email": "cccc",
				"tag": "good"
			},
			{
				"firstName": 1002.3,
				"age": 101
			}
		]
	},
	"lastly": {
		"yay": "final"
	}
}`

func testJSONFromString(s string) *JSON {
	return &JSON{
		JSON: testEncodingJSONFromString(s),
	}
}

func testEncodingJSONFromString(s string) *encoding.JSON {
	JSON, err := encoding.NewJSONFromString(s)
	if err != nil {
		panic(err)
	}
	return JSON
}

func TestNewJSONFromEncodingJSON(t *testing.T) {
	type args struct {
		j *encoding.JSON
	}
	tests := []struct {
		name string
		args args
		want *JSON
	}{
		{
			name: "1",
			args: args{
				j: testEncodingJSONFromString(basicJSON),
			},
			want: testJSONFromString(basicJSON),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJSONFromEncodingJSON(tt.args.j); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJSONFromEncodingJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewJSONFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    *JSON
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				s: basicJSON,
			},
			want: testJSONFromString(basicJSON),
		},

		{
			name: "2",
			args: args{
				s: invlaidJSON,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewJSONFromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJSONFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJSONFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewJSONFromBytes(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *JSON
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				b: []byte(basicJSON),
			},
			want: testJSONFromString(basicJSON),
		},

		{
			name: "2",
			args: args{
				b: []byte(invlaidJSON),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewJSONFromBytes(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJSONFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJSONFromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewJSONFromFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *JSON
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				"1213231k3j21kl3.dadsadasda",
			},
			wantErr: true,
		},
		{
			name: "2",
			args: args{
				"test_data",
			},
			want: testJSONFromString(strings.ReplaceAll(basicJSON, "\n", "\r\n")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewJSONFromFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJSONFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJSONFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSON_GetConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		j       *JSON
		args    args
		want    *JSON
		wantErr bool
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "loggy.programmers.0",
			},
			want: testJSONFromString(`{
				"firstName": "Brett",
				"lastName": "McLaughlin",
				"email": "aaaa",
				"tag": "good"
			}`),
			wantErr: false,
		},

		{
			name: "2",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "loggy.programmers.0.firstName",
			},
			wantErr: true,
		},
		{
			name: "3",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "loggy.programmers.0.1111",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.GetConfig(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSON.GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil && tt.want == nil {
				return
			}
			if got.String() != tt.want.String() {
				t.Errorf("JSON.GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSON_GetBoolOrDefaullt(t *testing.T) {
	type args struct {
		path         string
		defaultValue bool
	}
	tests := []struct {
		name string
		j    *JSON
		args args
		want bool
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "happy",
			},
			want: true,
		},

		{
			name: "2",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "immortal",
			},
		},
		{
			name: "3",
			j:    testJSONFromString(basicJSON),
			args: args{
				path:         "loggy.programmers.0",
				defaultValue: true,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.j.GetBoolOrDefaullt(tt.args.path, tt.args.defaultValue); got != tt.want {
				t.Errorf("JSON.GetBoolOrDefaullt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSON_GetInt64OrDefaullt(t *testing.T) {
	type args struct {
		path         string
		defaultValue int64
	}
	tests := []struct {
		name string
		j    *JSON
		args args
		want int64
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "age",
			},
			want: 100,
		},
		{
			name: "2",
			j:    testJSONFromString(basicJSON),
			args: args{
				path:         "arr",
				defaultValue: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.j.GetInt64OrDefaullt(tt.args.path, tt.args.defaultValue); got != tt.want {
				t.Errorf("JSON.GetInt64OrDefaullt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSON_GetFloat64OrDefaullt(t *testing.T) {
	type args struct {
		path         string
		defaultValue float64
	}
	tests := []struct {
		name string
		j    *JSON
		args args
		want float64
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "loggy.programmers.3.firstName",
			},
			want: 1002.3,
		},
		{
			name: "2",
			j:    testJSONFromString(basicJSON),
			args: args{
				path:         "arr",
				defaultValue: 10.1,
			},
			want: 10.1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.j.GetFloat64OrDefaullt(tt.args.path, tt.args.defaultValue); got != tt.want {
				t.Errorf("JSON.GetFloat64OrDefaullt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSON_GetStringOrDefaullt(t *testing.T) {
	type args struct {
		path         string
		defaultValue string
	}
	tests := []struct {
		name string
		j    *JSON
		args args
		want string
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "name.first",
			},
			want: "tom",
		},
		{
			name: "2",
			j:    testJSONFromString(basicJSON),
			args: args{
				path:         "arr",
				defaultValue: "test",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.j.GetStringOrDefaullt(tt.args.path, tt.args.defaultValue); got != tt.want {
				t.Errorf("JSON.GetStringOrDefaullt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJSON_GetConfigArray(t *testing.T) {
	forArray := func(s string, path string) []*encoding.JSON {
		a, err := testEncodingJSONFromString(s).GetArray(path)
		if err != nil {
			panic(err)
		}
		return a
	}

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		j       *JSON
		args    args
		want    []*encoding.JSON
		wantErr bool
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "vals",
			},
			want: forArray(basicJSON, "vals"),
		},
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "val",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.GetConfigArray(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSON.GetConfigArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("JSON.GetConfigArray() = %v, want %v", got, tt.want)
				return
			}

			for i := range got {
				if got[i].String() != tt.want[i].String() {
					t.Errorf("JSON.GetConfigArray() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestJSON_GetConfigMap(t *testing.T) {
	forMap := func(s string, path string) map[string]*encoding.JSON {
		m, err := testEncodingJSONFromString(s).GetMap(path)
		if err != nil {
			panic(err)
		}
		return m
	}

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		j       *JSON
		args    args
		want    map[string]*encoding.JSON
		wantErr bool
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "name",
			},
			want: forMap(basicJSON, "name"),
		},
		{
			name: "2",
			j:    testJSONFromString(basicJSON),
			args: args{
				path: "vals",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.j.GetConfigMap(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("JSON.GetConfigMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("JSON.GetConfigArray() = %v, want %v", got, tt.want)
				return
			}

			for k := range got {
				if got[k].String() != tt.want[k].String() {
					t.Errorf("JSON.GetConfigArray() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestJSON_CloneConfig(t *testing.T) {
	tests := []struct {
		name string
		j    *JSON
		want *JSON
	}{
		{
			name: "1",
			j:    testJSONFromString(basicJSON),
			want: testJSONFromString(basicJSON),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.j.CloneConfig()

			if got == tt.want {
				t.Errorf("JSON.CloneConfig() = %p, want %p", got, tt.want)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JSON.CloneConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
