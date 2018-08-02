package dotfile

import (
	"fmt"
	"io/ioutil"
	"os"
	// "os/exec"
	// "path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"

	"github.com/LEI/dot/utils"
)

var (
	extraFuncMap = template.FuncMap{
		// macos `defaults write` value formatter
		"toString": func(s interface{}) interface{} {
			str := ""
			switch t := s.(type) {
			// case bool, int, int32, int64, float32, float64:
			// 	str = fmt.Sprintf("%s", t)
			case bool:
				str = fmt.Sprintf("%v", t)
			case int, int32:
				str = fmt.Sprintf("%d", t)
			case float32, float64:
				str = fmt.Sprintf("%f", t)
			case string:
				str = shellEscape(t)
			case []string: // -array, -array-add
				for _, s := range t {
					if len(str) > 0 {
						str += " "
					}
					str += shellEscape(s)
				}
			case []interface{}:
				for _, s := range t {
					if len(str) > 0 {
						str += " "
					}
					str += shellEscape(s.(string))
				}
			// case map[string]interface{}:
			case map[interface{}]interface{}: // -dict, -dict-add
				for key, val := range t {
					if len(str) > 0 {
						str += " "
					}
					str += key.(string)
					switch v := val.(type) {
					case bool:
						str += fmt.Sprintf(" -bool %v", v)
					case int, int32, int64:
						str += fmt.Sprintf(" -int %d", v)
					case float32, float64:
						str += fmt.Sprintf(" -float %f", v)
					case string:
						str += fmt.Sprintf(" -string %s", shellEscape(v))
					default:
						fmt.Printf("unexpected default %s: %s\n", key, val)
					}
				}
			// -data, -date...
			default:
				str = fmt.Sprintf("%+v", s)
			}
			return str
		},
	}
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
			str, err := TemplateData(def.Name, d.Template, def, extraFuncMap)
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
		fmt.Printf("%s\n", strings.TrimRight(s, "\n"))
		// if err := execute(Shell, "-c", s); err != nil {
		// 	return err
		// }
		stdout, stderr, status := ExecCommand(Shell, "-c", s)
		if status != 0 {
			if stderr == "" {
				stderr = fmt.Sprintf(
					"defaults failed for `%s`: %s (code %d)",
					s,
					stderr,
					status)
			}
			return fmt.Errorf(stderr)
		}
		if stderr != "" {
			fmt.Fprintf(os.Stderr, stderr)
		}
		str := strings.TrimRight(stdout, "\n")
		if str != "" {
			fmt.Println(str)
		}
	}
	return nil
}
