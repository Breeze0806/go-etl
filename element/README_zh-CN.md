# go-etl数据类型说明

这个包主要定义了go-etl中的数据类型

## 记录

```
//Record 记录
type Record interface {
	fmt.Stringer

	Add(Column) error                      //新增列
	GetByIndex(i int) (Column, error)      //获取第i个列
	GetByName(name string) (Column, error) //获取列名为name的列
	Set(i int, c Column) error             //设置第i列
	ColumnNumber() int                     //获取列数
	ByteSize() int64                       //字节流大小
	MemorySize() int64                     //内存大小
}
```

## 数据类型转化

go-etl支持七种内部数据类型：

- `bigInt`：定点数(int64、int32、int16、int8、BigInt等)。
- `decimal`：浮点数(float32、float64、BigDecimal(无限精度)等)。
- `string`：字符串类型，底层不限长，使用通用字符集(Unicode)。
- `time`：日期类型。
- `bool`：布尔值。
- `bytes`：二进制，可以存放诸如MP3等非结构化数据。
- `json`：JSON类型，用于存储和处理JSON数据。

对应地，有`TimeColumnValue`、`BigIntColumnValue`、`DecimalColumnValue`、`BytesColumnValue`、`StringColumnValue`、`BoolColumnValue`和`JsonColumnValue`七种`ColumnValue`的实现。

这些`ColumnValue`提供一系列以`as`开头的数据类型转换转换方法。

```go
//ColumnValue 列值
type ColumnValue interface {
	fmt.Stringer

	Type() ColumnType                    //列类型
	IsNil() bool                         //是否为空
	AsBool() (bool, error)               //转化为布尔值
	AsBigInt() (*apd.BigInt, error)         //转化为整数
	AsDecimal() (*apd.Decimal, error) //转化为高精度实数
	AsString() (string, error)           //转化为字符串
	AsBytes() ([]byte, error)            //转化为字节流
	AsTime() (time.Time, error)          // 转化为时间
	AsJSON() (JSON, error)               //转化为JSON
}
```

在ColumnValue的基础上实现了下列方法

```go
//Column 列
type Column interface {
	ColumnValue
	AsInt64() (int64, error)     //转化为64位整数
	AsFloat64() (float64, error) //转化为64位实数
	Clone() (Column, error)      //克隆
	Cmp(Column) (int, error)     //比较, 1代表大于， 0代表相等， -1代表小于
	Name() string                //列名
	ByteSize() int64             //字节流大小
	MemorySize() int64           //内存大小
}
```

DataX的内部类型在实现上会选用不同的golang类型：
现在有两种实现方式，但是老方式在处理大树时存在性能上的问题，目前新实现方式还是beta版,未有较好实践验证。

+ 老的实现方式

| 内部类型 | 实现类型        | 备注                              |
| -------- | --------------- | --------------------------------- |
| time     | time.Time       |                                   |
| bigInt   | apd.BigInt      | 使用无限精度的大整数，保证不失真  |
| decimal  | apd.Decimal     | 用apd.Decimal表示，保证不失真 |
| bytes    | []byte          |                                   |
| string   | string          |                                   |
| bool     | bool            |                                   |
| json     | encoding.JSON   | 包装来自github.com/Breeze0806/go/encoding包的encoding.JSON |


+ 目前的实现方式

| 内部类型 | 实现类型        | 备注                              |
| -------- | --------------- | --------------------------------- |
| time     | time.Time       |                                   |
| bigInt   | BigIntNumber    | 使用Int64和BigIntStr交叉保存的方式，保证不失真  |
| decimal  | DecimalNumber   | 使用Float64, Int64，BigIntStr，DecimalStr和Decimal交叉保存的方式，保证不失真 |
| bytes    | []byte          |                                   |
| string   | string          |                                   |
| bool     | bool            |                                   |
| json     | DefaultJSON     | 包装encoding.JSON以实现element.JSON接口 |

这两种实现方式之间的差距主要在数值方面做出了调整，通过以下接口进行了整合：

