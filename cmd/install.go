package cmd

import (
	"context"
	"fmt"
	dot "github.com/LEI/dot/dotfile"
	"github.com/LEI/dot/git"
	"github.com/LEI/dot/prompt"
	"github.com/LEI/dot/role"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	// "strings"
)

var installCmd = &cobra.Command{
	Hidden: true,
	Use:    "install [flags]",
	Short:  "Install dotfiles",
	// Long:   ``,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			logger.Warnln("Extra arguments:", args)
			return cmd.Help()
		}
		return installCommand(Dot.Source, Dot.Target, Dot.Roles)
	},
}

var handlers = []ContextHandlerFunc{
	roleInit,
	roleInitRepo,
	roleInitConfig,
	roleInstallDirs,
	roleInstallLinks,
	roleInstallLines,
	roleDone,
}

type ContextHandler interface {
	Next(ctx context.Context, r *role.Role) error
}

type ContextHandlerFunc func(context.Context, *role.Role) error

func (h ContextHandlerFunc) Next(ctx context.Context, r *role.Role) error {
	return h(ctx, r)
}

func register(h ContextHandler) ContextHandler {
	return ContextHandlerFunc(func(ctx context.Context, r *role.Role) error {
		// ctx = newContext(ctx, r)
		return h.Next(ctx, r)
	})
}

// func chain(f ContextHandlerFunc, m ...func(ContextHandlerFunc) ContextHandlerFunc) ContextHandlerFunc {
// 	if len(m) == 0 {
// 		return f
// 	}
// 	return m[0](chain(f, m[1:cap(m)]...))
// }

// func handleGitRepo(f HandlerContextFunc) HandlerContextFunc {
// 	return func(ctx context.Context, *role.Role) error {
// 		fmt.Println("Handle", r.Name, "Git Repo")
// 		return nil
// 	}
// }

// type ContextAdapter struct {
// 	ctx context.Context
// 	handler ContextHandler
// }

func init() {
	RootCmd.AddCommand(installCmd)
	// installCmd.Flags().BoolVarP(&, "", "", , "")
	// Config.BindPFlags(installCmd.Flags())
}

func installCommand(source, target string, roles []*role.Role) error {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	err := installRoles(ctx, source, target, roles)
	if err != nil {
		return err
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

func installRoles(ctx context.Context, source, target string, roles []*role.Role) error {
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
			h := register(ContextHandlerFunc(f))
			ctx = context.WithValue(ctx, "role", r.Name)
			// logger = ctx.Value("logger").(*log.Logger)
			err := h.Next(ctx, r)
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
	return nil
}

func roleInit(ctx context.Context, r *role.Role) error {
	// Check platform
	ok := r.IsOs([]string{OS, OSTYPE})
	if !ok {
		return Skip
	}
	// Filter by name
	skip := len(filter) > 0
	for _, roleName := range filter {
		if roleName == r.Name {
			skip = false
			break
		}
	}
	if skip {
		return Skip
	}
	// logger.SetPrefix(r.Name+": ")
	return nil
}

func roleInitRepo(ctx context.Context, r *role.Role) error {
	dir := filepath.Join(r.Target, RolesDir, r.Name) // git.DefaultPath
	git.Https = https
	repo, err := git.New(r.Origin, dir)
	if err != nil {
		return err
	}
	repo.Name = r.Name
	err = repo.CloneOrPull()
	if err != nil {
		return err
	}
	if repo.Path != r.Source {
		r.Source = repo.Path
	}
	return nil
}

func roleInitConfig(ctx context.Context, r *role.Role) error {
	if r.Config == nil {
		r.Config = viper.New()
	}
	r.Config.SetConfigName(configName)
	r.Config.AddConfigPath(r.Source)
	err := r.Config.ReadInConfig()
	if err != nil { // && !os.IsNotExist(err)
		return err
	}
	cfgUsed := r.Config.ConfigFileUsed()
	if cfgUsed != "" {
		logger.Debugln("Using role config file:", cfgUsed)
	}
	return nil
}

func roleInstallDirs(ctx context.Context, r *role.Role) error {
	for _, d := range r.Dirs() {
		d.Path = os.ExpandEnv(d.Path)
		dir := filepath.Join(r.Target, d.Path)
		f, err := dot.NewDir(dir, 0755)
		if err != nil {
			return err
		}
		logger.Infof("$ mkdir -p %s\n", f)
	}
	return nil
}

func roleInstallLinks(ctx context.Context, r *role.Role) error {
	for _, l := range r.Links() {
		logger.Infof("- Symlink %s\n", l.Pattern)
		l.Pattern = os.ExpandEnv(l.Pattern)
		pattern := filepath.Join(r.Source, l.Pattern)
		files, err := dot.List(pattern, filterIgnored)
		if err != nil {
			return err
		}
		for _, f := range files {
			err := roleSymlink(ctx, r, f)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func roleSymlink(ctx context.Context, r *role.Role, f *dot.File) error {
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
	link, err := ln.Readlink()
	if err != nil && os.IsExist(err) {
		return err
	}
	if link != "" {
		msg := fmt.Sprintf("! %s is a link to %s, remove?", ln.Target(), link)
		if ok := prompt.Confirm(msg); ok {
			err := os.Remove(ln.Target())
			if err != nil {
				return err
			}
		}
	}
	fi, err := ln.DestInfo()
	if err != nil && os.IsExist(err) {
		return err
	}
	if fi != nil {
		if err := backupFile(ln.Target); err != nil {
			return err
		}
	}
	logger.Infof("$ ln -s %s %s\n", ln.Path(), ln.Target())
	err = os.Symlink(ln.Path(), ln.Target())
	if err != nil {
		return err
	}
	return nil
}

func backupFile(path string) bool {
	backup := path + ".backup"
	msg := fmt.Sprintf("! %s already exists, append .backup?", path)
	if ok := prompt.Confirm(msg); ok {
		err := os.Rename(path, backup)
		if err != nil {
			return err
		}
	}
	return nil
}

func filterIgnored(f *dot.File) bool {
	ignore, err := f.Match(DotIgnore...)
	if err != nil {
		logger.Error(err)
	}
	if ignore {
		logger.Debugf("# .ignore %s\n", f.Base())
		return false
	}
	return true
}

func roleInstallLines(ctx context.Context, r *role.Role) error {
	for _, l := range r.Lines() {
		logger.Infof("- Line in %s\n", l.File)
		l.File = os.ExpandEnv(l.File)
		l.File = filepath.Join(r.Target, l.File)
		changed, err := dot.LineInFile(l.File, l.Line)
		if err != nil {
			return err
		}
		if changed {
			logger.Infof("$ echo '%s' >> %s\n", l.Line, l.File)
		} else {
			logger.Infof("# line '%s' already in file %s\n", l.Line, l.File)
		}
	}
	return nil
}

func roleDone(ctx context.Context, r *role.Role) error {
	logger.Infoln("---", ctx.Value("role"))
	return nil
}
