package plugin

//Type 插件类型
type Type string

//插件类型枚举
var (
	Reader      Type = "reader"      //读取器
	Writer      Type = "writer"      //写入器
	Transformer Type = "transformer" //转化器
	Handler     Type = "handler"     //处理器
)

//NewType 新增类型
func NewType(s string) Type {
	return Type(s)
}

func (t Type) String() string {
	return string(t)
}

//IsValid 是否合法
func (t Type) IsValid() bool {
	switch t {
	case Reader, Writer, Transformer, Handler:
		return true
	}
	return false
}
