package role

import (
	"fmt"
	"github.com/LEI/dot/fileutil"
	"path/filepath"
	"os"
)

type Line struct {
	File string
	Line string
}

func (l *Line) String() string {
	return fmt.Sprintf("%s`%s`", l.File, l.Line)
}

func (l *Line) InFile(target string) error {
	l.File = os.ExpandEnv(l.File)
	l.File = filepath.Join(target, l.File)
	err := fileutil.LineInFile(l.File, l.Line)
	if err != nil {
		return err
	}
	return nil
}

func (r *Role) Lines() []*Line {
	p := r.Package
	if p == nil {
		p = &Package{}
	}
	r.Config.UnmarshalKey("line", &p.Line)
	r.Config.UnmarshalKey("lines", &p.Lines)
	if p.Line != nil {
		p.Lines = append(p.Lines, p.Line) // .(map[string]interface{})
		p.Line = nil
	}
	r.Config.Set("lines", p.Lines)
	r.Config.Set("line", p.Line)
	return p.Lines
}
