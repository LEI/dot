package cmd

import (
	"context"
	"fmt"
	"github.com/LEI/dot/prompt"
	"github.com/LEI/dot/fileutil"
	"github.com/LEI/dot/git"
	"github.com/LEI/dot/role"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"strings"
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
	DotCmd.AddCommand(installCmd)
	// installCmd.Flags().BoolVarP(&, "", "", , "")
	// Config.BindPFlags(installCmd.Flags())
}

func installCommand(source, target string, roles []*role.Role) error {
	var ctx context.Context
	var cancel context.CancelFunc

	ctx, cancel = context.WithCancel(context.Background())
	defer cancel()

	c := make(chan error, 1)
	go func() {
		c <- installRoles(ctx, source, target, roles)
	}()
	select {
	case <-ctx.Done():
		logger.Debugln("<-ctx.Done")
		<-c
		logger.Debugln("<-c")
		return ctx.Err()
	case err := <-c:
		logger.Debugln("err := <-c", err)
		return err
	}
	return nil
}

func installRoles(ctx context.Context, source, target string, roles []*role.Role) error {
	var handlers = []ContextHandlerFunc{
		roleInit,
		roleInitRepo,
		roleInitConfig,
		roleInstallDirs,
		roleInstallLinks,
		roleInstallLines,
		roleDone,
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
		// logger.Infof("- Create %s\n", d.Path)
		d.Path = os.ExpandEnv(d.Path)
		d.Path = filepath.Join(r.Target, d.Path)
		logger.Infof("$ mkdir -p %s\n", d.Path)
		err := os.MkdirAll(d.Path, 0755) // os.Remove(d.Path)
		if err != nil {
			return err
		}
	}
	return nil
}

func roleInstallLinks(ctx context.Context, r *role.Role) error {
	for _, l := range r.Links() {
		logger.Infof("- Symlink %s\n", l.Pattern)
		l.Pattern = os.ExpandEnv(l.Pattern)
		paths, err := filepath.Glob(filepath.Join(r.Source, l.Pattern))
		if err != nil {
			return err
		}
GLOB:
		for _, src := range paths {
			base := filepath.Base(src)
			for _, pattern := range IgnoreNames {
				ignore, err := filepath.Match(pattern, base)
				if err != nil {
					return err
				}
				if ignore {
					logger.Debugf("# ignore %s (filename)\n", base)
					continue GLOB
				}
			}
			fi, err := os.Stat(src)
			if err != nil {
				return err
			}
			switch {
				case l.Type == "directory" && !fi.IsDir(),
				l.Type == "file" && fi.IsDir():
				logger.Debugf("# ignore %s (filetype)\n", base)
				continue // GLOB
			}
			dst := strings.Replace(src, r.Source, r.Target, 1)

			fi, err = os.Lstat(dst)
			if err != nil && os.IsExist(err) {
				return err
			}
			if fi != nil && (fi.Mode()&os.ModeSymlink != 0) {
				link, err := os.Readlink(dst)
				if err != nil {
					return err
				}
				if link == src { // TODO os.SameFile?
					logger.Debugf("# ignore %s (already a link)\n", dst)
					// logger.Infof("Already linked: %s", src)
					return nil
				}
				// TODO check broken symlink?
				msg := fmt.Sprintf("! %s exists, linked to %s, replace with %s?", dst, link, src)
				if ok := prompt.Confirm(msg); ok {
					err := os.Remove(dst)
					if err != nil {
						return err
					}
				}
			} else if fi != nil {
				backup := dst + ".backup"
				msg := fmt.Sprintf("! %s exists, move to %s to link %s?", dst, backup, src)
				if ok := prompt.Confirm(msg); ok {
					err := os.Rename(dst, dst+".backup")
					if err != nil {
						return err
					}
				}
			}
			logger.Infof("$ ln -s %s %s\n", src, dst)
			err = os.Symlink(src, dst)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func roleInstallLines(ctx context.Context, r *role.Role) error {
	for _, l := range r.Lines() {
		logger.Infof("- Line in %s\n", l.File)
		l.File = os.ExpandEnv(l.File)
		l.File = filepath.Join(r.Target, l.File)
		changed, err := fileutil.LineInFile(l.File, l.Line)
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
