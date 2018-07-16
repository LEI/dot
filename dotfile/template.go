package dotfile

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	// "path"
	// "strings"
	"text/template"
)

// TemplateTask struct
type TemplateTask struct {
	Source, Target string
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
	tmpl, err := template.ParseGlob(t.Source)
	if err != nil {
		return "", err
	}
	tmpl = tmpl.Option("missingkey=zero")
	buf := &bytes.Buffer{}
	// env, err := GetEnv()
	// if err != nil {
	// 	return false, err
	// }
	// for k, v := range t.Vars {
	// 	fmt.Println("VAR", k, "=", v)
	// }
	for k, v := range t.Env {
		t.Vars[k] = v
	}
	if err = tmpl.Execute(buf, t.Vars); err != nil {
		return buf.String(), err
	}
	return buf.String(), nil
}

// Install template
func (t *TemplateTask) Install() error {
	if err := createBaseDir(t.Target); err != nil {
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
	// if err := removeBaseDir(t.Target); err != nil {
	// 	return err
	// }
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
	if str == string(b) {
		return false, nil
	}
	if DryRun {
		// if Verbose > 0 {
		// 	fmt.Println(str)
		// }
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
	if str != string(b) { // Mismatching content
		fmt.Printf("Warn: mismatching content %s\n", t.Target)
		return false, nil
	}
	if DryRun {
		return true, nil
	}
	if err := os.Remove(t.Target); err != nil {
		return false, err
	}
	return true, nil
}
