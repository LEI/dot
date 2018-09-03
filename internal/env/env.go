package env

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/LEI/dot/internal/shell"
)

var (
	original map[string]string

	// Fallback to default value if the environment variable is not defined
	// Example: ${VARIABLE:-fallback}
	// sustituteParameterRe = regexp.MustCompile(`^\${([a-zA-Z0-9_]+):-(.*)}$`)
	sustituteParameterRe = regexp.MustCompile(`^([a-zA-Z0-9_]+):-(.*)$`)
	// Command subsitutions with backticks are not supported
	// Example: $(command substitution)
	substituteCommandRe = regexp.MustCompile(`^\$\((.*)\)$`)
	// Used to handle quoted values
	quotedRe       = regexp.MustCompile(`^"(.*)"$`)
	singleQuotedRe = regexp.MustCompile(`^'(.*)'$`)
)

func init() {
	original = GetAll()
}

// // Init ...
// func Init() error {
// 	for k, v := range baseEnv {
// 		if err := os.Setenv(k, v); err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

// Set environment variable
func Set(k, v string) error {
	// if Verbose > 0 {
	// 	fmt.Printf("%s=%s\n", k, v)
	// }
	return os.Setenv(k, v)
}

// Unset environment variable
func Unset(k, v string) error {
	return os.Unsetenv(k)
}

// GetAll environment variables
func GetAll() map[string]string {
	env := make(map[string]string, 0)
	for _, i := range os.Environ() {
		k, v := Split(i)
		if _, ok := env[k]; !ok {
			env[k] = v
		}
	}
	return env
}

// Split "key=value" into two variables.
func Split(s string) (string, string) {
	sep := strings.Index(s, "=")
	k := strings.ToUpper(s[0:sep])
	v := s[sep+1:]
	if matches := quotedRe.FindStringSubmatch(k); len(matches) == 2 {
		v = matches[1]
	} else if matches := singleQuotedRe.FindStringSubmatch(k); len(matches) == 2 {
		v = matches[1]
	}
	return k, v
}

// Get environment variable
func Get(k string) string {
	// v, _ := GetEnv(k, GetAll())
	env := GetAll()
	v := env[k]
	return v
}

/* // ParseEnv variable
func ParseEnv(s string, env map[string]string) (k string, v string, ok bool) {
	// fmt.Printf(">>> ParseEnv(%#v)\n", s)
	// var defaultVal string
	// v, ok := env[k]
	// if !ok {
	// 	return v
	// }
	// v, err := TemplateEnv(k, v)
	// if err != nil {
	// 	return err
	// }
	// if matches := substituteCommandRe.FindStringSubmatch(key); len(matches) == 2 {
	// 	fmt.Println("111111", key, matches)
	// 	c := matches[1]
	// 	cmd := exec.Command(shell.Get(), "-c", c)
	// 	out, err := cmd.Output()
	// 	if err != nil {
	// 		fmt.Fprintf(os.Stderr, "failed to execute `%s`: %s\n", c, err)
	// 	}
	// 	fmt.Printf(">>> SUBST %s=%q\n", key, string(out))
	// 	if err == nil {
	// 		key = ""
	// 		val = strings.TrimRight(string(out), "\n")
	// 		// return v, true
	// 	}
	// }
	// if matches := sustituteParameterRe.FindStringSubmatch(key); len(matches) == 3 {
	// 	fmt.Println("22222222", key, matches)
	// 	key = matches[1]
	// 	fmt.Println(">fb expand", matches[2])
	// 	val = ExpandEnv(matches[2], env)
	// 	fmt.Printf(">>> FALLBACK %s=%q\n", key, val)
	// }
	if matches := sustituteParameterRe.FindStringSubmatch(s); len(matches) == 3 {
		k = matches[1]
		fmt.Println("OKKK", k)
		// if val, ok := env[k]; ok {
		// 	return k, val
		// }
		// v = matches[2]
		v = ExpandEnv(matches[2], env) // GetAll())
		// return k, v, false
		s = v // continue in case the fallback is a command
		fmt.Println("111", k, v)
	}
	if matches := substituteCommandRe.FindStringSubmatch(s); len(matches) == 2 {
		c := matches[1]
		cmd := exec.Command(shell.Get(), "-c", c)
		out, err := cmd.Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to execute `%s`: %s\n", c, err)
		}
		if err == nil {
			v = strings.TrimRight(string(out), "\n")
			// fmt.Printf("SUBST %s=%q\n", k, v)
			return k, v, true
		}
		fmt.Println("222", k, v)
	}
	// v, ok := env[k]
	// if !ok { // v == "" && val != ""
	// 	if defaultVal != "" {
	// 		v = defaultVal
	// 		ok = true
	// 	}
	// }
	// fmt.Printf("EXPANDED (%v) %s=%q\n", ok, k, v)
	return k, v, false
} */

// // SubstituteEnv variable returns key and default value
// func SubstituteEnv(s string, env map[string]string) (s string) {
// }

// Expand ...
func Expand(s string) string {
	// return os.ExpandEnv(s)
	// env := GetAll()
	// key, val, ok := ParseEnv(s, env)
	// if ; ok {
	// 	return val
	// }
	return os.Expand(s, Get)
}

// ExpandEnv variables and execute commands, or fallback to global env
func ExpandEnv(s string, env map[string]string) string {
	if matches := substituteCommandRe.FindStringSubmatch(s); len(matches) == 2 {
		c := matches[1]
		cmd := exec.Command(shell.Get(), "-c", c)
		// cmd.Env = os.Environ()
		out, err := cmd.Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to execute `%s`: %s\n", c, err)
		}
		if err == nil {
			v := strings.TrimRight(string(out), "\n")
			// fmt.Printf("SUBST %s=%q\n", k, v)
			return v
		}
	}
	expand := func(k string) string {
		var defaultValue string
		if matches := sustituteParameterRe.FindStringSubmatch(k); len(matches) == 3 {
			k = matches[1]
			// if val, ok := env[k]; ok {
			// 	return k, val
			// }
			// v = matches[2]
			defaultValue = ExpandEnv(matches[2], env) // GetAll())
		}
		if v, ok := env[k]; ok {
			return v
		}
		// fmt.Printf("GETENV %q\n", k)
		// if v, ok := GetEnv(k, env); ok {
		// 	return v
		// }
		v := Get(k)
		if v == "" && defaultValue != "" {
			v = defaultValue
		}
		return v
	}
	// key, val, ok := ParseEnv(s, env)
	// if ok { // command substituted
	// 	fmt.Println("SUBST", val)
	// 	return val
	// }
	// if key != "" {
	// 	fmt.Println("REPLACE S", s, key)
	// 	s = "$" + key
	// }
	// v := os.Expand(s, expand)
	// if v == "" { // key != ""
	// 	fmt.Println("FALLB", val)
	// 	return val // default value
	// }
	// fmt.Println("EXPANDED", v)
	// return v
	return os.Expand(s, expand)
}

// Lookup environment variable
func Lookup(k string) (string, bool) {
	return os.LookupEnv(k)
}

// Clear environment
func Clear() {
	os.Clearenv()
}

// Restore environment
func Restore(env map[string]string) error {
	Clear()
	for k, v := range env {
		if err := os.Setenv(k, v); err != nil {
			return err
		}
	}
	return nil
}

// RestoreOriginal environment
func RestoreOriginal(env map[string]string) error {
	return Restore(original)
}
