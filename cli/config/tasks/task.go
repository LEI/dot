package tasks

import (
	"fmt"
	"reflect"

	"github.com/LEI/dot/cli"
)

var (
	// Verbose ...
	Verbose bool

	// ErrSkip ...
	ErrSkip = fmt.Errorf("skip")
)

// Tasker interface
type Tasker interface {
	Check() error
	// Execute() error
	Install() error
	Remove() error
}

// Task struct
type Task struct {
	Tasker
	// execute bool
	toDo bool
}

// DoInstall ...
func (t *Task) DoInstall() bool {
	return !t.toDo
}

// DoRemove ...
func (t *Task) DoRemove() bool {
	return t.toDo
}

// TaskList ...
type TaskList []Tasker

// type Tasks struct {
// 	value []*Task
// }

// Check tasks
func Check(i interface{}) error {
	// tl := i.(*[]Tasker)
	tl, err := taskList(i)
	if err != nil {
		return err
	}
	errs := cli.Errors{}
	for _, t := range tl {
		if err := t.Check(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf(errs.Error())
	}
	return nil
}

// Install tasks
func Install(i interface{}) error {
	// tl := i.(TaskList)
	tl, err := taskList(i)
	if err != nil {
		return err
	}
	for _, t := range tl {
		err := t.Install()
		switch err {
		case nil:
			fallthrough
		case ErrSkip:
			// fmt.Println("skipped install", t)
			continue
		default:
			return err
		}
	}
	return nil
}

// Remove tasks
func Remove(i interface{}) error {
	// tl := i.(TaskList)
	tl, err := taskList(i)
	if err != nil {
		return err
	}
	for _, t := range tl {
		err := t.Remove()
		switch err {
		case nil:
			fallthrough
		case ErrSkip:
			// fmt.Println("skipped remove", t)
			continue
		default:
			return err
		}
	}
	return nil
}

// https://ahmet.im/blog/golang-take-slices-of-any-type-as-input-parameter/
func taskList(i interface{}) (TaskList, error) {
	val := reflect.ValueOf(i)
	if val.Kind() != reflect.Slice {
		return nil, fmt.Errorf("i (%s): %+v", val.Kind(), i)
	}
	// slice = val
	c := val.Len()
	tl := make(TaskList, c)
	// for i, v := range val {
	// 	tl[i] = v
	// }
	for i := 0; i < c; i++ {
		tl[i] = val.Index(i).Interface().(Tasker)
	}
	return tl, nil
}
