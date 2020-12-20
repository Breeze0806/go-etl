package config

import (
	"github.com/Breeze0806/go-etl/encoding"
)

type Json struct {
	*encoding.Json
}

func NewJsonFromEncodingJson(j *encoding.Json) *Json {
	return &Json{
		Json: j,
	}
}

func NewJsonFromString(s string) (*Json, error) {
	json, err := encoding.NewJsonFromString(s)
	if err != nil {
		return nil, err
	}
	return NewJsonFromEncodingJson(json), nil
}

func NewJsonFromBytes(b []byte) (*Json, error) {
	json, err := encoding.NewJsonFromBytes(b)
	if err != nil {
		return nil, err
	}
	return NewJsonFromEncodingJson(json), nil
}

func NewJsonFromFile(filename string) (*Json, error) {
	json, err := encoding.NewJsonFromFile(filename)
	if err != nil {
		return nil, err
	}
	return NewJsonFromEncodingJson(json), nil
}

func (j *Json) GetConfig(path string) (*Json, error) {
	json, err := j.GetJson(path)
	if err != nil {
		return nil, err
	}
	return NewJsonFromEncodingJson(json), nil
}

func (j *Json) GetBoolOrDefaullt(path string, defaultValue bool) bool {
	if v, err := j.GetBool(path); err == nil {
		return v
	}
	return defaultValue
}

func (j *Json) GetInt64OrDefaullt(path string, defaultValue int64) int64 {
	if v, err := j.GetInt64(path); err == nil {
		return v
	}
	return defaultValue
}

func (j *Json) GetFloat64OrDefaullt(path string, defaultValue float64) float64 {
	if v, err := j.GetFloat64(path); err == nil {
		return v
	}
	return defaultValue
}

func (j *Json) GetStringOrDefaullt(path string, defaultValue string) string {
	if v, err := j.Json.GetString(path); err == nil {
		return v
	}
	return defaultValue
}

func (j *Json) GetConfigArray(path string) ([]*Json, error) {
	a, err := j.Json.GetArray(path)
	if err != nil {
		return nil, err
	}

	var jsons []*Json

	for i := range a {
		jsons = append(jsons, NewJsonFromEncodingJson(a[i]))
	}

	return jsons, nil
}

func (j *Json) GetConfigMap(path string) (map[string]*Json, error) {
	m, err := j.Json.GetMap(path)
	if err != nil {
		return nil, err
	}

	jsons := make(map[string]*Json)

	for k, v := range m {
		jsons[k] = NewJsonFromEncodingJson(v)
	}
	return jsons, nil
}

func (j *Json) CloneConfig() *Json {
	return &Json{
		Json: j.Json.Clone(),
	}
}
