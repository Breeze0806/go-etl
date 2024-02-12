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

// Package config provides JSON configuration
package config

import (
	"github.com/Breeze0806/go/encoding"
)

// JSON JSON format configuration file
type JSON struct {
	*encoding.JSON
}

// NewJSONFromEncodingJSON gets JSON from encoded JSON j
func NewJSONFromEncodingJSON(j *encoding.JSON) *JSON {
	return &JSON{
		JSON: j,
	}
}

// NewJSONFromString gets json configuration from string s and returns an error if the json format is wrong
func NewJSONFromString(s string) (*JSON, error) {
	JSON, err := encoding.NewJSONFromString(s)
	if err != nil {
		return nil, err
	}
	return NewJSONFromEncodingJSON(JSON), nil
}

// NewJSONFromBytes gets json configuration from byte stream b and returns an error if the json format is wrong
func NewJSONFromBytes(b []byte) (*JSON, error) {
	JSON, err := encoding.NewJSONFromBytes(b)
	if err != nil {
		return nil, err
	}
	return NewJSONFromEncodingJSON(JSON), nil
}

// NewJSONFromFile gets json configuration from the file named filename
// And return an error if there is a json format error or a file read error
func NewJSONFromFile(filename string) (*JSON, error) {
	JSON, err := encoding.NewJSONFromFile(filename)
	if err != nil {
		return nil, err
	}
	return NewJSONFromEncodingJSON(JSON), nil
}

// GetConfig gets the value configuration file corresponding to the path. For the following json
//
//	{
//	 "a":{
//	    "b":[{
//	       c:"x"
//	     }]
//		}
//	}
//
// To access the x string, the access path for each layer of path is a, a.b, a.b.0, a.b.0.c
// If the corresponding path is not a json structure or does not exist, an error will be returned
func (j *JSON) GetConfig(path string) (*JSON, error) {
	JSON, err := j.GetJSON(path)
	if err != nil {
		return nil, err
	}
	return NewJSONFromEncodingJSON(JSON), nil
}

// GetBoolOrDefaullt gets the BOOL value corresponding to the path. For the following json
//
//	{
//	 "a":{
//	    "b":[{
//	       c:"x"
//	     }]
//		}
//	}
//
// To access the x string, the access path for each layer of path is a, a.b, a.b.0, a.b.0.c
// If the corresponding path is not int64 or does not exist, defaultValue will be returned
func (j *JSON) GetBoolOrDefaullt(path string, defaultValue bool) bool {
	if v, err := j.GetBool(path); err == nil {
		return v
	}
	return defaultValue
}

// GetInt64OrDefaullt gets the int64 value corresponding to the path. For the following json
//
//	{
//	 "a":{
//	    "b":[{
//	       c:"x"
//	     }]
//		}
//	}
//
// To access the x string, the access path for each layer of path is a, a.b, a.b.0, a.b.0.c
// If the corresponding path is not int64 or does not exist, defaultValue will be returned
func (j *JSON) GetInt64OrDefaullt(path string, defaultValue int64) int64 {
	if v, err := j.GetInt64(path); err == nil {
		return v
	}
	return defaultValue
}

// GetFloat64OrDefaullt gets the float64 value corresponding to the path. For the following json
//
//	{
//	 "a":{
//	    "b":[{
//	       c:"x"
//	     }]
//		}
//	}
//
// To access the x string, the access path for each layer of path is a, a.b, a.b.0, a.b.0.c
// If the corresponding path is not float64 or does not exist, defaultValue will be returned
func (j *JSON) GetFloat64OrDefaullt(path string, defaultValue float64) float64 {
	if v, err := j.GetFloat64(path); err == nil {
		return v
	}
	return defaultValue
}

// GetStringOrDefaullt gets the string value corresponding to the path. For the following json
//
//	{
//	 "a":{
//	    "b":[{
//	       c:"x"
//	     }]
//		}
//	}
//
// To access the x string, the access path for each layer of path is a, a.b, a.b.0, a.b.0.c
// If the corresponding path is not a string or does not exist, defaultValue will be returned
func (j *JSON) GetStringOrDefaullt(path string, defaultValue string) string {
	if v, err := j.JSON.GetString(path); err == nil {
		return v
	}
	return defaultValue
}

// GetConfigArray gets the configuration array corresponding to the path. For the following json
//
//	{
//	 "a":{
//	    "b":[{
//	       c:"x"
//	     }]
//		}
//	}
//
// To access the x string, the access path for each layer of path is a, a.b, a.b.0, a.b.0.c
// If the corresponding path is not a configuration array or does not exist, an error will be returned
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

// GetConfigMap gets the configuration mapping corresponding to the path. For the following json
//
//	{
//	 "a":{
//	    "b":[{
//	       c:"x"
//	     }]
//		}
//	}
//
// To access the x string, the access path for each layer of path is a, a.b, a.b.0, a.b.0.c
// If the corresponding path is not a configuration mapping or does not exist, an error will be returned
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

// CloneConfig clones the json configuration file
func (j *JSON) CloneConfig() *JSON {
	return &JSON{
		JSON: j.JSON.Clone(),
	}
}
