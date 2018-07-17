package dotfile

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	// "strings"
)

var (
	cacheDir string

	// cachePathSep = "=+"

	// osPathSep = string(os.PathSeparator)
)

func init() {
	cacheDir = os.ExpandEnv("$HOME/.cache/dot")
	_, err := CreateDir(cacheDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create directory: %s", cacheDir)
		os.Exit(1)
	}
}

// Cache ...
type Cache struct {
	Map map[string]string
}

// Get ...
func (c *Cache) Get(k string) string {
	// if c.Map == nil {
	// 	return ""
	// }
	key := cacheSerialize(k)
	// v, ok := c.Map[k]
	return c.Map[key]
}

// Validate ...
func (c *Cache) Validate(k, v string) error {
	cached := c.Get(k)
	if cached != "" && cached != cacheHashValue(v) {
		return fmt.Errorf("Mismatching cached file: %s", k)
	}
	return nil
}

// Put ...
func (c *Cache) Put(k, v string) error {
	if c.Map == nil {
		c.Map = make(map[string]string, 0)
	}
	key := cacheSerialize(k)
	val := cacheHashValue(v)
	c.Map[key] = val
	return nil
}

// WriteKey ...
func (c *Cache) WriteKey(k, v string) error {
	if err := c.Put(k, v); err != nil {
		return err
	}
	return c.Write()
}

// Write ...
func (c *Cache) Write() error {
	for k, v := range c.Map {
		f := filepath.Join(cacheDir, k)
		fmt.Println("Write cache", f, v)
		if err := ioutil.WriteFile(f, []byte(v), FileMode); err != nil {
			return err
		}
	}
	return nil
}

// Read ...
func (c *Cache) Read() error {
	// TODO
	return nil
}

func cacheSerialize(s string) string {
	// return strings.Replace(s, osPathSep, cachePathSep, -1)
	return url.QueryEscape(s) // url.QueryUnescape(f)
}

// Create hash from file content
func cacheHashValue(s string) string {
	h := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", string(h[:16]))
}
