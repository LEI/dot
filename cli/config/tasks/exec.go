package tasks

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/LEI/dot/cli/config/types"
	// "github.com/LEI/dot/internal/executils"
	// "github.com/LEI/dot/internal/ostype"
	"github.com/LEI/dot/system"
	"github.com/mitchellh/mapstructure"
)

var (
	// ExecDir ...
	ExecDir string

	defaultExecShell = "/bin/sh"
)

// Exec task
type Exec struct {
	Task
	Command     string
	Shell       string
	Action      string // install, remove
	types.HasOS `mapstructure:",squash"`
}

func (e *Exec) String() string {
	return fmt.Sprintf("exec[%s]", e.Command)
}

// Check copy task
func (e *Exec) Check() error {
	if e.Command == "" {
		return fmt.Errorf("exec: empty command")
	}
	if !e.CheckOS() { // len(e.OS) > 0 && !ostype.Has(e.OS...) {
		return fmt.Errorf("exec %s: only for %s", e.Command, e.OS)
	}
	// err := system.CheckExec(e.Command)
	// switch err {
	// // case system.ErrDone:
	// // 	e.Done() // Mark task as already executed
	// default:
	// 	return err
	// }
	return nil
}

// Install copy task
func (e *Exec) Install() error {
	if !e.CheckOS() { // len(e.OS) > 0 && !ostype.Has(e.OS...) {
		return ErrSkip
	}
	if e.Action != "" && e.Action != "install" {
		return ErrSkip
	}
	str := strings.TrimSuffix(e.Command, "\n")
	if !e.ShouldInstall() {
		if Verbose > 0 {
			fmt.Fprintf(Stdout, "# %s\n", str)
		}
		return ErrSkip
	}
	fmt.Fprintf(Stdout, "$ %s\n", str)
	// return system.Exec(e.Shell, e.Command)
	if system.DryRun {
		return nil
	}
	if e.Shell == "" {
		e.Shell = defaultExecShell
	}
	cmd := exec.Command(e.Shell, []string{"-c", e.Command}...)
	cmd.Stdout = Stdout
	cmd.Stderr = Stderr
	if ExecDir != "" {
		cmd.Dir = ExecDir
	}
	return cmd.Run()
}

// Remove copy task
func (e *Exec) Remove() error {
	if !e.CheckOS() { // len(e.OS) > 0 && !ostype.Has(e.OS...) {
		return ErrSkip
	}
	if e.Action != "" && e.Action != "remove" {
		return nil
	}
	str := strings.TrimSuffix(e.Command, "\n")
	if !e.ShouldRemove() {
		if Verbose > 0 {
			fmt.Fprintf(Stdout, "# %s\n", str)
		}
		return ErrSkip
	}
	fmt.Fprintf(Stdout, "$ %s\n", str)
	// return system.Exec(e.Shell, e.Command)
	if system.DryRun {
		return nil
	}
	if e.Shell == "" {
		e.Shell = defaultExecShell
	}
	cmd := exec.Command(e.Shell, []string{"-c", e.Command}...)
	cmd.Stdout = Stdout
	cmd.Stderr = Stderr
	if ExecDir != "" {
		cmd.Dir = ExecDir
	}
	return cmd.Run()
}

// Commands task slice
type Commands []*Exec

func (commands *Commands) String() string {
	// s := ""
	// for i, c := range *commands {
	// 	s += fmt.Sprintf("%s", c)
	// 	if i > 0 {
	// 		s += "\n"
	// 	}
	// }
	// return s
	return fmt.Sprintf("%s", *commands)
}

// Parse copy tasks
func (commands *Commands) Parse(i interface{}) error {
	cc := &Commands{}
	s, err := types.NewSliceMap(i)
	if err != nil {
		return err
	}
	for _, v := range *s {
		e := &Exec{}
		switch val := v.(type) {
		case string:
			// e = &Exec{Command: val}
			e.Command = val
		case *types.Map:
			mapstructure.Decode(val, &e)
		case interface{}:
			e = val.(*Exec)
		default:
			return fmt.Errorf("invalid exec map: %+v", val)
		}
		// fmt.Printf("COMMAND [%+v] = [%+v]\n", v, *e)
		cc.Add(*e)
	}
	*commands = *cc
	return nil
}

// Add a command to execute
func (commands *Commands) Add(c Exec) {
	*commands = append(*commands, &c)
}
