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
	"github.com/LEI/dot/pkg/ostype"
	"github.com/LEI/dot/system"
)

var (
	sudo bool

	// Manager pkg
	Manager *Cmd

	managers = map[string]*Cmd{
		// https://wiki.alpinelinux.org/wiki/Alpine_Linux_package_management
		"apk": {
			Sudo: true,
			Bin:  "apk",
			Acts: map[string]string{
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
			Acts: map[string]string{
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
			Acts: map[string]string{
				"install": "install",
				"remove":  "uninstall",
			},
			Opts: []*Opt{
				{
					Args: []string{"--quiet"},
				},
			},
		},
		"cask": {
			Bin: "brew",
			Acts: map[string]string{
				"install": "install",
				"remove":  "uninstall",
			},
			Opts: []*Opt{
				{
					Args: []string{"cask"},
				},
			},
		},
		// https://www.archlinux.org/pacman/pacman.8.html
		"pacman": {
			Sudo: true,
			Bin:  "pacman",
			Acts: map[string]string{
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
			Acts: map[string]string{
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
			Acts: map[string]string{
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

// func init() {
// 	Detect()
// }

// Detect default package manager
func Detect() (c *Cmd) {
	switch runtime.GOOS {
	case "darwin":
		c = managers["brew"]
	case "linux":
		switch true {
		case executable("apk"):
			c = managers["apk"]
		case executable("apt-get"):
			c = managers["apt-get"]
		case executable("pacman"):
			c = managers["pacman"]
		case executable("yum"):
			c = managers["yum"]
		default:
			fmt.Fprintf(os.Stderr, "no package manager for OS: %s", ostype.List)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "no package manager for OS %s", runtime.GOOS)
		os.Exit(1)
	}
	// if c == nil {
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

func isRoot() bool {
	return os.Geteuid() == 0
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
	var c *Cmd
	if manager == "" {
		c = Detect()
	} else {
		var ok bool
		c, ok = managers[manager]
		if !ok {
			return fmt.Errorf(
				"%s: invalid package manager to %s %s",
				manager,
				action,
				name,
			)
		}
	}
	pkgs := strings.Split(name, " ")
	opts = append(pkgs, opts...)
	cmdArgs, err := c.Build(action, opts...)
	if err != nil {
		return err
	}
	bin := c.Bin
	// Switch binary for sudo
	if c.Sudo && bin != "sudo" && !isRoot() {
		cmdArgs = append([]string{bin}, cmdArgs...)
		bin = "sudo"
	}
	fmt.Printf("$ %s %s\n", bin, strings.Join(cmdArgs, " "))
	if system.DryRun {
		return nil
	}
	cmd := exec.Command(bin, cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// func NewCmd(t string)

// Cmd ...
type Cmd struct {
	Sudo bool
	Bin  string            // Package manager binary path
	Acts map[string]string // Command actions map
	Opts []*Opt            // General pkg manager options
	// ActOpts []*Opt         // Action options
	// types.HasOS `mapstructure:",squash"` // OS   map[string][]string // Platform options
	// types.HasIf `mapstructure:",squash"` // If   map[string][]string // Conditional opts
	Init func() error // Install or prepare bin
}

// Opt ...
type Opt struct {
	Args        interface{} // *parsers.Slice
	types.HasOS             // `mapstructure:",squash"`
	types.HasIf             // `mapstructure:",squash"`
}

// Add ...
func (cmd *Cmd) Add(opt *Opt) ([]string, error) {
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
func (cmd *Cmd) Build(a string, slice ...string) ([]string, error) {
	opts := []string{}

	// // General manager options
	// if len(cmd.Opts) == 0 && !ostype.Has("alpine") {
	// 	cmd.Opts = append(cmd.Opts, &Opt{Args: []string{"--noconfirm"}})
	// }

	// Package manager action
	action, ok := cmd.Acts[strings.ToLower(a)]
	if !ok {
		return []string{}, fmt.Errorf("unknown pkg action: %s", a)
	}
	opts = append(opts, action)

	// Action options
	for _, a := range cmd.Opts {
		add, err := cmd.Add(a)
		if err != nil {
			return opts, err
		}
		if len(add) > 0 {
			opts = append(opts, add...)
		}
	}

	// Insert package names and extra options
	opts = append(opts, slice...)

	// After action
	// for _, a := range cmd.ActOpts {
	// 	add, err := cmd.Add(a)
	// 	if err != nil {
	// 		return opts, err
	// 	}
	// 	if len(add) > 0 {
	// 		opts = append(opts, add...)
	// 	}
	// }
	return opts, nil
}
