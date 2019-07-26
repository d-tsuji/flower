package repository

type StatusType string

const (
	WaitExecute StatusType = "0"
	Executable  StatusType = "1"
	Running     StatusType = "2"
	Completed   StatusType = "3"
	Error       StatusType = "9"
)
