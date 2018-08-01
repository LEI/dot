package dotfile

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/LEI/dot/parsers"
)

const (
	pacaptURL = "https://github.com/icy/pacapt/raw/ng/pacapt"

	pacaptBin = "/usr/local/bin/pacapt"
)

var (
	sudo bool

	pkgCommands = map[string]PkgCmd{
		"pacapt": {
			Bin: "pacapt",
			// Opts: []string{"--noconfirm"},
			Opts: []*PkgOpt{
				// FIXME prepend "alpine": {"--no-cache"},
				{
					Args: "--noconfirm",
				},
			},
			Acts: map[string]string{
				"install": "-S",
				"remove":  "-R",
			},
			Args: []*PkgOpt{
				{
					Args: []string{"--needed", "--noprogressbar"},
					OS:   []string{"archlinux"},
				},
				{
					Args: []string{"--no-install-suggests", "--no-install-recommends", "--quiet"},
					OS:   []string{"debian"},
				},
				{
					Args: []string{"--quiet"},
					If:   []string{"{{and (eq .Verbose 0) (hasOS \"!alpine\")}}"},
					// "{{eq .Verbose 0}}": {"--quiet"},
				},
			},
			Init: func() error {
				return downloadFromURL(pacaptURL, pacaptBin, 0755)
				// execute("sudo", "chmod", "+x", pacaptBin)
			},
		},
		"pacman": {
			Bin: "pacman",
			Opts: []*PkgOpt{
				{
					Args: []string{"--noconfirm", "--needed", "--noprogressbar"},
				},
			},
			Acts: map[string]string{
				"install": "-S",
				"remove":  "-R",
			},
			Args: []*PkgOpt{
				{
					Args: []string{"--quiet"},
					If:   []string{"{{eq .Verbose 0}}"},
				},
			},
		},
		"cask": {
			Bin: "brew",
			Opts: []*PkgOpt{
				{
					Args: []string{"cask"},
				},
			},
			Acts: map[string]string{
				"install": "install",
				"remove":  "uninstall",
			},
		},
	}
)

// PkgCmd ...
type PkgCmd struct {
	Bin  string            // Package manager binary path
	Opts []*PkgOpt         // General pkg manager options
	Acts map[string]string // Command actions map
	Args []*PkgOpt         // General pkg manager options
	// OS   map[string][]string // Platform options
	// If   map[string][]string // Conditional opts
	Init func() error // Install or prepare bin
	done bool
}

// Build command arguments
func (cmd *PkgCmd) Build(a string, slice ...string) ([]string, error) {
	args := []string{}

	// General manager options
	if len(cmd.Opts) == 0 && !HasOSType("alpine") {
		cmd.Opts = append(cmd.Opts, &PkgOpt{Args: []string{"--noconfirm"}})
	}

	// Before action
	for _, a := range cmd.Opts {
		add, err := cmd.Add(a)
		if err != nil {
			return args, err
		}
		if len(add) > 0 {
			args = append(args, add...)
		}
	}

	// Package manager action
	action, ok := cmd.Acts[strings.ToLower(a)]
	if !ok {
		return []string{}, fmt.Errorf("unknown pkg action: %s", a)
	}
	args = append(args, action)

	// Insert package names and extra options
	args = append(args, slice...)

	// After action
	for _, a := range cmd.Args {
		add, err := cmd.Add(a)
		if err != nil {
			return args, err
		}
		if len(add) > 0 {
			args = append(args, add...)
		}
	}

	// // Platform specific options
	// for p, opt := range cmd.OS {
	// 	patterns := strings.Split(p, ",")
	// 	for _, p := range patterns {
	// 		if HasOSType(p) {
	// 			args = append(args, opt...)
	// 		}
	// 	}
	// }
	// // Conditional options
	// pkgVarsMap := map[string]interface{}{
	// 	"DryRun":  DryRun,
	// 	"Verbose": Verbose,
	// 	"OS":      OS,
	// }
	// pkgFuncMap := template.FuncMap{
	// 	"hasOS": HasOSType,
	// }
	// for tpl, opt := range cmd.If {
	// 	str, err := TemplateData(cmd.Bin, tpl, pkgVarsMap, pkgFuncMap)
	// 	if err != nil {
	// 		return "", err
	// 	}
	// 	if str == "true" {
	// 		args = append(args, opt...)
	// 	}
	// }
	return args, nil
}

