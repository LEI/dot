package dotfile

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"text/template"
)

var (
	// Shell ...
	Shell = "bash"

	// OS ...
	OS = runtime.GOOS

	osTypes []string

	originalEnv map[string]string

	extraEnv = map[string]string{
		"OS": OS,
	}
)

func init() {
	osTypes = GetOSTypes()
	originalEnv = GetEnv()
}

// InitEnv ...
func InitEnv() error {
	if err := os.Setenv("OS", OS); err != nil {
		return err
	}
	return nil
}

// GetEnv ...
func GetEnv() map[string]string {
	env := make(map[string]string, 0)
	for _, i := range os.Environ() {
		sep := strings.Index(i, "=")
		k := i[0:sep]
		v := i[sep+1:]
		env[k] = v
	}
	// for k, v := range dotEnv {
	// 	prefix := ""
	// 	if _, ok := env[k]; !ok {
	// 		env[k] = v
	// 	} else {
	// 		prefix = "# (SKIPPED) "
	// 	}
	// 	cfgLogger.Debugf("%s%s=%+v", prefix, k, v)
	// }
	return env
}

// TemplateEnv ...
func TemplateEnv(k, v string) (string, error) {
	if v == "" {
		return v, nil
	}
	templ, err := template.New(k).Option("missingkey=zero").Parse(v)
	if err != nil {
		return v, err
	}
	buf := &bytes.Buffer{}
	err = templ.Execute(buf, GetEnv())
	if err != nil {
		return v, err
	}
	v = buf.String()
	return v, nil
}

// SetEnv ...
func SetEnv(k, v string) error {
	v, err := TemplateEnv(k, v)
	if err != nil {
		return err
	}
	fmt.Printf("%s=%s\n", k, v)
	return os.Setenv(k, v)
}

// ExpandEnv ...
func ExpandEnv(s string, envs ...map[string]string) string {
	// TODO for _, e := range envs {
	// 	s = os.Expand(s, e)
	// }
	s = os.ExpandEnv(s)
	return s
}

// HasOSType ...
func HasOSType(s ...string) bool {
	return HasOne(s, osTypes)
}

// HasOne ...
func HasOne(in []string, list []string) bool {
	for _, a := range in {
		for _, b := range list {
			if b == a {
				return true
			}
		}
	}
	return false
}

// GetOSTypes (name, family, distrib...)
func GetOSTypes() []string {
	types := []string{OS}

	// Add OS family
	c := exec.Command(Shell, "-c", "cat /etc/*-release")
	stdout, _ := c.StdoutPipe()
	// stderr, _ := c.StderrPipe()
	c.Start()
	scanner := bufio.NewScanner(stdout)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
		v := strings.TrimLeft(m, "ID=")
		if m != v {
			types = append(types, v)
			break
		}
	}
	c.Wait()

	OSTYPE, ok := os.LookupEnv("OSTYPE")
	if ok && OSTYPE != "" {
		types = append(types, OSTYPE)
	} else { // !ok || OSTYPE == ""
		// fmt.Printf("OSTYPE='%s' (%v)\n", OSTYPE, ok)
		out, err := exec.Command(Shell, "-c", "printf '%s' \"$OSTYPE\"").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
		}
		if len(out) > 0 {
			OSTYPE = string(out)
			o := strings.Split(OSTYPE, ".")
			if len(o) > 0 {
				types = append(types, o[0])
			}
			types = append(types, OSTYPE)
		}
	}
	if OSTYPE == "" {
		fmt.Println("OSTYPE is not set or empty")
	}
	return types
}
