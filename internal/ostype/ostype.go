package ostype

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

// NAME="CentOS Linux"
// VERSION="7 (Core)"
// ID="centos"
// ID_LIKE="rhel fedora"
// VERSION_ID="7"
// PRETTY_NAME="CentOS Linux 7 (Core)"

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"gopkg.in/go-ini/ini.v1"
)

var (
	// List stores the list of OS types
	List []string

	// release *Release

	// releasePattern is used to find release files
	releasePattern = "/etc/*-release"
)

func init() {
	List = Get()
}

// Release ...
type Release struct {
	ID         string `ini:"ID"`
	IDLike     string `ini:"ID_LIKE"`
	Name       string `ini:"NAME"`
	PrettyName string `ini:"PRETTY_NAME"`
	Version    string `ini:"VERSION"`
	VersionID  string `ini:"VERSION_ID"`
	// HomeURL string `ini:"HOME_URL"`
	// SupportURL string `ini:"SUPPORT_URL"`
	// BugReportURL string `ini:"BUG_REPORT_URL"`
	DistribID          string `ini:"DISTRIB_ID"`
	DistribRelease     string `ini:"DISTRIB_RELEASE"`
	DistribCodename    string `ini:"DISTRIB_CODENAME"`
	DistribDescription string `ini:"DISTRIB_DESCRIPTION"`
}

// NewRelease read release files as INI
func NewRelease() *Release {
	release := &Release{}
	paths, err := filepath.Glob(releasePattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", releasePattern, err)
		os.Exit(1)
	}
	for _, p := range paths {
		if err := ini.MapTo(&release, p); err != nil {
			// fmt.Fprintf(os.Stderr, "%s: %s\n", p, err)
			continue // return err
		}
	}
	return release
}

// Get OS types: name, release, family, distrib...
func Get() []string {
	types := []string{runtime.GOOS}
	release := NewRelease()
	name := strings.ToLower(release.Name)
	id := strings.ToLower(release.ID)
	if name != "" && id != "" && isNum(id) {
		types = append(types, name+id)
	} else if id != "" {
		types = append(types, id)
	} else if name != "" {
		types = append(types, name)
	}
	if release.IDLike != "" {
		for _, id := range strings.Fields(release.IDLike) {
			types = append(types, id)
		}
	}
	if release.DistribCodename != "" {
		types = append(types, release.DistribCodename)
	}
	types = append(types, parseEnvVar("OSTYPE")...)
	return types
}

func isNum(v string) bool {
	_, err := strconv.Atoi(v)
	return err == nil
}

func parseEnvVar(name string) []string {
	types := make([]string, 0)
	if o, ok := os.LookupEnv(name); ok && o != "" {
		types = append(types, o)
	} else { // !ok || s == ""
		// fmt.Printf("%s='%s' (%v)\n", name, s, ok)
		out, err := exec.Command("bash", "-c", "printf '%s' \"$"+name+"\"").Output()
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

// Has OS type
func Has(s ...string) bool {
	// ok, _ := matches(s, List)
	// return ok
	ok, err := matches(s, List)
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
