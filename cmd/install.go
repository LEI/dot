package cmd

import (
	"fmt"
	dot "github.com/LEI/dot/dotfile"
	"github.com/LEI/dot/prompt"
	"github.com/LEI/dot/role"
	"github.com/spf13/cobra"
	"os"
	"path"
	"strings"
)

var installCmd = &cobra.Command{
	Hidden:  true,
	Use:     "install [flags]",
	Aliases: []string{"i"},
	Short:   "Install dotfiles",
	// Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		dot.DryRun = dryrun
		if len(args) > 0 {
			logger.Warnln("Extra arguments:", args)
			return cmd.Help()
		}
		return installRoles()
	},
}

// type RoleHandler interface { Next(r *role.Role) error }
// type RoleHandlerFunc func(*role.Role) error

func init() {
	RootCmd.AddCommand(installCmd)
	// installCmd.Flags().BoolVarP(&, "", "", , "")
	// Config.BindPFlags(installCmd.Flags())
}

func installRoles() error {
	var handlers = []func(*role.Role) error{
		validateRole,
		initGitRepo,
		initRoleConfig,
		installDirs,
		installLinks,
		installLines,
	}
	return apply(handlers...)
}

func installDirs(r *role.Role) error {
	var prefix string
	for _, d := range r.Dirs() {
		d.Path = os.ExpandEnv(d.Path)
		path := path.Join(r.Target, d.Path)
		logger.Debugf("Create directory %s\n", path)
		created, err := dot.CreateDir(path, 0755)
		if err != nil {
			return err
		}
		if created {
			prefix = "$"
		} else {
			prefix = "#"
		}
		logger.Infof("%s mkdir -p %s\n", prefix, path)
	}
	return nil
}

func installLinks(r *role.Role) error {
	var prefix string
	for _, l := range r.Links() {
		logger.Debugf("Symlink %s\n", l.Pattern)
		l.Pattern = os.ExpandEnv(l.Pattern)
		pattern := path.Join(r.Source, l.Pattern)
		paths, err := dot.List(pattern, filterIgnored, only(l.Type))
		if err != nil {
			return err
		}
		for _, source := range paths {
			target := strings.Replace(source, r.Source, r.Target, 1)
			linked, err := dot.InstallSymlink(source, target, removeOrBackup)
			if err != nil {
				return err
			}
			if linked {
				prefix = "$"
			} else {
				prefix = "#"
			}
			logger.Infof("%s ln -s %s %s\n", prefix, source, target)
		}
	}
	return nil
}

func removeOrBackup(path string, link string) (bool, error) {
	if link != "" {
		msg := fmt.Sprintf("> %s is a link to %s, remove?", target, link)
		if ok := prompt.Confirm(msg); !ok {
			return false, nil
		}
		if !dot.DryRun {
			err := os.Remove(target)
			if err != nil {
				return false, err
			}
		}
	} else {
		new := path + ".backup"
		msg := fmt.Sprintf("> %s already exists, backup?", path)
		if ok := prompt.Confirm(msg); !ok {
			return false, nil
		}
		if dot.DryRun {
			err := os.Rename(target, new)
			if err != nil {
				return false, err
			}
		}
	}
	return true, nil
}

func installLines(r *role.Role) error {
	var prefix string
	for _, l := range r.Lines() {
		logger.Debugf("Line in %s\n", l.File)
		l.File = os.ExpandEnv(l.File)
		l.File = path.Join(r.Target, l.File)
		changed, err := dot.LineInFile(l.File, l.Line)
		if err != nil {
			return err
		}
		if changed {
			prefix = "$"
		} else {
			prefix = "#"
		}
		logger.Infof("%s echo '%s' >> %s\n", prefix, l.Line, l.File)
	}
	return nil
}
