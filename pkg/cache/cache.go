package cache

// https://www.opsdash.com/blog/persistent-key-value-store-golang.html
// https://github.com/dgraph-io/badger

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/LEI/dot/pkg/homedir"
	"github.com/rapidloop/skv"
)

var (
	// DirMode ...
	DirMode os.FileMode = 0755

	defaultCacheDir = ".cache"
	defaultFileExt  = "db"
	homeDir         string
	// store *Store // *svk.KVStore
)

// Store ...
type Store struct {
	dir string
	ext string
	*skv.KVStore
}

// New cache store
func New(dir string) (*Store, error) {
	// if dir == "" {
	// 	return fmt.Errorf("cache: missing directory name")
	// }
	s := &Store{
		dir: dir,
		ext: defaultFileExt,
	}
	err := s.Init()
	return s, err
}

// Init cache directory
func Init(path string) error {
	fi, err := os.Stat(path)
	if err != nil && os.IsExist(err) {
		return err
	}
	if fi != nil {
		return nil
	}
	fmt.Printf("Creating cache directory '%s'\n", path)
	if err := os.MkdirAll(path, DirMode); err != nil {
		return err
	}
	return nil
}

// BaseDir cache
func BaseDir() string {
	return filepath.Join(homedir.Get(), defaultCacheDir)
}

// Init cache store
func (s *Store) Init() error {
	return Init(s.Dir())
}

// Open inits a key value store
func (s *Store) Open(name string) error {
	if name == "" {
		return fmt.Errorf("cache: missing file name")
	}
	path := s.Path(name)
	// Open the store, e.g. "/var/lib/mywebapp/sessions.db"
	store, err := skv.Open(path)
	s.KVStore = store
	if err != nil {
		return err
	}
	return nil
}

// Path ...
func (s *Store) Path(name string) string {
	if s.ext != "" {
		name = fmt.Sprintf("%s.%s", name, s.ext)
	}
	return filepath.Join(s.Dir(), name)
}

// Dir cache store
func (s *Store) Dir() string {
	return filepath.Join(BaseDir(), s.dir)
}

// Save ...
func (s *Store) Save(dst string) error {
	b, err := ioutil.ReadFile(dst)
	if err != nil { // && os.IsExist(err) {
		return err
	}
	// c := string(b)
	// if err := dotCache.Put(dst, c); err != nil {
	// 	return true, err
	// }
	return s.Put(dst, b)
}

// // Forget ...
// func (s *Store) Forget(dst string) error {
// 	return s.Delete(dst)
// }

// Compare ...
func (s *Store) Compare(dst string) (bool, error) {
	var b []byte
	if err := s.Get(dst, &b); err != nil {
		return false, err
	}
	a, err := ioutil.ReadFile(dst)
	if err != nil { // && os.IsExist(err) {
		return false, err
	}
	ok := string(a) == string(b)
	// fmt.Println("COMPARING", dst, "->", ok)
	// fmt.Printf("<<<<\n%s====\n%s>>>>\n", string(a), string(b))
	return ok, nil
}
