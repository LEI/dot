package git

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/LEI/dot/internal/cli"
)

var (
	// Scheme git, https or ssh
	Scheme = "https" // "git://"

	// Host git
	Host = "github.com"

	// User git
	User *url.Userinfo // = url.User("git")
)

// Repo ...
type Repo struct {
	Dir string
	URL *url.URL

	branch string // default: master
	remote string // default: origin
}

// NewRepo git
func NewRepo(u *url.URL, repo, dir string) (*Repo, error) {
	if u == nil {
		u = &url.URL{}
		// Scheme     string
		// Opaque     string    // encoded opaque data
		// User       *Userinfo // username and password information
		// Host       string    // host or host:port
		// Path       string    // path (relative paths may omit leading slash)
		// RawPath    string    // encoded path hint (see EscapedPath method)
		// ForceQuery bool      // append a query ('?') even if RawQuery is empty
		// RawQuery   string    // encoded query values, without '?'
		// Fragment   string    // fragment for references, without '#'
	}
	if dir == "" {
		return nil, fmt.Errorf("missing repo dir")
	}
	// strings.Contains(dir, "/")
	// && string(dir[0]) != "/"
	// && string(dir[0]) != "~"
	if repo == "" && !filepath.IsAbs(dir) {
		repo = dir // ParseURL(dir)
	}
	// format,
	// &URL{user, host, repo},
	// return &Remote{proto, user, host}
	r := &Repo{
		Dir:    dir,
		remote: defaultRemote,
		branch: defaultBranch,
	}
	// fmt.Println("URL parse repo", repo)
	// fmt.Printf("URL: %+v\n", u)
	repoURL, err := ParseURL(u, repo)
	if err != nil {
		return r, err
	}
	r.URL = repoURL
	if r.URL.String() == "" {
		return r, fmt.Errorf("missing repo url")
	}
	return r, nil
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
		return fmt.Errorf("uncommitted changes in %s:\n%s", r.Dir, stdout)
	}
	return nil
}

// Clone repo
func (r *Repo) Clone() error {
	args := []string{"clone", r.URL.String()}
	if r.Dir != "" {
		args = append(args, r.Dir)
	}
	if r.branch != "" {
		args = append(args, "--branch", r.branch)
	}
	if r.remote != "" {
		args = append(args, "--origin", r.remote)
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
	if DryRun {
		fmt.Printf("DRY-RUN: %s %s\n", GitBin, cli.FormatArgs(args))
		return nil
	}
	// status := r.ExecStatus(args...)
	// if status != 0 {
	//     return fmt.Errorf("git clone %s failed with exit code %d", r.URL, status)
	// }
	stdout, stderr, err := git(args...)
	if stderr != "" { // && Verbose > 0 {
		fmt.Fprintln(Stderr, stderr)
	}
	if stdout != "" { // && Verbose > 0 {
		fmt.Fprintln(Stdout, stdout)
	}
	if err != nil {
		return fmt.Errorf("unable to clone %s in %s: %s", r.URL, r.Dir, err)
	}
	return nil
}

// Pull repo
func (r *Repo) Pull() error {
	args := []string{"pull", r.remote, r.branch}
	if DryRun {
		args = append(args, "--dry-run")
	}
	if Quiet {
		args = append(args, "--quiet")
	}
	// if Verbose > 0 {
	// 	fmt.Println("git pull", r.remote, r.branch)
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

// ParseURL ...
func ParseURL(u *url.URL, repo string) (*url.URL, error) {
	if u.Scheme == "" && Scheme != "" {
		// fmt.Println("ParseURL", repo, "set Scheme", Scheme)
		u.Scheme = Scheme
	}
	if u.Host == "" && Host != "" { // u.Opaque == ""
		// fmt.Println("ParseURL", repo, "set Host", Host)
		u.Host = Host
	}
	if u.User != nil && u.User.String() == "" && User != nil {
		// fmt.Println("ParseURL", repo, "set User", User.String())
		u.User = User // url.User(username)
	}
	if repo != "" && !strings.HasSuffix(repo, ".git") {
		repo += ".git"
	}
	return u.Parse(repo)
}
