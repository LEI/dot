package cmd

import (
	"bytes"
	"fmt"
	// "io/ioutil"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

// Repo ...
type Repo struct {
	Path, URL      string
	remote, branch string
}

var (
	online bool

	defaultRemote = "origin"
	defaultBranch = "master"

	// ErrDirtyRepo ...
	ErrDirtyRepo = fmt.Errorf("Dirty repository")

	// ErrNetworkUnreachable ...
	ErrNetworkUnreachable = fmt.Errorf("Network unreachable")
)

func init() {
	online = networkReachable()
}

// NewRepo ...
func NewRepo(p, url string) *Repo {
	if !strings.Contains(url, "https://") {
		url = "https://github.com/" + url
	}
	r := &Repo{Path: p, URL: url}
	if r.remote == "" {
		r.remote = defaultRemote
	}
	if r.branch == "" {
		r.branch = defaultBranch
	}
	return r
}

func (r *Repo) checkRepo() error {
	args := []string{"-C", r.Path, "diff-index", "--quiet", "HEAD"}
	c := exec.Command("git", args...)
	err := c.Run()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// fmt.Fprintf(os.Stderr, "Uncommited changes in '%s'", r.Path)
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok && status == 256 {
				return ErrDirtyRepo
			}
		}
		fmt.Fprintf(os.Stderr, "%s: %s\n", r.Path, err)
		return err
	}
	return nil
}

// Pull ...
func (r *Repo) Pull() error {
	if !online {
		return ErrNetworkUnreachable
	}
	args := []string{"-C", r.Path, "pull", r.remote, r.branch, "--quiet"}
	// fmt.Printf("git %s\n", strings.Join(args, " "))
	if Verbose {
		fmt.Printf("git pull %s %s\n", r.remote, r.branch)
	}
	c := exec.Command("git", args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}
	return nil
}

// Clone ...
func (r *Repo) Clone() error {
	if !online {
		return ErrNetworkUnreachable
	}
	if _, err := os.Stat(r.Path); err == nil {
		return r.checkRemote()
	}
	args := []string{"clone", r.URL, r.Path, "--recursive", "--quiet"}
	// fmt.Printf("git %s\n", strings.Join(args, " "))
	if Verbose {
		fmt.Printf("git clone %s %s\n", r.URL, r.Path)
	}
	c := exec.Command("git", args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}
	// actual := strings.TrimRight(out, "\n")
	return nil
}

func networkReachable() bool {
	timeout := time.Duration(1 * time.Second)
	_, err := net.DialTimeout("tcp", "github.com:443", timeout)
	// fmt.Println(err)
	return err == nil
}

func (r *Repo) checkRemote() error {
	args := []string{"-C", r.Path, "config", "--local", "--get", "remote." + r.remote + ".url"}
	// fmt.Printf("git %s\n", strings.Join(args, " "))
	c := exec.Command("git", args...)
	stdout, err := c.StdoutPipe()
	if err != nil {
		return err
	}
	// stderr, err := c.StderrPipe()
	// if err != nil {
	// 	return err
	// }
	if err := c.Start(); err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	buf.ReadFrom(stdout)
	out := buf.String()
	// out, _ := ioutil.ReadAll(stdout)
	// fmt.Printf("stdout: %s\n", out)
	// outErr, _ := ioutil.ReadAll(stderr)
	// fmt.Printf("stderr: %s\n", outErr)
	if err := c.Wait(); err != nil {
		return err
	}
	url := strings.TrimRight(out, "\n")
	// fmt.Println(parseRepo(r.URL), parseRepo(url))
	if parseRepo(r.URL) != parseRepo(url) {
		return fmt.Errorf("Remote mismatch: url is '%s' but existing repo has '%s'", r.URL, url)
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