### Number接口
```golang
//NumberConverter 数字转化器
type NumberConverter interface {
	ConvertBigIntFromInt(i int64) (num BigIntNumber)
	ConvertDecimalFromFloat(f float64) (num DecimalNumber)
	ConvertBigInt(s string) (num BigIntNumber, err error)
	ConvertDecimal(s string) (num DecimalNumber, err error)
}

//Number 数字
type Number interface {
	Bool() (bool, error)
	String() string
}

//BigIntNumber 高精度整数
type BigIntNumber interface {
	Number

	Int64() (int64, error)
	Decimal() DecimalNumber
	CloneBigInt() BigIntNumber
	AsBigInt() *apd.BigInt
}

//DecimalNumber 高精度实数
type DecimalNumber interface {
	Number

	Float64() (float64, error)
	BigInt() BigIntNumber
	CloneDecimal() DecimalNumber
	AsDecimal() *apd.Decimal
}
```
主要实现了NumberConverter的Converter(目前的实现方式)和OldConverter(老的实现方式)，Converter比OldConverter性能更好, 通过number_bench_test.go的测试结果如下：
```
go test -bench=.
BenchmarkConverter_ConvertBigIntFromInt-32                      1000000000               0.1849 ns/op          0 B/op          0 allocs/op
BenchmarkOldConverter_ConvertBigIntFromInt-32                   1000000000               0.4963 ns/op          0 B/op          0 allocs/op
BenchmarkConverter_ConvertDecimalFromloat-32                    1000000000               0.1903 ns/op          0 B/op          0 allocs/op
BenchmarkOldConverter_ConvertDecimalFromFloat-32                 4170295               292.9 ns/op            32 B/op          1 allocs/op
BenchmarkConverter_ConvertBigInt_Int64-32                       18324470                61.40 ns/op           31 B/op          2 allocs/op
BenchmarkOldConverter_ConvertBigInt_Int64-32                    14466127                73.18 ns/op           55 B/op          3 allocs/op
BenchmarkCoventor_ConvertBigInt_large_number-32                   138158              8414 ns/op           11520 B/op        288 allocs/op
BenchmarkOldCoventor_ConvertBigInt_large_number-32                 51420             22713 ns/op           25344 B/op        648 allocs/op
BenchmarkConverter_ConvertDecimal_Int64-32                      12897178                83.36 ns/op           31 B/op          2 allocs/op
BenchmarkOldConverter_ConvertDecimal_Int64-32                    9517863               117.4 ns/op            63 B/op          3 allocs/op
BenchmarkConverter_ConvertDecimal_Float64-32                     5664470               204.9 ns/op           147 B/op          5 allocs/op
BenchmarkOldConverter_ConvertDecimal_Float64-32                  2939163               403.6 ns/op           262 B/op          8 allocs/op
BenchmarkConverter_ConvertDecimal-32                              328171              3620 ns/op            2280 B/op        144 allocs/op
BenchmarkOldConverter_ConvertDecimal-32                           131666              7630 ns/op            3424 B/op        216 allocs/op
BenchmarkConverter_ConvertDecimal_large_number-32                 164592              7290 ns/op            5184 B/op        144 allocs/op
BenchmarkOldConverter_ConvertDecimal_large_number-32               36433             31969 ns/op           29376 B/op        720 allocs/op
BenchmarkConverter_ConvertDecimal_Exp-32                           41289             27639 ns/op           20736 B/op        720 allocs/op
BenchmarkOldConverter_ConvertDecimal_Exp-32                        49203             24741 ns/op           17280 B/op        576 allocs/op
BenchmarkDecimal_Decmial_String-32                               8836029               119.8 ns/op           112 B/op          4 allocs/op
BenchmarkDecimal_DecmialStr_String-32                           1000000000               0.2572 ns/op          0 B/op          0 allocs/op
BenchmarkDecimal_Float64_String-32                               3808023               300.3 ns/op            88 B/op          3 allocs/op
BenchmarkDecimal_Int64_String-32                                52006051                22.84 ns/op           24 B/op          1 allocs/op
BenchmarkDecimal_BigInt_String-32                               16544946                63.10 ns/op           48 B/op          2 allocs/op
BenchmarkDecimal_BigIntStr_String-32                            1000000000               0.2596 ns/op          0 B/op          0 allocs/op
```
另外，如果遇到问题可以通过修改number.go中_DefaultNumberConverter的取值回到老的实现方式

### JSON接口

```go
// JSON JSON接口
type JSON interface {
	ToString() string
	ToBytes() []byte
	Clone() JSON
}
// JSONConverter JSON转换器接口
type JSONConverter interface {
	ConvertFromString(s string) (json JSON, err error)
	ConvertFromBytes(b []byte) (json JSON, err error)
}
```

`JSON`接口提供将JSON值转换为字符串和字节的方法，以及克隆JSON值的功能。

`JSONConverter`提供将字符串和字节转换为JSON值的功能。

JSON类型实现包装了`github.com/Breeze0806/go/encoding`包，以提供强大的JSON解析和操作能力。

### 类型转换关系表

类型之间相互转换的关系如下：

| from\to | time                                             | bigInt                               | decimal                       | bytes                                          | string                                         | bool                                                         | json                                                           |
| ------- | ------------------------------------------------ | ------------------------------------ | ----------------------------- | ---------------------------------------------- | ---------------------------------------------- | ------------------------------------------------------------ | ------------------------------------------------------------ |
| time    | -                                                | 不支持                               | 不支持                        | 支持指定时间格式的转化（一般支持默认时间格式） | 支持指定时间格式的转化（一般支持默认时间格式） | 不支持                                                       | 不支持                                                       |
| bigInt  | 不支持                                           | -                                    | 支持                          | 支持                                           | 支持                                           | 支持非0的转化成ture,0转化成false                             | 不支持                                                       |
| decimal | 不支持                                           | 取整，直接截取整数部分               | -                             | 支持                                           | 支持                                           | 支持非0的转化成ture,0转化成false                             | 不支持                                                       |
| bytes   | 仅支持指定时间格式的转化（一般支持默认时间格式） | 实数型以及科学性计数法字符串会被取整 | 实数型以及科学性计数法字符串  | -                                              | 支持                                           | 支持"1"," t", "T", "TRUE", "true", "True"转化为true，"0", "f"," F", "FALSE", "false", "False"转化为false | 支持（解析为JSON）                                     |
| string  | 仅支持指定时间格式的转化（一般支持默认时间格式） | 实数型以及科学性计数法字符串会被取整 | 实数型以及科学性计数法字符串  | 支持                                           | -                                              | 支持"1", "t", "T", "TRUE", "true", "True"转化为"true"，"0", "f", "F", "FALSE", "false", "False"转化为"false" | 支持（解析为JSON）                                     |
| bool    | 不支持                                           | ture转化为1，false转化为0            | ture转化为1.0，false转化为0.0 | true转化为"true"，false转化为"false"           | true转化为"true"，false转化为"false"           | -                                                            | 不支持                                                       |
| json    | 不支持                                           | 不支持                               | 不支持                        | 支持（输出为JSON字节）                               | 支持（输出为JSON字符串）                               | 不支持                                                   | -                                                            |

**注：默认时间格式为`2006-01-02 15:04:05.999999999Z07:00`**

此表提供了time、bigInt、decimal、bytes、string、bool和json不同格式之间数据类型转换的概述，包括支持哪些转换、不支持哪些转换，以及与每个转换相关的特定行为或限制。