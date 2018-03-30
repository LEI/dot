package dotlib

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"
)

// TemplateTask struct
type TemplateTask struct {
	Source, Target string
	Env            map[string]string
	// Task
}

// Install template
func (t *TemplateTask) Install() error {
	_, f := path.Split(t.Source)
	dst := path.Join(t.Target, strings.TrimSuffix(f, ".tpl"))
	changed, err := Template(t.Source, dst, t.Env)
	if err != nil {
		return err
	}
	prefix := "# "
	if changed {
		prefix = ""
	}
	for k, v := range t.Env {
		fmt.Printf("%s=\"%s\"\n", k, v)
	}
	// fmt.Printf("%senvsubst < %s | tee %s\n", prefix, t.Source, dst)
	fmt.Printf("%stemplate %s -> %s\n", prefix, t.Source, dst)
	return nil
}

// Template task
func Template(src, dst string, env map[string]string) (bool, error) {
	tmpl, err := template.ParseGlob(src)
	if err != nil {
		return false, err
	}
	tmpl = tmpl.Option("missingkey=zero")
	buf := &bytes.Buffer{}
	// env, err := GetEnv()
	// if err != nil {
	// 	return false, err
	// }
	if err = tmpl.Execute(buf, env); err != nil {
		return false, err
	}
	str := buf.String()
	b, err := ioutil.ReadFile(dst)
	if err != nil && os.IsExist(err) {
		return false, err
	}
	if str == string(b) {
		return false, nil
	}
	if DryRun {
		return true, nil
	}
	if err := ioutil.WriteFile(dst, []byte(str), FileMode); err != nil {
		return false, err
	}
	return true, nil
}
