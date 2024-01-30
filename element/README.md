# go-etl Data Type Descriptions

This package primarily defines the data types used in go-etl.

## Record

```go
// Record represents a data record.
type Record interface {
 fmt.Stringer

 Add(Column) error                      // Adds a new column.
 GetByIndex(i int) (Column, error)      // Retrieves the column at the specified index.
 GetByName(name string) (Column, error) // Retrieves the column with the specified name.
 Set(i int, c Column) error             // Sets the column at the specified index.
 ColumnNumber() int                     // Returns the number of columns.
 ByteSize() int64                       // Returns the size of the record in bytes.
 MemorySize() int64                     // Returns the memory usage of the record.
}
```

## Data Type Conversions

go-etl supports six internal data types:

- `bigInt`: Fixed-point numbers (int64, int32, int16, int8, BigInt, etc.).
- `decimal`: Floating-point numbers (float32, float64, BigDecimal (unlimited precision), etc.).
- `string`: String type, with unlimited length and using a universal character set (Unicode).
- `time`: Date and time type.
- `bool`: Boolean value.
- `bytes`: Binary data, which can store unstructured data such as MP3 files.

Correspondingly, there are six implementations of `ColumnValue`: `TimeColumnValue`, `BigIntColumnValue`, `DecimalColumnValue`, `BytesColumnValue`, `StringColumnValue`, and `BoolColumnValue`.

These `ColumnValue` interfaces provide a series of data type conversion methods that start with `as`.

```go
// ColumnValue represents a value in a column.
type ColumnValue interface {
 fmt.Stringer

 Type() ColumnType                    // Returns the column type.
 IsNil() bool                         // Checks if the value is nil.
 AsBool() (bool, error)               // Converts the value to a boolean.
 AsBigInt() (*big.Int, error)         // Converts the value to a big integer.
 AsDecimal() (decimal.Decimal, error) // Converts the value to a decimal with unlimited precision.
 AsString() (string, error)           // Converts the value to a string.
 AsBytes() ([]byte, error)            // Converts the value to a byte array.
 AsTime() (time.Time, error)          // Converts the value to a time.
}
```

Based on the `ColumnValue` interface, the following methods are implemented:

```go
// Column represents a data column.
type Column interface {
 ColumnValue
 AsInt64() (int64, error)     // Converts the value to a 64-bit integer.
 AsFloat64() (float64, error) // Converts the value to a 64-bit floating point number.
 Clone() (Column, error)      // Clones the column.
 Cmp(Column) (int, error)     // Compares the column with another column. Returns 1 if greater, 0 if equal, -1 if less.
 Name() string                // Returns the name of the column.
 ByteSize() int64             // Returns the size of the column in bytes.
 MemorySize() int64           // Returns the memory usage of the column.
}
```

The internal types of DataX are implemented using different Golang types:
Currently, there are two implementation approaches, but the older approach has performance issues when dealing with large datasets. The new implementation is still in beta and has not been thoroughly validated through practical use.

+ Older Implementation Approach

| Internal Type | Implementation Type | Notes                                                        |
| ------------- | ------------------- | ------------------------------------------------------------ |
| time          | time.Time           |                                                              |
| bigInt        | big.Int             | Uses arbitrary-precision integers to ensure no loss of precision. |
| decimal       | decimal.Decimal     | Represented using decimal.Decimal to ensure no loss of precision. |
| bytes         | []byte              |                                                              |
| string        | string              |                                                              |
| bool          | bool                |                                                              |

+ Current Implementation Approach

| Internal Type | Implementation Type | Notes                                                        |
| ------------- | ------------------- | ------------------------------------------------------------ |
| time          | time.Time           |                                                              |
| bigInt        | BigIntNumber        | Uses a hybrid approach of storing values as Int64 and BigIntStr to ensure no loss of precision. |
| decimal       | DecimalNumber       | Uses a hybrid approach of storing values as Float64, Int64, BigIntStr, DecimalStr, and Decimal to ensure no loss of precision. |
| bytes         | []byte              |                                                              |
| string        | string              |                                                              |
| bool          | bool                |                                                              |

The gap between these two implementation methods mainly lies in numerical adjustments, which are integrated through the following interfaces:


```golang
// NumberConverter: Digital Converter
type NumberConverter interface {
 ConvertBigIntFromInt(i int64) (num BigIntNumber)
 ConvertDecimalFromFloat(f float64) (num DecimalNumber)
 ConvertBigInt(s string) (num BigIntNumber, err error)
 ConvertDecimal(s string) (num DecimalNumber, err error)
}

// Number: Represents a numeric value
type Number interface {
 Bool() (bool, error)
 String() string
}

// BigIntNumber: Represents a high-precision integer
type BigIntNumber interface {
 Number

 Int64() (int64, error)
 Decimal() DecimalNumber
 CloneBigInt() BigIntNumber
 AsBigInt() *big.Int
}

// DecimalNumber: Represents a high-precision decimal number
type DecimalNumber interface {
 Number

 Float64() (float64, error)
 BigInt() BigIntNumber
 CloneDecimal() DecimalNumber
 AsDecimal() decimal.Decimal
}
```
The main implementations are Converter (the current implementation method) and OldConverter (the previous implementation method). Converter outperforms OldConverter in terms of performance. The test results from `number_bench_test.go` are as follows:


