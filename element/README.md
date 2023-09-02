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

go-etl支持六种内部数据类型：

- `bigInt`：定点数(int64、int32、int16、int8、BigInt等)。
- `decimal`：浮点数(float32、float63、BigDecimal(无限精度)等)。
- `string`：字符串类型，底层不限长，使用通用字符集(Unicode)。
- `time`：日期类型。
- `bool`：布尔值。
- `bytes`：二进制，可以存放诸如MP3等非结构化数据。

对应地，有`TimeColumnValue`、`BigIntColumnValue`、`DecimalColumnValue`、`BytesColumnValue`、`StringColumnValue`和`BoolColumnValue`六种`ColumnValue`的实现。

这些`ColumnValue`提供一系列以`as`开头的数据类型转换转换方法。

```go
//ColumnValue 列值
type ColumnValue interface {
	fmt.Stringer

	Type() ColumnType                    //列类型
	IsNil() bool                         //是否为空
	AsBool() (bool, error)               //转化为布尔值
	AsBigInt() (*big.Int, error)         //转化为整数
	AsDecimal() (decimal.Decimal, error) //转化为高精度实数
	AsString() (string, error)           //转化为字符串
	AsBytes() ([]byte, error)            //转化为字节流
	AsTime() (time.Time, error)          // 转化为时间
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
| bigInt   | big.Int         | 使用无限精度的大整数，保证不失真  |
| decimal  | decimal.Decimal | 用decimal.Decimal表示，保证不失真 |
| bytes    | []byte          |                                   |
| string   | string          |                                   |
| bool     | bool            |                                   |


+ 目前的实现方式

| 内部类型 | 实现类型        | 备注                              |
| -------- | --------------- | --------------------------------- |
| time     | time.Time       |                                   |
| bigInt   | BigIntNumber    | 使用Int64和BigIntStr交叉保存的方式，保证不失真  |
| decimal  | DecimalNumber   | 使用Float64, Int64，BigIntStr，DecimalStr和Decimal交叉保存的方式，保证不失真 |
| bytes    | []byte          |                                   |
| string   | string          |                                   |
| bool     | bool            |                                   |

这两种实现方式之间的差距主要在数值方面做出了调整，通过以下接口进行了整合：
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
	AsBigInt() *big.Int
}

//DecimalNumber 高精度实数
type DecimalNumber interface {
	Number

	Float64() (float64, error)
	BigInt() BigIntNumber
	CloneDecimal() DecimalNumber
	AsDecimal() decimal.Decimal
}
```
主要实现了NumberConverter的Converter(目前的实现方式)和OldConverter(老的实现方式)，Converter比OldConverter性能更好, 通过number_bench_test.go的测试结果如下：
```
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
另外，如果遇到问题可以通过修改number.go中_DefaultNumberConverter的取值回到老的实现方式

类型之间相互转换的关系如下：

| from\to | time                                             | bigInt                               | decimal                       | bytes                                          | string                                         | bool                                                         |
| ------- | ------------------------------------------------ | ------------------------------------ | ----------------------------- | ---------------------------------------------- | ---------------------------------------------- | ------------------------------------------------------------ |
| time    | -                                                | 不支持                               | 不支持                        | 支持指定时间格式的转化（一般支持默认时间格式） | 支持指定时间格式的转化（一般支持默认时间格式） | 不支持                                                       |
| bigInt  | 不支持                                           | -                                    | 支持                          | 支持                                           | 支持                                           | 支持非0的转化成ture,0转化成false                             |
| decimal | 不支持                                           | 取整，直接截取整数部分               | -                             | 支持                                           | 支持                                           | 支持非0的转化成ture,0转化成false                             |
| bytes   | 仅支持指定时间格式的转化（一般支持默认时间格式） | 实数型以及科学性计数法字符串会被取整 | 实数型以及科学性计数法字符串  | -                                              | 支持                                           | 支持"1"," t", "T", "TRUE", "true", "True"转化为true，"0", "f"," F", "FALSE", "false", "False"转化为false |
| string  | 仅支持指定时间格式的转化（一般支持默认时间格式） | 实数型以及科学性计数法字符串会被取整 | 实数型以及科学性计数法字符串  | 支持                                           | -                                              | 支持"1", "t", "T", "TRUE", "true", "True"转化为"true"，"0", "f", "F", "FALSE", "false", "False"转化为false |
| bool    | 不支持                                           | ture转化为1，false转化为0            | ture转化为1.0，false转化为0.0 | true转化为"true"，false转化为"false"           | true转化为"true"，false转化为"false"           | -                                                            |

**注：默认时间格式为2006-01-02 15:04:05.999999999Z07:00**