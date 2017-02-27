package git

import (
	"fmt"
	"strings"
)

var (
	User = "git"
	Host = "github.com"
)

type Remote struct {
	Name string
	URL  string
}

func NewRemote(name string, url string) *Remote {
	r := &Remote{Name: name}
	r.SetUrl(url)
	return r
}

func (r *Remote) String() string {
	return fmt.Sprintf("%s %s", r.Name, r.URL)
}

func (r *Remote) Pull(path string) error {
	cmd := []string{"-C", path, "pull"}
	// "--git-dir", dir+"/.git", "--work-tree", dir,
	if Quiet {
		cmd = append(cmd, "--quiet")
	}
	cmd = append(cmd, r.Name) // repo.Branch
	err := Exec(cmd...)
	if err != nil {
		return err
	}
	return nil
}

func (r *Remote) SetUrl(url string) *Remote {
	switch {
	case strings.HasPrefix(url, "git@"),
		strings.HasPrefix(url, "git://"),
		// strings.HasPrefix(url, "http://"),
		strings.HasPrefix(url, "https://"),
		strings.HasPrefix(url, "ssh://"):
		r.URL = url
	default:
		if Https {
			r.URL = fmt.Sprintf("https://%s/%s", Host, url)
		} else {
			r.URL = fmt.Sprintf("%s@%s:%s", User, Host, url)
		}
	}
	if !strings.HasSuffix(url, ".git") {
		r.URL += ".git"
	}
	return r
}
