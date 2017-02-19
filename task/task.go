package task

import fmt (
	"fmt"
)

// type ResourceProvider interface {
// 	Clone(res *Res) error
// 	Update(res *Res) error
// }

type Tasker interface {
	Check() bool
	Sync() error
}

func (t *Task) String() string {
	return fmt.Sprintf("%+v", t)
}
