package main

// TODO: url parsing & support hg?
// https://github.com/sourcegraph/go-vcs
// https://github.com/libgit2/git2go

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/LEI/dot/dotfile"
	"github.com/LEI/dot/utils"
)

// Repo ...
type Repo struct {
	Path, URL      string
	Remote, Branch string
}

var (
	online bool

	defaultRemote = "origin"
	defaultBranch = "master"

	repoFmt = "https://github.com/%s.git"

	// ErrNoGitDir ...
	ErrNoGitDir = fmt.Errorf("no .git directory")

	// ErrDirtyRepo ...
	ErrDirtyRepo = fmt.Errorf("dirty repository")

	// ErrNetworkUnreachable ...
	ErrNetworkUnreachable = fmt.Errorf("network unreachable")
)

func init() {
	online = networkReachable()
}

// NewRepo ...
func NewRepo(p, url string) *Repo {
	if url != "" && !strings.Contains(url, "https://") {
		url = fmt.Sprintf(repoFmt, url)
		fmt.Println("NewRepo URL:", url)
	} else if url == "" && strings.Contains(p, "/") && string(p[0]) != "/" && string(p[0]) != "~" {
		url = fmt.Sprintf(repoFmt, p)
		fmt.Println("NewRepo URL:", url)
	}
	r := &Repo{
		Path:   p,
		URL:    url,
		Remote: defaultRemote,
		Branch: defaultBranch,
	}
	return r
}

func (r *Repo) checkRepo() error {
	if !isGitDir(r.Path) {
		return ErrNoGitDir
	}
	args := []string{"-C", r.Path, "diff-index", "--quiet", "HEAD"}
	_, stderr, status := dotfile.ExecCommand("git", args...)
	if status == 1 {
		return ErrDirtyRepo
	} else if status != 0 {
		return fmt.Errorf("check repo unknown error: %s", stderr)
	}
	// c := exec.Command("git", args...)
	// err := c.Run()
	// if err != nil {
	// 	if exitError, ok := err.(*exec.ExitError); ok {
	// 		// fmt.Fprintf(os.Stderr, "Uncommitted changes in '%s'", r.Path)
	// 		if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
	// 			if status.ExitStatus() == 1 {
	// 				return ErrDirtyRepo
	// 			}
	// 		}
	// 	}
	// 	fmt.Fprintf(os.Stderr, "%s: %s\n", r.Path, err)
	// 	return err
	// }
	return nil
}

// Pull ...
func (r *Repo) Pull() error {
	if !online {
		return ErrNetworkUnreachable
	}
	args := []string{"-C", r.Path, "pull", r.Remote, r.Branch}
	switch true {
	case Verbose == 0:
		args = append(args, "--quiet")
	// case Verbose == 1:
	// 	fmt.Printf("git pull %s %s\n", r.Remote, r.Branch)
	case Verbose > 1:
		fmt.Printf("git %s\n", strings.Join(args, " "))
	}
	c := exec.Command("git", args...)
	if Verbose > 0 {
		c.Stdout = os.Stdout
	}
	c.Stderr = os.Stderr
	return c.Run()
}

// Clone ...
func (r *Repo) Clone() error {
	if !online {
		return ErrNetworkUnreachable
	}
	if utils.Exist(r.Path) {
		return r.checkRemote()
	}
	args := []string{"clone", "--recursive"}
	if r.Branch != "" {
		args = append(args, "--branch", r.Branch)
	}
	args = append(args, r.URL, r.Path)
	switch true {
	case Verbose == 0:
		args = append(args, "--quiet")
	// case Verbose == 1:
	// 	fmt.Printf("git clone %s %s\n", r.URL, r.Path)
	case Verbose > 1:
		fmt.Printf("git %s\n", strings.Join(args, " "))
	}
	c := exec.Command("git", args...)
	if Verbose > 0 {
		c.Stdout = os.Stdout
	}
	c.Stderr = os.Stderr
	return c.Run()
}

func networkReachable() bool {
	timeout := time.Duration(1 * time.Second)
	_, err := net.DialTimeout("tcp", "github.com:443", timeout)
	// fmt.Println(err)
	return err == nil
}

func (r *Repo) checkRemote() error {
	urlKey := "remote." + r.Remote + ".url"
	args := []string{"-C", r.Path, "config", "--local", "--get", urlKey}
	// fmt.Printf("git %s\n", strings.Join(args, " "))
	stdout, stderr, status := dotfile.ExecCommand("git", args...)
	if status != 0 {
		return fmt.Errorf(stderr)
	}
	url := strings.TrimRight(stdout, "\n")
	if parseRepo(r.URL) != parseRepo(url) {
		return fmt.Errorf("remote mismatch: url is '%s' but existing repo has '%s'", r.URL, url)
	}
	return nil
}

func parseRepo(str string) string {
	str = strings.TrimSuffix(str, ".git")
	str = strings.Replace(str, ":", "/", 1)
	parts := strings.Split(str, "/")
	if len(parts) > 1 {
		return strings.Join(parts[len(parts)-2:], "/")
	}
	return str
}

func isGitDir(s string) bool {
	return utils.Exist(filepath.Join(s, ".git"))
}