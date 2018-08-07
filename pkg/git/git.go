package git

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/LEI/dot/pkg/executils"
)

var (
	cloneDepth = 1
	defaultBranch = "master"
	defaultRemote = "origin"
	quiet bool
	repoFmt = "https://github.com/%s.git"
	verbose bool
)

// Repo ...
type Repo struct {
	Dir string
	URL string
	Branch string
	Remote string
}

// NewRepo ...
func NewRepo(dir, url string) (*Repo, error) {
	repo := &Repo{
		Remote: defaultRemote,
		Branch: defaultBranch,
	}
	if dir == "" {
		return repo, fmt.Errorf("missing repo dir")
	}
	if url != "" && !strings.Contains(url, "https://") {
		url = fmt.Sprintf(repoFmt, url)
		// fmt.Println("NewRepo URL:", url)
	} else if url == "" && strings.Contains(dir, "/") && string(dir[0]) != "/" && string(dir[0]) != "~" {
		url = fmt.Sprintf(repoFmt, dir)
		// fmt.Println("NewRepo URL:", url)
	}
	if url == "" {
		return repo, fmt.Errorf("missing repo url")
	}
	repo.Dir = dir
	repo.URL = url
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

// Exec repo
func (r *Repo) Exec(args ...string) (string, string, int) {
	stdout, stderr, status := executils.Execute("git", args...)
	// out := strings.TrimRight(string(stdout), "\n")
	// err := strings.TrimRight(string(stderr), "\n")
	out := strings.TrimRight(string(stdout), "\n")
	err := strings.TrimRight(string(stderr), "\n")
	return out, err, status
}

// Status repo
func (r *Repo) Status() error {
	args := []string{"status", "--porcelain"}
	if r.Dir != "" {
	    args = append([]string{"-C", r.Dir}, args...)
	}
	stdout, stderr, status := r.Exec(args...)
	str := strings.TrimRight(stdout, "\n")
	err := strings.TrimRight(stderr, "\n")
	if status == 1 {
		return fmt.Errorf("%s: not a git directory", r.Dir)
	} else if status != 0 {
		return fmt.Errorf("%s: git status exit code %d", err, status)
	}
	if str != "" {
		// fmt.Println(str)
		return fmt.Errorf("%s: dirty git directory -> %s", r.Dir, str)
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
	if quiet {
		args = append(args, "--quiet")
	}
	// if verbose {
	// 	fmt.Println("git clone", r.URL, r.Dir)
	// }
	stdout, stderr, status := r.Exec(args...)
	if status != 0 {
	    return fmt.Errorf("git clone exit status %d: %s", status, stderr)
	}
	if stdout != "" && verbose {
		fmt.Println(stdout)
	}
	if stderr != "" && verbose {
		fmt.Fprintln(os.Stderr, stderr)
	}
	return nil
}

// Pull repo
func (r *Repo) Pull() error {
	args := []string{"pull", r.Remote, r.Branch}
	if r.Dir != "" {
	    args = append([]string{"-C", r.Dir}, args...)
	}
	if quiet {
		args = append(args, "--quiet")
	}
	// if verbose {
	// 	fmt.Println("git pull", r.Remote, r.Branch)
	// }
	stdout, stderr, status := r.Exec(args...)
	if status != 0 {
	    return fmt.Errorf("git pull exit status %d: %s", status, stderr)
	}
	if stdout != "" && verbose {
		fmt.Println(stdout)
	}
	if stderr != "" && verbose {
		fmt.Fprintln(os.Stderr, stderr)
	}
	return nil
}
