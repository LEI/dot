package task

import (
	// "fmt"
)

// type ResourceProvider interface {
// 	Clone(res *Res) error
// 	Update(res *Res) error
// }

type Context interface {
	Check() bool
	Sync() error
}

// func (t *Task) String() string {
// 	return fmt.Sprintf("%+v", t)
// }
