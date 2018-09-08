package host

// https://github.com/shirou/gopsutil/blob/master/host/host.go

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/LEI/dot/internal/env"
	"github.com/LEI/dot/internal/shell"
)

var (
	// OSTypes stores the list of OS types
	OSTypes []string
)

func init() {
	OSTypes = GetOSTypes()
}

// GetOSTypes types (name, release, family, distrib...).
func GetOSTypes() []string {
	types := []string{runtime.GOOS}
	types = append(types, NewRelease().Parse()...)
	types = append(types, parseEnvVar("OSTYPE")...)
	return types
}

// HasOS checks at least one given OS type matches current host.
func HasOS(s ...string) bool {
	// ok, _ := matches(s, OSTypes)
	// return ok
	ok, err := matches(s, OSTypes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", s, err)
		os.Exit(1)
	}
	return ok
}

func matches(in []string, list []string) (bool, error) {
	if len(in) == 0 {
		return false, nil
	}
	nn := 0
	for _, pattern := range in {
		negated := pattern[0] == '!'
		if negated {
			pattern = pattern[1:]
		}
		for _, str := range list {
			matched, err := regexp.MatchString(pattern, str)
			if err != nil {
				// pattern error
				return false, err
			}
			if negated && matched {
				return false, nil
			}
			if matched {
				return true, nil
			}
		}
		if negated {
			nn++
			// return true, nil
		}
	}
	if nn == len(in) { // && nn > 0
		return true, nil
	}
	return false, nil
}

func parseEnvVar(name string) []string {
	types := make([]string, 0)
	if o, ok := env.Lookup(name); ok && o != "" {
		types = append(types, o)
	} else { // !ok || s == ""
		// Fallback to shell invocation
		// fmt.Printf("%s='%s' (%v)\n", name, s, ok)
		out, err := exec.Command(shell.Get(), "-c", "printf '%s' \"$"+name+"\"").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s error: %s\n", name, err)
		}
		if len(out) > 0 {
			s := string(out)
			o := strings.Split(s, ".")
			if len(o) > 0 && o[0] != s {
				types = append(types, o[0])
			}
			types = append(types, s)
		}
	}
	return types
}
