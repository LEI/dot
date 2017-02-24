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
	"strings"
)

var RemoveEmpty bool

var removeCmd = &cobra.Command{
	Use:     "remove [flags]",
	Aliases: []string{"rm"},
	Short:   "Remove dotfiles",
	// Long:   ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		dot.DryRun = dryrun
		dot.RemoveEmptyFile = RemoveEmpty
		if len(args) > 0 {
			logger.Warnln("Extra arguments:", args)
			return cmd.Help()
		}
		return removeRoles()
	},
}

func init() {
	RootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolVarP(&RemoveEmpty, "empty", "", RemoveEmpty, "Remove empty files")
	// removeCmd.Flags().BoolVarP(&, "", "", , "")
	// Config.BindPFlags(removeCmd.Flags())
}

func removeRoles() error {
	var handlers = []func(*role.Role) error{
		validateRole,
		initGitRepo, // TODO update role.Path without git pull?
		initRoleConfig,
		removeLinks,
		removeLines,
		removeDirs,
	}
	return apply(handlers...)
}

func removeDirs(r *role.Role) error {
	var prefix string
	for _, d := range r.Dirs() {
		d.Path = os.ExpandEnv(d.Path)
		path := path.Join(r.Target, d.Path)
		logger.Debugf("Remove directory %s\n", d.Path)

		di, err := dot.ReadDir(path)
		if err != nil {
			if os.IsNotExist(err) {
				logger.Infof("# rmdir %s\n", path)
				continue
			}
			return err
		}
		if len(di) > 0 {
			logger.Warnf("%s is not empty\n", path)
			break
		}
		prefix = "#"
		if RemoveEmpty || prompt.Confirm("> Remove empty directory %s?", path) {
			removed, err := dot.RemoveDir(path)
			if err != nil {
				return err
			}
			if removed {
				prefix = "$"
			}
		}
		logger.Infof("%s rmdir %s\n", prefix, path)
	}
	return nil
}

func removeLinks(r *role.Role) error {
	var prefix string
	for _, l := range r.Links() {
		logger.Debugf("Unlink %s\n", l.Pattern)
		l.Pattern = os.ExpandEnv(l.Pattern)
		pattern := path.Join(r.Source, l.Pattern)
		paths, err := dot.List(pattern, filterIgnored, only(l.Type))
		if err != nil {
			return err
		}
		for _, source := range paths {
			target := strings.Replace(source, r.Source, r.Target, 1)
			removed, err := dot.RemoveSymlink(source, target)
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
