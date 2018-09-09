package git

// https://github.com/src-d/go-git
// https://github.com/graymeta/stow

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/LEI/dot/internal/shell"
)

var (
	min = &Version{2, 11}
)

// Version semantic
type Version struct {
	major int
	minor int
	// patch string
}

func (v *Version) String() string {
	return fmt.Sprintf("%d.%d", v.major, v.minor)
}

var (
	// // ErrDirtyRepo ...
	// ErrDirtyRepo = fmt.Errorf("dirty repository")

	// DryRun flag
	DryRun bool

	// Force ignores dirty repo
	Force bool

	// Quiet flag
	Quiet bool

	// Verbose level
	Verbose int

	// GitBin path
	GitBin = "git"

	// Stdout ...
	Stdout io.Writer = os.Stdout
	// Stderr ...
	Stderr io.Writer = os.Stderr

	cloneDepth    = 1
	defaultBranch = "master"
	defaultRemote = "origin"
)

func init() {
	// TODO check executable git before version
	if err := checkGitVersion(); err != nil {
		fmt.Fprintf(Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func checkGitVersion() error {
	buf, err := gitCombined("--version")
	if err != nil {
		return err
	}
	out := string(buf)
	out = strings.TrimPrefix(out, "git version ")
	out = strings.TrimSuffix(out, "\n")
	// fmt.Println("GIT_VERSION", out)
	// if out == "" {
	// 	return fmt.Errorf("%s: unable to parse git version", str)
	// }
	parts := strings.Split(out, ".")
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return err
	}
	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return err
	}
	if major < min.major || (major == min.major && minor < min.minor) {
		return fmt.Errorf("git version %s is required", min)
	}
	return nil
}

func git(args ...string) (string, string, error) {
	var stdout, stderr bytes.Buffer
	// _, err := shell.Exec(nil, &stdout, &stderr, GitBin, args...)
	cmd := exec.Command(GitBin, args...)
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout
	cmd.Stdin = os.Stdin
	if Verbose > 0 {
		fmt.Fprintln(Stdout, "exec:", GitBin, shell.FormatArgs(args))
	}
	err := cmd.Run()
	outstr := strings.TrimSuffix(stdout.String(), "\n")
	errstr := strings.TrimSuffix(stderr.String(), "\n")
	return outstr, errstr, err
}

func gitCombined(args ...string) (string, error) {
	cmd := exec.Command(GitBin, args...)
	if Verbose > 1 {
		fmt.Fprintln(Stdout, "exec:", GitBin, shell.FormatArgs(args))
	}
	buf, err := cmd.CombinedOutput()
	str := strings.TrimSuffix(string(buf), "\n")
	return str, err
}
