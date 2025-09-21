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

package element

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"testing"

	"github.com/cockroachdb/apd/v3"
)

var (
	benchDecimal = &Decimal{
		value: apd.New(math.MaxInt64, -9),
	}

	benchDecimalStr = &DecimalStr{
		value:  strconv.FormatInt(math.MaxInt64, 10),
		intLen: 9,
	}

	benchInt64 = &Int64{
		value: math.MaxInt64,
	}

	benchBigInt = &BigInt{
		value: apd.NewBigInt(math.MaxInt64),
	}

	benchBigIntStr = &BigIntStr{
		value: strconv.FormatInt(math.MaxInt64, 10),
	}
)

func BenchmarkConverter_ConvertFromBigInt(b *testing.B) {
	rng := rand.New(rand.NewSource(0xdead1337))
	in := make([]int64, b.N)
	for i := range in {
		in[i] = int64(rng.Intn(math.MaxInt64))
	}
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = testNumConverter.ConvertBigIntFromInt(in[i])
	}
}

func BenchmarkOldConverter_ConvertFromBigInt(b *testing.B) {
	rng := rand.New(rand.NewSource(0xdead1337))
	in := make([]int64, b.N)
	for i := range in {
		in[i] = int64(rng.Intn(math.MaxInt64))
	}
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = testOldNumConverter.ConvertBigIntFromInt(in[i])
	}
}

func BenchmarkConverter_ConvertDecimalFromloat(b *testing.B) {
	rng := rand.New(rand.NewSource(0xdead1337))
	in := make([]float64, b.N)
	for i := range in {
		in[i] = rng.NormFloat64() * 10e20
	}
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = testNumConverter.ConvertDecimalFromFloat(in[i])
	}
}

func BenchmarkOldConverter_ConvertDecimalFromFloat(b *testing.B) {
	rng := rand.New(rand.NewSource(0xdead1337))
	in := make([]float64, b.N)
	for i := range in {
		in[i] = rng.NormFloat64() * 10e20
	}
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_ = testOldNumConverter.ConvertDecimalFromFloat(in[i])
	}
}

func BenchmarkConverter_ConvertBigInt_Int64(b *testing.B) {
	rng := rand.New(rand.NewSource(0xdead1337))
	in := make([]int64, b.N)
	for i := range in {
		in[i] = int64(rng.Intn(math.MaxInt64))
	}
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		in := strconv.FormatInt(in[i], 10)
		_, _ = testNumConverter.ConvertBigInt(in)
	}
}

func BenchmarkOldConverter_ConvertBigInt_Int64(b *testing.B) {
	rng := rand.New(rand.NewSource(0xdead1337))
	in := make([]int64, b.N)
	for i := range in {
		in[i] = int64(rng.Intn(math.MaxInt64))
	}
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		in := strconv.FormatInt(in[i], 10)
		_, _ = testOldNumConverter.ConvertBigInt(in)
	}
}

func BenchmarkCoventor_ConvertBigInt_large_number(b *testing.B) {
	count := 72
	prices := make([]string, 0, count)
	for i := 1; i <= count; i++ {
		prices = append(prices, "93233720368547758079223372036854775807")
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, p := range prices {
			_, _ = testNumConverter.ConvertBigInt(p)
		}
	}
}

func BenchmarkOldCoventor_ConvertBigInt_large_number(b *testing.B) {
	count := 72
	prices := make([]string, 0, count)
	for i := 1; i <= count; i++ {
		prices = append(prices, "93233720368547758079223372036854775807")
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, p := range prices {
			_, _ = testOldNumConverter.ConvertBigInt(p)
		}
	}
}

func BenchmarkConverter_ConvertDecimal_Int64(b *testing.B) {
	rng := rand.New(rand.NewSource(0xdead1337))
	in := make([]int64, b.N)
	for i := range in {
		in[i] = int64(rng.Intn(math.MaxInt64))
	}
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		in := strconv.FormatInt(in[i], 10)
		_, _ = testNumConverter.ConvertDecimal(in)
	}
}

