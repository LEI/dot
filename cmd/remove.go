package cmd

import (
	// "fmt"
	// "github.com/LEI/dot/config"
	dot "github.com/LEI/dot/dotfile"
	// "github.com/LEI/dot/fileutil"
	// "github.com/LEI/dot/log"
	"github.com/LEI/dot/prompt"
	"github.com/LEI/dot/role"
	"github.com/spf13/cobra"
	"os"
	"path"
	// "strings"
)

var RemoveEmpty bool

var removeCmd = &cobra.Command{
	Use:     "remove [flags]",
	Aliases: []string{"rm"},
	Short:   "Remove dotfiles",
	// Long:   ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			logger.Warnln("Extra arguments:", args)
			return cmd.Help()
		}
		return removeRoles(Dot.Source, Dot.Target, Dot.Roles)
	},
}

func init() {
	RootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolVarP(&RemoveEmpty, "empty", "", RemoveEmpty, "Remove empty files")
	// removeCmd.Flags().BoolVarP(&, "", "", , "")
	// Config.BindPFlags(removeCmd.Flags())
}

func removeRoles(source, target string, roles []*role.Role) error {
	dot.RemoveEmptyFile = RemoveEmpty
	var handlers = []func(*role.Role) error{
		validateRole,
		initGitRepo, // TODO update role.Path without git pull?
		initRoleConfig,
		removeLinks,
		removeLines,
		removeDirs,
	}
ROLES:
	for _, r := range roles {
		r, err := r.New(source, target)
		if err != nil {
			return err
		}
		for _, f := range handlers {
			err := f(r)
			if err != nil {
				switch err {
				case Skip:
					continue ROLES
				}
				return err
			}
		}
	}
	return nil
}

func removeDirs(r *role.Role) error {
	for _, d := range r.Dirs() {
		d.Path = os.ExpandEnv(d.Path)
		dir := path.Join(r.Target, d.Path)
		logger.Debugf("Remove directory %s\n", d.Path)

		d, err := os.Open(dir)
		if err != nil {
			if os.IsNotExist(err) {
				logger.Infof("# rmdir %s\n", dir)
				continue
			}
			return err
		}
		defer d.Close()
		di, err := d.Readdir(-1)
		if err != nil {
			return err
		}
		if len(di) > 0 {
			logger.Warnf("%s is not empty\n", dir)
			break
		}
		if RemoveEmpty || prompt.Confirm("Remove empty directory %s?", dir) {
			err := os.Remove(dir)
			if err != nil {
				return err
			}
		}
		logger.Infof("$ rmdir %s\n", dir)
	}
	return nil
}

func removeLinks(r *role.Role) error {
	for _, l := range r.Links() {
		logger.Debugf("Unlink %s\n", l.Pattern)
		l.Pattern = os.ExpandEnv(l.Pattern)
		pattern := path.Join(r.Source, l.Pattern)
		files, err := dot.List(pattern, filterIgnored)
		if err != nil {
			return err
		}
		for _, f := range files {
			isDir, err := f.IsDir()
			if err != nil {
				return err
			}
			switch {
			case l.Type == "directory" && !isDir:
				logger.Debugf("Ignore directory %s\n", f.Base())
				continue
			case l.Type == "file" && isDir:
				logger.Debugf("Ignore file %s\n", f.Base())
				continue
			}
			ln := dot.NewLink(f.Path(), f.Replace(r.Source, r.Target))
			linked, err := ln.IsLinked()
			if err != nil {
				return err
			}
			if !linked {
				logger.Infof("# rm %s\n", ln.Target())
				continue
			}
			err = os.Remove(ln.Target())
			if err != nil {
				return err
			}
			logger.Infof("$ rm %s\n", ln.Target())
		}
	}
	return nil
}

func removeLines(r *role.Role) error {
	var prefix string
	for _, l := range r.Lines() {
		logger.Debugf("Line out %s\n", l.File)
		l.File = os.ExpandEnv(l.File)
		l.File = path.Join(r.Target, l.File)
		changed, err := dot.LineOutFile(l.File, l.Line)
		if err != nil {
			return err
		}
		if changed {
			prefix = "$"
		} else {
			prefix = "#"
		}
		logger.Infof("%s grep -v '%s' > %s\n", prefix, l.Line, l.File)
	}
	return nil
}
