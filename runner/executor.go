package runner

// Executor is a struct for executing tasks.
type executor struct {
	params map[string]string
}

// NewExecutor creates a new executor.
func NewExecutor(m map[string]string) *executor {
	return &executor{params: m}
}
