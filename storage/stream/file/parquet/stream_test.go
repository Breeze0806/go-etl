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

package parquet

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Breeze0806/go-etl/config"
	"github.com/Breeze0806/go-etl/element"
)

func TestWriter_Write(t *testing.T) {
	tmpDir := os.TempDir()

	type args struct {
		columns  []element.Column
		out      *config.JSON
		in       *config.JSON
		filename string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "basic_types",
			args: args{
				columns: []element.Column{
					// Bytes field
					element.NewDefaultColumn(element.NewBytesColumnValue([]byte("hello world")), "bytes_field", 0),
					// String field
					element.NewDefaultColumn(element.NewStringColumnValue("test string"), "string_field", 0),
					// BigInt field
					element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(123123123), "bigint_field", 0),
					// Boolean field
					element.NewDefaultColumn(element.NewBoolColumnValue(false), "bool_field", 0),
					// Decimal field
					element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat32(142.12312312312), "decimal_field", 0),
				},
				filename: filepath.Join(tmpDir, "basic_types.parquet"),
				in:       testJSONFromString(`{"column":[{"name":"bytes_field","type":"string"},{"name":"string_field","type":"string"},{"name":"bigint_field","type":"bigInt"},{"name":"bool_field","type":"bool"},{"name":"decimal_field","type":"decimal"}]}`),
				out:      testJSONFromString(`{"column":[{"name":"bytes_field","type":"string"},{"name":"string_field","type":"string"},{"name":"bigint_field","type":"bigInt"},{"name":"bool_field","type":"bool"},{"name":"decimal_field","type":"decimal"}]}`),
			},
		},
		{
			name: "all_types",
			args: args{
				columns: []element.Column{
					// String field
					element.NewDefaultColumn(element.NewStringColumnValue("example text"), "text_field", 0),
					// Bytes field
					element.NewDefaultColumn(element.NewBytesColumnValue([]byte{0x00, 0x01, 0x02, 0xFF}), "binary_field", 0),
					// Int64 field
					element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(9223372036854775807), "int64_field", 0),
					// Int32 field (using BigInt with smaller value)
					element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(2147483647), "int32_field", 0),
					// Boolean field
					element.NewDefaultColumn(element.NewBoolColumnValue(true), "boolean_field", 0),
					// Float field
					element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat32(3.14159), "float_field", 0),
					// Double field
					element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(2.718281828459045), "double_field", 0),
					// Timestamp field (as int64)
					element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(time.Now().Unix()), "timestamp_field", 0),
				},
				filename: filepath.Join(tmpDir, "all_types.parquet"),
				in:       testJSONFromString(`{"column":[{"name":"text_field","type":"string"},{"name":"binary_field","type":"string"},{"name":"int64_field","type":"bigInt"},{"name":"int32_field","type":"bigInt"},{"name":"boolean_field","type":"bool"},{"name":"float_field","type":"decimal"},{"name":"double_field","type":"decimal"},{"name":"timestamp_field","type":"bigInt"}]}`),
				out:      testJSONFromString(`{"column":[{"name":"text_field","type":"string"},{"name":"binary_field","type":"string"},{"name":"int64_field","type":"bigInt"},{"name":"int32_field","type":"bigInt"},{"name":"boolean_field","type":"bool"},{"name":"float_field","type":"decimal"},{"name":"double_field","type":"decimal"},{"name":"timestamp_field","type":"bigInt"}]}`),
			},
		},
		{
			name: "edge_cases",
			args: args{
				columns: []element.Column{
					// Empty string
					element.NewDefaultColumn(element.NewStringColumnValue(""), "empty_string", 0),
					// Null/empty bytes
					element.NewDefaultColumn(element.NewBytesColumnValue([]byte{}), "empty_bytes", 0),
					// Zero values
					element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(0), "zero_int", 0),
					element.NewDefaultColumn(element.NewBoolColumnValue(false), "false_bool", 0),
					element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(0.0), "zero_float", 0),
					// Large values
					element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(-9223372036854775808), "min_int64", 0),
					element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(9223372036854775807), "max_int64", 0),
				},
				filename: filepath.Join(tmpDir, "edge_cases.parquet"),
				in:       testJSONFromString(`{"column":[{"name":"empty_string","type":"string"},{"name":"empty_bytes","type":"string"},{"name":"zero_int","type":"bigInt"},{"name":"false_bool","type":"bool"},{"name":"zero_float","type":"decimal"},{"name":"min_int64","type":"bigInt"},{"name":"max_int64","type":"bigInt"}]}`),
				out:      testJSONFromString(`{"column":[{"name":"empty_string","type":"string"},{"name":"empty_bytes","type":"string"},{"name":"zero_int","type":"bigInt"},{"name":"false_bool","type":"bool"},{"name":"zero_float","type":"decimal"},{"name":"min_int64","type":"bigInt"},{"name":"max_int64","type":"bigInt"}]}`),
			},
		},
		{
			name: "long_text",
			args: args{
				columns: []element.Column{
					// Long string
					element.NewDefaultColumn(element.NewStringColumnValue("This is a very long string that contains a lot of text to test how the parquet writer handles longer content. It should be able to store and retrieve this text without any issues."), "long_text", 0),
					// Multiple binary values
					element.NewDefaultColumn(element.NewBytesColumnValue([]byte("binary data 1")), "binary_1", 0),
					element.NewDefaultColumn(element.NewBytesColumnValue([]byte("binary data 2 with more content")), "binary_2", 0),
					// Mixed numeric values
					element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(1000000), "large_int", 0),
					element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(999.999), "large_float", 0),
				},
				filename: filepath.Join(tmpDir, "long_text.parquet"),
				in:       testJSONFromString(`{"column":[{"name":"long_text","type":"string"},{"name":"binary_1","type":"string"},{"name":"binary_2","type":"string"},{"name":"large_int","type":"bigInt"},{"name":"large_float","type":"decimal"}]}`),
				out:      testJSONFromString(`{"column":[{"name":"long_text","type":"string"},{"name":"binary_1","type":"string"},{"name":"binary_2","type":"string"},{"name":"large_int","type":"bigInt"},{"name":"large_float","type":"decimal"}]}`),
			},
		},
		{
			name: "timestamp_data",
			args: args{
				columns: []element.Column{
					// Current timestamp
					element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(time.Now().Unix()), "current_timestamp", 0),
					// Past timestamp
					element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).Unix()), "past_timestamp", 0),
					// Future timestamp
					element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(time.Date(2030, 12, 31, 23, 59, 59, 0, time.UTC).Unix()), "future_timestamp", 0),
					// String representation of time
					element.NewDefaultColumn(element.NewStringColumnValue(time.Now().Format(time.RFC3339)), "time_string", 0),
				},
				filename: filepath.Join(tmpDir, "timestamp_data.parquet"),
				in:       testJSONFromString(`{"column":[{"name":"current_timestamp","type":"bigInt"},{"name":"past_timestamp","type":"bigInt"},{"name":"future_timestamp","type":"bigInt"},{"name":"time_string","type":"string"}]}`),
				out:      testJSONFromString(`{"column":[{"name":"current_timestamp","type":"bigInt"},{"name":"past_timestamp","type":"bigInt"},{"name":"future_timestamp","type":"bigInt"},{"name":"time_string","type":"string"}]}`),
			},
		},
		{
			name: "boolean_variations",
			args: args{
				columns: []element.Column{
					element.NewDefaultColumn(element.NewBoolColumnValue(true), "true_value", 0),
					element.NewDefaultColumn(element.NewBoolColumnValue(false), "false_value", 0),
					element.NewDefaultColumn(element.NewStringColumnValue("true"), "string_true", 0),
					element.NewDefaultColumn(element.NewStringColumnValue("false"), "string_false", 0),
					element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(1), "int_true", 0),
					element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(0), "int_false", 0),
				},
				filename: filepath.Join(tmpDir, "boolean_variations.parquet"),
				in:       testJSONFromString(`{"column":[{"name":"true_value","type":"bool"},{"name":"false_value","type":"bool"},{"name":"string_true","type":"string"},{"name":"string_false","type":"string"},{"name":"int_true","type":"bigInt"},{"name":"int_false","type":"bigInt"}]}`),
				out:      testJSONFromString(`{"column":[{"name":"true_value","type":"bool"},{"name":"false_value","type":"bool"},{"name":"string_true","type":"string"},{"name":"string_false","type":"string"},{"name":"int_true","type":"bigInt"},{"name":"int_false","type":"bigInt"}]}`),
			},
		},
		{
			name: "numeric_precision",
			args: args{
				columns: []element.Column{
					// Integer values
					element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(1), "int_1", 0),
					element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(123), "int_123", 0),
					element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(12345), "int_12345", 0),
					// Float values with different precision
					element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat32(1.1), "float_1_1", 0),
					element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat32(1.12345), "float_1_12345", 0),
					element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(1.123456789012345), "double_precise", 0),
					// Negative values
					element.NewDefaultColumn(element.NewBigIntColumnValueFromInt64(-123), "negative_int", 0),
					element.NewDefaultColumn(element.NewDecimalColumnValueFromFloat(-1.23456789), "negative_float", 0),
				},
				filename: filepath.Join(tmpDir, "numeric_precision.parquet"),
				in:       testJSONFromString(`{"column":[{"name":"int_1","type":"bigInt"},{"name":"int_123","type":"bigInt"},{"name":"int_12345","type":"bigInt"},{"name":"float_1_1","type":"decimal"},{"name":"float_1_12345","type":"decimal"},{"name":"double_precise","type":"decimal"},{"name":"negative_int","type":"bigInt"},{"name":"negative_float","type":"decimal"}]}`),
				out:      testJSONFromString(`{"column":[{"name":"int_1","type":"bigInt"},{"name":"int_123","type":"bigInt"},{"name":"int_12345","type":"bigInt"},{"name":"float_1_1","type":"decimal"},{"name":"float_1_12345","type":"decimal"},{"name":"double_precise","type":"decimal"},{"name":"negative_int","type":"bigInt"},{"name":"negative_float","type":"decimal"}]}`),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer os.Remove(tt.args.filename)
			var creator Creator
			out, err := creator.Create(tt.args.filename)
			if err != nil {
				t.Fatal(err)
			}
			defer out.Close()
			w, err := out.Writer(tt.args.out)
			if err != nil {
				t.Fatal(err)
			}
			defer w.Close()
			defer w.Flush()
			r := element.NewDefaultRecord()
			for _, c := range tt.args.columns {
				r.Add(c)
			}
			err = w.Write(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("writer.Write() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
