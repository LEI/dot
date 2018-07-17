package dotfile

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	// "path"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

var (
	tplFuncMap = template.FuncMap{
		"lcFirst": func (s string) string {
			for i, v := range s {
				return string(unicode.ToLower(v)) + s[i+1:]
			}
			return ""
		},
		"ucFirst": func (s string) string {
			for i, v := range s {
				return string(unicode.ToUpper(v)) + s[i+1:]
			}
			return ""
		},
	}
)

// TemplateTask struct
type TemplateTask struct {
	// parsers.Tpl
	Source, Target string
	Ext            string // `default:"tpl"`
	Env            map[string]string
	Vars           map[string]interface{}
	Task
}

// Do ...
func (t *TemplateTask) Do(a string) error {
	return do(t, a)
}

// Parse template file
func (t *TemplateTask) Parse() (string, error) {
	_, name := filepath.Split(t.Source)
	tmpl, err := template.New(name).Option("missingkey=zero").Funcs(tplFuncMap).ParseGlob(t.Source)
	// b, err := ioutil.ReadFile(t.Source)
	// c := string(b) // Template contents
	// if err != nil && os.IsExist(err) {
	// 	return c, err
	// }
	// tmpl, err := template.New(name).Option("missingkey=zero").Funcs(tplFuncMap).Parse(c)
	if err != nil {
		return "", err
	}
	buf := &bytes.Buffer{}
	data := make(map[string]interface{}, 0)
	// Global environment variables
	// for k, v := range GetEnv() {
	// 	data[k] = v
	// }
	// Custom application environment
	for k, v := range baseEnv {
		k = strings.ToTitle(k)
		v, err := TemplateEnv(k, v)
		if err != nil {
			return "", err
		}
		data[k] = v
	}
	// Specific role environment
	for k, v := range t.Env {
		k = strings.ToTitle(k)
		v, err := TemplateEnv(k, v)
		if err != nil {
			return "", err
		}
		data[k] = v
	}
	// Extra variables (not string only)
	for k, v := range t.Vars {
		data[k] = v
	}
	if err = tmpl.Execute(buf, data); err != nil {
		return buf.String(), err
	}
	return buf.String(), nil
}

// Install template
func (t *TemplateTask) Install() error {
	if err := createBaseDir(t.Target); err != nil && err != ErrDirShouldExist {
		return err
	}
	changed, err := Template(t) // t.Source, dst, t.Env
	if err != nil {
		return err
	}
	prefix := ""
	if !changed {
		prefix = "# "
	}
	/*
		vars := []string{}
		for k, v := range t.Env { // + dotEnv
			// fmt.Printf("%s=\"%s\"\n", k, v)
			vars = append(vars, fmt.Sprintf("%s: %s", k, v))
		}
	*/
	// envsubst
	// fmt.Printf("%senvsubst < %s | tee %s\n", prefix, t.Source, dst)
	// fmt.Printf("%sgotpl %s <<< '%s' | tee %s\n", prefix, t.Source, strings.Join(vars, "\n"), t.Target)

	// TODO? github.com/tsg/gotpl with option missingkey=zero
	// fmt.Printf("%sgotpl %s <<'EOF' | tee %s\n%s\nEOF\n", prefix, t.Source, t.Target, strings.Join(vars, "\n"))
	fmt.Printf("%stpl %s -> %s\n", prefix, t.Source, t.Target)
	return nil
}

// Remove template
func (t *TemplateTask) Remove() error {
	changed, err := Untemplate(t)
	if err != nil {
		return err
	}
	prefix := ""
	if !changed {
		prefix = "# "
	}
	/*for k, v := range t.Env { // + dotEnv
		fmt.Printf("%s=\"%s\"\n", k, v)
	}*/
	fmt.Printf("%srm %s\n", prefix, t.Target)
	if RemoveEmptyDirs {
		if err := removeBaseDir(t.Target); err != nil {
			return err
		}
	}
	return nil
}

// Template task
func Template(t *TemplateTask) (bool, error) {
	str, err := t.Parse()
	if err != nil {
		return false, err
	}
	b, err := ioutil.ReadFile(t.Target)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	c := string(b) // Current contents
	if str == c {
		return false, nil
	} else if str != c && c != "" {
		// TODO: cache checksum of previous run to compare
		// or ask for user confirmation to remove the file
		diff := t.Source // TODO diff
		return false, fmt.Errorf("# /!\\ Template content mismatch: %s\n%s", t.Target, diff)
	}
	if Verbose > 1 {
		fmt.Printf("---START---\n%s\n----END----\n", str)
	}
	if DryRun {
		return true, nil
	}
	if err := ioutil.WriteFile(t.Target, []byte(str), FileMode); err != nil {
		return false, err
	}
	return true, nil
}

// Untemplate task
func Untemplate(t *TemplateTask) (bool, error) {
	str, err := t.Parse()
	if err != nil {
		return false, err
	}
	b, err := ioutil.ReadFile(t.Target)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if len(b) == 0 { // Empty file
		return false, nil
	}
	c := string(b) // Current contents
	if str != c && c != "" {
		return false, fmt.Errorf("# /!\\ Template content mismatch: %s", t.Target)
	}
	if DryRun {
		return true, nil
	}
	if err := os.Remove(t.Target); err != nil {
		return false, err
	}
	return true, nil
}
