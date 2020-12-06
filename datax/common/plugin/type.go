package plugin

type Type string

var (
	Reader      Type = "reader"
	Writer      Type = "writer"
	Transformer Type = "transformer"
	Handler     Type = "handler"
)

func NewType(s string) Type {
	return Type(s)
}

func (t Type) String() string {
	return string(t)
}

func (t Type) IsValid() bool {
	switch t {
	case Reader, Writer, Transformer, Handler:
		return true
	}
	return false
}
