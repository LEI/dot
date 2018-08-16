package system

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/LEI/dot/internal/comp"
	"gopkg.in/yaml.v2"
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
		"expand": os.ExpandEnv,
		"escape": func(s interface{}) interface{} {
			str, ok := s.(string)
			if !ok {
				return s
			}
			return shellEscape(str)
		},
	}
)

func shellEscape(s string) string {
	if !strings.Contains(s, " ") && !strings.Contains(s, "\"") {
		return s
	}
	return strconv.Quote(s)
}

// CheckTemplate ... (verify/validate)
func CheckTemplate(src, dst string, data map[string]interface{}) error {
	if !Exists(src) {
		// return ErrIsNotExist
		return fmt.Errorf("%s: no such file to template to %s", src, dst)
	}
	if !Exists(dst) {
		// Stop here if the target does not exist
		return nil
	}
	b, err := parseTpl(src, data)
	if err != nil {
		return err
	}
	// Compare to cached content
	ok, err := store.CompareFile(dst)
	if err != nil {
		return err
	}
	if ok {
		fmt.Println("cache matched, remove", dst, "?")
		return nil // ErrFileExist
	}
	// Compare to new content
	ok, err = comp.RegularFile(dst, b)
	if err != nil {
		return err
	}
	if !ok {
		fmt.Println("DIFF OF", dst)
		if err := printDiff(dst, b); err != nil {
			return err
		}
		return ErrFileExist
	}
	return ErrTemplateAlreadyExist
}

// CreateTemplate ...
func CreateTemplate(src, dst string, data map[string]interface{}) (err error) {
	content, err := parseTpl(src, data)
	if err != nil {
		return err
	}
	if DryRun {
		return nil
	}
	if err := ioutil.WriteFile(dst, content, FileMode); err != nil {
		return err
	}
	return store.Put(dst, content)
}

// RemoveTemplate ...
func RemoveTemplate(src, dst string, data map[string]interface{}) error {
	// content, err := parseTpl(src, data)
	// if err != nil {
	// 	return err
	// }
	// if DryRun {
	// 	return nil
	// }
	return Remove(dst)
}

// ParseTemplate ...
func ParseTemplate(file string) (data map[string]interface{}, err error) {
	if strings.HasPrefix(file, "~/") {
		file = filepath.Join(os.Getenv("HOME"), file[2:])
	}
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return data, nil
		}
		return data, err
	}
	if err := yaml.Unmarshal(bytes, &data); err != nil {
		return data, err
	}
	return data, nil
}

// TemplateData ...
func TemplateData(k, v string, data interface{}, funcMaps ...template.FuncMap) (string, error) {
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

func parseTpl(src string, data map[string]interface{}) ([]byte, error) {
	buf := &bytes.Buffer{}
	_, name := filepath.Split(src)
	tmpl, err := template.New(name).Option("missingkey=zero").Funcs(tplFuncMap).ParseGlob(src)
	if err != nil {
		return buf.Bytes(), err
	}
	if err != nil {
		return buf.Bytes(), err
	}
	if err = tmpl.Execute(buf, data); err != nil {
		return buf.Bytes(), err
	}
	return buf.Bytes(), nil
}

func printDiff(path string, content []byte) error {
	// stdout, stderr, status := ExecCommand("")
	diffCmd := exec.Command("diff", path, "-")
	// --side-by-side --suppress-common-lines
	stdin, err := diffCmd.StdinPipe()
	if err != nil {
		return err
	}
	defer stdin.Close()
	diffCmd.Stdout = os.Stdout
	diffCmd.Stderr = os.Stderr
	fmt.Println("START DIFF", path)
	if err := diffCmd.Start(); err != nil {
		return err
	}
	io.WriteString(stdin, string(content))
	// fmt.Println("WAIT")
	stdin.Close()
	diffCmd.Wait()
	fmt.Println("END DIFF", path)
	return nil
}
