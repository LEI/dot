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
func PacInstall(args ...string) error {
	return pacDo("install", args...)
}

// PacRemove ...
func PacRemove(args ...string) error {
	return pacDo("remove", args...)
}

func pacDo(a string, args ...string) error {
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
	pacArgs := []string{}
	switch strings.ToLower(a) {
	case "install":
		pacArgs = append(pacArgs, "-S")
	// case "remove":
	// 	pacArgs = append(pacArgs, "-R")
	default:
		fmt.Println("abort pacDo")
		return nil
	}
	// if HasOSType("darwin") {
	pacArgs = append(pacArgs, "--noconfirm")
	// }
	pacArgs = append(pacArgs, args...)
	if sudo {
		pacArgs = append([]string{pacBin}, pacArgs...)
		return execute("sudo", pacArgs...)
	}
	return execute(pacBin, pacArgs...)
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
