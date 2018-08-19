package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	// "text/template"

	// "github.com/LEI/dot/cli/config/tasks"

	"github.com/LEI/dot/internal/ostype"
)

var (
	// DryRun mode
	DryRun bool

	// Upgrade pacakges
	Upgrade bool

	// ErrExist ...
	ErrExist = fmt.Errorf("package already present")

	// ErrUnknown ...
	ErrUnknown = fmt.Errorf("unable to determine if package is present")

	// Manager pkg
	Manager *Pm

	managers = map[string]*Pm{
		"apk":     apk,
		"apt-get": aptGet,
		"brew":    brew,
		"cask":    brewCask,
		"pacman":  pacman,
		"yaourt":  yaourt,
		"yum":     yum,
	}
)

// Pm package manager
type Pm struct {
	Sudo        bool        // Prefix Bin with sudo if not root
	Bin         string      // Path to package manager binary
	Sub         []string    // Sub command and main options
	Install     interface{} // Install command name
	Remove      interface{} // Remove command name
	DryRun      []string    // Check mode, do not run if absent
	Opts        []string    // Common pkg manager options
	InstallOpts []string    // Install pkg manager options
	RemoveOpts  []string    // Remove pkg manager options
	// ActOpts []*Opt         // Action options
	// types.HasOS `mapstructure:",squash"` // OS   map[string][]string // Platform options
	// types.HasIf `mapstructure:",squash"` // If   map[string][]string // Conditional opts
	Env  map[string]string
	Init func() error               // Install or prepare bin
	Has  func(string) (bool, error) // Search install package
	done bool
}

// NewPm ...
func NewPm(name string) (*Pm, error) {
	m := &Pm{}
	if name == "" {
		m = Detect()
	} else {
		var ok bool
		m, ok = managers[name]
		if !ok {
			return m, fmt.Errorf("%s: invalid package manager name", name)
		}
	}
	return m, nil
}

// BuildOptions constructs the command arguments.
func (m *Pm) BuildOptions(a string, in ...string) ([]string, error) {
	opts := []string{}

	// // General manager options
	// if len(m.Opts) == 0 && !ostype.Has("alpine") {
	// 	m.Opts = append(m.Opts, &Opt{Args: []string{"--noconfirm"}})
	// }

	// Sub command and general options
	if len(m.Sub) > 0 {
		opts = append(opts, m.Sub...)
	}

	// Package manager action
	action, err := m.GetAction(a, in[0], in[1:]...)
	if err != nil {
		return opts, err
	}
	opts = append(opts, action)
	// Insert common and action specific options
	opts = append(opts, m.Opts...)
	switch a {
	case "install":
		opts = append(opts, m.InstallOpts...)
	case "remove":
		opts = append(opts, m.RemoveOpts...)
	}
	// Append package name and extra options
	opts = append(opts, in...)
	return opts, nil
}

// GetAction constructs the manager command for a given package.
func (m *Pm) GetAction(name, pkgName string, pkgOpts ...string) (string, error) {
	action := strings.ToLower(name)
	var i interface{}
	switch action {
	case "install":
		i = m.Install
	case "remove":
		i = m.Remove
	default:
		return action, fmt.Errorf("invalid pkg action: %s", action)
	}
	// act, ok := m.Acts[a]
	// if !ok {
	// 	return []string{}, fmt.Errorf("invalid pkg action: %s", a)
	// }
	switch a := i.(type) {
	case string:
		action = a
	// case []string:
	case func(*Pm, string, ...string) string:
		action = a(m, pkgName, pkgOpts...)
	default:
		return action, fmt.Errorf("%s: unknown pkg manager", a)
	}
	if action == "" {
		return action, fmt.Errorf("empty action %+v", m)
	}
	return action, nil
}

// func init() {
// 	Detect()
// }

// Detect default package manager
func Detect() (m *Pm) {
	switch runtime.GOOS {
	case "darwin":
		m = managers["brew"]
	case "linux":
		switch true {
		case executable("apk"):
			m = managers["apk"]
		case executable("apt-get"):
			m = managers["apt-get"]
		case executable("pacman"):
			m = managers["pacman"]
		case executable("yum"):
			m = managers["yum"]
		default:
			fmt.Fprintf(os.Stderr, "no package manager for OS: %s", ostype.List)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "no package manager for OS %s", runtime.GOOS)
		os.Exit(1)
	}
	// if m == nil {
	// 	os.Exit(128)
	// }
	return
}

func executable(bin string) bool {
	c := fmt.Sprintf("command -v %s", bin)
	cmd := exec.Command("sh", "-c", c)
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr
	// err := cmd.Run()
	out, err := cmd.CombinedOutput()
	// if err != nil {
	// 	return false
	// }
	return err == nil && len(out) > 0
}

// Has ...
func Has(manager, name string, opts ...string) (bool, error) {
	m, err := NewPm(manager)
	if err != nil {
		return false, err
	}
	if m.Has == nil {
		return false, ErrUnknown
	}
	return m.Has(name)
}

// Install ...
func Install(manager, name string, opts ...string) error {
	// fmt.Printf("%s %s\n", cmd.Bin, strings.Join(cmdArgs, " "))
	// stdout, stderr, status := ExecCommand(cmd.Bin, cmdArgs...)
	// str := strings.TrimRight(stdout, "\n")
	// // Quickfix centos yum
	// if status == 1 && stderr == "Error: Nothing to do\n" {
	// 	return str, nil
	// }
	// if status != 0 {
	// 	return str, fmt.Errorf(stderr)
	// }
	return execute(manager, "install", name, opts...)
}

// Remove ...
func Remove(manager, name string, opts ...string) error {
	return execute(manager, "remove", name, opts...)
}

// Exec ...
func execute(manager, action, name string, opts ...string) error {
	bin, opts, err := Init(manager, action, name, opts...)
	if err != nil {
		return err
	}
	// Get manager again to check package presence
	m, err := NewPm(manager)
	if err != nil {
		return err
	}
	for k, v := range m.Env {
		o := os.Getenv(k)
		if o != v {
			defer os.Setenv(k, o)
		}
		os.Setenv(k, v)
	}
	if action == "install" && m.Has != nil {
		ok, err := m.Has(name)
		if err != nil {
			return err
		}
		if ok {
			return ErrExist
		}
	}
	return execManagerCommand(m, bin, opts...)
}

func execManagerCommand(m *Pm, name string, args ...string) error {
	// fmt.Printf("$ %s %s\n", name, strings.Join(args, " "))
	if DryRun {
		if len(m.DryRun) == 0 {
			return nil
		}
		// Append check mode options and run
		args = append(args, m.DryRun...)
	}
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func execCommand(name string, args ...string) error {
	// fmt.Printf("$ %s %s\n", name, strings.Join(args, " "))
	if DryRun {
		return nil
	}
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Init can return ErrExist if the package is already installed
func Init(manager, action, name string, opts ...string) (string, []string, error) {
	m, err := NewPm(manager)
	if err != nil {
		return "", []string{}, err
	}
	// TODO forbid opts in name and/or multiple package names
	pkgs := strings.Split(name, " ")
	name = pkgs[0]
	opts = append(pkgs, opts...)
	opts, err = m.BuildOptions(action, opts...)
	if err != nil {
		return m.Bin, opts, err
	}
	bin, opts, err := getBin(m, opts)
	if err != nil {
		return bin, opts, err
	}
	if !m.done && m.Init != nil {
		if err := m.Init(); err != nil {
			return bin, opts, err
		}
		m.done = true
	}
	return bin, opts, nil
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
