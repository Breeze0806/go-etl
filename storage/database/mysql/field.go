package mysql

import (
	"database/sql"
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

func (f *Field) BindVar(i int) string {
	return "?"
}

func (f *Field) Select() string {
	return Quoted(f.Name())
}

func (f *Field) Type() database.FieldType {
	return NewFieldType(f.ColumnType())
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

func NewFieldType(columnType *sql.ColumnType) *FieldType {
	return &FieldType{
		BaseFieldType: database.NewBaseFieldType(columnType),
	}
}

func (f *FieldType) GoType() database.GoType {
	switch f.DatabaseTypeName() {
	//由于存在非负整数，如果直接变为对应的int类型，则会导致转化错误
	case "MEDIUMINT", "INT", "BIGINT", "SMALLINT", "TINYINT",
		"TEXT", "LONGTEXT", "MEDIUMTEXT", "TINYTEXT", "CHAR", "VARCHAR",
		//TIME存在负数无法正常转化，YEAR就是TINYINT
		//todo: test year
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
	*database.BaseScanner
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
		data, ok := src.([]byte)
		if !ok {
			return fmt.Errorf("src is %v(%T), but is not []byte", src, src)
		}
		if cv, err = element.NewBigIntColumnValueFromString(string(data)); err != nil {
			return
		}
	case "BLOB", "LONGBLOB", "MEDIUMBLOB", "BINARY", "TINYBLOB", "VARBINARY":
		data, ok := src.([]byte)
		if !ok {
			return fmt.Errorf("src is %v(%T), but is not []byte", src, src)
		}
		cv = element.NewBytesColumnValue(data)
	case "DATE", "DATETIME", "TIMESTAMP":
		data, ok := src.(time.Time)
		if !ok {
			return fmt.Errorf("src is %v(%T), but is not []byte", src, src)
		}
		cv = element.NewTimeColumnValue(data)
	case "TEXT", "LONGTEXT", "MEDIUMTEXT", "TINYTEXT", "CHAR", "VARCHAR", "TIME":
		data, ok := src.([]byte)
		if !ok {
			return fmt.Errorf("src is %v(%T), but is not []byte", src, src)
		}
		cv = element.NewStringColumnValue(string(data))
	case "DOUBLE", "FLOAT", "DECIMAL":
		data, ok := src.([]byte)
		if !ok {
			return fmt.Errorf("src is %v(%T), but is not []byte", src, src)
		}
		if cv, err = element.NewDecimalColumnValueFromString(string(data)); err != nil {
			return
		}
	}
	s.SetColumn(element.NewDefaultColumn(cv, s.f.Name(), byteSize))
	return
}
