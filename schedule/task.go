package schedule

//Task 任务
type Task interface {
	Do() error //同步执行
}

//AsyncTask 异步任务
type AsyncTask interface {
	Do() error   //同步执行
	Post() error //后续通知
}
