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

//go:build go1.18
// +build go1.18

package element

import (
	"testing"
)

func FuzzConverterConvertDecimal(f *testing.F) {
	for _, t := range testTableDecimalStr {
		f.Add(t.short)
	}

	f.Fuzz(func(t *testing.T, number string) {
		num1, err1 := testNumConverter.ConvertDecimal(number)
		num2, err2 := testOldNumConverter.ConvertDecimal(number)
		if err1 == nil && err2 != nil {
			t.Fatalf("input: %v err1: %v err2: %v", number, err1, err2)
		}
		if err1 != nil && err2 == nil {
			t.Fatalf("input: %v err1: %v err2:%v", number, err1, err2)
		}
		if err1 == nil && err2 == nil {
			if num1.String() != num2.String() {
				t.Fatalf("input: %v num1: %v num2: %v", number, num1.String(), num2.String())
			}
		}
	})
}
