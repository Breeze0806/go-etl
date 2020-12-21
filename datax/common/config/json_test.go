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

func testJsonFromString(s string) *Json {
	return &Json{
		Json: testEncodingJsonFromString(s),
	}
}

func testEncodingJsonFromString(s string) *encoding.Json {
	json, err := encoding.NewJsonFromString(s)
	if err != nil {
		panic(err)
	}
	return json
}

func TestNewJsonFromEncodingJson(t *testing.T) {
	type args struct {
		j *encoding.Json
	}
	tests := []struct {
		name string
		args args
		want *Json
	}{
		{
			name: "1",
			args: args{
				j: testEncodingJsonFromString(basicJSON),
			},
			want: testJsonFromString(basicJSON),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJsonFromEncodingJson(tt.args.j); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJsonFromEncodingJson() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewJsonFromString(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    *Json
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				s: basicJSON,
			},
			want: testJsonFromString(basicJSON),
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
			got, err := NewJsonFromString(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJsonFromString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJsonFromString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewJsonFromBytes(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *Json
		wantErr bool
	}{
		{
			name: "1",
			args: args{
				b: []byte(basicJSON),
			},
			want: testJsonFromString(basicJSON),
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
			got, err := NewJsonFromBytes(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJsonFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJsonFromBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewJsonFromFile(t *testing.T) {
	type args struct {
		filename string
	}
	tests := []struct {
		name    string
		args    args
		want    *Json
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
			want: testJsonFromString(strings.ReplaceAll(basicJSON, "\n", "\r\n")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewJsonFromFile(tt.args.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewJsonFromFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJsonFromFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_GetConfig(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		j       *Json
		args    args
		want    *Json
		wantErr bool
	}{
		{
			name: "1",
			j:    testJsonFromString(basicJSON),
			args: args{
				path: "loggy.programmers.0",
			},
			want: testJsonFromString(`{
				"firstName": "Brett",
				"lastName": "McLaughlin",
				"email": "aaaa",
				"tag": "good"
			}`),
			wantErr: false,
		},

		{
			name: "2",
			j:    testJsonFromString(basicJSON),
			args: args{
				path: "loggy.programmers.0.firstName",
			},
			wantErr: true,
		},
		{
			name: "3",
			j:    testJsonFromString(basicJSON),
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
				t.Errorf("Json.GetConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got == nil && tt.want == nil {
				return
			}
			if got.String() != tt.want.String() {
				t.Errorf("Json.GetConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_GetBoolOrDefaullt(t *testing.T) {
	type args struct {
		path         string
		defaultValue bool
	}
	tests := []struct {
		name string
		j    *Json
		args args
		want bool
	}{
		{
			name: "1",
			j:    testJsonFromString(basicJSON),
			args: args{
				path: "happy",
			},
			want: true,
		},

		{
			name: "2",
			j:    testJsonFromString(basicJSON),
			args: args{
				path: "immortal",
			},
		},
		{
			name: "3",
			j:    testJsonFromString(basicJSON),
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
				t.Errorf("Json.GetBoolOrDefaullt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_GetInt64OrDefaullt(t *testing.T) {
	type args struct {
		path         string
		defaultValue int64
	}
	tests := []struct {
		name string
		j    *Json
		args args
		want int64
	}{
		{
			name: "1",
			j:    testJsonFromString(basicJSON),
			args: args{
				path: "age",
			},
			want: 100,
		},
		{
			name: "2",
			j:    testJsonFromString(basicJSON),
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
				t.Errorf("Json.GetInt64OrDefaullt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_GetFloat64OrDefaullt(t *testing.T) {
	type args struct {
		path         string
		defaultValue float64
	}
	tests := []struct {
		name string
		j    *Json
		args args
		want float64
	}{
		{
			name: "1",
			j:    testJsonFromString(basicJSON),
			args: args{
				path: "loggy.programmers.3.firstName",
			},
			want: 1002.3,
		},
		{
			name: "2",
			j:    testJsonFromString(basicJSON),
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
				t.Errorf("Json.GetFloat64OrDefaullt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_GetStringOrDefaullt(t *testing.T) {
	type args struct {
		path         string
		defaultValue string
	}
	tests := []struct {
		name string
		j    *Json
		args args
		want string
	}{
		{
			name: "1",
			j:    testJsonFromString(basicJSON),
			args: args{
				path: "name.first",
			},
			want: "tom",
		},
		{
			name: "2",
			j:    testJsonFromString(basicJSON),
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
				t.Errorf("Json.GetStringOrDefaullt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJson_GetConfigArray(t *testing.T) {
	forArray := func(s string, path string) []*encoding.Json {
		a, err := testEncodingJsonFromString(s).GetArray(path)
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
		j       *Json
		args    args
		want    []*encoding.Json
		wantErr bool
	}{
		{
			name: "1",
			j:    testJsonFromString(basicJSON),
			args: args{
				path: "vals",
			},
			want: forArray(basicJSON, "vals"),
		},
		{
			name: "1",
			j:    testJsonFromString(basicJSON),
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
				t.Errorf("Json.GetConfigArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("Json.GetConfigArray() = %v, want %v", got, tt.want)
				return
			}

			for i := range got {
				if got[i].String() != tt.want[i].String() {
					t.Errorf("Json.GetConfigArray() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestJson_GetConfigMap(t *testing.T) {
	forMap := func(s string, path string) map[string]*encoding.Json {
		m, err := testEncodingJsonFromString(s).GetMap(path)
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
		j       *Json
		args    args
		want    map[string]*encoding.Json
		wantErr bool
	}{
		{
			name: "1",
			j:    testJsonFromString(basicJSON),
			args: args{
				path: "name",
			},
			want: forMap(basicJSON, "name"),
		},
		{
			name: "2",
			j:    testJsonFromString(basicJSON),
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
				t.Errorf("Json.GetConfigMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("Json.GetConfigArray() = %v, want %v", got, tt.want)
				return
			}

			for k := range got {
				if got[k].String() != tt.want[k].String() {
					t.Errorf("Json.GetConfigArray() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestJson_CloneConfig(t *testing.T) {
	tests := []struct {
		name string
		j    *Json
		want *Json
	}{
		{
			name: "1",
			j:    testJsonFromString(basicJSON),
			want: testJsonFromString(basicJSON),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.j.CloneConfig()

			if got == tt.want {
				t.Errorf("Json.CloneConfig() = %p, want %p", got, tt.want)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Json.CloneConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
