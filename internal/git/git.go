package git

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/LEI/dot/cli/config/tasks"
	"github.com/LEI/dot/internal/shell"
	"github.com/LEI/dot/system"
)

const (
	minVer = 2
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

	cloneDepth    = 1
	defaultBranch = "master"
	defaultRemote = "origin"
	repoFmt       = "https://github.com/%s.git"
)

func init() {
	if err := checkGitVersion(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
	Stdout = os.Stdout
	Stdout = os.Stderr
}

func checkGitVersion() error {
	cmd := exec.Command("git", "--version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	str := string(out)
	ver := strings.TrimPrefix(str, "git version ")
	// fmt.Println("GIT_VERSION", ver)
	// if ver == "" {
	// 	return fmt.Errorf("%s: unable to parse git version", str)
	// }
	parts := strings.Split(ver, ".")
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}
	if major < minVer {
		return fmt.Errorf("git version %s is required", string(minVer))
	}
	return nil
}

// Repo ...
type Repo struct {
	Dir    string
	URL    string
	Branch string
	Remote string
}

// func quiet() bool {
// 	return tasks.Verbose == 0
// }

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
		Dir:    dir,
		URL:    parseURL(dir, url),
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

// // ExecStatus repo command
// func (r *Repo) ExecStatus(args ...string) int {
// 	return shell.Run(GitBin, args...)
// }

// Exec repo command
func (r *Repo) Exec(args ...string) (string, string, error) {
	var stdout, stderr *bytes.Buffer // io.Writer
	_, err := shell.Exec(nil, stdout, stderr, GitBin, args...)
	return stdout.String(), stderr.String(), err
}

// Status repo
func (r *Repo) Status() error {
	args := []string{"status", "--porcelain"}
	if r.Dir != "" {
		args = append([]string{"-C", r.Dir}, args...)
	}
	stdout, stderr, err := r.Exec(args...)
	if err != nil {
		// return fmt.Errorf("%s: not a git directory", r.Dir)
		return fmt.Errorf("git status %s: %s", r.Dir, err)
	}
	if stderr != "" {
		fmt.Fprintf(os.Stderr, stderr)
	}
	if stdout != "" && !Force {
		// ErrDirtyRepo
		return fmt.Errorf("Uncommitted changes in %s:\n%s", r.Dir, stdout)
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
	if tasks.Verbose == 0 {
		args = append(args, "--quiet")
	}
	// if tasks.Verbose > 0 {
	// 	fmt.Println("git clone", r.URL, r.Dir)
	// }
	// status := r.ExecStatus(args...)
	// if status != 0 {
	//     return fmt.Errorf("git clone %s failed with exit code %d", r.URL, status)
	// }
	stdout, stderr, err := r.Exec(args...)
	if err != nil {
		return fmt.Errorf("Unable to clone %s in %s:\n%s", r.URL, r.Dir, err)
		// return fmt.Errorf(stderr)
	}
	if stderr != "" && tasks.Verbose > 0 {
		fmt.Fprintln(Stderr, stderr)
	}
	if stdout != "" && tasks.Verbose > 0 {
		fmt.Fprintf(Stdout, "%s\n", stdout)
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
	if tasks.Verbose == 0 {
		args = append(args, "--quiet")
	}
	// if tasks.Verbose > 0 {
	// 	fmt.Println("git pull", r.Remote, r.Branch)
	// }
	// status := r.ExecStatus(args...)
	// if status != 0 {
	//     return fmt.Errorf("git clone %s failed with exit code %d", r.URL, status)
	// }
	stdout, stderr, err := r.Exec(args...)
	if err != nil {
		// '{{.URL}}': Could not resolve host: {{.Host}}
		// ErrNetworkUnreachable
		if Force && strings.HasPrefix(stderr, "fatal: unable to access") {
			return nil
		}
		// return fmt.Errorf("Unable to pull %s in %s:\n%s", r.URL, r.Dir, stderr)
		return err
	}
	if stderr != "" { // && tasks.Verbose > 0 {
		fmt.Fprintln(Stderr, stderr)
	}
	if stdout != "" && tasks.Verbose > 0 {
		fmt.Fprintf(Stdout, "%s\n", stdout)
	}
	return nil
}
