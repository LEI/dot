package dotfile

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"text/template"
)

// PkgTask struct
type PkgTask struct {
	Name string
	Opts []string
	Type string
	Sudo bool
	Task
}

// PkgType ...
type PkgType struct {
	Bin  string // Package manager binary path
	Opts []string // General pkg manager options
	Acts map[string]string // Command actions map
	OS   map[string][]string // Platform options
	If   map[string][]string // Conditional opts
	Init func() error // Install or prepare bin
	done bool
}

const (
	pacaptURL = "https://github.com/icy/pacapt/raw/ng/pacapt"

	pacaptBin = "/usr/local/bin/pacapt"
)

var (
	sudo bool

	pkgTypes = map[string]*PkgType{
		"pacapt": {
			Bin: "pacapt",
			// Opts: []string{"--noconfirm"},
			Acts: map[string]string{
				"install": "-S",
				"remove":  "-R",
			},
			OS: map[string][]string{
				// FIXME prepend "alpine": {"--no-cache"},
				"archlinux": {"--needed", "--noprogressbar"},
				"debian": {"--no-install-suggests", "--no-install-recommends", "--quiet"},
			},
			// https://golang.org/pkg/text/template/#hdr-Functions
			If: map[string][]string{
				// "{{eq .Verbose 0}}": {"--quiet"},
				"{{and (eq .Verbose 0) (hasOS \"!alpine\")}}": {"--quiet"},
			},
			Init: func() error {
				return downloadFromURL(pacaptURL, pacaptBin, 0755)
				// execute("sudo", "chmod", "+x", pacaptBin)
			},
		},
		"pacman": {
			Bin:  "pacman",
			Opts: []string{"--noconfirm", "--needed", "--noprogressbar"},
			Acts: map[string]string{
				"install": "-S",
				"remove":  "-R",
			},
			If: map[string][]string{
				"{{eq .Verbose 0}}": {"--quiet"},
			},
		},
		"cask": {
			Bin:  "brew",
			Opts: []string{"cask"},
			Acts: map[string]string{
				"install": "install",
				"remove":  "uninstall",
			},
		},
	}
)

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
	pt, ok := pkgTypes[t.Type]
	if !ok {
		return "", fmt.Errorf("unknown pkg type: %s", t.Type)
	}
	bin := pt.Bin
	if bin == "" {
		return "", fmt.Errorf("missing pkg bin: %+v", t)
	}
	if !pt.done && pt.Init != nil {
		if err := pt.Init(); err != nil {
			return "", err
		}
		pt.done = true
	}
	pacArgs := []string{}
	// General manager options
	if len(pt.Opts) == 0 && !HasOSType("alpine") {
		pt.Opts = append(pt.Opts, "--noconfirm")
	}
	pacArgs = append(pacArgs, pt.Opts...)
	action, ok := pt.Acts[strings.ToLower(a)]
	if !ok {
		return "", fmt.Errorf("unknown pkg action: %s", a)
	}
	pacArgs = append(pacArgs, action)
	// Insert package names and extra options
	pacArgs = append(pacArgs, args...)
	// Platform specific options
	for p, opt := range pt.OS {
		patterns := strings.Split(p, ",")
		for _, p := range patterns {
			if HasOSType(p) {
				pacArgs = append(pacArgs, opt...)
			}
		}
	}
	// Conditional options
	pkgVarsMap := map[string]interface{}{
		"DryRun":  DryRun,
		"Verbose": Verbose,
		"OS":      OS,
	}
	pkgFuncMap := template.FuncMap{
		"hasOS": HasOSType,
	}
	for tpl, opt := range pt.If {
		str, err := TemplateData(pt.Bin, tpl, pkgVarsMap, pkgFuncMap)
		if err != nil {
			return "", err
		}
		if str == "true" {
			pacArgs = append(pacArgs, opt...)
		}
	}
	// Switch binary for sudo
	if t.Sudo {
		pacArgs = append([]string{bin}, pacArgs...)
		bin = "sudo"
	}
	fmt.Printf("%s %s\n", bin, strings.Join(pacArgs, " "))
	if DryRun {
		return "", nil
	}
	stdout, stderr, status := ExecCommand(bin, pacArgs...)
	// Quickfix centos yum
	if status == 1 && stderr == "Error: Nothing to do\n" {
		return stdout, nil
	}
	if status != 0 {
		return stdout, fmt.Errorf(stderr)
	}
	return stdout, nil
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
