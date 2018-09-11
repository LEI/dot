package dot

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// https://docs.ansible.com/ansible/latest/modules/get_url_module.html
// https://github.com/LEI/dot-php/blob/master/install-composer.sh

// Download from a given URL
func getURL(url, dst string, perm os.FileMode) error {
	if dst == "" {
		tokens := strings.Split(url, "/")
		dst = tokens[len(tokens)-1]
	}

	fi, err := os.Stat(dst)
	if err != nil && os.IsExist(err) {
		return err
	}
	if fi != nil && !os.IsNotExist(err) {
		// fmt.Fprintf(os.Stderr, "%s: already exists!\n", fi.Name())
		// return nil
		return fmt.Errorf("%s: destination is already a file", dst)
	}
	fmt.Printf("Downloading %q...\n", url)
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
	fmt.Printf("%d bytes downloaded into %s\n", n, dst)
	return nil
}
