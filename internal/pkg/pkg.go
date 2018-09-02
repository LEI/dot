package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	// "text/template"

	// "github.com/LEI/dot/cli/config/tasks"

	"github.com/LEI/dot/internal/host"
	"github.com/LEI/dot/internal/shell"
)

var (
	// DryRun mode
	DryRun bool

	// Upgrade pacakges
	Upgrade bool

	// ErrExist ...
	ErrExist = fmt.Errorf("package already present")

	// Manager pkg
	Manager *Pm

	managers = map[string]*Pm{
		"apk":     apk,      // alpine
		"apt-cyg": aptCyg,   // cygwin
		"apt-get": aptGet,   // debian
		"brew":    brew,     // homebrew
		"cask":    brewCask, // darwin
		"choco":   choco,    // windows
		// "cpan":    cpan,     // perl
		"gem":    gem,    // ruby
		"npm":    npm,    // node
		"pacman": pacman, // archlinux
		"pip":    pip,    // python
		"pip2":   pip2,
		"pip3":   pip3,
		"termux": termux, // android pkg
		"yaourt": yaourt, // arch user repository
		"yum":    yum,    // centos
	}
)

type hasFunc func(*Pm, []string) (bool, error)

// Pm package manager
type Pm struct {
	Shell       string
	Sudo        bool        // Prefix Bin with sudo if not root
	Bin         string      // Path to package manager binary
	Sub         []string    // Sub command and main options
	Install     interface{} // Install command name
	Remove      interface{} // Remove command name
	Opts        []string    // Common pkg manager options
	InstallOpts []string    // Install pkg manager options
	RemoveOpts  []string    // Remove pkg manager options
	DryRunOpts  []string    // Check mode, do not run if absent
	// ActOpts []*Opt         // Action options
	// types.HasOS `mapstructure:",squash"` // OS   map[string][]string // Platform options
	// types.HasIf `mapstructure:",squash"` // If   map[string][]string // Conditional opts
	Env  map[string]string // Execution environment variables
	Init func(*Pm) error   // Install manager or prepare bin
	Has  hasFunc           // Search local packages
	done bool              // Initialized (TODO: use sync.Once instead?)
}

// NewPm ...
func NewPm(name string) (m *Pm, err error) {
	// m := &Pm{}
	if name == "" {
		name, err = Detect()
		if err != nil {
			return m, err
		}
	}
	var ok bool
	m, ok = managers[name]
	if !ok {
		return m, fmt.Errorf("%s: invalid package manager name", name)
	}
	if m == nil {
		return m, fmt.Errorf("unable to detect package manager %s", name)
	}
	return m, nil
}

// Build can return ErrExist if the package is already installed
func (m *Pm) Build(action string, pkgs []string, opts ...string) (string, []string, error) {
	// input := strings.Fields(name)
	// if len(input) == 0 { ... }
	//opts = append(pkgs, opts...)
	opts, err := m.BuildOptions(action, pkgs, opts...)
	if err != nil {
		return m.Bin, opts, err
	}
	bin, args, err := getBin(m, opts)
	if err != nil {
		return bin, args, err
	}
	// m.Init(m)
	return bin, args, nil
}

// BuildOptions constructs the command arguments.
func (m *Pm) BuildOptions(a string, pkgs []string, opts ...string) ([]string, error) {
	a = strings.ToLower(a)
	s := []string{}

	// // General manager options
	// if len(m.Opts) == 0 && !host.HasOS("alpine") {
	// 	m.Opts = append(m.Opts, &Opt{Args: []string{"--noconfirm"}})
	// }

	// Sub command and general options
	if len(m.Sub) > 0 {
		s = append(s, m.Sub...)
	}
	// Package manager action
	action, err := m.GetAction(a, pkgs...)
	if err != nil {
		return s, err
	}
	s = append(s, action)
	// Insert common and action specific options
	s = append(s, m.Opts...)
	switch a {
	case "install":
		s = append(s, m.InstallOpts...)
	case "remove":
		s = append(s, m.RemoveOpts...)
	}
	// Append package name and extra options
	s = append(s, pkgs...)
	s = append(s, opts...)
	return s, nil
}

// GetAction constructs the manager command for a given package.
func (m *Pm) GetAction(name string, input ...string) (string, error) {
	var i interface{}
	switch name {
	case "sync":
		fallthrough
	case "install":
		i = m.Install
	case "remove":
		i = m.Remove
	default:
		// Unknown action
		return name, fmt.Errorf("invalid pkg action: %s", name)
	}
	// act, ok := m.Acts[a]
	// if !ok {
	// 	return []string{}, fmt.Errorf("invalid pkg action: %s", a)
	// }
	switch a := i.(type) {
	case string:
		name = a
	// case []string:
	case func(*Pm, ...string) string:
		pkgs := pkgOnly(input)
		if len(pkgs) == 0 {
			return name, fmt.Errorf("empty name %+v", input)
		}
		name = a(m, pkgs...)
	default:
		return name, fmt.Errorf("%s: unknown pkg manager", a)
	}
	if name == "" {
		return name, fmt.Errorf("empty action for package manager %+v", m)
	}
	return name, nil
}

func pkgOnly(input []string) []string {
	pkgs := []string{}
	for _, s := range input {
		if !strings.HasPrefix(s, "-") {
			pkgs = append(pkgs, s)
		}
	}
	return pkgs
}

