package dot

import (
	"bytes"
	"fmt"
	// "io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Repo ...
type Repo struct {
	Path, URL string
	remote, branch string
}

var (
	defaultRemote = "origin"
	defaultBranch = "master"
)

// NewRepo ...
func NewRepo(p, url string) *Repo {
	r := &Repo{Path: p, URL: url}
	if r.remote == "" {
		r.remote = defaultRemote
	}
	if r.branch == "" {
		r.branch = defaultBranch
	}
	return r
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
		log.Fatalf("Remote mismatch: url is '%s' but existing repo has '%s'\n", r.URL, url)
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

func (r *Repo) checkRepo() error {
	args := []string{"-C", r.Path, "diff-index", "--quiet", "HEAD"}
	c := exec.Command("git", args...)
	err := c.Run()
	if err != nil { // fmt.Fprintf(os.Stderr)
		return fmt.Errorf("Uncommited changes in '%s'", r.Path)
	}
	return nil
}

func (r *Repo) pullRepo() error {
	args := []string{"-C", r.Path, "pull", r.remote, r.branch, "--quiet"}
	// // fmt.Printf("git %s\n", strings.Join(args, " "))
	// fmt.Printf("git pull %s %s\n", r.remote, r.branch)
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
	args := []string{"clone", r.URL, r.Path, "--recursive", "--quiet"}
	// // fmt.Printf("git %s\n", strings.Join(args, " "))
	// fmt.Printf("git clone %s %s\n", r.URL, r.Path)
	c := exec.Command("git", args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	if err := c.Run(); err != nil {
		return err
	}
	// actual := strings.TrimRight(out, "\n")
	return nil
}
