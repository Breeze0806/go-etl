// Copyright 2020 the go-etl Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//Package config 提供JSON配置
package config

import (
	"github.com/Breeze0806/go/encoding"
)

//JSON JSON格式配置文件
type JSON struct {
	*encoding.JSON
}

//NewJSONFromEncodingJSON 从编码JSON j中获取JSON
func NewJSONFromEncodingJSON(j *encoding.JSON) *JSON {
	return &JSON{
		JSON: j,
	}
}

//NewJSONFromString 从字符串s 获取json配置 ，并在json格式错误时返回错误
func NewJSONFromString(s string) (*JSON, error) {
	JSON, err := encoding.NewJSONFromString(s)
	if err != nil {
		return nil, err
	}
	return NewJSONFromEncodingJSON(JSON), nil
}

//NewJSONFromBytes 从字节流b 获取json配置 ，并在json格式错误时返回错误
func NewJSONFromBytes(b []byte) (*JSON, error) {
	JSON, err := encoding.NewJSONFromBytes(b)
	if err != nil {
		return nil, err
	}
	return NewJSONFromEncodingJSON(JSON), nil
}

//NewJSONFromFile 从文件名为filename的文件中获取json配置
//并在json格式错误或者读取文件错误时返回错误
func NewJSONFromFile(filename string) (*JSON, error) {
	JSON, err := encoding.NewJSONFromFile(filename)
	if err != nil {
		return nil, err
	}
	return NewJSONFromEncodingJSON(JSON), nil
}

//GetConfig 获取path路径对应的值配置文件,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
//如果path对应的不是json结构或者不存在，就会返回错误
func (j *JSON) GetConfig(path string) (*JSON, error) {
	JSON, err := j.GetJSON(path)
	if err != nil {
		return nil, err
	}
	return NewJSONFromEncodingJSON(JSON), nil
}

//GetBoolOrDefaullt 获取path路径对应的BOOL值,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
//如果path对应的不是int64或者不存在，就会返回defaultValue
func (j *JSON) GetBoolOrDefaullt(path string, defaultValue bool) bool {
	if v, err := j.GetBool(path); err == nil {
		return v
	}
	return defaultValue
}

//GetInt64OrDefaullt 获取path路径对应的int64值,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
//如果path对应的不是int64或者不存在，就会返回defaultValue
func (j *JSON) GetInt64OrDefaullt(path string, defaultValue int64) int64 {
	if v, err := j.GetInt64(path); err == nil {
		return v
	}
	return defaultValue
}

//GetFloat64OrDefaullt 获取path路径对应的float64值,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
//如果path对应的不是float64或者不存在，就会返回defaultValue
func (j *JSON) GetFloat64OrDefaullt(path string, defaultValue float64) float64 {
	if v, err := j.GetFloat64(path); err == nil {
		return v
	}
	return defaultValue
}

//GetStringOrDefaullt 获取path路径对应的字符串值,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
//如果path对应的不是字符串或者不存在，就会返回defaultValue
func (j *JSON) GetStringOrDefaullt(path string, defaultValue string) string {
	if v, err := j.JSON.GetString(path); err == nil {
		return v
	}
	return defaultValue
}

//GetConfigArray 获取path路径对应的配置数组,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
//如果path对应的不是配置数组或者不存在，就会返回错误
func (j *JSON) GetConfigArray(path string) ([]*JSON, error) {
	a, err := j.JSON.GetArray(path)
	if err != nil {
		return nil, err
	}

	var JSONs []*JSON

	for i := range a {
		JSONs = append(JSONs, NewJSONFromEncodingJSON(a[i]))
	}

	return JSONs, nil
}

//GetConfigMap 获取path路径对应的配置映射,对于下列json
//{
//  "a":{
//     "b":[{
//        c:"x"
//      }]
//	}
//}
//要访问到x字符串 path每层的访问路径为a,a.b,a.b.0，a.b.0.c
//如果path对应的不是配置映射或者不存在，就会返回错误
func (j *JSON) GetConfigMap(path string) (map[string]*JSON, error) {
	m, err := j.JSON.GetMap(path)
	if err != nil {
		return nil, err
	}

	JSONs := make(map[string]*JSON)

	for k, v := range m {
		JSONs[k] = NewJSONFromEncodingJSON(v)
	}
	return JSONs, nil
}

//CloneConfig 克隆json配置文件
func (j *JSON) CloneConfig() *JSON {
	return &JSON{
		JSON: j.JSON.Clone(),
	}
}
