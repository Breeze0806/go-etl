# 类二维表文件存储开发者指南

类二维表文件存储是一个对类二维表文件流进行抽象的框架，这个框架可以支持各类类二维表文件格式的读取和写入

## 输入文件流

```go
// Opener 用于打开一个输入流的打开器
type Opener interface {
	Open(filename string) (stream InStream, err error) //打开文件名filename的输入流
}

// InStream 输入流
type InStream interface {
	Rows(conf *config.JSON) (rows Rows, err error) //获取行读取器
	Close() (err error)                            //关闭输入流
}

// Rows 行读取器
type Rows interface {
	Next() bool                                  //获取下一行，如果没有返回false，有返回true
	Scan() (columns []element.Column, err error) //扫描出每一行的列
	Error() error                                //获取下一行的错误
	Close() error                                //关闭行读取器
}
```

InStream输入流可以通过传入json配置文件来获取行读取器Rows，在Rows把一个行数据转化成一条记录，可以参考csv包的实现。另外，需要通过把Opener注册。

```go
func init() {
	var opener Opener
	file.RegisterOpener("csv", &opener)
}
```

## 输出文件流

```go
// Creator 创建输出流的创建器
type Creator interface {
	Create(filename string) (stream OutStream, err error) //创建名为filename的输出流
}

// OutStream 输出流
type OutStream interface {
	Writer(conf *config.JSON) (writer StreamWriter, err error) //创建写入器
	Close() (err error)                                        //关闭输出流
}

// StreamWriter 输出流写入器
type StreamWriter interface {
	Write(record element.Record) (err error) //写入记录
	Flush() (err error)                      //刷新至文件
	Close() (err error)                      //关闭输出流写入器
}
```

OutStream输入流可以通过传入json配置文件来获取输出流写入器StreamWriter，在StreamWriter把成一条记录转化一个行数据,可以参考csv包的实现。另外，需要通过把Creator注册。

```go
func init() {
	var creator Creator
	file.RegisterCreator("csv", &creator)
}
```

