package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	// "text/template"

	// "github.com/LEI/dot/cli/config/tasks"
	"github.com/LEI/dot/cli/config/types"
	"github.com/LEI/dot/internal/ostype"
	"github.com/LEI/dot/system"
)

var (
	// Upgrade pacakges
	Upgrade bool

	// Manager pkg
	Manager *Pm

	managers = map[string]*Pm{
		// https://wiki.alpinelinux.org/wiki/Alpine_Linux_package_management
		"apk": {
			Sudo: true,
			Bin:  "apk",
			Acts: map[string]interface{}{
				"install": "add",
				"remove":  "del",
			},
			Opts: []*Opt{
				{
					Args: []string{
						"--no-cache",
						"--quiet",
						"--update",
					},
				},
			},
		},
		// https://manpages.debian.org/stretch/apt/apt-get.8.en.html
		"apt-get": {
			Sudo: true,
			Bin:  "apt-get",
			Acts: map[string]interface{}{
				"install": "install",
				"remove":  "remove",
			},
			Opts: []*Opt{
				{
					// Args: []string{"-qqy"},
					Args: []string{
						"--assume-yes",
						"--no-install-recommends",
						"--no-install-suggests",
						"--quiet",
						"--quiet",
					},
				},
			},
		},
		// https://docs.brew.sh/Manpage
		"brew": {
			Bin: "brew",
			Acts: map[string]interface{}{
				"install": func(m *Pm, in []string) string {
					// TODO filter strings.HasPrefix(in, "-")?
					opts := append([]string{"ls", "--versions"}, in...)
					err := exec.Command("brew", opts...).Run()
					if err == nil && Upgrade {
						return "upgrade"
					}
					return "install"
				},
				"remove": "uninstall",
			},
			Opts: []*Opt{
				{
					Args: []string{"--quiet"},
				},
			},
			Env: map[string]string{
				// "HOMEBREW_NO_ANALYTICS": "1",
				"HOMEBREW_NO_AUTO_UPDATE": "1",
				// "HOMEBREW_NO_EMOJI": "1",
			},
			Init: func() error {
				return execute("brew", "update", "--quiet")
			},
		},
		"cask": {
			Bin: "brew",
			Sub: []string{"cask"},
			Acts: map[string]interface{}{
				"install": "install",
				"remove":  "uninstall",
			},
		},
		// https://www.archlinux.org/pacman/pacman.8.html
		"pacman": {
			Sudo: true,
			Bin:  "pacman",
			Acts: map[string]interface{}{
				"install": "--sync",   // -S
				"remove":  "--remove", // -R
			},
			Opts: []*Opt{
				{
					Args: []string{
						"--needed",
						"--noconfirm",
						"--noprogressbar",
						"--quiet",
						"--refresh",    // -y
						"--sysupgrade", // -u
					},
				},
				// {
				// 	Args: []string{"--quiet"},
				// 	// If:   []string{"{{eq .Verbose 0}}"},
				// 	HasIf: types.HasIf{If: []string{"{{eq .Verbose 0}}"}},
				// },
			},
		},
		// https://archlinux.fr/man/yaourt.8.html
		"yaourt": {
			// Sudo: false,
			Bin: "yaourt",
			Acts: map[string]interface{}{
				"install": "--sync",   // -S
				"remove":  "--remove", // -R
			},
			Opts: []*Opt{
				{
					Args: []string{
						"--noconfirm",
						// "--sysupgrade", // -u
					},
				},
			},
		},
		"yum": {
			Sudo: true,
			Bin:  "yum",
			Acts: map[string]interface{}{
				"install": "install",
				"remove":  "remove",
			},
			Opts: []*Opt{
				{
					Args: []string{
						"--assumeyes",
						// "--error=0",
						"--quiet",
					},
				},
			},
		},
	}
)

// Pm ...
type Pm struct {
	Sudo bool
	Bin  string                 // Package manager binary path
	Sub  []string               // Sub commands
	Acts map[string]interface{} // Command actions map
	Opts []*Opt                 // General pkg manager options
	// ActOpts []*Opt         // Action options
	// types.HasOS `mapstructure:",squash"` // OS   map[string][]string // Platform options
	// types.HasIf `mapstructure:",squash"` // If   map[string][]string // Conditional opts
	Env  map[string]string
	Init func() error // Install or prepare bin
	done bool
}

