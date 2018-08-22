package git

import (
	"fmt"
	"strconv"
	"strings"
)

// Repo ...
type Repo struct {
	Dir    string
	URL    string
	Branch string
	Remote string
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

func parseURL(dir, url string) string {
	if url != "" && !strings.Contains(url, "https://") {
		url = fmt.Sprintf(remoteURLFormat, url)
		// fmt.Println("NewRepo URL:", url)
	} else if url == "" && strings.Contains(dir, "/") && string(dir[0]) != "/" && string(dir[0]) != "~" {
		url = fmt.Sprintf(remoteURLFormat, dir)
		// fmt.Println("NewRepo URL:", url)
	}
	return url
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

// Exec git command
func (r *Repo) Exec(args ...string) (string, string, error) {
	if r.Dir != "" {
		args = append([]string{"-C", r.Dir}, args...)
	}
	return git(args...)
}

// Status repo
func (r *Repo) Status() error {
	args := []string{"status", "--porcelain"}
	stdout, stderr, err := r.Exec(args...)
	if err != nil {
		// return fmt.Errorf("%s: not a git directory", r.Dir)
		return fmt.Errorf("git status %s: %s", r.Dir, err)
	}
	if stderr != "" {
		fmt.Fprintf(Stderr, stderr)
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
	if r.Branch != "" {
		args = append(args, "--branch", r.Branch)
	}
	if r.Remote != "" {
		args = append(args, "--origin", r.Remote)
	}
	if cloneDepth > 0 {
		args = append(args, "--depth", strconv.Itoa(cloneDepth))
	}
	if Quiet {
		args = append(args, "--quiet")
	}
	// if Verbose > 0 {
	// 	fmt.Println("git clone", r.URL, r.Dir)
	// }
	// status := r.ExecStatus(args...)
	// if status != 0 {
	//     return fmt.Errorf("git clone %s failed with exit code %d", r.URL, status)
	// }
	stdout, stderr, err := git(args...)
	if err != nil {
		return fmt.Errorf("Unable to clone %s in %s:\n%s", r.URL, r.Dir, err)
		// return fmt.Errorf(stderr)
	}
	if stderr != "" { // && Verbose > 0 {
		fmt.Fprintln(Stderr, stderr)
	}
	if stdout != "" && Verbose > 0 {
		fmt.Fprintln(Stdout, stdout)
	}
	return nil
}

// Pull repo
func (r *Repo) Pull() error {
	args := []string{"pull", r.Remote, r.Branch}
	if DryRun {
		args = append(args, "--dry-run")
	}
	if Quiet {
		args = append(args, "--quiet")
	}
	// if Verbose > 0 {
	// 	fmt.Println("git pull", r.Remote, r.Branch)
	// }
	// status := r.ExecStatus(args...)
	// if status != 0 {
	//     return fmt.Errorf("git clone %s failed with exit code %d", r.URL, status)
	// }
	stderr, stdout, err := r.Exec(args...)
	if err != nil {
		if Force && strings.HasPrefix(stderr, "fatal: unable to access") {
			return nil
		}
		return err
	}
	if stderr != "" { // && Verbose > 0 {
		fmt.Fprintln(Stderr, stderr)
	}
	if stdout != "" && Verbose > 0 {
		fmt.Fprintln(Stdout, stdout)
	}
	// stdout, stderr, err := r.Exec(args...)
	// if err != nil {
	// 	// '{{.URL}}': Could not resolve host: {{.Host}}
	// 	// ErrNetworkUnreachable
	// 	if Force && strings.HasPrefix(stderr, "fatal: unable to access") {
	// 		return nil
	// 	}
	// 	// return fmt.Errorf("Unable to pull %s in %s:\n%s", r.URL, r.Dir, stderr)
	// 	return err
	// }
	// if stderr != "" { // && Verbose > 0 {
	// 	fmt.Fprintln(Stderr, stderr)
	// }
	// if stdout != "" && Verbose > 0 {
	// 	fmt.Fprintf(Stdout, "%s\n", stdout)
	// }
	return nil
}