```plaintext
BenchmarkConverter_ConvertFromBigInt-4                	34292768	        40.13 ns/op	       8 B/op	       0 allocs/op
BenchmarkOldConverter_ConvertFromBigInt-4             	19314712	        58.69 ns/op	      16 B/op	       1 allocs/op
BenchmarkConverter_ConvertDecimalFromloat-4           	100000000	        15.74 ns/op	       8 B/op	       0 allocs/op
BenchmarkOldConverter_ConvertDecimalFromFloat-4       	 1654504	       725.8 ns/op	      48 B/op	       2 allocs/op
BenchmarkConverter_ConvertBigInt_Int64-4              	 5020077	       230.0 ns/op	      39 B/op	       2 allocs/op
BenchmarkOldConverter_ConvertBigInt_Int64-4           	 2232102	       627.3 ns/op	     111 B/op	       5 allocs/op
BenchmarkCoventor_ConvertBigInt_large_number-4        	   50010	     21211 ns/op	    8064 B/op	     216 allocs/op
BenchmarkOldCoventor_ConvertBigInt_large_number-4     	   23709	     51818 ns/op	    9216 B/op	     360 allocs/op
BenchmarkConverter_ConvertDecimal_Int64-4             	 3830624	       312.6 ns/op	      39 B/op	       2 allocs/op
BenchmarkOldConverter_ConvertDecimal_Int64-4          	 1995441	       611.4 ns/op	     116 B/op	       4 allocs/op
BenchmarkConverter_ConvertDecimal_Float64-4           	 1707649	       671.4 ns/op	     178 B/op	       5 allocs/op
BenchmarkOldConverter_ConvertDecimal_Float64-4        	 1229505	       991.1 ns/op	     191 B/op	       6 allocs/op
BenchmarkConverter_ConvertDecimal-4                   	   80113	     15009 ns/op	    2280 B/op	     144 allocs/op
BenchmarkOldConverter_ConvertDecimal-4                	   56880	     26496 ns/op	    4608 B/op	     288 allocs/op
BenchmarkConverter_ConvertDecimal_large_number-4      	   45754	     22387 ns/op	    5184 B/op	     144 allocs/op
BenchmarkOldConverter_ConvertDecimal_large_number-4   	   16726	     69543 ns/op	   13248 B/op	     432 allocs/op
BenchmarkConverter_ConvertDecimal_Exp-4               	   15516	     86355 ns/op	   18432 B/op	     648 allocs/op
BenchmarkOldConverter_ConvertDecimal_Exp-4            	   17992	     56777 ns/op	   11520 B/op	     432 allocs/op
BenchmarkDecimal_Decmial_String-4                     	 3443062	       361.0 ns/op	      88 B/op	       5 allocs/op
BenchmarkDecimal_DecmialStr_String-4                  	1000000000	         0.6694 ns/op	       0 B/op	       0 allocs/op
BenchmarkDecimal_Float64_String-4                     	 5254669	       260.7 ns/op	      48 B/op	       2 allocs/op
BenchmarkDecimal_Int64_String-4                       	13537401	        89.62 ns/op	      24 B/op	       1 allocs/op
BenchmarkDecimal_BigInt_String-4                      	 4664106	       247.4 ns/op	      56 B/op	       3 allocs/op
BenchmarkDecimal_BigIntStr_String-4                   	1000000000	         0.6873 ns/op	       0 B/op	       0 allocs/op
```
Additionally, if any issues arise, you can revert to the old implementation by modifying the `_DefaultNumberConverter` value in `number.go`.

The relationship between the types and their conversions is as follows:

| from\to | time                                                   | bigInt                                   | decimal                           | bytes                                                      | string                                                     | bool                                                                |
| --- | -------------------------------------------------- | ------------------------------------ | ------------------------------- | -------------------------------------------------- | -------------------------------------------------- | ----------------------------------------------------------- |
| time    | -                                                    | Not supported                        | Not supported                   | Supports conversion of specified time formats (generally supports default time format) | Supports conversion of specified time formats (generally supports default time format) | Not supported                                                   |
| bigInt  | Not supported                                      | -                                    | Supported                       | Supported                                                  | Supported                                                  | Converts non-zero values to true, and zero to false               |
| decimal | Not supported                                      | Rounds to the nearest integer, truncating the decimal part | -                                 | Supported                                                  | Supported                                                  | Converts non-zero values to true, and zero to false               |
| bytes   | Only supports conversion of specified time formats (generally supports default time format) | Real numbers and scientific notation strings are rounded | Real numbers and scientific notation strings | -                                                      | Supported                                                  | Supports conversion of "1", "t", "T", "TRUE", "true", "True" to true, and "0", "f", "F", "FALSE", "false", "False" to false |
| string  | Only supports conversion of specified time formats (generally supports default time format) | Real numbers and scientific notation strings are rounded | Real numbers and scientific notation strings | Supported                                                  | -                                                        | Supports conversion of "1", "t", "T", "TRUE", "true", "True" to "true", and "0", "f", "F", "FALSE", "false", "False" to false |
| bool    | Not supported                                      | true converts to 1, false converts to 0                 | true converts to 1.0, false converts to 0.0 | true converts to "true", false converts to "false" | true converts to "true", false converts to "false" | -                                                           |

**Note: The default time format is 2006-01-02 15:04:05.999999999Z07:00**

This table provides an overview of data type conversions between different formats, including time, bigInt, decimal, bytes, string, and bool. It specifies which conversions are supported and which are not, as well as any specific behavior or limitations associated with each conversion.