package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func New(name string, path string) *Repo {
	return &Repo{
		Name: name,
		Branch: "master",
		Path: path,
		Remotes: make(map[string]*Remote, 0),
	}
}

func NewRemote(name string, url string) *Remote {
	return &Remote{
		Name: name,
		URL: UrlScheme(url),
	}
}

func UrlScheme(str string) string {
	url := "https://github.com/" + str + ".git"
	switch {
	case strings.HasPrefix(str, "git@github.com"),
		strings.HasPrefix(str, "git://"),
		strings.HasPrefix(str, "http://"),
		strings.HasPrefix(str, "https://"),
		strings.HasPrefix(str, "ssh://"):
		url = str
	}
	return url
}

type Repo struct {
	Name    string
	Branch  string
	Path    string
	Remotes map[string]*Remote
}

type Remote struct {
	Name string
	URL  string
}

func (repo *Repo) String() string {
	return fmt.Sprintf("%s@%s", repo.Name, repo.Branch)
}

func (repo *Repo) AddRemote(name string, url string) *Repo {
	repo.Remotes[name] = NewRemote(name, url)
	return repo
}

func (repo *Repo) IsCloned() bool {
	_, err := os.Stat(repo.Path)
	if err != nil { // && os.IsNotExist(err)
		return false
	}
	return true
}

func (repo *Repo) Clone() error {
	for _, remote := range repo.Remotes {
		cmd := exec.Command("git", "clone", remote.URL, repo.Path)
		out, err := cmd.CombinedOutput()
		if len(out) > 0 {
			fmt.Printf("%s: %s", repo.Name, out)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *Repo) Pull(args ...string) error {
	for _, remote := range repo.Remotes {
		pull := []string{"-C", repo.Path, "pull"}
		if len(args) == 0 {
			args = []string{remote.Name, repo.Branch}
		}
		pull = append(pull, args...)
		cmd := exec.Command("git", pull...)
		// "--git-dir", dir+"/.git",
		// "--work-tree", dir,
		out, err := cmd.CombinedOutput()
		if len(out) > 0 {
			fmt.Printf("%s: %s", repo.Name, out)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
