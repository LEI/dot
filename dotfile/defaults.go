package dotfile

import (
	"fmt"
	"io/ioutil"
	// "os"
	// "os/exec"
	// "path/filepath"

	"gopkg.in/yaml.v2"

	"github.com/LEI/dot/utils"
)

// Defaults ...
type Defaults struct {
	Template string
	Defaults map[string]map[string]Def // []*Default
	Commands []string
}

// Def ...
type Def struct {
	App    string
	Domain string
	Name   string
	Type   string
	Value  interface{}
	Sudo   bool
}

// Read ...
func (d *Defaults) Read(s string) error {
	if s == "" {
		return nil
	}
	if !utils.Exist(s) {
		return nil
	}
	cfg, err := ioutil.ReadFile(s)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(cfg, &d)
}

// Parse ...
func (d *Defaults) Parse() error {
	for a, b := range d.Defaults {
		// fmt.Printf("Defaults: %s\n (%d)\n", d.Template, len(D.Commands))
		for name, def := range b {
			def.App = a
			def.Name = name
			// s := fmt.Sprintf("%s %s %s %s\n", def.Domain, def.Name, def.Type, def.Value)
			str, err := TemplateData(def.Name, d.Template, def)
			if err != nil {
				return err
			}
			// fmt.Printf("[%s] %s\n", c, str)
			if def.Sudo {
				str = "sudo " + str
			}
			d.Commands = append(d.Commands, str)
		}
	}
	return nil
}

// Exec ...
func (d *Defaults) Exec() error {
	for _, s := range d.Commands {
		fmt.Printf("%s\n", s)
		if err := execute(Shell, "-c", s); err != nil {
			return err
		}
	}
	return nil
}
