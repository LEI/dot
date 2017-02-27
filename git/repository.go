package git

import (
	"fmt"
	"os"
	"path/filepath"
	// "strings"
)

var (
	DefaultBranch   = "master"
	DefaultRemote   = "origin"
	DefaultCloneDir string
)

type Repository struct {
	Name    string
	Branch  string
	Path    string
	Remotes map[string]*Remote
	GitDir  string
}

// type Options map[string]interface{}

func NewRepo(name string, clonePath string, remoteUrl string) (*Repository, error) {
	repo := &Repository{
		Name:    name,
		Branch:  DefaultBranch,
		Path:    clonePath,
		Remotes: make(map[string]*Remote, 0),
	}
	if remoteUrl != "" {
		repo.AddRemote(DefaultRemote, remoteUrl)
	} else {
		return repo, fmt.Errorf("Empty remote url in %s", repo)
	}
	return repo, nil
}

func (repo *Repository) String() string {
	return fmt.Sprintf("%s@%s [%s] %+v", repo.Name, repo.Branch, repo.Path, repo.Remotes)
}

func (repo *Repository) WorkTree() string {
	if DefaultCloneDir == "" && repo.Path == "" {
		fmt.Printf("Warning: %s\n", "No default git clone path")
	} else if repo.Path == "" {
		repo.Path = filepath.Join(DefaultCloneDir, repo.Name)
	}
	return repo.Path
}

func (repo *Repository) GetGitDir() string {
	if repo.GitDir == "" {
		repo.GitDir = filepath.Join(repo.WorkTree(), ".git")
	}
	return repo.GitDir
}

func (repo *Repository) AddRemote(name string, url string) *Repository {
	repo.Remotes[name] = NewRemote(name, url)
	return repo

}

func (repo *Repository) Status() error {
	repo.WorkTree()
	err := Exec("-C", repo.Path, "status")
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) CloneOrPull() error {
	repo.WorkTree()
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
	return nil
}

func (repo *Repository) IsCloned() bool {
	fi, err := os.Stat(repo.GetGitDir())
	if err != nil && os.IsExist(err) {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if fi == nil {
		return false
	}
	return true
}

func (repo *Repository) Clone() error {
	repo.WorkTree()
	for _, remote := range repo.Remotes {
		cmd := []string{"clone"}
		if Quiet {
			cmd = append(cmd, "--quiet")
		}
		cmd = append(cmd, remote.URL)
		cmd = append(cmd, repo.Path)
		err := Exec(cmd...)
		if err != nil {
			return err
		}
	}
	return nil
}

func (repo *Repository) Pull() error {
	repo.WorkTree()
	for _, remote := range repo.Remotes {
		err := remote.Pull(repo.Path)
		if err != nil {
			return err
		}
	}
	return nil
}
