package element

type Recode interface {
	Add(Column) error
	GetByIndex(i int) (Column, error)
	GetByName(name string) (Column, error)
	Set(i int, c Column) error
	ColumnNumber() int
	ByteSize() int64
	MemorySize() int64
}

type DefaultRecode struct {
	names      []string
	columns    map[string]Column
	byteSize   int64
	memorySize int64
}

func NewDefaultRecode() *DefaultRecode {
	return &DefaultRecode{
		names:   make([]string, 0),
		columns: make(map[string]Column),
	}
}

func (r *DefaultRecode) Add(c Column) error {
	r.names = append(r.names, c.Name())
	if _, ok := r.columns[c.Name()]; ok {
		return ErrColumnExist
	}
	r.columns[c.Name()] = c
	r.incSize(c)
	return nil
}

func (r *DefaultRecode) GetByIndex(i int) (Column, error) {
	if i >= len(r.names) || i < 0 {
		return nil, ErrIndexOutOfRange
	}
	if v, ok := r.columns[r.names[i]]; ok {
		return v, nil
	}
	return nil, ErrColumnNotExist
}

func (r *DefaultRecode) GetByName(name string) (Column, error) {
	if v, ok := r.columns[name]; ok {
		return v, nil
	}
	return nil, ErrColumnNotExist
}

func (r *DefaultRecode) Set(i int, c Column) error {
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

func (r *DefaultRecode) ColumnNumber() int {
	return len(r.columns)
}

func (r *DefaultRecode) ByteSize() int64 {
	return r.byteSize
}

func (r *DefaultRecode) MemorySize() int64 {
	return r.memorySize
}

func (r *DefaultRecode) incSize(c Column) {
	r.byteSize += c.ByteSize()
	r.memorySize += c.MemorySize()
}

func (r *DefaultRecode) decSize(c Column) {
	r.byteSize -= c.ByteSize()
	r.memorySize -= c.MemorySize()
}
