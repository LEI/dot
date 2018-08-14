package ostype

import (
	"fmt"
	"gopkg.in/go-ini/ini.v1"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var (
	// OS ...
	OS = runtime.GOOS

	// List of OS types
	List []string

	// Verbose ...
	Verbose bool

	release Release
)

func init() {
	List = Get()
	// fmt.Printf("OS types: %+v\n", List)
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
	r := parseRelease()
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
	types = append(types, parseEnvVar("OSTYPE")...)
	return types
}

func isNum(v string) bool {
	_, err := strconv.Atoi(v)
	return err == nil
}

// Has OS type
func Has(s ...string) bool {
	return matches(s, List)
}

func matches(in []string, list []string) bool {
	nn := 0
	for _, pattern := range in {
		negated := pattern[0] == '!'
		if negated {
			pattern = pattern[1:]
		}
		for _, str := range list {
			matched, err := regexp.MatchString(pattern, str)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s: %s\n", pattern, err)
				os.Exit(1)
			}
			if negated && matched {
				// fmt.Println("NEGATED", pattern, str, in, list)
				return false
			}
			if matched {
				// fmt.Println("MATCHED", pattern, str, in, list)
				return true
			}
		}
		if negated {
			nn++
			// fmt.Println("MATCHED NEGATED", pattern, in, list)
			// return true
		}
	}
	if nn == len(in) && nn > 0 {
		// fmt.Println("MATCHED NEGATED", in, list)
		return true
	}
	// fmt.Println("NOMATCH", in, list)
	return false
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

// Read release file as INI
func parseRelease() Release {
	pattern := "/etc/*-release"
	paths, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", pattern, err)
		os.Exit(1)
	}
	for _, p := range paths {
		if err := ini.MapTo(&release, p); err != nil {
			// fmt.Fprintf(os.Stderr, "%s: %s\n", p, err)
			continue
			// return err
		}
		// cmd := exec.Command("cat", p)
		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr
		// cmd.Run()
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