// func init() {
// 	Detect()
// }

// Detect default package manager
func Detect() (name string, err error) {
	switch runtime.GOOS {
	case "android": // executable("pkg"):
		name = "termux"
	case "darwin":
		name = "brew"
	case "linux":
		switch {
		case executable("apk"):
			name = "apk"
		case executable("apt-get"):
			name = "apt-get"
		case executable("pacman"):
			name = "pacman"
		case executable("yum"):
			name = "yum"
		}
	case "windows": // executable("choco"):
		name = "apt-cyg" // "choco"
	}
	if name == "" {
		return "", fmt.Errorf("no package manager for OS %s (%s)", runtime.GOOS, host.OSTypes)
	}
	return
}

func executable(name string) bool {
	// c := fmt.Sprintf("command -v %s", name)
	// cmd := exec.Command(shell.Get(), "-c", c)
	// out, err := cmd.CombinedOutput()
	// return err == nil && len(out) > 0
	out, err := exec.LookPath(name)
	return err == nil && len(out) > 0
}

// Has ...
func Has(manager string, pkgs []string, opts ...string) (bool, error) {
	m, err := NewPm(manager)
	if err != nil {
		return false, err
	}
	if m == nil {
		return false, fmt.Errorf(manager, "no pkg manager", manager)
	}
	if m.Has == nil {
		return false, nil
		// ErrUnknown = fmt.Errorf("unable to determine if package is present")
	}
	return m.Has(m, pkgs)
}

// Install ...
func Install(manager string, pkgs []string, opts ...string) error {
	// fmt.Printf("%s %s\n", cmd.Bin, shell.FormatArgs(cmdArgs))
	// stdout, stderr, status := ExecCommand(cmd.Bin, cmdArgs...)
	// str := strings.TrimRight(stdout, "\n")
	// // Quickfix centos yum
	// if status == 1 && stderr == "Error: Nothing to do\n" {
	// 	return str, nil
	// }
	// if status != 0 {
	// 	return str, fmt.Errorf(stderr)
	// }
	return execute(manager, "install", pkgs, opts...)
}

// Remove ...
func Remove(manager string, pkgs []string, opts ...string) error {
	return execute(manager, "remove", pkgs, opts...)
}

// Exec ...
func execute(manager, action string, pkgs []string, opts ...string) error {
	m, err := NewPm(manager)
	if err != nil {
		return err
	}
	if m == nil {
		return fmt.Errorf(manager, "no pkg manager", manager)
	}
	bin, args, err := m.Build(action, pkgs, opts...)
	if err != nil {
		return err
	}
	// First run initialisation
	// TODO: before execute
	if !m.done && m.Init != nil {
		if err := m.Init(m); err != nil {
			return err
		}
		m.done = true
	}
	if err := m.checkIfExist(pkgs); err != nil {
		return err
	}
	return execManagerCommand(m, bin, args...)
}

// TODO: return *exec.Cmd
func execManagerCommand(m *Pm, bin string, args ...string) error {
	if DryRun {
		if len(m.DryRunOpts) == 0 {
			return nil
		}
		// Append check mode options and run
		args = append(args, m.DryRunOpts...)
	}
	var cmd *exec.Cmd
	if m.Shell != "" {
		s := shell.Get()
		// if Verbose { fmt.Println("Using shell:", s) }
		c := fmt.Sprintf("%s %s", bin, shell.FormatArgs(args))
		cmd = exec.Command(s, "-c", c)
	} else {
		cmd = exec.Command(bin, args...)
	}
	return execWithEnv(m.Env, cmd)
}

// Run a command with custom env vars, even if DryRun is enabled.
func execWithEnv(env map[string]string, cmd *exec.Cmd) error {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Env = os.Environ()
	for k, v := range env {
		// o := os.Getenv(k)
		// if o == v {
		// 	continue
		// }
		// defer os.Setenv(k, o)
		// if err := os.Setenv(k, v); err != nil {
		// 	return err
		// }

		// e := fmt.Sprintf("%s=%s", k, v)
		// fmt.Println(e)
		cmd.Env = append(cmd.Env, k+"="+v)
	}
	// if err := cmd.Run(); err != nil {
	// 	if !m.AllowFailure {
	// 		return err
	// 	}
	// 	fmt.Fprintf(os.Stderr, "$ %s %s: %s", bin, shell.FormatArgs(args), err) // return err
	// }
	return cmd.Run()
}

func execCommand(name string, args ...string) error {
	fmt.Printf("$ %s %s\n", name, shell.FormatArgs(args))
	if DryRun {
		return nil
	}
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

// Return ErrExist if packages are present (installed)
func (m *Pm) checkIfExist(pkgs []string) error {
	// Abort if not check function is available
	if m.Has != nil {
		return nil
	}
	// // Ignore presence check if install is a func
	// if _, ok := m.Install.(hasFunc); ok {
	// 	return nil
	// }
	ok, err := m.Has(m, pkgs)
	if err != nil {
		return err
	}
	if ok {
		return ErrExist
	}
	return nil
}

func getBin(m *Pm, opts []string) (string, []string, error) {
	bin := m.Bin
	// Switch binary for sudo
	if m.Sudo && bin != "sudo" && !isRoot() {
		opts = append([]string{bin}, opts...)
		bin = "sudo"
	}
	return bin, opts, nil
}

func isRoot() bool {
	return os.Geteuid() == 0
}
