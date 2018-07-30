package dotfile

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	// PACAPTURL pacapt download URL
	PACAPTURL = "https://github.com/icy/pacapt/raw/ng/pacapt"
	// PACAPT pacapt bin
	PACAPT = "/usr/local/bin/pacapt"
	// PACMAN pacman bin
	PACMAN = "pacman"
)

var (
	sudo   bool
)

func has(p string) bool {
	path, err := exec.LookPath(p)
	if err != nil {
		return false
	}
	return len(path) > 0
}

// PacInstall ...
func PacInstall(args ...string) (string, error) {
	// pacBin := pacBin()
	// stdout, _, status := ExecCommand(pacBin, append([]string{"-Qs"}, args...)...)
	// if status == 0 && len(stdout) > 0 {
	// 	fmt.Println("Error: Nothing to do")
	// }
	bin, opts := pac("install", args...)
	if bin == "" {
		return "", nil
	}
	fmt.Printf("%s %s\n", bin, strings.Join(opts, " "))
	stdout, stderr, status := ExecCommand(bin, opts...)
	// Quickfix centos yum
	if status == 1 && stderr == "Error: Nothing to do\n" {
		return stdout, nil
	}
	if status != 0 {
		return stdout, fmt.Errorf(stderr)
	}
	return stdout, nil
	// return "", execute(bin, opts...)
	// str, err := str, execute(bin, opts...)
	// if err != nil {
	// 	// pacapt -Syu
	// 	return str, err
	// }
	// return str, nil
}

// PacRemove ...
func PacRemove(args ...string) (string, error) {
	bin, opts := pac("remove", args...)
	if bin == "" {
		return "", nil
	}
	fmt.Printf("%s %s\n", bin, strings.Join(opts, " "))
	return "", nil // TODO execute(bin, opts...)
}

func pacBin() (pacBin string) {
	if has(PACMAN) {
		// Arch Linux
		pacBin = PACMAN
	} else {
		// Unices
		pacBin = PACAPT
		downloadFromURL(PACAPTURL, PACAPT, 0755)
		// execute("sudo", "chmod", "+x", PACAPT)
	}
	return
}

func pac(a string, args ...string) (string, []string) {
	pacBin := pacBin()
	pa := []string{}
	switch strings.ToLower(a) {
	case "install":
		pa = append(pa, "-S")
	// case "remove":
	// 	pa = append(pa, "-R")
	default:
		fmt.Println("abort pac", a)
		return "", args
	}
	// pacman
	if HasOSType("archlinux") {
		pa = append(pa, "--needed", "--noprogressbar")
	}
	// pacman, apt...
	if !HasOSType("darwin") {
		pa = append(pa, "--noconfirm")
	}
	pa = append(pa, args...)
	if Verbose == 0 {
		pa = append(pa, "--quiet")
	}
	if sudo {
		pa = append([]string{pacBin}, pa...)
		return "sudo", pa // execute("sudo", pa...)
	}
	return pacBin, pa // execute(pacBin, pa...)
}

func downloadFromURL(url, dst string, perm os.FileMode) {
	if dst == "" {
		tokens := strings.Split(url, "/")
		dst = tokens[len(tokens)-1]
	}

	fi, err := os.Stat(dst)
	if err != nil && os.IsExist(err) {
		log.Fatal(err)
	}
	if fi != nil && !os.IsNotExist(err) {
		fmt.Println("Already exists:", dst)
		return
	}

	fmt.Println("Downloading", url, "to", dst)

	output, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, perm)
	if err != nil {
		fmt.Println("Error while creating", dst, "-", err)
		return
	}
	defer output.Close()

	response, err := http.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}
	defer response.Body.Close()

	n, err := io.Copy(output, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return
	}

	fmt.Println(n, "bytes downloaded.")
}
