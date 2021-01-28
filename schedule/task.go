package schedule

//Task 任务
type Task interface {
	Do() error //同步执行
}
