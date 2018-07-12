package dotfile

import (
	// "bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"text/template"

	"github.com/jessevdk/go-flags"
)

// Release ...
type Release struct {
	ID         string `ini-name:"ID"`
	IDLike     string `ini-name:"ID_LIKE"`
	Name       string `ini-name:"NAME"`
	PrettyName string `ini-name:"PRETTY_NAME"`
	Version    string `ini-name:"VERSION"`
	VersionID  string `ini-name:"VERSION_ID"`
	// HomeURL string `ini-name:"HOME_URL"`
	// SupportURL string `ini-name:"SUPPORT_URL"`
	// BugReportURL string `ini-name:"BUG_REPORT_URL"`
	DistribID string `ini-name:"DISTRIB_ID"`
	DistribRelease string `ini-name:"DISTRIB_RELEASE"`
	DistribCodename string `ini-name:"DISTRIB_CODENAME"`
	DistribDescription string `ini-name:"DISTRIB_DESCRIPTION"`
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
	fmt.Printf("OS types:\n%+v\n", strings.Join(osTypes[:], "\n"))
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
	// if isNum(r.ID) {
	// 	if r.Name != "" {
	// 		types = append(types, r.Name)
	// 		if r.ID != "" {
	// 			types = append(types, r.Name+r.ID)
	// 		}
	// 	}
	// }
	if r.ID != "" && !isNum(r.ID) {
		types = append(types, r.Name)
	} else if r.ID != "" {
		types = append(types, r.ID)
	} else if r.Name != "" {
		types = append(types, strings.ToLower(r.Name))
	}
	if r.IDLike != "" {
		types = append(types, r.IDLike)
	}
	if r.DistribCodename != "" {
		types = append(types, r.DistribCodename)
	}
	types = append(types, parseOSTypes()...)
	return types
}

func isNum(v string) bool {
	_, err := strconv.Atoi(v)
	return err == nil
}

// /etc/os-release

// PRETTY_NAME="Debian GNU/Linux 9 (stretch)"
// NAME="Debian GNU/Linux"
// VERSION_ID="9"
// VERSION="9 (stretch)"
// ID=debian

// NAME="Ubuntu"
// VERSION="14.04.5 LTS, Trusty Tahr"
// ID=ubuntu
// ID_LIKE=debian
// PRETTY_NAME="Ubuntu 14.04.5 LTS"
// VERSION_ID="14.04"

// /etc/lsb-release

// DISTRIB_ID=Ubuntu
// DISTRIB_RELEASE=14.04
// DISTRIB_CODENAME=trusty
// DISTRIB_DESCRIPTION="Ubuntu 14.04.5 LTS"
//

// Read release files as INI
func parseReleases() Release {
	pattern := "/etc/*-release"
	paths, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", pattern, err)
		os.Exit(1)
	}
	for _, p := range paths {
		parser := flags.NewParser(&release, flags.IgnoreUnknown)
		ini := flags.NewIniParser(parser)
		// ini.ParseAsDefaults = true
		err := ini.ParseFile(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s(ini): %s\n", p, err)
			os.Exit(1)
		}
		// if Verbose {
		// 	fmt.Printf("%s:\n%+v\n", p, release)
		// }
		// execute("cat", p)
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
