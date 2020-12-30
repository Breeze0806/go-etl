package element

type Record interface {
	Add(Column) error
	GetByIndex(i int) (Column, error)
	GetByName(name string) (Column, error)
	Set(i int, c Column) error
	ColumnNumber() int
	ByteSize() int64
	MemorySize() int64
}

var singleTerminateRecord = &TerminateRecord{}

func GetTerminateRecord() Record {
	return singleTerminateRecord
}

type TerminateRecord struct{}

func (t *TerminateRecord) Add(Column) error {
	return nil
}

func (t *TerminateRecord) GetByIndex(i int) (Column, error) {
	return nil, nil
}

func (t *TerminateRecord) GetByName(name string) (Column, error) {
	return nil, nil
}

func (t *TerminateRecord) Set(i int, c Column) error {
	return nil
}

func (t *TerminateRecord) ColumnNumber() int {
	return 0
}

func (t *TerminateRecord) ByteSize() int64 {
	return 0
}

func (t *TerminateRecord) MemorySize() int64 {
	return 0
}

type DefaultRecord struct {
	names      []string
	columns    map[string]Column
	byteSize   int64
	memorySize int64
}

func NewDefaultRecord() *DefaultRecord {
	return &DefaultRecord{
		names:   make([]string, 0),
		columns: make(map[string]Column),
	}
}

func (r *DefaultRecord) Add(c Column) error {
	if _, ok := r.columns[c.Name()]; ok {
		return ErrColumnExist
	}
	r.names = append(r.names, c.Name())
	r.columns[c.Name()] = c
	r.incSize(c)
	return nil
}

func (r *DefaultRecord) GetByIndex(i int) (Column, error) {
	if i >= len(r.names) || i < 0 {
		return nil, ErrIndexOutOfRange
	}
	if v, ok := r.columns[r.names[i]]; ok {
		return v, nil
	}
	return nil, ErrColumnNotExist
}

func (r *DefaultRecord) GetByName(name string) (Column, error) {
	if v, ok := r.columns[name]; ok {
		return v, nil
	}
	return nil, ErrColumnNotExist
}

func (r *DefaultRecord) Set(i int, c Column) error {
	if i >= len(r.names) || i < 0 {
		return ErrIndexOutOfRange
	}

	if v, ok := r.columns[r.names[i]]; ok {
		r.decSize(v)
	}
	r.names[i] = c.Name()
	r.columns[c.Name()] = c
	r.incSize(c)
	return nil
}

func (r *DefaultRecord) ColumnNumber() int {
	return len(r.columns)
}

func (r *DefaultRecord) ByteSize() int64 {
	return r.byteSize
}

func (r *DefaultRecord) MemorySize() int64 {
	return r.memorySize
}

func (r *DefaultRecord) incSize(c Column) {
	r.byteSize += c.ByteSize()
	r.memorySize += c.MemorySize()
}

func (r *DefaultRecord) decSize(c Column) {
	r.byteSize -= c.ByteSize()
	r.memorySize -= c.MemorySize()
}
