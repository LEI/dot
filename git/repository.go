package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type Repository struct {
	Name    string
	Branch  string
	Path    string
	Remotes map[string]*Remote
}

func NewRepository(name string, path string, remotes ...*Remote) *Repository {
	return &Repository{
		Name: name,
		Branch: "master",
		Path: path,
		// Remotes: make(map[string]*Remote, 0),
		Remotes: remotes
	}
}

func (repo *Repository) String() string {
	return fmt.Sprintf("%s@%s", repo.Name, repo.Branch)
}

func (repo *Repository) AddRemote(name string, url string) *Repository {
	repo.Remotes[name] = NewRemote(name, url)
	return repo
}

func (repo *Repository) CloneOrPull() error {
	if repo != nil {
		if repo.IsCloned() {
			err := repo.Pull()
			if err != nil {
				return err
			}
		} else {
			err := repo.Clone()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (repo *Repository) IsCloned() bool {
	_, err := os.Stat(repo.Path)
	if err != nil {
		if os.IsExist(err) {
			panic(err)
		}
		return false
	}
	return true
}

func (repo *Repository) Clone() error {
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

func (repo *Repository) Pull(args ...string) error {
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
