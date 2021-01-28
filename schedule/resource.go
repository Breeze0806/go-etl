package schedule

//Resource 资源
type Resource interface {
	Close() error //关闭释放资源
}

//MappedResource 可映射资源
type MappedResource interface {
	Resource

	Key() string //关键字
}
