package helpers

import (
	"bytes"
	// "fmt"
	"io/ioutil"
	"os"
	"text/template"
)

var ()

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
	if err := ioutil.WriteFile(dst, []byte(str), FileMode); err != nil {
		return false, err
	}
	return true, nil
}
