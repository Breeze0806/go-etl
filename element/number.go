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
	"github.com/cockroachdb/apd/v3"
)

// Number   Numeric value
type Number interface {
	Bool() (bool, error)
	String() string
}

// BigIntNumber   High-precision integer
type BigIntNumber interface {
	Number

	Int64() (int64, error)
	Decimal() DecimalNumber
	CloneBigInt() BigIntNumber
	AsBigInt() *apd.BigInt
}

// DecimalNumber   High-precision decimal
type DecimalNumber interface {
	Number

	Float64() (float64, error)
	BigInt() BigIntNumber
	CloneDecimal() DecimalNumber
	AsDecimal() *apd.Decimal
}
