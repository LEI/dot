package dot

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	// TODO: allow custom duration
	timeout = time.Duration(8 * time.Second)
)

// Copy task
type Copy struct {
	Task   `mapstructure:",squash"` // Action, If, OS
	Source string
	Target string
	Mode   os.FileMode
}

func (t *Copy) String() string {
	s := fmt.Sprintf("%s:%s", t.Source, t.Target)
	switch Action {
	case "install":
		if isRemote(t.Source) {
			s = fmt.Sprintf("curl -sSL %q -o %q", t.Source, tildify(t.Target))
		} else {
			s = fmt.Sprintf("cp %s %s", tildify(t.Source), tildify(t.Target))
		}
	case "remove":
		s = fmt.Sprintf("rm %s", tildify(t.Target))
	}
	return s
}

// Init task
func (t *Copy) Init() error {
	// ...
	return nil
}

// RemoteString command
func (t *Copy) RemoteString() string {
	s := fmt.Sprintf("curl -sSL %q -o %s", t.Source, tildify(t.Target))
	if t.Mode != 0 {
		// s += fmt.Sprintf("\nchmod %o %q", t.Mode, t.Dest)
		s += fmt.Sprintf("; chmod %o $_", t.Mode)
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
func (t *Copy) Status() error {
	exists, err := copyExists(t.Source, t.Target)
	if err != nil {
		return err
	}
	if exists {
		return ErrExist
	}
	return nil
}

// Do task
func (t *Copy) Do() error {
	if err := t.Status(); err != nil {
		switch err {
		case ErrExist, ErrSkip:
			return nil
		default:
			return err
		}
	}
	if isRemote(t.Source) {
		if t.Mode == 0 {
			t.Mode = defaultFileMode
		}
		err := getURL(t.Source, t.Target, t.Mode)
		if err != nil {
			return err
		}
		// tmpfile, err := ioutil.TempFile("", filepath.Basename(t.Source))
		// if err != nil {
		// 	return err
		// }
		// defer os.Remove(tmpfile.Name())
		// if _, err := tmpfile.Write(content); err != nil {
		// 	log.Fatal(err)
		// }
		// if err := tmpfile.Close(); err != nil {
		// 	log.Fatal(err)
		// }
		return nil
	}
	return copyFile(t.Source, t.Target)
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
func (t *Copy) Undo() error {
	if err := t.Status(); err != nil {
		switch err {
		case ErrExist:
			// continue
		case ErrSkip:
			return nil
		default:
			return err
		}
	}
	return os.Remove(t.Target)
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
		if rf.Length == -1 {
			// Ignore failed HEAD request
			return false, nil
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

// TODO: ask overwrite confirmation if remote file changed
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
	// if r.Length == -1 {
	// 	return false, fmt.Errorf("%s: no content length", r.URL)
	// }
	if fi.Size() != r.Length {
		// fmt.Println("DIFFERENT SIZE", fi.Size(), "->", r.Length)
		return false, nil
	}
	if fi.ModTime().After(r.Date) {
		// fmt.Println("DIFFERENT DATE", fi.ModTime(), "->", r.Date)
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
	client := http.Client{
		Timeout: timeout,
		// Fixes 403 forbidden for some github raw files
		// https://stackoverflow.com/a/42185713/7796750
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	//fmt.Println("> HEAD", r.URL)
	resp, err := client.Head(r.URL)
	if err != nil {
		return r, err
	}

	// req, err := http.NewRequest("HEAD", url, nil) // http.NoBody
	// if err != nil {
	// 	return r, err
	// }
	// /*
	// req.Header.Set("Accept", "*/*")
	// req.Header.Set("Content-Type", "text/plain; charset=utf-8")
	// req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; rv:31.0) Gecko/20100101 Firefox/31.0")
	// */
	// resp, err := client.Do(req)
	// if err != nil {
	// 	return r, err
	// }
	// defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return r, err
	}
	content := string(body)

	if resp.StatusCode >= 400 {
		return r, fmt.Errorf("HEAD %s %s (%+v)", r.URL, resp.Status, content)
	} else if resp.StatusCode != 200 && resp.StatusCode != 302 {
		fmt.Fprintf(os.Stderr, "HEAD %s %s (%+v)\n", r.URL, resp.Status, content)
	}
	r.Length = resp.ContentLength
	if resp.ContentLength == -1 && content != "" {
		fmt.Fprintf(os.Stderr, "HEAD %s %s: no content length but got: %+v\n", r.URL, resp.Status, content)
	}
	// // contentLen := resp.Header.Get("Content-Length")
	// if r.Length, err = strconv.ParseInt(contentLen, 10, 64); err != nil {
	// 	return r, err
	// }
	r.Date, err = time.Parse(time.RFC1123, resp.Header.Get("Date"))
	if err != nil {
		return r, err
	}
	r.Etag = resp.Header.Get("Etag")
	return r, nil
}

// remoteFileCompare TODO not implemented
// func remoteFileCompare(src, dst string) (bool, error) {
// 	return true, nil
// }
