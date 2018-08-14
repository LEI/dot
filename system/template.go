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
	content, err := parseTpl(src, data)
	if err != nil {
		return err
	}
	if Exists(dst) {
		_, ok, err := CompareFileContent(dst, content)
		if err != nil {
			return err
		}
		if ok {
			return ErrTemplateAlreadyExist
		}
		fmt.Println("DIFF", dst)
		if err := printDiff(dst, content); err != nil {
			return err
		}
	}
	return nil
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
	if err := ioutil.WriteFile(dst, []byte(content), FileMode); err != nil {
		return err
	}
	return nil
}

// RemoveTemplate ...
func RemoveTemplate(src, dst string, data map[string]interface{}) error {
	// content, err := parseTpl(src, data)
	// if err != nil {
	// 	return err
	// }
	if DryRun {
		return nil
	}
	fmt.Println("TODO RemoveTemplate", src, dst)
	// if err := os.Remove(dst); err != nil {
	// 	return err
	// }
	return nil
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

// CompareFileContent ...
func CompareFileContent(s, str string) (string, bool, error) {
	fi, err := os.Open(s)
	if err != nil {
		return "", false, err
	}
	// if err != nil && os.IsExist(err) {
	// 	return "", false, err
	// }
	// if fi == nil {
	// 	return "", false, nil
	// }
	defer fi.Close()
	stat, err := fi.Stat()
	if err != nil && os.IsExist(err) {
		return "", false, err
	}
	if stat != nil && !stat.Mode().IsRegular() {
		return "", false, fmt.Errorf("not a regular file: %s (%q)", stat.Name(), stat.Mode().String())
	}
	// b, err := ioutil.ReadFile(s)
	// if err != nil && os.IsExist(err) {
	// 	return false, err
	// }
	b, err := ioutil.ReadAll(fi)
	content := string(b)
	if err != nil {
		return content, false, err
	}
	// fmt.Println("COMPARED FILE CONTENT", s, len(str), "vs", len(content), "->", content == str)
	return content, content == str, nil
}

// // TemplateData ...
// func TemplateData(k, v string, data interface{}, funcMaps ...template.FuncMap) (string, error) {
// 	if v == "" {
// 		return v, nil
// 	}
// 	tmpl := template.New(k).Option("missingkey=zero")
// 	tmpl.Funcs(tplFuncMap)
// 	for _, funcMap := range funcMaps {
// 		tmpl.Funcs(funcMap)
// 	}
// 	tmpl, err := tmpl.Parse(v)
// 	if err != nil {
// 		return v, err
// 	}
// 	buf := &bytes.Buffer{}
// 	err = tmpl.Execute(buf, data)
// 	if err != nil {
// 		return v, err
// 	}
// 	v = buf.String()
// 	return v, nil
// }

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

func printDiff(s, content string) error {
	// stdout, stderr, status := ExecCommand("")
	diffCmd := exec.Command("diff", s, "-")
	// --side-by-side --suppress-common-lines
	stdin, err := diffCmd.StdinPipe()
	if err != nil {
		return err
	}
	defer stdin.Close()
	diffCmd.Stdout = os.Stdout
	diffCmd.Stderr = os.Stderr
	fmt.Println("START DIFF", s)
	if err := diffCmd.Start(); err != nil {
		return err
	}
	io.WriteString(stdin, content)
	// fmt.Println("WAIT")
	stdin.Close()
	diffCmd.Wait()
	fmt.Println("END DIFF", s)
	return nil
}
