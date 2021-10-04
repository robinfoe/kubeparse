package task

type Task interface {
	proceed() error
}
