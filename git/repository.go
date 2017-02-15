package git

import (
	"fmt"
	"github.com/LEI/dot/fileutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	DefaultBranch = "master"
	DefaultRemote = "origin"
	DefaultClonePath = filepath.Join(os.Getenv("HOME"), ".dot")
	// TODO init() viper.Get("target")
)

type Repository struct {
	Name   string
	Branch string
	Path   string // Git work tree
	GitDir string
	Remotes map[string]*Remote
}

func NewRepository(spec string /*, path string, remotes ...*Remote*/) (*Repository, error) {
	repo := &Repository{
		Name:    spec,
		Branch:  DefaultBranch,
		Remotes: make(map[string]*Remote, 0),
	}
	remoteUrl := ""
	if strings.HasPrefix(spec, string(os.PathSeparator)) { // filepath.IsAbs(spec)
		exists, err := fileutil.Exists(spec)
		if err != nil {
			return repo, err
		}
		if !exists {
			return repo, fmt.Errorf("%s: No such repository\n", spec)
		}
		// TODO find out branch
		repo.Path = spec
		repo.Name = filepath.Dir(repo.Path)
		// return repo, nil
	} else if strings.Contains(spec, "=") {
		parts := strings.Split(spec, "=")
		if len(parts) != 2 {
			return repo, fmt.Errorf("%s: Invalid repository spec\n", spec)
		}
		repo.Name = parts[0]
		remoteUrl = parts[1]
	} else if strings.Contains(spec, "/") {
		remoteUrl = spec
	} else {
		return repo, fmt.Errorf("%s: Unknown repository spec\n", spec)
	}
	if remoteUrl != "" {
		repo.AddRemote(DefaultRemote, remoteUrl)
	}
	if repo.Path == "" {
		repo.Path = filepath.Join(DefaultClonePath, repo.Name)
	}
	repo.GitDir = filepath.Join(repo.Path, ".git")
	return repo, nil
}

func (repo *Repository) String() string {
	return fmt.Sprintf("%s@%s [%s] %+v", repo.Name, repo.Branch, repo.Path, repo.Remotes)
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
	exists, err := fileutil.Exists(repo.GitDir) // repo.Path
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return exists
}

func (repo *Repository) Clone() error {
	for _, remote := range repo.Remotes {
		fmt.Println("git", "clone", remote.URL, repo.Path)
		cmd := exec.Command("git", "clone", remote.URL, repo.Path)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
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
		fmt.Printf("git %s\n", strings.Join(pull, " "))
		cmd := exec.Command("git", pull...)
		// "--git-dir", dir+"/.git",
		// "--work-tree", dir,
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}
