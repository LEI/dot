package cmd

import (
	"bytes"
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
	// SuggestFor: []string{},
	Short:   "Install dotfiles",
	// Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
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
	// DotCmd.AddCommand(installCmd)
	// installCmd.Flags().BoolVarP(&, "", "", , "")
	// Config.BindPFlags(installCmd.Flags())
}

func installRoles() error {
	var handlers = []func(*role.Role) error{
		func(r *role.Role) error {
			logger.Infof("## Installing %s...\n", r.Title())
			return r.Repo.CloneOrPull()
		},
		initRoleConfig,
		do("pre", "install"),
		installDirs,
		installLinks,
		installLines,
		installTemplates,
		do("post", "install"),
	}
	return apply(Dot.Roles, handlers...)
}

func installDirs(r *role.Role) error {
	var prefix string
	for _, d := range r.GetDirs() {
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
	var targetDir, prefix string
	for _, l := range r.GetLinks() {
		logger.Debugf("Symlink %s\n", l.Path)
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
			linked, err := dot.SyncLink(source, target, removeOrBackup)
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
	for _, l := range r.GetLines() {
		logger.Debugf("Line in %s\n", l.Path)
		l.Path = os.ExpandEnv(l.Path)
		l.Path = path.Join(r.Target, l.Path)
		changed, err := dot.LineInFile(l.Path, l.Line)
		if err != nil {
			return err
		}
		if changed {
			prefix = "$"
		} else {
			prefix = "#"
		}
		logger.Infof("%s echo '%s' >> %s\n", prefix, l.Line, l.Path)
	}
	return nil
}

func installTemplates(r *role.Role) error {
	var prefix string
	for _, t := range r.GetTemplates() {
		logger.Debugf("Template %s\n", t.Path)
		vars, err := r.GetEnv()
		if err != nil {
			return err
		}
		for k, v := range vars {
			logger.Debugf("%s=\"%s\"\n", k, v)
			if v != "" {
				err = os.Setenv(k, v)
				if err != nil {
					return err
				}
			} else {
				logger.Warnf("Empty variable: %s", k)
			}
		}
		t.Path = os.ExpandEnv(t.Path)
		pattern := path.Join(r.Source, t.Path)
		source := path.Clean(pattern)
		target := path.Join(r.Target, strings.TrimSuffix(t.Path, ".tpl"))
		tmpl, err := template.ParseGlob(pattern)
		tmpl = tmpl.Option("missingkey=zero")
		if err != nil {
			return err
		}
		buf := &bytes.Buffer{}
		err = tmpl.Execute(buf, role.Env())
		if err != nil {
			return err
		}
		str := buf.String()
		changed, err := dot.WriteString(target, str)
		if err != nil {
			return err
		}
		if changed {
			prefix = "$"
		} else {
			prefix = "#"
		}
		logger.Infof("%s envsubst < %s | tee %s\n", prefix, source, target)
	}
	return nil
}