// Add ...
func (cmd *PkgCmd) Add(opt *PkgOpt) ([]string, error) {
	args := []string{}
	// Check platform
	if len(opt.OS) > 0 {
		hasOS := false
		for _, o := range opt.OS {
			if HasOSType(o) {
				hasOS = true
				break
			}
		}
		if !hasOS {
			// continue
			return args, nil
		}
	}
	// Check condition template
	pkgVarsMap := map[string]interface{}{
		"DryRun":  DryRun,
		"Verbose": Verbose,
		"OS":      OS,
	}
	pkgFuncMap := template.FuncMap{
		"hasOS": HasOSType,
	}
	for _, cond := range opt.If {
		str, err := TemplateData(cmd.Bin, cond, pkgVarsMap, pkgFuncMap)
		if err != nil {
			return args, err
		}
		if str != "true" {
			// continue
			return args, nil
		}
	}
	// for _, a := range *opt.Args {
	// 	args = append(args, a)
	// }
	optArgs, err := parsers.NewSlice(opt.Args)
	if err != nil {
		return args, err
	}
	args = append(args, *optArgs...)
	return args, err
}

// PkgOpt ...
type PkgOpt struct {
	Args interface{} // *parsers.Slice
	OS   []string
	// https://golang.org/pkg/text/template/#hdr-Functions
	If []string
}

// NewPkgCmd ...
func NewPkgCmd(t string, args ...string) (*PkgCmd, error) {
	cmd, ok := pkgCommands[t]
	if !ok {
		return &cmd, fmt.Errorf("unknown pkg type: %s", t)
	}
	bin := cmd.Bin
	if bin == "" {
		return &cmd, fmt.Errorf("missing pkg bin: %+v", t)
	}
	if !cmd.done && cmd.Init != nil {
		if err := cmd.Init(); err != nil {
			return &cmd, err
		}
		cmd.done = true
	}
	return &cmd, nil
}

// PkgTask struct
type PkgTask struct {
	Name string
	Opts []string
	Type string
	Sudo bool
	Task
}

// Status package
func (t *PkgTask) Status() bool {
	return true
}

// Do package task
func (t *PkgTask) Do(a string) (string, error) {
	return do(t, a)
}

// List package
func (t *PkgTask) List() (string, error) {
	str := fmt.Sprintf("Pkg: %s %+v", t.Name, t.Opts)
	// t.Sudo
	return str, nil
}

// Install package
func (t *PkgTask) Install() (string, error) {
	args := strings.Split(t.Name, " ")
	for _, a := range t.Opts {
		args = append(args, a) // "--"+
	}
	return t.Exec("install", args...) // PacInstall(args...)
}

// Remove package
func (t *PkgTask) Remove() (string, error) {
	args := strings.Split(t.Name, " ")
	return t.Exec("remove", args...) // PacRemove(args...)
}

// Exec ...
func (t *PkgTask) Exec(a string, args ...string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("no package specified")
	}
	if t.Type == "" {
		t.Type = "pacapt"
	}
	// switch t.Type {
	// case "pacman", "cask":
	// 	break
	// case "pac":
	// 	fallthrough
	// case "pacapt":
	// 	fallthrough
	// default:
	// 	t.Type = "pacapt"
	// }
	if t.Type == "pacapt" && has("pacman") {
		t.Type = "pacman"
	}
	cmd, err := NewPkgCmd(t.Type)
	if err != nil {
		return "", err
	}
	cmdArgs, err := cmd.Build(a, args...)
	if err != nil {
		return "", err
	}
	// Switch binary for sudo
	if t.Sudo {
		cmdArgs = append([]string{cmd.Bin}, cmdArgs...)
		cmd.Bin = "sudo"
	}
	fmt.Printf("%s %s\n", cmd.Bin, strings.Join(cmdArgs, " "))
	if DryRun {
		return "", nil
	}
	stdout, stderr, status := ExecCommand(cmd.Bin, cmdArgs...)
	str := strings.TrimRight(stdout, "\n")
	// Quickfix centos yum
	if status == 1 && stderr == "Error: Nothing to do\n" {
		return str, nil
	}
	if status != 0 {
		return str, fmt.Errorf(stderr)
	}
	return str, nil
}

func has(p string) bool {
	path, err := exec.LookPath(p)
	if err != nil {
		return false
	}
	return len(path) > 0
}

func downloadFromURL(url, dst string, perm os.FileMode) error {
	if dst == "" {
		tokens := strings.Split(url, "/")
		dst = tokens[len(tokens)-1]
	}

	fi, err := os.Stat(dst)
	if err != nil && os.IsExist(err) {
		return err
	}
	if fi != nil && !os.IsNotExist(err) {
		return nil // fmt.Errorf("already exists: %s", dst)
	}
	fmt.Println("Downloading", url, "to", dst)
	output, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		return fmt.Errorf("error while creating %s: %s", dst, err)
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error while downloading %s: %s", url, err)
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		return fmt.Errorf("error while copying to %s: %s", dst, err)
	}
	fmt.Println(n, "bytes downloaded.")
	return nil
}
