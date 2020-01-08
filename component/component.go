// component is the implementation of the method that is executed as a task in the workflow
package component

type component struct {
	params map[string]string
}

// NewComponent creates a new component.
func NewComponent(m map[string]string) *component {
	return &component{params: m}
}
