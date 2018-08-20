package dot

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	"github.com/LEI/dot/internal/env"
)

var (
	tplFuncMap = template.FuncMap{
		// https://github.com/hashicorp/consul-template/blob/de2ebf4/template_functions.go#L727-L901
		"add": func(a, b int) int {
			return a + b
		},
		"title": strings.Title,
		"lcFirst": func(s string) string {
			for i, v := range s {
				return string(unicode.ToLower(v)) + s[i+1:]
			}
			return ""
		},
		"ucFirst": func(s string) string {
			for i, v := range s {
				return string(unicode.ToUpper(v)) + s[i+1:]
			}
			return ""
		},
		"expand": os.ExpandEnv, // TODO custom env
		"escape": func(s interface{}) interface{} {
			str, ok := s.(string)
			if !ok {
				return s
			}
			// shellEscape
			if !strings.Contains(str, " ") && !strings.Contains(str, "\"") {
				return str
			}
			return strconv.Quote(str)
		},
	}
)

// Tpl task
type Tpl struct {
	Task        `mapstructure:",squash"` // Action, If, OS
	Source      string
	Target      string
	Env         map[string]string
	Vars        map[string]interface{}
	IncludeVars string `mapstructure:"include_vars"`
}

func (t *Tpl) String() string {
	return fmt.Sprintf("%s:%s", t.Source, t.Target)
}

// Type task name
func (t *Tpl) Type() string {
	return "tpl" // template
}

// DoString string
func (t *Tpl) DoString() string {
	return fmt.Sprintf("gotpl %s %s", t.Source, t.Target)
}

// UndoString string
func (t *Tpl) UndoString() string {
	return fmt.Sprintf("rm %s", t.Target)
}

// Status check task
func (t *Tpl) Status() error {
	data, err := t.Data()
	if err != nil {
		return err
	}
	exists, err := tplExists(t.Source, t.Target, data)
	if err != nil {
		return err
	}
	if exists {
		return ErrAlreadyExist
	}
	return nil
}

// Do task
func (t *Tpl) Do() error {
	if err := t.Status(); err != nil {
		switch err {
		case ErrAlreadyExist, ErrSkip:
			return nil
		default:
			return err
		}
	}
	data, err := t.Data()
	if err != nil {
		return err
	}
	content, err := parseTpl(t.Source, data)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(t.Target, []byte(content), defaultFileMode); err != nil {
		return err
	}
	return nil
}

// Undo task
func (t *Tpl) Undo() error {
	if err := t.Status(); err != nil {
		switch err {
		case ErrSkip:
			return nil
		case ErrAlreadyExist:
			// continue
		default:
			return err
		}
	}
	return os.Remove(t.Target)
}

// Data ...
func (t *Tpl) Data() (map[string]interface{}, error) {
	data := make(map[string]interface{}, 0)
	// Global environment variables
	// and custom application baseEnv
	for k, v := range env.GetAll() {
		data[k] = v
	}
	// Specific role environment
	for k, v := range t.Env {
		// k = strings.ToTitle(k)
		ev, err := buildTplEnv(k, v)
		if err != nil {
			return data, err
		}
		// fmt.Println("$ ENV", k, "=", v)
		data[k] = ev
	}
	// Extra variables (not string only)
	for k, v := range t.Vars {
		// fmt.Println("$ VAR", k, "=", v)
		data[k] = v
	}
	return data, nil
}

// templateExists returns true if the template is the same.
func tplExists(src, dst string, data map[string]interface{}) (bool, error) {
	if !exists(src) {
		return false, fmt.Errorf("%s: no such file or directory (to tpl %s)", src, dst)
	}
	if !exists(dst) {
		// Stop here if the target does not exist
		return false, nil
	}
	// TODO compare file contents
	return true, nil
}

func parseTpl(src string, data map[string]interface{}) (string, error) {
	_, name := filepath.Split(src)
	tmpl, err := template.New(name).Option("missingkey=zero").Funcs(tplFuncMap).ParseGlob(src)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	if err != nil {
		return buf.String(), err
	}
	if err = tmpl.Execute(buf, data); err != nil {
		return buf.String(), err
	}
	return buf.String(), nil
}

// buildTpl ...
func buildTpl(k, v string, data interface{}, funcMaps ...template.FuncMap) (string, error) {
	if v == "" {
		return v, nil
	}
	tmpl := template.New(k).Option("missingkey=zero")
	tmpl.Funcs(tplFuncMap)
	for _, funcMap := range funcMaps {
		tmpl.Funcs(funcMap)
	}
	tmpl, err := tmpl.Parse(v)
	if err != nil {
		return v, err
	}
	buf := &bytes.Buffer{}
	err = tmpl.Execute(buf, data)
	if err != nil {
		return v, err
	}
	v = buf.String()
	return v, nil
}

// buildTplEnv ...
func buildTplEnv(k, v string) (string, error) {
	return buildTpl(k, v, env.GetAll())
}
