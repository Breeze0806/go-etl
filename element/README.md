# element包说明
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
	AsInt8() (int8, error)       //转化为8位整数
	AsInt16() (int16, error)     //转化为16位整数
	AsInt32() (int32, error)     //转化为32位整数
	AsInt64() (int64, error)     //转化为64位整数
	AsFloat32() (float32, error) //转化为32位实数
	AsFloat64() (float64, error) //转化为64位实数
	Clone() (Column, error)      //克隆
	Cmp(Column) (int, error)     //比较, 1代表大于， 0代表相等， -1代表小于
	Name() string                //列名
	ByteSize() int64             //字节流大小
	MemorySize() int64           //内存大小
}
```

DataX的内部类型在实现上会选用不同的golang类型：

| 内部类型 | 实现类型        | 备注                              |
| -------- | --------------- | --------------------------------- |
| time     | time.Time       |                                   |
| bigInt   | big.Int         | 使用无限精度的大整数，保证不失真  |
| decimal  | decimal.Decimal | 用decimal.Decimal表示，保证不失真 |
| bytes    | []byte          |                                   |
| string   | string          |                                   |
| bool     | bool            |                                   |

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