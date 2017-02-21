package dotfile

import (
	"fmt"
	"os"
)

type Link struct {
	File
	target string
	lstat os.FileInfo
}

func NewLink(src string, dst string) *Link {
	return &Link{
		File: File{path: src},
		target: dst,
	}
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
	fi, err := os.Lstat(l.path)
	l.lstat = fi
	return l.lstat, err
}

func (l *Link) IsLink() (bool, error) {
	fi, err := l.Lstat()
	if err != nil {
		return false, err
	}
	if IsSymlink(fi) {
		return true, nil
	}
	return false, nil
}
