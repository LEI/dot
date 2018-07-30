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
	pacBin string
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
func PacInstall(slice ...string) (string, error) {
	bin, args := pac("install", slice...)
	if bin == "" {
		return "", nil
	}
	str := fmt.Sprintf("%s %s", bin, strings.Join(args, " "))
	return str, execute(bin, args...)
	// str, err := str, execute(bin, args...)
	// if err != nil {
	// 	// pacapt -Syu
	// 	return str, err
	// }
	// return str, nil
}

// PacRemove ...
func PacRemove(slice ...string) (string, error) {
	bin, args := pac("remove", slice...)
	if bin == "" {
		return "", nil
	}
	return fmt.Sprintf("%s %s", bin, args), execute(bin, args...)
}

func pac(a string, args ...string) (string, []string) {
	// Init
	if has(PACMAN) {
		// Arch Linux
		pacBin = PACMAN
	} else {
		// Unices
		pacBin = PACAPT
		downloadFromURL(PACAPTURL, PACAPT, 0755)
		// execute("sudo", "chmod", "+x", PACAPT)
	}
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
