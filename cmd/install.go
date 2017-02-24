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
	"text/template"
)

var installCmd = &cobra.Command{
	Hidden:  true,
	Use:     "install [flags]",
	Aliases: []string{"i"},
	Short:   "Install dotfiles",
	// Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		// dot.DryRun = dryrun
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
		// do("pre", "install"),
		installDirs,
		installLinks,
		installLines,
		// installTemplates,
		// do("post", "install"),
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
		logger.Debugf("Symlink %s\n", l.Path)
		l.Path = os.ExpandEnv(l.Path)
		glob := path.Join(r.Source, l.Path)
		paths, err := dot.List(glob, filterIgnored, only(l.Type))
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
		msg := fmt.Sprintf("> %s is a link to %s, remove?", path, link)
		if ok := prompt.Confirm(msg); !ok {
			return false, nil
		}
		if !dot.DryRun {
			err := os.Remove(path)
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
		if !dot.DryRun {
			err := os.Rename(path, new)
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

func installTemplates(r *role.Role) error {
	var prefix string
	for _, t := range r.Templates() {
		logger.Debugf("Template %s\n", t.Path)
		t.Path = os.ExpandEnv(t.Path)
		glob := path.Join(r.Source, t.Path)
		// dest := path.Join(r.Target, t.Path)
		tmpl, err := template.ParseGlob(glob) // template.Must()
		if err != nil {
			return err
		}
		vars := map[string]string{}
		for k, v := range t.Vars {
			k = strings.Replace(k, "_", " ", -1)
			k = strings.Title(k)
			k = strings.Replace(k, " ", "", -1)
			v = os.ExpandEnv(v)
			// t, err := template.New(k).Parse(v)
			// if err != nil {
			// 	return err
			// }
			// err = t.Execute(os.Stdout, nil)
			// if err != nil {
			// 	return err
			// }
			vars[k] = v
			logger.Infof("         %s -> %s\n", k, v)
		}
		if !dot.DryRun {
			err := tmpl.Execute(os.Stdout, vars)
			if err != nil {
				return err
			}
		}
		prefix = "@"
		logger.Infof("%s template %s\n", prefix, t.Path)
	}
	return nil
}
