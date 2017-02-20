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
	DefaultPath   string
	PathSep       = string(os.PathSeparator)
)

type Repository struct {
	Name    string
	Branch  string
	Path    string
	Remotes map[string]*Remote
	gitdir  string
}

func NewRepository(spec string /*, clonePath string, remotes ...*Remote*/) (*Repository, error) {
	name, dir, url, err := ParseSpec(spec)
	if err != nil {
		return nil, err
	}
	repo := &Repository{
		Name:    name,
		Branch:  DefaultBranch,
		Path:    dir,
		Remotes: make(map[string]*Remote, 0),
	}
	if url != "" {
		repo.AddRemote(DefaultRemote, url)
	}
	// remoteUrl := ""
	// if strings.HasPrefix(spec, string(os.PathSeparator)) { // filepath.IsAbs(spec)
	// 	exists, err := fileutil.Exists(spec)
	// 	if err != nil {
	// 		return repo, err
	// 	}
	// 	if !exists {
	// 		return repo, fmt.Errorf("%s: No such repository\n", spec)
	// 	}
	// 	// TODO find out branch
	// 	repo.Path = spec
	// 	repo.Name = filepath.Dir(repo.Path)
	// 	// return repo, nil
	// } else if strings.Contains(spec, "=") {
	// 	parts := strings.Split(spec, "=")
	// 	if len(parts) != 2 {
	// 		return repo, fmt.Errorf("%s: Invalid repository spec\n", spec)
	// 	}
	// 	repo.Name = parts[0]
	// 	remoteUrl = parts[1]
	// } else if strings.Contains(spec, "/") {
	// 	remoteUrl = spec
	// } else {
	// 	return repo, fmt.Errorf("%s: Unknown repository spec\n", spec)
	// }
	return repo, nil
}

func (repo *Repository) String() string {
	return fmt.Sprintf("%s@%s [%s] %+v", repo.Name, repo.Branch, repo.Path, repo.Remotes)
}

func (repo *Repository) WorkTree() string {
	if DefaultPath == "" && repo.Path == "" {
		fmt.Printf("Warning: %s\n", "No default git clone path")
	} else if repo.Path == "" {
		repo.Path = filepath.Join(DefaultPath, repo.Name)
	}
	return repo.Path
}

func (repo *Repository) GitDir() string {
	if repo.gitdir == "" {
		repo.gitdir = filepath.Join(repo.WorkTree(), ".git")
	}
	return repo.gitdir
}

func (repo *Repository) AddRemote(name string, url string) *Repository {
	repo.Remotes[name] = NewRemote(name, url)
	return repo
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
	exists, err := fileutil.Exists(repo.GitDir())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	return exists
}

func (repo *Repository) Clone() error {
	repo.WorkTree()
	for _, remote := range repo.Remotes {
		// fmt.Println("git", "clone", remote.URL, repo.Path)
		cmd := exec.Command("git", "clone", "--quiet", remote.URL, repo.Path)
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
	repo.WorkTree()
	for _, remote := range repo.Remotes {
		pull := []string{"-C", repo.Path, "pull"}
		if len(args) == 0 {
			args = []string{remote.Name} // , repo.Branch}
		}
		pull = append(pull, args...)
		// fmt.Printf("git %s\n", strings.Join(pull, " "))
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

// name=user/repo
// user/repo
func ParseSpec(str string) (string, string, string, error) {
	var nameSep = "="
	var name = str
	var dir string
	var url string
	var err error
	if strings.HasPrefix(str, PathSep) {
		exists, err := fileutil.Exists(str)
		if err != nil || !exists {
			return name, dir, url, err
		}
		dir = str
		name = filepath.Dir(dir)
	} else if strings.Contains(str, PathSep) {
		if strings.Contains(str, nameSep) {
			parts := strings.Split(str, nameSep)
			if len(parts) != 2 {
				return name, dir, url, fmt.Errorf("Invalid spec: '%s'", str)
			}
			name = parts[0]
			url = parts[1]
		} else {
			// name = filepath.Base(str)
			url = str
		}
	} else {
		err = fmt.Errorf("Unknown git origin: '%s'", str)
	}
	return name, dir, url, err
}
