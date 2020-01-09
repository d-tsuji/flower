// Package component is the implementation of the method that is executed as a task in the workflow.
//
// The task corresponds to a method of one component structure.
// Task methods are implemented as follows:
//
//		func (c *component) EchoRandomTimeSleep() error {
//			randTime := random.Intn(5) + 1
//
//			fmt.Printf("[component] start EchoRandomTimeSleep. (%v second sleep)\n", randTime)
//			time.Sleep(time.Duration(randTime) * time.Second)
//			fmt.Printf("[component] finish EchoRandomTimeSleep\n")
//
//			return nil
//		}
//
// Register the component method implemented as a task in "ms_task_definition" as a workflow.
// The data structure of "ms_task_definition" is as follows. The *component method is specified
// in the "ms_task_definition" program column.
//
//		CREATE TABLE IF NOT EXISTS ms_task_definition (
//		task_id                  varchar(256) NOT NULL
//		,   task_seq                 numeric NOT NULL
//		,   program                  varchar(256) NOT NULL
//		,   task_priority            numeric NOT NULL
//		,   param1_key               varchar(1024)
//		,   param1_value             varchar(1024)
//		,   param2_key               varchar(1024)
//		,   param2_value             varchar(1024)
//		,   param3_key               varchar(1024)
//		,   param3_value             varchar(1024)
//		,   param4_key               varchar(1024)
//		,   param4_value             varchar(1024)
//		,   param5_key               varchar(1024)
//		,   param5_value             varchar(1024)
//
//		,   PRIMARY KEY (task_id, task_seq)
//		);
//
// To register a task of EchoRandomTimeSleep, add a record as follows.
//
//		INSERT INTO ms_task_definition(task_id, task_seq, program) VALUES ('example', 1, 'EchoRandomTimeSleep');
package component

type component struct {
	params map[string]string
}

// NewComponent creates a new component.
func NewComponent(m map[string]string) *component {
	return &component{params: m}
}
