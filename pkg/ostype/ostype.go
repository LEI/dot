package ostype

import (
	"fmt"
	"gopkg.in/go-ini/ini.v1"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/LEI/dot/pkg/sliceutils"
)

var (
	// OS ...
	OS = runtime.GOOS

	// Types ...
	Types []string

	// Verbose ...
	Verbose bool

	release Release
)

func init() {
	Types = Get()
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

// Get OS types: name, release, family, distrib...
func Get() []string {
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
	name := strings.ToLower(r.Name)
	id := strings.ToLower(r.ID)
	if name != "" && id != "" && isNum(id) {
		types = append(types, name+id)
	} else if id != "" {
		types = append(types, id)
	} else if name != "" {
		types = append(types, name)
	}
	if r.IDLike != "" {
		for _, id := range strings.Split(r.IDLike, " ") {
			types = append(types, id)
		}
	}
	if r.DistribCodename != "" {
		types = append(types, r.DistribCodename)
	}
	types = append(types, parseOSTypes()...)
	return types
}

// Has ...
func Has(s ...string) bool {
	return sliceutils.Matches(s, Types)
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

// NAME="CentOS Linux"
// VERSION="7 (Core)"
// ID="centos"
// ID_LIKE="rhel fedora"
// VERSION_ID="7"
// PRETTY_NAME="CentOS Linux 7 (Core)"

// Read release files as INI
func parseReleases() Release {
	pattern := "/etc/*-release"
	paths, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", pattern, err)
		os.Exit(1)
	}
	for _, p := range paths {
		if err := ini.MapTo(release, p); err != nil {
			continue
			// return err
		}
	}
	// for _, p := range paths {
	// 	parser := flags.NewParser(&release, flags.IgnoreUnknown)
	// 	ini := flags.NewIniParser(parser)
	// 	// ini.ParseAsDefaults = true
	// 	err := ini.ParseFile(p)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 		// b, err := ioutil.ReadFile(p)
	// 		// if err != nil {
	// 		// 	fmt.Fprintf(os.Stderr, "Error while reading file %s: %s\n", p, err)
	// 		// 	return release
	// 		// }
	// 		// str := string(b)
	// 		// if str == "" {
	// 		// 	fmt.Fprintf(os.Stderr, "Empty release file: %s\n", p)
	// 		// 	return release
	// 		// }
	// 		// lines := strings.Split(str, "\n")
	// 		// if len(lines) != 2 || lines[0] == "" {
	// 		// 	fmt.Fprintf(os.Stderr, "Unexpected release file %s:---\n%s\n---\n", p, str)
	// 		// 	return release
	// 		// }
	// 		// if release.PrettyName == "" {
	// 		// 	release.PrettyName = lines[0]
	// 		// }
	// 		// if release.ID == "" {
	// 		// 	release.ID = strings.Split(lines[0], " ")[0]
	// 		// }
	// 		continue
	// 	}
	// 	if Verbose > 1 {
	// 		fmt.Printf("%s:\n%+v\n", p, release)
	// 	}
	// 	if Verbose > 2 {
	// 		fmt.Println(p)
	// 		execute("cat", p)
	// 	}
	// }
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
			if len(o) > 0 && o[0] != s {
				types = append(types, o[0])
			}
			types = append(types, s)
		}
	}
	return types
}