// Opt ...
type Opt struct {
	Args        interface{} // *parsers.Slice
	types.HasOS             // `mapstructure:",squash"`
	types.HasIf             // `mapstructure:",squash"`
}

// NewMngr ...
func NewMngr(name string) (*Pm, error) {
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

// Add ...
func (m *Pm) Add(opt *Opt) ([]string, error) {
	args := []string{}
	// Check platform
	if !opt.CheckOS() {
		// continue
		return args, nil
	}
	if !opt.CheckIf() {
		// continue
		return args, nil
	}
	// Check condition template
	// pkgVarsMap := map[string]interface{}{
	// 	// "DryRun":  DryRun,
	// 	// "Verbose": Verbose,
	// 	// "OS":      OS,
	// }
	// pkgFuncMap := template.FuncMap{
	// 	// "hasOS": HasOSType,
	// }
	// for _, cond := range opt.If {
	// 	str, err := TemplateData(cmd.Bin, cond, pkgVarsMap, pkgFuncMap)
	// 	if err != nil {
	// 		return args, err
	// 	}
	// 	stdout, stderr, status := ExecCommand(Shell, "-c", str)
	// 	stdout = strings.TrimRight(stdout, "\n")
	// 	stderr = strings.TrimRight(stderr, "\n")
	// 	// if stdout != "" {
	// 	// 	fmt.Printf("stdout: %s\n", stdout)
	// 	// }
	// 	if stderr != "" {
	// 		fmt.Fprintf(os.Stderr, "stderr: %s\n", stderr)
	// 	}
	// 	if status != 0 {
	// 		// continue
	// 		return args, nil
	// 	}
	// }
	// optArgs, err := tasks.NewSlice(opt.Args)
	// if err != nil {
	// 	return args, err
	// }
	// args = append(args, *optArgs...)
	switch o := opt.Args.(type) {
	case string:
		args = append(args, o)
	case []string:
		args = append(args, o...)
	default:
		return args, fmt.Errorf("todo: opt args %+v", opt)
	}
	// return args, err
	return args, nil
}

// Build command arguments
func (m *Pm) Build(a string, in ...string) ([]string, error) {
	opts := []string{}

	// // General manager options
	// if len(m.Opts) == 0 && !ostype.Has("alpine") {
	// 	m.Opts = append(m.Opts, &Opt{Args: []string{"--noconfirm"}})
	// }

	// Sub command
	if len(m.Sub) > 0 {
		opts = append(opts, m.Sub...)
	}

	// Package manager action
	act, ok := m.Acts[strings.ToLower(a)]
	if !ok {
		return []string{}, fmt.Errorf("unknown pkg action: %s", a)
	}
	var action string
	switch A := act.(type) {
	case string:
		action = A
	// case []string:
	case func(m *Pm, in []string) string:
		action = A(m, in)
	default:
		return opts, fmt.Errorf("%s: unknown pkg manager", A)
	}
	if action == "" {
		return opts, fmt.Errorf("empty action %+v", m)
	}
	opts = append(opts, action)

	// Action options
	for _, a := range m.Opts {
		add, err := m.Add(a)
		if err != nil {
			return opts, err
		}
		if len(add) > 0 {
			opts = append(opts, add...)
		}
	}

	// Insert package names and extra options
	opts = append(opts, in...)

	// After action
	// for _, a := range m.ActOpts {
	// 	add, err := m.Add(a)
	// 	if err != nil {
	// 		return opts, err
	// 	}
	// 	if len(add) > 0 {
	// 		opts = append(opts, add...)
	// 	}
	// }
	return opts, nil
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
	return Exec("install", manager, name, opts...)
}

// Remove ...
func Remove(manager, name string, opts ...string) error {
	return Exec("remove", manager, name, opts...)
}

// Exec ...
func Exec(action, manager, name string, opts ...string) error {
	m, err := NewMngr(manager)
	if err != nil {
		return err
	}
	pkgs := strings.Split(name, " ")
	opts = append(pkgs, opts...)
	opts, err = m.Build(action, opts...)
	if err != nil {
		return err
	}
	bin, opts, err := getBin(m, opts)
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
	if !m.done && m.Init != nil {
		if err := m.Init(); err != nil {
			return err
		}
		m.done = true
	}
	return execute(bin, opts...)
}

func execute(name string, args ...string) error {
	fmt.Printf("$ %s %s\n", name, strings.Join(args, " "))
	if system.DryRun {
		return nil
	}
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
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
