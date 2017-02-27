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
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var RemoveEmpty bool

var removeCmd = &cobra.Command{
	Use:     "remove [flags]",
	Aliases: []string{"rm"},
	SuggestFor: []string{"delete", "uninstall"},
	Short:   "Remove dotfiles",
	// Long:   ``,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		// dot.DryRun = dryrun
		dot.RemoveEmptyFile = RemoveEmpty
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			logger.Warnln("Extra arguments:", args)
			return cmd.Help()
		}
		return removeRoles()
	},
}

func init() {
	// DotCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolVarP(&RemoveEmpty, "empty", "", RemoveEmpty, "Remove empty files")
	// removeCmd.Flags().BoolVarP(&, "", "", , "")
	// Config.BindPFlags(removeCmd.Flags())
}

func removeRoles() error {
	var handlers = []func(*role.Role) error{
		func(r *role.Role) error {
			logger.Infof("## Removing %s...\n", r.Title())
			return nil
		},
		initRoleConfig,
		do("pre", "remove"),
		removeLinks,
		removeLines,
		removeTemplates,
		removeDirs,
		do("post", "remove"),
	}
	return apply(Dot.Roles, handlers...)
}

func removeDirs(r *role.Role) error {
	for _, d := range r.GetDirs() {
		prefix := "#"
		d.Path = os.ExpandEnv(d.Path)
		path := path.Join(r.Target, d.Path)
		logger.Debugf("Remove directory %s\n", d.Path)
		changed, err := removeDir(path)
		if err != nil {
			return err
		}
		if changed {
			prefix = "$"
		}
		logger.Infof("%s rmdir %s\n", prefix, path)
	}
	return nil
}

func removeDir(path string) (bool, error) {
	di, err := dot.ReadDir(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
			// continue
		}
		return false, err
	}
	if len(di) > 0 {
		logger.Warnf("%s is not empty\n", path)
		return false, nil
	}
	if RemoveEmpty || prompt.Confirm("> Remove empty directory %s?", path) {
		removed, err := dot.RemoveDir(path)
		return removed, err
	}
	return false, nil
}

func removeLinks(r *role.Role) error {
	var targetDir, prefix string
	for _, l := range r.GetLinks() {
		logger.Debugf("Unlink %s\n", l.Path)
		l.Path = os.ExpandEnv(l.Path)
		if strings.Contains(l.Path, ":") {
			s := strings.Split(l.Path, ":")
			if len(s) != 2 {
				logger.Errorf("%s: Invalid link path", l.Path)
			}
			l.Path = s[0]
			targetDir = s[1]
		}
		pattern := path.Join(r.Source, l.Path)
		paths, err := dot.List(pattern, filterIgnored, only(l.Type))
		if err != nil {
			return err
		}
		for _, source := range paths {
			s := r.Source
			t := r.Target
			if targetDir != "" {
				s = path.Dir(source)
				t = path.Join(t, targetDir)
			}
			target := strings.Replace(source, s, t, 1)
			removed, err := dot.RemoveLink(source, target)
			if err != nil {
				return err
			}
			if removed {
				prefix = "$"
			} else {
				prefix = "#"
			}
			logger.Infof("%s rm %s\n", prefix, target)
		}
	}
	return nil
}

func removeLines(r *role.Role) error {
	var prefix string
	for _, l := range r.GetLines() {
		logger.Debugf("Line in %s\n", l.Path)
		l.Path = os.ExpandEnv(l.Path)
		l.Path = path.Join(r.Target, l.Path)
		changed, err := dot.LineOutFile(l.Path, l.Line)
		if err != nil {
			return err
		}
		if changed {
			prefix = "$"
		} else {
			prefix = "#"
		}
		// grep -v 'line' "file" > "tmpfile" && mv "tmpfile" "file"
		logger.Infof("%s sed -i '/%s/d' %s\n", prefix, l.Line, l.Path)
	}
	return nil
}

func removeTemplates(r *role.Role) error {
	var prefix string
	for _, t := range r.GetTemplates() {
		logger.Debugf("Template %s\n", t.Path)
		t.Path = os.ExpandEnv(t.Path)
		// pattern := path.Join(r.Source, t.Path)
		// source := path.Clean(pattern)
		target := path.Join(r.Target, strings.TrimSuffix(t.Path, ".tpl"))
		_, err := ioutil.ReadFile(target)
		if err != nil && os.IsExist(err) {
			return err
		}
		if err != nil && os.IsNotExist(err) {
			prefix = "#"
		} else {
			prefix = "$"
			err := os.Remove(target)
			if err != nil {
				return err
			}
		}
		logger.Infof("%s rm %s\n", prefix, target)
	}
	return nil
}
