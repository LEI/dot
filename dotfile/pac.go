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

func init() {
	if has(PACMAN) {
		// Arch Linux
		pacBin = PACMAN
	} else {
		// Unices
		pacBin = PACAPT
		downloadFromURL(PACAPTURL, PACAPT, 0755)
		// execute("sudo", "chmod", "+x", PACAPT)
	}
}

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
	pa := []string{} // pacapt args
	opts := ""
	if Verbose == 0 {
		opts+= "q"
	}
	switch strings.ToLower(a) {
	case "install":
		pa = append(pa, "-S" + opts)
		break
	// case "remove":
	// 	pa = append(pa, "-R" + opts)
	// 	break
	default:
		fmt.Println("abort pacDo")
		return nil
	}
	if ok := HasOne([]string{"darwin"}, GetOSTypes()); !ok {
		pa = append(pa, "--noconfirm")
	}
	pa = append(pa, args...)
	if sudo {
		pa = append([]string{pacBin}, pa...)
		return execute("sudo", pa...)
	}
	return execute(pacBin, pa...)
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
