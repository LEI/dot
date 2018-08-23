package dot

// https://github.com/LEI/dot/blob/go-flags/dotfile/template.go

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
	yaml "gopkg.in/yaml.v2"
)

var (
	defaultTemplateExt = "tpl"

	// Number of context lines in difflib output
	diffContextLines = 3

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
	Ext         string // Template extenstion (default: tpl)
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

// Prepare template task
func (t *Tpl) Prepare() error {
	if t.Ext == "" {
		t.Ext = defaultTemplateExt
	}
	if t.Target != "" && t.Ext != "" && strings.HasSuffix(t.Target, "."+t.Ext) {
		t.Target = strings.TrimSuffix(t.Target, "."+t.Ext)
	}
	// Already done in role ParseTpls
	// if t.Vars == nil {
	// 	t.Vars = map[string]interface{}{}
	// }
	// for k, v := range t.Env {
	// 	// ...
	// }
	if t.IncludeVars != "" {
		// Included variables override existing tpl.Vars keys
		inclVars, err := includeVars(t.IncludeVars) // os.ExpandEnv?
		if err != nil {
			return err
		}
		for k, v := range inclVars {
			// if val, ok := t.Vars[k]; !ok {
			// 	return fmt.Errorf("include vars %s: %s=%v already set to %v", t.IncludeVars, k, v, val)
			// }
			t.Vars[k] = v
		}
	}
	return nil
}

// DoString string
func (t *Tpl) DoString() string {
	return fmt.Sprintf("gotpl %s %s", tildify(t.Source), tildify(t.Target))
}

// UndoString string
func (t *Tpl) UndoString() string {
	return fmt.Sprintf("rm %s", tildify(t.Target))
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
	e := env.GetAll() // map[string]string{}
	// for k, v := range env.GetAll() {
	// 	data[k] = v
	// }
	// Specific role environment variables (uppercase key)
	for k, v := range t.Env {
		k = strings.ToUpper(k)
		ev, err := buildTplEnv(k, v, e)
		if err != nil {
			return data, err
		}
		fmt.Printf("$ export %s=%q\n", k, ev)
		data[k] = ev // e[k] = ev
	}
	// Extra variables (not only strings)
	for k, v := range t.Vars {
		// if k == "Env" ...
		if val, ok := v.(string); ok && val != "" {
			ev, err := buildTplEnv(k, val, e)
			if err != nil {
				return data, err
			}
			v = ev
		}
		// fmt.Printf("# var %s = %+v\n", k, v)
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
	content, err := parseTpl(src, data)
	if err != nil {
		return false, err
	}
	// TODO compare file content and ask confirmation
	// printDiff(dst, content)
	diff, err := getDiff(src, dst, content)
	if err == nil {
		fmt.Println(strings.TrimSuffix(diff, "\n"))
	}
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
func buildTplEnv(k, v string, envs ...map[string]string) (string, error) {
	environ := env.GetAll()
	for _, e := range envs {
		for k, v := range e {
			environ[k] = v
		}
	}
	return buildTpl(k, v, environ)
}

func includeVars(file string) (vars map[string]interface{}, err error) {
	if strings.HasPrefix(file, "~/") {
		file = filepath.Join(homeDir, file[2:])
	}
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return vars, nil
		}
		return vars, err
	}
	if err := yaml.Unmarshal(bytes, &vars); err != nil {
		return vars, err
	}
	return vars, nil
}
