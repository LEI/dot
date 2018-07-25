package dotfile

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"

	"github.com/LEI/dot/utils"
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
func (t *TemplateTask) Do(a string) (string, error) {
	return do(t, a)
}

// Install template
func (t *TemplateTask) Install() (string, error) {
	// if utils.Exist(dst) {
	// 	return "", nil
	// }
	if err := createBaseDir(t.Target); err != nil && err != ErrDirShouldExist {
		return "", err
	}
	data, err := t.Data()
	if err != nil {
		return "", err
	}
	changed, err := Template(t.Source, t.Target, data)
	if err != nil {
		return "", err
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
	return fmt.Sprintf("%stpl %s -> %s", prefix, t.Source, t.Target), nil
}

// Remove template
func (t *TemplateTask) Remove() (string, error) {
	data, err := t.Data()
	if err != nil {
		return "", err
	}
	changed, err := Untemplate(t.Source, t.Target, data)
	if err != nil {
		return "", err
	}
	prefix := ""
	if !changed {
		prefix = "# "
	}
	/*for k, v := range t.Env { // + dotEnv
		fmt.Printf("%s=\"%s\"\n", k, v)
	}*/
	if RemoveEmptyDirs {
		if err := removeBaseDir(t.Target); err != nil {
			return "", err
		}
	}
	return fmt.Sprintf("%srm %s", prefix, t.Target), nil
}

// Data ...
func (t *TemplateTask) Data() (map[string]interface{}, error) {
	data := make(map[string]interface{}, 0)
	// Global environment variables
	// and custom application baseEnv
	for k, v := range GetEnv() {
		data[k] = v
	}
	// Specific role environment
	for k, v := range t.Env {
		k = strings.ToTitle(k)
		v, err := TemplateEnv(k, v)
		if err != nil {
			return data, err
		}
		data[k] = v
	}
	// Extra variables (not string only)
	for k, v := range t.Vars {
		data[k] = v
	}
	return data, nil
}

// Template task
func Template(src, dst string, data map[string]interface{}) (bool, error) {
	content, err := parseTpl(src, data)
	if err != nil {
		return false, err
	}
	if utils.Exist(dst) {
		c, ok, err := utils.CompareFileContent(dst, content)
		if err != nil {
			return false, err
		}
		if ok {
			return false, nil
		}
		valid, err := dotCache.Validate(dst, c)
		if err != nil {
			return false, err
		}
		if !valid {
			overwrite, err := tplOverwrite(src, dst, content)
			// if err != nil || !overwrite {
			// 	return overwrite, err
			// }
			if err != nil {
				return false, err
			}
			if !overwrite {
				return false, nil
			}
		}
		// if !overwrite {
		// 	// if err := dotCache.Put(dst, content); err != nil {
		// 	// 	return false, err
		// 	// }
		// 	return false, nil
		// }
		// if !ok {
		// 	// overwrite, err := tplOverwrite(src, dst, content)
		// 	// if err != nil || !overwrite {
		// 	// 	return overwrite, err
		// 	// }
		// 	return false, fmt.Errorf("different template target: %s", dst)
		// }
	}
	if DryRun {
		return true, nil
	}
	// fmt.Println("------------------- xxx", content, "xxx")
	// fmt.Println("------------------- yyy", c, "yyy")
	// fmt.Println("WRITE", content, "ENDWRITE")
	if err := ioutil.WriteFile(dst, []byte(content), FileMode); err != nil {
		return false, err
	}
	if err := dotCache.Put(dst, content); err != nil {
		return true, err
	}
	return true, nil
}

// Untemplate task
func Untemplate(src, dst string, data map[string]interface{}) (bool, error) {
	if !utils.Exist(dst) {
		if err := dotCache.Del(dst); err != nil {
			return false, err
		}
		return false, nil
	}
	content, err := parseTpl(src, data)
	if err != nil {
		return false, err
	}
	c, ok, err := utils.CompareFileContent(dst, content)
	if err != nil {
		return false, err
	}
	if !ok {
		valid, err := dotCache.Validate(dst, c)
		if err != nil {
			return false, err
		}
		if !valid {
			overwrite, err := tplOverwrite(src, dst, content)
			// if err != nil || !overwrite {
			// 	return overwrite, err
			// }
			if err != nil {
				return false, err
			}
			if !overwrite {
				return false, nil
			}
		}
	}
	if DryRun {
		return true, nil
	}
	if err := os.Remove(dst); err != nil {
		return false, err
	}
	if err := dotCache.Del(dst); err != nil {
		return false, err
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

func askOverwriteTpl(src, dst, content string) (bool, error) {
	c, ok, err := utils.CompareFileContent(dst, content)
	if err != nil {
		return false, err
	}
	if ok {
		// FIXME untemplate?
		fmt.Println("COMPARED OK", dst)
		return true, nil // false, nil
	}
	// b, err := ioutil.ReadFile(dst)
	// if err != nil && os.IsExist(err) {
	// 	return content, false, err
	// }
	// c := string(b) // Current file content
	// if content == c { // Same file content
	// 	return false, fmt.Errorf("same file content for %s, should be handled by CheckFile", dst)
	// } else if content != c && c != "" { // Target changed
	valid, err := dotCache.Validate(dst, c)
	if err != nil {
		return false, err
	}
	if !valid {
		overwrite, err := tplOverwrite(src, dst, content)
		// if err != nil || !overwrite {
		// 	return overwrite, err
		// }
		if err != nil {
			return false, err
		}
		if !overwrite {
			return false, nil
		}
	} // else if content != c && c == "" && OverwriteEmptyFiles {}
	// if Verbose > 1 {
	// 	fmt.Printf("---START---\n%s\n----END----\n", content)
	// }
	return true, nil
}

func tplOverwrite(src, dst, content string) (bool, error) {
	if err := printDiff(dst, content); err != nil {
		return false, err
	}
	q := fmt.Sprintf("Overwrite existing template target: %s", dst)
	if !AskConfirmation(q) {
		// diff := src // TODO diff?
		// return content, false, fmt.Errorf("# /!\\ Template content mismatch: %s\n%s", dst, diff)
		fmt.Fprintf(os.Stderr, "Skipping template %s because its target exists: %s", src, dst)
		return false, nil
	}
	// if DryRun {
	// 	return true, nil
	// }
	if err := Backup(dst); err != nil {
		return false, err
	}
	return true, nil
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
