package element

//Record 记录
type Record interface {
	Add(Column) error                      //新增列
	GetByIndex(i int) (Column, error)      //获取第i个列
	GetByName(name string) (Column, error) //获取列名为name的列
	Set(i int, c Column) error             //设置第i列
	ColumnNumber() int                     //获取列数
	ByteSize() int64                       //字节流大小
	MemorySize() int64                     //内存大小
}

var singleTerminateRecord = &TerminateRecord{}

//GetTerminateRecord 获取终止记录
func GetTerminateRecord() Record {
	return singleTerminateRecord
}

//TerminateRecord 终止记录
type TerminateRecord struct{}

//Add 空方法
func (t *TerminateRecord) Add(Column) error {
	return nil
}

//GetByIndex 空方法
func (t *TerminateRecord) GetByIndex(i int) (Column, error) {
	return nil, nil
}

//GetByName 空方法
func (t *TerminateRecord) GetByName(name string) (Column, error) {
	return nil, nil
}

//Set 空方法
func (t *TerminateRecord) Set(i int, c Column) error {
	return nil
}

//ColumnNumber 空方法
func (t *TerminateRecord) ColumnNumber() int {
	return 0
}

//ByteSize 空方法
func (t *TerminateRecord) ByteSize() int64 {
	return 0
}

//MemorySize 空方法
func (t *TerminateRecord) MemorySize() int64 {
	return 0
}

//DefaultRecord 默认记录
type DefaultRecord struct {
	names      []string          //列名数组
	columns    map[string]Column //列映射
	byteSize   int64             //字节流大小
	memorySize int64             //内存大小
}

//NewDefaultRecord 创建默认记录
func NewDefaultRecord() *DefaultRecord {
	return &DefaultRecord{
		names:   make([]string, 0),
		columns: make(map[string]Column),
	}
}

//Add 新增列c,若列c已经存在，就会报错
func (r *DefaultRecord) Add(c Column) error {
	if _, ok := r.columns[c.Name()]; ok {
		return ErrColumnExist
	}
	r.names = append(r.names, c.Name())
	r.columns[c.Name()] = c
	r.incSize(c)
	return nil
}

//GetByIndex 获取第i列,若索引i超出范围或者不存在，就会报错
func (r *DefaultRecord) GetByIndex(i int) (Column, error) {
	if i >= len(r.names) || i < 0 {
		return nil, ErrIndexOutOfRange
	}
	if v, ok := r.columns[r.names[i]]; ok {
		return v, nil
	}
	return nil, ErrColumnNotExist
}

//GetByName 获取列名为name的列,若列名为name的列不存在，就会报错
func (r *DefaultRecord) GetByName(name string) (Column, error) {
	if v, ok := r.columns[name]; ok {
		return v, nil
	}
	return nil, ErrColumnNotExist
}

//Set 设置第i列,若索引i超出范围，就会报错
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

//ColumnNumber 列数量
func (r *DefaultRecord) ColumnNumber() int {
	return len(r.columns)
}

//ByteSize 字节流大小
func (r *DefaultRecord) ByteSize() int64 {
	return r.byteSize
}

//MemorySize 内存大小
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
