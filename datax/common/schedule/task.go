package schedule

type Task interface {
	Do() error
}
