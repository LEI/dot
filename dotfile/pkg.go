package dotfile

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var (
	// OS ...
	OS = runtime.GOOS

	// Shell ...
	Shell = "bash"

	osTypes []string
)

func init() {
	osTypes = getOSTypes()
}

// InitEnv ...
func InitEnv() error {
	if err := os.Setenv("OS", OS); err != nil {
		return err
	}
	return nil
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

// List of OS name and family/type
func getOSTypes() []string {
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
