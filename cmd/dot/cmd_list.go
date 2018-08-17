package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ListOptions collects all options for the list command.
type ListOptions struct {
	ListLong bool
	// Host      string
	// Tags      restic.TagLists
	// Paths     []string
	// Recursive bool
}

var listOptions ListOptions

var cmdList = &cobra.Command{
	Use:     "list [flags]", //  [snapshotID] [dir...]
	Aliases: []string{"ls"},
	Short:   "List managed files",
	Long: `
The "list" command lists roles and their tasks.
`,
	DisableAutoGenTag: true,
	Args:              cobra.NoArgs,
	RunE:              runList,
}

func init() {
	cmdRoot.AddCommand(cmdList)

	flags := cmdList.Flags()
	flags.BoolVarP(&listOptions.ListLong, "long", "l", false, "use a long listing format showing size and mode")
	// flags.StringVarP(&listOptions.Host, "host", "H", "", "only consider snapshots for this `host`, when no snapshot ID is given")
	// flags.Var(&listOptions.Tags, "tag", "only consider snapshots which include this `taglist`, when no snapshot ID is given")
	// flags.StringArrayVar(&listOptions.Paths, "path", nil, "only consider snapshots which include this (absolute) `path`, when no snapshot ID is given")
	// flags.BoolVar(&listOptions.Recursive, "recursive", false, "include files in subfolders of the listed directories")
}

func runList(cmd *cobra.Command, args []string) error {
	// if len(args) == 0 && opts.Host == "" && len(opts.Tags) == 0 && len(opts.Paths) == 0 {
	// 	return errors.Fatal("Invalid arguments, either give one or more snapshot IDs or set filters.")
	// }
	for _, r := range globalConfig.Roles {
		if globalOptions.Quiet {
			fmt.Println(r.Name)
			continue
		}
		// fmt.Printf("%+v\n", r)
		fmt.Println(r)
	}

	// // extract any specific directories to walk
	// var dirs []string
	// if len(args) > 1 {
	// 	dirs = args[1:]
	// 	for _, dir := range dirs {
	// 		if !strings.HasPrefix(dir, "/") {
	// 			return errors.Fatal("All path filters must be absolute, starting with a forward slash '/'")
	// 		}
	// 	}
	// }

	// withinDir := func(nodepath string) bool {
	// 	if len(dirs) == 0 {
	// 		return true
	// 	}

	// 	for _, dir := range dirs {
	// 		// we're within one of the selected dirs, example:
	// 		//   nodepath: "/test/foo"
	// 		//   dir:      "/test"
	// 		if fs.HasPathPrefix(dir, nodepath) {
	// 			return true
	// 		}
	// 	}
	// 	return false
	// }

	// approachingMatchingTree := func(nodepath string) bool {
	// 	if len(dirs) == 0 {
	// 		return true
	// 	}

	// 	for _, dir := range dirs {
	// 		// the current node path is a prefix for one of the
	// 		// directories, so we're interested in something deeper in the
	// 		// tree. Example:
	// 		//   nodepath: "/test"
	// 		//   dir:      "/test/foo"
	// 		if fs.HasPathPrefix(nodepath, dir) {
	// 			return true
	// 		}
	// 	}
	// 	return false
	// }

	// repo, err := OpenRepository(gopts)
	// if err != nil {
	// 	return err
	// }

	// if err = repo.LoadIndex(gopts.ctx); err != nil {
	// 	return err
	// }

	// ctx, cancel := context.WithCancel(gopts.ctx)
	// defer cancel()
	// for sn := range FindFilteredSnapshots(ctx, repo, opts.Host, opts.Tags, opts.Paths, args[:1]) {
	// 	Verbosef("snapshot %s of %v filtered by %v at %s):\n", sn.ID().Str(), sn.Paths, dirs, sn.Time)

	// 	err := walker.Walk(ctx, repo, *sn.Tree, nil, func(nodepath string, node *restic.Node, err error) (bool, error) {
	// 		if err != nil {
	// 			return false, err
	// 		}
	// 		if node == nil {
	// 			return false, nil
	// 		}

	// 		if withinDir(nodepath) {
	// 			// if we're within a dir, print the node
	// 			Printf("%s\n", formatNode(nodepath, node, lsOptions.ListLong))

	// 			// if recursive listing is requested, signal the walker that it
	// 			// should continue walking recursively
	// 			if opts.Recursive {
	// 				return false, nil
	// 			}
	// 		}

	// 		// if there's an upcoming match deeper in the tree (but we're not
	// 		// there yet), signal the walker to descend into any subdirs
	// 		if approachingMatchingTree(nodepath) {
	// 			return false, nil
	// 		}

	// 		// otherwise, signal the walker to not walk recursively into any
	// 		// subdirs
	// 		if node.Type == "dir" {
	// 			return false, walker.SkipNode
	// 		}
	// 		return false, nil
	// 	})

	// 	if err != nil {
	// 		return err
	// 	}
	// }
	return nil
}
