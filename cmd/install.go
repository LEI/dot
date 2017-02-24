package cmd

import (
	"fmt"
	dot "github.com/LEI/dot/dotfile"
	"github.com/LEI/dot/prompt"
	"github.com/LEI/dot/role"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var installCmd = &cobra.Command{
	Hidden:  true,
	Use:     "install [flags]",
	Aliases: []string{"i"},
	Short:   "Install dotfiles",
	// Long:    ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			logger.Warnln("Extra arguments:", args)
			return cmd.Help()
		}
		return installRoles(Dot.Source, Dot.Target, Dot.Roles)
	},
}

// type RoleHandler interface { Next(r *role.Role) error }
// type RoleHandlerFunc func(*role.Role) error

func init() {
	RootCmd.AddCommand(installCmd)
	// installCmd.Flags().BoolVarP(&, "", "", , "")
	// Config.BindPFlags(installCmd.Flags())
}

func installRoles(source, target string, roles []*role.Role) error {
	// var ctx context.Context
	// var cancel context.CancelFunc

	// ctx, cancel = context.WithCancel(context.Background())
	// defer cancel()

	var handlers = []func(*role.Role) error{
		validateRole,
		initGitRepo,
		initRoleConfig,
		// beforeInstall,
		installDirs,
		installLinks,
		installLines,
		// afterInstall,
	}

ROLES:
	for _, r := range roles {
		r, err := r.New(source, target)
		if err != nil {
			return err
		}
		// h := &ContextAdapter{
		// 	handler: middleware(ContextHandlerFunc(handler)),
		// }
		// a := chain(handler, register(ContextHandlerFunc(checkHandler)))
		for _, f := range handlers {
			// fmt.Println("Handler #", i, f)
			// h := register(ContextHandlerFunc(f))
			// ctx = context.WithValue(ctx, "role", r.Name)
			err := f(r)
			if err != nil {
				switch err {
				case Skip:
					continue ROLES
				}
				return err
			}
		}
		// gh := register(ContextHandlerFunc(handler))
		// gh.Next(ctx, r)

		// err := installRole(ctx, r, func(rol *role.Role) error {
	}
	// c := make(chan error, 1)
	// go func() {
	// 	c <- installRoles(ctx, source, target, roles)
	// }()
	// select {
	// case <-ctx.Done():
	// 	<-c
	// 	return ctx.Err()
	// case err := <-c:
	// 	return err
	// }
	return nil
}

func installDirs(r *role.Role) error {
	for _, d := range r.Dirs() {
		d.Path = os.ExpandEnv(d.Path)
		dir := path.Join(r.Target, d.Path)
		logger.Debugf("Create directory %s\n", dir)
		fi, err := os.Stat(dir)
		if err == nil {
			if fi.IsDir() {
				logger.Infof("# mkdir -p %s\n", dir)
				return nil
			}
			// return &os.PathError{"dir", f.path, syscall.ENOTDIR}
		}
		f, err := dot.NewDir(dir, 0755)
		if err != nil {
			return err
		}
		logger.Infof("$ mkdir -p %s\n", f.Path())
	}
	return nil
}

func installLinks(r *role.Role) error {
	for _, l := range r.Links() {
		logger.Debugf("Symlink %s\n", l.Pattern)
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
				logger.Debugf("# ignore directory %s\n", f.Base())
				continue
			case l.Type == "file" && isDir:
				logger.Debugf("# ignore file %s\n", f.Base())
				continue
			}
			ln := dot.NewLink(f.Path(), f.Replace(r.Source, r.Target))
			linked, err := ln.IsLinked()
			if err != nil {
				return err
			}
			if linked {
				logger.Infof("# ln -s %s %s\n", ln.Path(), ln.Target())
				continue
			}
			err = roleSymlink(ln)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func roleSymlink(ln *dot.Link) error {
	if ln.IsLink() {
		err := readSymlink(ln)
		if err != nil {
			return err
		}
	}
	fi, err := ln.Dstat()
	if err != nil && os.IsExist(err) {
		return err
	}
	if fi != nil {
		moved, err := backupFile(ln.Target())
		if err != nil {
			return err
		}
		if !moved {
			logger.Warn("Ignore existing link: %s", ln.Target())
			return nil
		}
	}
	err = os.Symlink(ln.Path(), ln.Target())
	if err != nil {
		return err
	}
	logger.Infof("$ ln -s %s %s\n", ln.Path(), ln.Target())
	return nil
}

func readSymlink(ln *dot.Link) error {
	link, err := ln.Readlink()
	if err != nil && os.IsExist(err) {
		return err
	}
	if link != "" {
		msg := fmt.Sprintf("> %s is a link to %s, remove?", ln.Target(), link)
		if ok := prompt.Confirm(msg); ok {
			err := os.Remove(ln.Target())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func backupFile(path string) (bool, error) {
	backup := path + ".backup"
	msg := fmt.Sprintf("> %s already exists, backup?", path)
	if ok := prompt.Confirm(msg); !ok {
		return false, nil
	}
	err := os.Rename(path, backup)
	if err != nil {
		return false, err
	}
	return true, nil
}

func filterIgnored(f *dot.File) bool {
	ignore, err := f.BaseMatch(DotIgnore...)
	if err != nil {
		logger.Error(err)
	}
	if ignore {
		logger.Debugf("Ignore %s\n", f.Base())
		return false
	}
	return true
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
