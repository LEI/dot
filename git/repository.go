package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"github.com/LEI/dot/fileutil"
)

type Repository struct {
	Name    string
	Branch  string
	Path    string
	Remotes map[string]*Remote
}

func NewRepository(spec string/*, path string, remotes ...*Remote*/) (*Repository, error) {
	repo := &Repository{Name: spec, Branch: "master"}
	remote := ""
	if strings.HasPrefix(spec, string(os.PathSeparator)) { // filepath.IsAbs(spec)
		if !fileutil.Exists(spec) {
			repo.Path = spec
			// TODO find out branch
			return repo, nil
		} else {
			return repo, fmt.Errorf("%s: No such repository\n", spec)
		}
	}
	if strings.Contains(spec, "=") {
		parts := strings.Split(spec, "=")
		if len(parts) != 2 {
			return repo, fmt.Errorf("%s: Invalid repository spec\n", spec)
		}
		repo.Name = parts[0]
		remote = parts[1]
	}
	if remote != "" {
		repo.AddRemote("origin", remote)
	}

	return repo, nil
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
