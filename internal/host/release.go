package host

import (
	"fmt"
	"os"
	"path/filepath"
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
	name := strings.ToLower(r.Name)
	id := strings.ToLower(r.ID)
	if name != "" && id != "" && isNum(id) {
		s = append(s, name+id)
	} else if id != "" {
		s = append(s, id)
	} else if name != "" {
		s = append(s, name)
	}
	if r.IDLike != "" {
		for _, id := range strings.Fields(r.IDLike) {
			s = append(s, id)
		}
	}
	if r.DistribCodename != "" {
		s = append(s, r.DistribCodename)
	}
	return s
}

func isNum(v string) bool {
	_, err := strconv.Atoi(v)
	return err == nil
}