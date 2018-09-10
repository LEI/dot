package dot

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Copy task
type Copy struct {
	Task   `mapstructure:",squash"` // Action, If, OS
	Source string
	Target string
	Mode   os.FileMode
}

func (c *Copy) String() string {
	s := fmt.Sprintf("%s:%s", c.Source, c.Target)
	switch Action {
	case "install":
		if isRemote(c.Source) {
			s = fmt.Sprintf("curl -sSL %q -o %q", c.Source, tildify(c.Target))
		} else {
			s = fmt.Sprintf("cp %s %s", tildify(c.Source), tildify(c.Target))
		}
	case "remove":
		s = fmt.Sprintf("rm %s", tildify(c.Target))
	}
	return s
}

// RemoteString command
func (c *Copy) RemoteString() string {
	s := fmt.Sprintf("curl -sSL %q -o %s", c.Source, tildify(c.Target))
	if c.Mode != 0 {
		// s += fmt.Sprintf("\nchmod %o %q", c.Mode, c.Dest)
		s += fmt.Sprintf("; chmod %o $_", c.Mode)
	}
	return s
}

// func (h *Hook) buildCommandString() error {
// 	if h.Command != "" {
// 		return fmt.Errorf("%+v: invalid hook", h)
// 	}
// 	return nil
// }

// Status check task
func (c *Copy) Status() error {
	exists, err := copyExists(c.Source, c.Target)
	if err != nil {
		return err
	}
	if exists {
		return ErrExist
	}
	return nil
}

// Do task
func (c *Copy) Do() error {
	if err := c.Status(); err != nil {
		switch err {
		case ErrExist, ErrSkip:
			return nil
		default:
			return err
		}
	}
	if isRemote(c.Source) {
		if c.Mode == 0 {
			c.Mode = defaultFileMode
		}
		err := getURL(c.Source, c.Target, c.Mode)
		if err != nil {
			return err
		}
		return nil
	}
	return copyFile(c.Source, c.Target)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err = io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}

// Undo task
func (c *Copy) Undo() error {
	if err := c.Status(); err != nil {
		switch err {
		case ErrExist:
			// continue
		case ErrSkip:
			return nil
		default:
			return err
		}
	}
	return os.Remove(c.Target)
}

func isRemote(s string) bool {
	return strings.HasPrefix(s, "http://") ||
		strings.HasPrefix(s, "https://")
}

// copyExists returns true if the file source and target have the same content.
func copyExists(src, dst string) (bool, error) {
	if isRemote(src) {
		rf, err := newRemoteFile(src)
		// ok, err := remoteFileExists(src, dst)
		if err != nil {
			return false, err
		}
		if !exists(dst) {
			// Stop here if the target does not exist
			return false, nil
		}
		// if !ok {
		// 	return false, &url.Error{Op: "copy", URL: src, Err: ErrNotExist}
		// }
		return rf.Compare(dst)
	}
	if !exists(src) {
		// fmt.Errorf("%s: no such file to copy to %s", src, dst)
		return false, &os.PathError{Op: "copy", Path: src, Err: ErrNotExist}
	}
	if !exists(dst) {
		// Stop here if the target does not exist
		return false, nil
	}
	return fileCompare(src, dst)
}

// fileExists returns true if the name exists and is a not a directory.
func fileExists(name string) bool {
	f, err := os.Open(name)
	if err != nil {
		return false
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return !fi.IsDir()
}

// fileCompare TODO read in chunks
func fileCompare(p1, p2 string) (bool, error) {
	a, err := ioutil.ReadFile(p1)
	if err != nil {
		return false, err
	}
	b, err := ioutil.ReadFile(p2)
	if err != nil {
		return false, err
	}
	return bytes.Equal(a, b), nil
	// return bytes.Compare(a, b) == 0, nil
}

type remoteFile struct {
	URL string
	// Length of the file contents
	Length int64
	// Date last modified time
	Date time.Time
	// ETag or entity tag, is an opaque token that identifies a version of
	// the component served by a particular URL. The token can be anything
	// enclosed in quotes; often it's an md5 hash of the content, or the content's
	// VCS version number.
	Etag string
}

// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/If-Modified-Since
func (r *remoteFile) Compare(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()
	fi, err := f.Stat()
	// b, err := ioutil.ReadFile(name) // size := len(b)
	if err != nil {
		return false, err
	}
	if fi.Size() != r.Length || fi.ModTime().After(r.Date) {
		// TODO: confirm overwrite
		// fmt.Println("mismatch size", fi.Size() != r.Length, fi.Size(), r.Length)
		// fmt.Println("mismatch date", fi.ModTime().After(r.Date), fi.ModTime(), r.Date)
		return false, nil
	}
	// fmt.Println("Etag", r.Etag)
	// h := sha1.New()
	// if _, err = io.Copy(h, f); err != nil {
	// 	return false, err
	// }
	// hash := hex.EncodeToString(h.Sum(nil))
	// fmt.Println("SHA1", hash)
	return true, nil
}

func newRemoteFile(url string) (*remoteFile, error) {
	r := &remoteFile{URL: url}
	timeout := time.Duration(3 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	resp, err := client.Head(r.URL)
	if err != nil {
		return r, err
	}
	if r.Length, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64); err != nil {
		return r, err
	}
	r.Date, err = time.Parse(time.RFC1123, resp.Header.Get("Date"))
	if err != nil {
		return r, err
	}
	r.Etag = resp.Header.Get("Etag")
	if resp.StatusCode != 200 {
		return r, fmt.Errorf("%s: head request returned %+v", r.URL, resp.Status)
	}
	return r, nil
}

// remoteFileCompare TODO not implemented
// func remoteFileCompare(src, dst string) (bool, error) {
// 	return true, nil
// }
