package dotfile

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
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
	Opts []string // General manager options
	Acts map[string]string // Action map
	OS   map[string][]string // Platform options
	If   map[string][]string // Conditional opts
	Init func() error // Install or prepare bin
	init bool
}

const (
	// PACAPTURL pacapt download URL
	PACAPTURL = "https://github.com/icy/pacapt/raw/ng/pacapt"

	// PACAPT bin
	PACAPT = "/usr/local/bin/pacapt"

	// PACMAN bin
	PACMAN = "pacman"
)

var (
	sudo bool

	pkgTypes = map[string]*PkgType{
		"pacapt": {
			Bin: PACAPT,
			Opts: []string{"--noconfirm"},
			Acts: map[string]string{
				"install": "-S",
				"remove":  "-R",
			},
			OS: map[string][]string{
				"archlinux": {"--needed", "--noprogressbar"},
				"debian": {"--no-install-suggests", "--no-install-recommends", "--quiet"},
			},
			If: map[string][]string{
				"eq .Verbose 0": {"--quiet"},
			},
			Init: func() error {
				return downloadFromURL(PACAPTURL, PACAPT, 0755)
				// execute("sudo", "chmod", "+x", PACAPT)
			},
		},
		"pacman": {
			Bin:  PACMAN,
			Opts: []string{"--noconfirm", "--needed", "--noprogressbar"},
			Acts: map[string]string{
				"install": "-S",
				"remove":  "-R",
			},
			If: map[string][]string{
				"eq .Verbose 0": {"--quiet"},
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
	if t.Type == "pacapt" && has(PACMAN) {
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
	if !pt.init && pt.Init != nil {
		if err := pt.Init(); err != nil {
			return "", err
		}
		pt.init = true
	}
	// General manager options
	opts := pt.Opts
	// Action (install, remove)
	action, ok := pt.Acts[strings.ToLower(a)]
	if !ok {
		return "", fmt.Errorf("unknown pkg action: %s", a)
	}
	opts = append(opts, action)
	// Packages and options
	opts = append(opts, args...)

	args = opts

	// Platform specific options
	for o, opts := range pt.OS {
		patterns := strings.Split(o, ",")
		for _, p := range patterns {
			if HasOSType(p) {
				args = append(args, opts...)
			}
		}
	}
	// Conditional options
	pkgTplMap := map[string]interface{}{
		"DryRun":  DryRun,
		"Verbose": Verbose,
		"OS":      OS,
	}
	for cond, opts := range pt.If {
		str, err := TemplateData(args[0], "{{"+cond+"}}", pkgTplMap)
		if err != nil {
			return "", err
		}
		if str == "true" {
			args = append(args, opts...)
		}
	}
	// Switch binary for sudo
	if t.Sudo {
		args = append([]string{bin}, args...)
		bin = "sudo"
	}
	fmt.Printf("%s %s\n", bin, strings.Join(args, " "))
	if DryRun {
		return "", nil
	}
	stdout, stderr, status := ExecCommand(bin, args...)
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
