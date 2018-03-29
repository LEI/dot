package cmd

import (
	// "fmt"
	// "github.com/LEI/dot/config"
	dot "github.com/LEI/dot/dotfile"
	// "github.com/LEI/git"
	"github.com/LEI/dot/role"
	// "io/ioutil"
	"os"
	"os/exec"
	// "path"
	// "time"
)

func apply(roles []*role.Role, handlers ...func(*role.Role) error) error {
ROLES:
	for i := range roles {
		for _, f := range handlers {
			err := f(roles[i])
			if err != nil {
				switch err {
				case dot.Skip:
					continue ROLES
				}
				return err
			}
		}
	}
	return nil
}

func do(state string, action string) func(*role.Role) error {
	key := state + "_" + action
	return func(r *role.Role) error {
		//logger.Infof("do -> CONFIG %v\n", r.Config)
		
		command := r.Config.GetString(key)
		if command == "" {
			return nil
		}
		logger.Infof("$ %s\n", command)
		if dot.DryRun {
			return nil
		}
		cmd := exec.Command(shell, "-c", command)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
}

func only(t string) dot.FileHandler {
	return func(path string, fi os.FileInfo) error {
		switch {
		// case fi == nil:
		// 	logger.Debugf("# ignore %s (does not exist)\n", fi.Name(), t)
		// 	return dot.Skip
		case t == "directory" && !fi.IsDir(),
			t == "file" && fi.IsDir():
			logger.Debugf("# ignore %s (not a %s)\n", fi.Name(), t)
			return dot.Skip
			// case t == "":
			// 	logger.Errorf("! Invalid type: %s\n", t)
			// default:
			// 	logger.Debugln("default", fi.Name(), t)
		}
		return nil
	}
}

func filterIgnored(path string, fi os.FileInfo) error {
	ignore, err := dot.Match(path, DotExclude...)
	if err != nil {
		return err
	}
	if ignore {
		logger.Debugf("Ignore %s\n", fi.Name())
		return dot.Skip
	}
	ignore, err = dot.Match(fi.Name(), DotIgnore...)
	if err != nil {
		return err
	}
	if ignore {
		logger.Debugf("Ignore %s\n", fi.Name())
		return dot.Skip
	}
	return nil
}
