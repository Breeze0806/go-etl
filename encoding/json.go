package encoding

import (
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type Json struct {
	res gjson.Result
}

func NewJsonFromString(s string) (*Json, error) {
	if !gjson.Valid(s) {
		return nil, fmt.Errorf("%v is not valid json", s)
	}
	return newJsonFromString(s), nil
}

func newJsonFromString(s string) *Json {
	return &Json{
		res: gjson.Parse(s),
	}
}

func NewJsonFromBytes(b []byte) (*Json, error) {
	if !gjson.ValidBytes(b) {
		return nil, fmt.Errorf("%v is not valid json", string(b))
	}
	return newJsonFromBytes(b), nil
}

func newJsonFromBytes(b []byte) *Json {
	return &Json{
		res: gjson.ParseBytes(b),
	}
}

func NewJsonFromFile(filename string) (*Json, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read file %v fail. err： %v", filename, err)
	}
	return NewJsonFromBytes(data)
}

func (j *Json) GetJson(path string) (*Json, error) {
	res, err := j.getResult(path)
	if err != nil {
		return nil, err
	}
	if res.Type != gjson.JSON {
		return nil, fmt.Errorf("path(%v) is not json", path)
	}

	return &Json{
		res: res,
	}, nil
}

func (j *Json) GetBool(path string) (bool, error) {
	res, err := j.getResult(path)
	if err != nil {
		return false, err
	}
	switch res.Type {
	case gjson.False:
		return false, nil
	case gjson.True:
		return true, nil
	}
	return false, fmt.Errorf("path(%v) is not bool", path)
}

func (j *Json) GetInt64(path string) (int64, error) {
	res, err := j.getResult(path)
	if err != nil {
		return 0, err
	}
	switch res.Type {
	case gjson.Number:
		v, err := strconv.ParseInt(res.String(), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("path(%v) is not int64. val: %v", res.String(), err)
		}
		return v, nil
	}
	return 0, fmt.Errorf("path(%v) is not bool", path)
}

func (j *Json) GetFloat64(path string) (float64, error) {
	res, err := j.getResult(path)
	if err != nil {
		return 0, err
	}
	switch res.Type {
	case gjson.Number:
		v, err := strconv.ParseFloat(res.String(), 64)
		if err != nil {
			return 0, fmt.Errorf("path(%v) is not float64. val: %v", res.String(), err)
		}
		return v, nil
	}
	return 0, fmt.Errorf("path(%v) is not bool", path)
}

func (j *Json) GetString(path string) (string, error) {
	res, err := j.getResult(path)
	if err != nil {
		return "", err
	}
	switch res.Type {
	case gjson.String:
		return res.String(), nil
	}
	return "", fmt.Errorf("path(%v) is not string", path)
}

func (j *Json) GetArray(path string) ([]*Json, error) {
	res, err := j.getResult(path)
	if err != nil {
		return nil, err
	}
	switch {
	case res.IsArray():
		var jsons []*Json
		a := res.Array()
		for _, v := range a {
			jsons = append(jsons, &Json{res: v})
		}
		return jsons, nil
	}
	return nil, fmt.Errorf("path(%v) is not array", path)
}

func (j *Json) GetMap(path string) (map[string]*Json, error) {
	res, err := j.getResult(path)
	if err != nil {
		return nil, err
	}
	switch {
	case res.IsObject():
		jsons := make(map[string]*Json)
		m := res.Map()
		for k, v := range m {
			jsons[k] = &Json{res: v}
		}
		return jsons, nil
	}
	return nil, fmt.Errorf("path(%v) is not map", path)
}

func (j *Json) String() string {
	return j.res.String()
}

func (j *Json) IsArray(path string) bool {
	return j.res.Get(path).IsArray()
}

func (j *Json) IsNumber(path string) bool {
	return j.res.Get(path).Type == gjson.Number
}

func (j *Json) IsJson(path string) bool {
	return j.res.Get(path).IsObject()
}

func (j *Json) IsBool(path string) bool {
	switch j.res.Get(path).Type {
	case gjson.False, gjson.True:
		return true
	}
	return false
}

func (j *Json) IsString(path string) bool {
	return j.res.Get(path).Type == gjson.String
}

func (j *Json) IsNull(path string) bool {
	return j.res.Get(path).Type == gjson.Null
}

func (j *Json) Exists(path string) bool {
	return j.res.Get(path).Exists()
}

func (j *Json) Set(path string, v interface{}) error {
	s, err := sjson.Set(j.String(), path, v)
	if err != nil {
		return fmt.Errorf("path(%v) set fail. err: %v", path, err)
	}
	j.fromString(s)
	return nil
}

func (j *Json) SetRawBytes(path string, b []byte) error {
	return j.SetRawString(path, string(b))
}

func (j *Json) SetRawString(path string, s string) error {
	ns, err := sjson.SetRaw(j.String(), path, s)
	if err != nil {
		return fmt.Errorf("path(%v) set fail. err: %v", path, err)
	}
	j.fromString(ns)
	return nil
}

func (j *Json) Remove(path string) error {
	s, err := sjson.Delete(j.String(), path)
	if err != nil {
		return fmt.Errorf("path(%v) remove fail. err: %v", path, err)
	}
	j.fromString(s)
	return nil
}

func (j *Json) FromString(s string) error {
	new, err := NewJsonFromString(s)
	if err != nil {
		return err
	}
	j.res = new.res
	return nil
}

func (j *Json) FromBytes(b []byte) error {
	new, err := NewJsonFromBytes(b)
	if err != nil {
		return err
	}
	j.res = new.res
	return nil
}

func (j *Json) FromFile(filename string) error {
	new, err := NewJsonFromFile(filename)
	if err != nil {
		return err
	}
	j.res = new.res
	return nil
}

func (j *Json) Clone() *Json {
	return &Json{
		res: j.res,
	}
}

func (j *Json) MarshalJSON() ([]byte, error) {
	return []byte(j.String()), nil
}

func (j *Json) fromString(s string) {
	new := newJsonFromString(s)
	j.res = new.res
}

func (j *Json) getResult(path string) (gjson.Result, error) {
	res := j.res.Get(path)
	if res.Exists() {
		return res, nil
	}
	return gjson.Result{}, fmt.Errorf("path(%v) does not exist", path)
}
