package parsers

import (
	"fmt"
	"reflect"
)

// Cmd ...
type Cmd struct {
	Command string
	Shell   string
	OS     Slice
	Action string // install, remove
}

// NewCmd ...
func NewCmd(i interface{}) (*Cmd, error) {
	cmd := &Cmd{}
	if i == nil {
		return cmd, fmt.Errorf("trying to add nil cmd: %+v", i)
	}
	if val, ok := i.(string); ok {
		cmd.Command = val
	} else if val, ok := i.(Cmd); ok {
		*cmd = val
	} else if val, ok := i.(map[interface{}]interface{}); ok {
		// Get command
		cmdName, ok := val["command"].(string)
		if !ok {
			return cmd, fmt.Errorf("missing cmd command: %+v", val)
		}
		cmd.Command = cmdName
		cmdShell, ok := val["shell"].(string)
		if ok && cmdShell != "" {
			cmd.Shell = cmdShell
		}
		cmdOS, err := NewSlice(val["os"])
		if err != nil {
			return cmd, err
		}
		cmd.OS = *cmdOS
		cmdAction, ok := val["action"].(string)
		if ok {
			cmd.Action = cmdAction
		}
	} else {
		return cmd, fmt.Errorf("unable to assert Cmd: %+v", i)
	}
	return cmd, nil
}

// Commands ...
type Commands []*Cmd

// Add ...
func (p *Commands) Add(i interface{}) error {
	cmd, err := NewCmd(i)
	if err != nil {
		return err
	}
	*p = append(*p, cmd)
	return nil
}

// UnmarshalYAML ...
func (p *Commands) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var i interface{}
	if err := unmarshal(&i); err != nil {
		return err
	}
	switch val := i.(type) {
	case []string:
		for _, v := range val {
			if err := p.Add(v); err != nil {
				return err
			}
		}
	case []interface{}:
		for _, v := range val {
			if err := p.Add(v); err != nil {
				return err
			}
		}
	default:
		t := reflect.TypeOf(val)
		T := t.Elem()
		if t.Kind() == reflect.Map {
			T = reflect.MapOf(t.Key(), t.Elem())
		}
		return fmt.Errorf("unable to unmarshal packages (%s) into struct: %+v", T, val)
	}
	return nil
}
