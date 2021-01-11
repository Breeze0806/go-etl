package mysql

import "github.com/Breeze0806/go-etl/storage/database"

type Table struct {
	*database.BaseTable
}

func NewTable(b *database.BaseTable) database.Table {
	return &Table{
		BaseTable: database.NewBaseTable(b.Instance(), b.Schema(), b.Name()),
	}
}

func (t *Table) Quoted() string {
	return Quoted(t.Instance()) + "." + Quoted(t.Name())
}

func (t *Table) String() string {
	return t.Quoted()
}

func (t *Table) AddField(baseField *database.BaseField) {
	t.AppendField(NewField(baseField))
}
