package mysql

import (
	"fmt"
	"time"

	"github.com/Breeze0806/go-etl/element"
	"github.com/Breeze0806/go-etl/storage/database"
)

type Field struct {
	*database.BaseField
}

func NewField(bf *database.BaseField) *Field {
	return &Field{
		BaseField: bf,
	}
}

func (f *Field) Quoted() string {
	return Quoted(f.Name())
}

func (f *Field) BindVar(_ int) string {
	return "?"
}

func (f *Field) Select() string {
	return Quoted(f.Name())
}

func (f *Field) Type() database.FieldType {
	return NewFieldType(f.FieldType())
}

func (f *Field) Scanner() database.Scanner {
	return NewScanner(f)
}

func (f *Field) Valuer(c element.Column) database.Valuer {
	return database.NewGoValuer(f, c)
}

type FieldType struct {
	*database.BaseFieldType
	goType database.GoType
}

func NewFieldType(fieldType database.FieldType) *FieldType {
	return &FieldType{
		BaseFieldType: database.NewBaseFieldType(fieldType),
	}
}

func (f *FieldType) GoType() database.GoType {
	switch f.DatabaseTypeName() {
	//由于存在非负整数，如果直接变为对应的int类型，则会导致转化错误
	//TIME存在负数无法正常转化，YEAR就是TINYINT
	//todo: test YEAR
	case "MEDIUMINT", "INT", "BIGINT", "SMALLINT", "TINYINT",
		"TEXT", "LONGTEXT", "MEDIUMTEXT", "TINYTEXT", "CHAR", "VARCHAR",
		"TIME", "YEAR",
		"DECIMAL":
		return database.GoTypeString
	case "BLOB", "LONGBLOB", "MEDIUMBLOB", "BINARY", "TINYBLOB", "VARBINARY":
		return database.GoTypeBytes
	case "DOUBLE", "FLOAT":
		return database.GoTypeFloat64
	case "DATE", "DATETIME", "TIMESTAMP":
		return database.GoTypeTime
	}
	return database.GoTypeUnknow
}

type Scanner struct {
	f *Field
	database.BaseScanner
}

func NewScanner(f *Field) *Scanner {
	return &Scanner{
		f: f,
	}
}

func (s *Scanner) Scan(src interface{}) (err error) {
	var cv element.ColumnValue
	//todo: byteSize is 0, fix it
	var byteSize int
	switch s.f.Type().DatabaseTypeName() {
	//todo: test year
	case "MEDIUMINT", "INT", "BIGINT", "SMALLINT", "TINYINT", "YEAR":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilBigIntColumnValue()
		case []byte:
			if cv, err = element.NewBigIntColumnValueFromString(string(data)); err != nil {
				return
			}
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeBigInt)
		}
	case "BLOB", "LONGBLOB", "MEDIUMBLOB", "BINARY", "TINYBLOB", "VARBINARY":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilBytesColumnValue()
		case []byte:
			cv = element.NewBytesColumnValue(data)
		default:
			return fmt.Errorf("src is %v(%T),but not %v", src, src, element.TypeBytes)
		}
	case "DATE", "DATETIME", "TIMESTAMP":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilTimeColumnValue()
		case time.Time:
			cv = element.NewTimeColumnValue(data)
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeTime)
		}
	case "TEXT", "LONGTEXT", "MEDIUMTEXT", "TINYTEXT", "CHAR", "VARCHAR", "TIME":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilStringColumnValue()
		case []byte:
			cv = element.NewStringColumnValue(string(data))
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeString)
		}
	case "DOUBLE", "FLOAT", "DECIMAL":
		switch data := src.(type) {
		case nil:
			cv = element.NewNilDecimalColumnValue()
		case []byte:
			if cv, err = element.NewDecimalColumnValueFromString(string(data)); err != nil {
				return
			}
		default:
			return fmt.Errorf("src is %v(%T), but not %v", src, src, element.TypeDecimal)
		}
	}
	s.SetColumn(element.NewDefaultColumn(cv, s.f.Name(), byteSize))
	return
}
