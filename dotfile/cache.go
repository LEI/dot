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
	// CacheDir ...
	CacheDir string

	// ClearCache  ...
	ClearCache bool // = true

	// ErrCacheKeyNotFound ...
	ErrCacheKeyNotFound = fmt.Errorf("Cache entry not found")

	dotCache = &Cache{Map: map[string]string{}}
)

func init() {
	CacheDir = os.ExpandEnv("$HOME/.cache/dot")
	_, err := CreateDir(CacheDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create directory: %s", CacheDir)
		os.Exit(1)
	}
	if ClearCache {
		if err := dotCache.Clear(); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to clear cache: %s", CacheDir)
			os.Exit(1)
		}
	} else if _, err := dotCache.Read(); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read cache: %s", CacheDir)
		os.Exit(1)
	}
}

// Cache ...
type Cache struct {
	// FIXME: allow different cache types
	// (copies, templates...) TODO init
	// Dir string
	Map map[string]string
}

// New ...
func (c *Cache) New() *Cache {
	return &Cache{Map: map[string]string{}}
}

// Init ...
func (c *Cache) Init() *Cache {
	if c.Map == nil {
		*c = *c.New()
	}
	return c
}

// Get ...
func (c *Cache) Get(k string) (string, error) {
	c.Init()
	// key := cacheSerialize(k)
	v, ok := c.Map[k]
	if !ok {
		return v, ErrCacheKeyNotFound
	}
	return v, nil
}

// Validate ...
func (c *Cache) Validate(k, v string) (bool, error) {
	cached, err := c.Get(k)
	if err != nil && err != ErrCacheKeyNotFound {
		return false, err
	}
	// if cached == "" {
	// 	fmt.Println("VALIDATED EMPTY CACHE", k, cacheHashValue(v))
	// 	return true, nil
	// }
	if err == ErrCacheKeyNotFound {
		fmt.Println("CACHE KEY NOT FOUND", k)
		return true, nil
	}
	return cached == cacheHashValue(v), nil
}

// Put ...
func (c *Cache) Put(k, v string) error {
	c.Init()
	key := cacheSerialize(k)
	val := cacheHashValue(v)
	c.Map[key] = val // (*c)
	file := filepath.Join(CacheDir, key)
	// fmt.Println("Write cache", file, val)
	return ioutil.WriteFile(file, []byte(val), FileMode)
}

// Del ...
func (c *Cache) Del(k string) error {
	c.Init()

	nc := c.New()
	for key, val := range c.Map {
		if key != k {
			nc.Map[key] = val
		}
	}
	c = nc

	key := cacheSerialize(k)
	file := filepath.Join(CacheDir, key)
	return os.Remove(file)
}

// Read ...
func (c *Cache) Read() (map[string]string, error) {
	c.Init()
	list, err := c.List()
	if err != nil {
		return c.Map, err
	}
	for _, f := range list {
		b, err := ioutil.ReadFile(f)
		if err != nil {
			return c.Map, err
		}
		_, n := filepath.Split(f)
		k, err := cacheUnserialize(n)
		if err != nil {
			return c.Map, err
		}
		v := string(b)
		c.Map[k] = v // (*c)
	}
	// fmt.Println("Read cache", c)
	return c.Map, nil
}

// Write ...
func (c *Cache) Write() error {
	for k, v := range c.Map {
		err := c.Put(k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// List ...
func (c *Cache) List() ([]string, error) {
	p := filepath.Join(CacheDir, "*")
	return filepath.Glob(p)
}

// Clear ...
func (c *Cache) Clear() error {
	c = c.New()
	list, err := c.List()
	if err != nil {
		return err
	}
	for _, f := range list {
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}

func cacheSerialize(s string) string {
	return url.QueryEscape(s)
}

func cacheUnserialize(s string) (string, error) {
	s, err := url.QueryUnescape(s)
	if err != nil {
		return s, err
	}
	return s, err
}

// Create hash from file content
func cacheHashValue(s string) string {
	h := md5.Sum([]byte(s))
	return fmt.Sprintf("%x", string(h[:16]))
}
