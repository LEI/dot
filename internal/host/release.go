package host

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	ini "gopkg.in/go-ini/ini.v1"
)

var (
	// release *Release

	// releasePattern is used to find release files
	releasePattern = "/etc/*-release"
)

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

// NewRelease parses release files as INI.
func NewRelease() *Release {
	r := &Release{}
	paths, err := filepath.Glob(releasePattern)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", releasePattern, err)
		os.Exit(1)
	}
	for _, p := range paths {
		if err := ini.MapTo(&r, p); err != nil {
			// fmt.Fprintf(os.Stderr, "%s: %s\n", p, err)
			continue // return err
		}
	}
	return r
}

// Parse release information into a slice of strings.
func (r *Release) Parse() (s []string) {
	s = append(s, parseRelease(r)...)
	if r.IDLike != "" {
		for _, id := range strings.Fields(r.IDLike) {
			s = append(s, id)
		}
	}
	if !isNum(r.Version) {
		re := regexp.MustCompile(`\((\w+)\)`)
		if matches := re.FindStringSubmatch(r.Version); len(matches) > 1 {
			// s = append(s, strings.ToLower(matches[1]))
			s = append(s, matches[1:]...)
		}
	}
	s = append(s, parseDistrib(r)...)
	return s
}

func parseRelease(r *Release) (s []string) {
	name := strings.ToLower(r.Name)
	id := strings.ToLower(r.ID)
	if name != "" && id != "" && isNum(id) {
		s = append(s, name+id)
	} else if id != "" {
		s = append(s, id)
		if !isNum(id) && isNum(r.VersionID) {
			s = append(s, id+r.VersionID)
		}
	} else if name != "" {
		s = append(s, name)
	}
	return s
}

func parseDistrib(r *Release) (s []string) {
	if r.DistribCodename != "" {
		s = append(s, r.DistribCodename)
	}
	if r.DistribID != "" {
		did := strings.ToLower(r.DistribID)
		if r.DistribRelease != "" {
			s = append(s, r.DistribRelease)
			parts := strings.Split(r.DistribRelease, ".")
			if len(parts) > 1 && isNum(parts[0]) {
				s = append(s, did+parts[0])
			}
		}
	}
	return s
}

func isNum(v string) bool {
	_, err := strconv.Atoi(v)
	return err == nil
}