func BenchmarkOldConverter_ConvertDecimal_Int64(b *testing.B) {
	rng := rand.New(rand.NewSource(0xdead1337))
	in := make([]int64, b.N)
	for i := range in {
		in[i] = int64(rng.Intn(math.MaxInt64))
	}
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		in := strconv.FormatInt(in[i], 10)
		_, _ = testOldNumConverter.ConvertDecimal(in)
	}
}

func BenchmarkConverter_ConvertDecimal_Float64(b *testing.B) {
	rng := rand.New(rand.NewSource(0xdead1337))
	in := make([]float64, b.N)
	for i := range in {
		in[i] = rng.NormFloat64() * 10e20
	}
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		in := strconv.FormatFloat(in[i], 'f', -1, 64)
		_, _ = testNumConverter.ConvertDecimal(in)
	}
}

func BenchmarkOldConverter_ConvertDecimal_Float64(b *testing.B) {
	rng := rand.New(rand.NewSource(0xdead1337))
	in := make([]float64, b.N)
	for i := range in {
		in[i] = rng.NormFloat64() * 10e20
	}
	b.ReportAllocs()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		in := strconv.FormatFloat(in[i], 'f', -1, 64)
		_, _ = testOldNumConverter.ConvertDecimal(in)
	}
}

func BenchmarkConverter_ConvertDecimal(b *testing.B) {
	count := 72
	prices := make([]string, 0, count)
	for i := 1; i <= count; i++ {
		prices = append(prices, fmt.Sprintf("%d.%d", i*100, i))
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, p := range prices {
			d, err := testNumConverter.ConvertDecimal(p)
			if err != nil {
				b.Log(d)
				b.Error(err)
			}
		}
	}
}

func BenchmarkOldConverter_ConvertDecimal(b *testing.B) {
	count := 72
	prices := make([]string, 0, count)
	for i := 1; i <= count; i++ {
		prices = append(prices, fmt.Sprintf("%d.%d", i*100, i))
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, p := range prices {
			d, err := testOldNumConverter.ConvertDecimal(p)
			if err != nil {
				b.Log(d)
				b.Error(err)
			}
		}
	}
}

func BenchmarkConverter_ConvertDecimal_large_number(b *testing.B) {
	count := 72
	prices := make([]string, 0, count)
	for i := 1; i <= count; i++ {
		prices = append(prices, "9323372036854775807.9223372036854775807")
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, p := range prices {
			d, err := testNumConverter.ConvertDecimal(p)
			if err != nil {
				b.Log(d)
				b.Error(err)
			}
		}
	}
}

func BenchmarkOldConverter_ConvertDecimal_large_number(b *testing.B) {
	count := 72
	prices := make([]string, 0, count)
	for i := 1; i <= count; i++ {
		prices = append(prices, "9323372036854775807.9223372036854775807")
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, p := range prices {
			d, err := testOldNumConverter.ConvertDecimal(p)
			if err != nil {
				b.Log(d)
				b.Error(err)
			}
		}
	}
}

func BenchmarkConverter_ConvertDecimal_Exp(b *testing.B) {
	count := 72
	prices := make([]string, 0, count)
	for i := 1; i <= count; i++ {
		prices = append(prices, "9323372036854775807.922e123456")
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, p := range prices {
			d, err := testNumConverter.ConvertDecimal(p)
			if err != nil {
				b.Log(d)
				b.Error(err)
			}
		}
	}
}

func BenchmarkOldConverter_ConvertDecimal_Exp(b *testing.B) {
	count := 72
	prices := make([]string, 0, count)
	for i := 1; i <= count; i++ {
		prices = append(prices, "9323372036854775807.922e12356")
	}

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, p := range prices {
			d, err := testOldNumConverter.ConvertDecimal(p)
			if err != nil {
				b.Log(d)
				b.Error(err)
			}
		}
	}
}

func BenchmarkDecimal_Decmial_String(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = benchDecimal.String()
	}
}

func BenchmarkDecimal_DecmialStr_String(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = benchDecimalStr.String()
	}
}

func BenchmarkDecimal_Int64_String(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = benchInt64.String()
	}
}

func BenchmarkDecimal_BigInt_String(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = benchBigInt.String()
	}
}

func BenchmarkDecimal_BigIntStr_String(b *testing.B) {

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = benchBigIntStr.String()
	}
}
