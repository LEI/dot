package main

import (
	"github.com/spf13/cobra"
)

var cmdInit = &cobra.Command{
	Use:     "init",
	Aliases: []string{},
	Short:   "Initialize a new repository",
	Long: `
The "init" command initializes a new repository.
`,
	DisableAutoGenTag: true,
	Args:              cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runInit(globalOptions, args)
	},
}

func init() {
	cmdRoot.AddCommand(cmdInit)
}

func runInit(gopts GlobalOptions, args []string) error {
	// if gopts.Repo == "" {
	// 	return fmt.Errorf("Please specify repository location (-r)")
	// }

	// be, err := create(gopts.Repo, gopts.extended)
	// if err != nil {
	// 	return fmt.Errorff("create repository at %s failed: %v\n", gopts.Repo, err)
	// }

	// gopts.password, err = ReadPasswordTwice(gopts,
	// 	"enter password for new repository: ",
	// 	"enter password again: ")
	// if err != nil {
	// 	return err
	// }

	// s := repository.New(be)

	// err = s.Init(gopts.ctx, gopts.password)
	// if err != nil {
	// 	return fmt.Errorff("create key in repository at %s failed: %v\n", gopts.Repo, err)
	// }

	// Verbosef("created dot repository %v at %s\n", s.Config().ID[:10], gopts.Repo)
	// Verbosef("\n")
	// Verbosef("Please note that knowledge of your password is required to access\n")
	// Verbosef("the repository. Losing your password means that your data is\n")
	// Verbosef("irrecoverably lost.\n")

	return nil
}
