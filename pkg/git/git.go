package git

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/LEI/dot/cli/config/tasks"
	"github.com/LEI/dot/pkg/executils"
	"github.com/LEI/dot/system"
)

var (
	// // ErrDirtyRepo ...
	// ErrDirtyRepo = fmt.Errorf("dirty repository")

	// Force ...
	Force bool

	// GitBin path
	GitBin = "git"

	// Stdout ...
	Stdout io.Writer
	// Stderr ...
	Stderr io.Writer

	cloneDepth = 1
	defaultBranch = "master"
	defaultRemote = "origin"
	repoFmt = "https://github.com/%s.git"
)

func init() {
	Stdout = os.Stdout
	Stdout = os.Stderr
}

// Repo ...
type Repo struct {
	Dir string
	URL string
	Branch string
	Remote string
}

func quiet() bool {
	return !tasks.Verbose
}

func parseURL(dir, url string) string {
	if url != "" && !strings.Contains(url, "https://") {
		url = fmt.Sprintf(repoFmt, url)
		// fmt.Println("NewRepo URL:", url)
	} else if url == "" && strings.Contains(dir, "/") && string(dir[0]) != "/" && string(dir[0]) != "~" {
		url = fmt.Sprintf(repoFmt, dir)
		// fmt.Println("NewRepo URL:", url)
	}
	return url
}

// NewRepo ...
func NewRepo(dir, url string) (*Repo, error) {
	repo := &Repo{
		Dir: dir,
		URL: parseURL(dir, url),
		Remote: defaultRemote,
		Branch: defaultBranch,
	}
	if repo.Dir == "" {
		return repo, fmt.Errorf("missing repo dir")
	}
	if repo.URL == "" {
		return repo, fmt.Errorf("missing repo url")
	}
	return repo, nil
}

// SetDir ...
func (r *Repo) SetDir(dir string) *Repo {
	r.Dir = dir
	return r
}

// SetURL ...
func (r *Repo) SetURL(url string) *Repo {
	r.URL = url
	return r
}

// Exec repo command
func (r *Repo) Exec(args ...string) (int) {
	return executils.Execute(GitBin, args...)
}

// ExecBuf repo command
func (r *Repo) ExecBuf(args ...string) (string, string, int) {
	stdOut, stdErr, status := executils.ExecuteBuf(GitBin, args...)
	out := strings.TrimRight(string(stdOut), "\n")
	err := strings.TrimRight(string(stdErr), "\n")
	return out, err, status
}

// Status repo
func (r *Repo) Status() error {
	args := []string{"status", "--porcelain"}
	if r.Dir != "" {
	    args = append([]string{"-C", r.Dir}, args...)
	}
	stdOut, stdErr, status := r.ExecBuf(args...)
	if status == 1 {
		return fmt.Errorf("%s: not a git directory", r.Dir)
	} else if status != 0 {
		return fmt.Errorf("%s: git status exit code %d", stdErr, status)
	}
	if stdOut != "" && !Force {
		// ErrDirtyRepo
		return fmt.Errorf("Uncommitted changes in %s:\n%s", r.Dir, stdOut)
	}
	return nil
}

// Clone repo
func (r *Repo) Clone() error {
	args := []string{"clone", r.URL}
	if r.Dir != "" {
		args = append(args, r.Dir)
	}
	if r.Branch != "" {
		args = append(args, "--branch", r.Branch)
	}
	if r.Remote != "" {
		args = append(args, "--origin", r.Remote)
	}
	if cloneDepth > 0 {
		args = append(args, "--depth", strconv.Itoa(cloneDepth))
	}
	if !tasks.Verbose {
		args = append(args, "--quiet")
	}
	// if tasks.Verbose {
	// 	fmt.Println("git clone", r.URL, r.Dir)
	// }
	// status := r.Exec(args...)
	// if status != 0 {
	//     return fmt.Errorf("git clone %s failed with exit code %d", r.URL, status)
	// }
	stdOut, stdErr, status := r.ExecBuf(args...)
	if status != 0 {
	    return fmt.Errorf(stdErr)
	}
	if stdErr != "" && tasks.Verbose {
		fmt.Fprintln(Stderr, stdErr)
	}
	if stdOut != "" && tasks.Verbose {
		fmt.Fprintf(Stdout, "%s\n", stdOut)
	}
	return nil
}

// Pull repo
func (r *Repo) Pull() error {
	args := []string{"pull", r.Remote, r.Branch}
	if r.Dir != "" {
	    args = append([]string{"-C", r.Dir}, args...)
	}
	if system.DryRun {
		args = append(args, "--dry-run")
	}
	if !tasks.Verbose {
		args = append(args, "--quiet")
	}
	// if tasks.Verbose {
	// 	fmt.Println("git pull", r.Remote, r.Branch)
	// }
	// status := r.Exec(args...)
	// if status != 0 {
	//     return fmt.Errorf("git clone %s failed with exit code %d", r.URL, status)
	// }
	stdOut, stdErr, status := r.ExecBuf(args...)
	if status != 0 {
		// '{{.URL}}': Could not resolve host: {{.Host}}
		// ErrNetworkUnreachable
		if Force && strings.HasPrefix(stdErr, "fatal: unable to access") {
			return nil
		}
		return fmt.Errorf(stdErr)
	}
	if stdErr != "" { // && tasks.Verbose {
		fmt.Fprintln(Stderr, stdErr)
	}
	if stdOut != "" && tasks.Verbose {
		fmt.Fprintf(Stdout, "%s\n", stdOut)
	}
	return nil
}
