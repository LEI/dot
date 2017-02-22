package dotfile

import (
	"fmt"
	"log"
	"os"
)

type Link struct {
	File
	target string
	lstat  os.FileInfo
}

func NewLink(src string, dst string) *Link {
	return &Link{
		File:   File{path: src},
		target: dst,
	}
}

func IsSymlink(fi os.FileInfo) bool {
	return fi != nil && fi.Mode()&os.ModeSymlink != 0
}

func (l *Link) String() string {
	return fmt.Sprintf("Link[%s][%s]", l.path, l.target)
}

func (l *Link) Target() string {
	return l.target
}

func (l *Link) SetTarget(dst string) {
	l.target = dst
}

func (l *Link) Info() (os.FileInfo, error) {
	if l.lstat != nil {
		return l.lstat, nil
	}
	fi, err := l.Lstat()
	l.lstat = fi
	return l.lstat, err
}

func (l *Link) Lstat() (os.FileInfo, error) {
	fi, err := os.Lstat(l.target)
	l.lstat = fi
	return l.lstat, err
}

func (l *Link) DestInfo() (os.FileInfo, error) {
	fi, err := os.Stat(l.target)
	return fi, err
}

func (l *Link) IsLink() bool {
	fi, err := l.Lstat()
	if err != nil && os.IsExist(err) {
		log.Fatal(err)
	}
	if IsSymlink(fi) {
		return true
	}
	return false
}

func (l *Link) IsLinked() (bool, error) {
	if !l.IsLink() {
		return false, nil
	}
	real, err := l.Readlink()
	if err != nil {
		return false, err
	}
	if real == l.path {
		return true, nil
	}
	return false, nil
}

func (l *Link) Readlink() (string, error) {
	path, err := os.Readlink(l.target)
	return path, err
}
