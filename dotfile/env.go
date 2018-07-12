package dotfile

import (
	// "bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/jessevdk/go-flags"
)

// Release ...
type Release struct {
	ID string `ini-name:"ID"` // debian
	Name string `ini-name:"NAME"` // Debian GNU/Linux
	PrettyName string `ini-name:"PRETTY_NAME"` // Debian GNU/Linux 9 (stretch)
	Version string `ini-name:"VERSION"` // 9 (stretch)
	VersionID string `ini-name:"VERSION_ID"` // 9
	// HomeURL string `ini-name:"HOME_URL"`
	// SupportURL string `ini-name:"SUPPORT_URL"`
	// BugReportURL string `ini-name:"BUG_REPORT_URL"`
}

var (
	// OS ...
	OS = runtime.GOOS

	release Release

	osTypes []string

	originalEnv map[string]string

	extraEnv = map[string]string{
		"OS": OS,
	}
)

func init() {
	osTypes = GetOSTypes()
	// fmt.Printf("OS types:\n%+v\n", strings.Join(osTypes[:], "\n"))
	originalEnv = GetEnv()
	// fmt.Printf("Original env: %+v\n", originalEnv)
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

// GetOSTypes OS name, release, family, distrib...
func GetOSTypes() []string {
	types := []string{OS}
	r := parseReleases()
	if r.Name != "" {
		types = append(types, r.Name)
		if r.ID != "" {
			types = append(types, r.Name + r.ID)
		}
	}
	types = append(types, parseOSTypes()...)
	return types
}

// Read release files as INI
func parseReleases() Release {
	paths, err := filepath.Glob("/etc/*-release")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	for _, p := range paths {
		parser := flags.NewParser(&release, flags.IgnoreUnknown)
		ini := flags.NewIniParser(parser)
		// ini.ParseAsDefaults = true
		err := ini.ParseFile(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(1)
		}
		// if Verbose {
		fmt.Printf("%s:\n%+v\n", p, release)
		// }
	}
	return release
}

func parseOSTypes() []string {
	types := make([]string, 0)
	if o, ok := os.LookupEnv("OSTYPE"); ok && o != "" {
		types = append(types, o)
	} else { // !ok || s == ""
		// fmt.Printf("OSTYPE='%s' (%v)\n", s, ok)
		out, err := exec.Command("bash", "-c", "printf '%s' \"$OSTYPE\"").Output()
		if err != nil {
			fmt.Fprintf(os.Stderr, "# OSTYPE error: %s\n", err)
		}
		if len(out) > 0 {
			s := string(out)
			o := strings.Split(s, ".")
			if len(o) > 0 {
				types = append(types, o[0])
			}
			types = append(types, s)
		}
	}
	return types
}
