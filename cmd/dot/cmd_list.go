package main

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/LEI/dot/internal/dot"
	"github.com/spf13/cobra"
)

// Options for the list command.
type listOptions struct {
	quiet  bool
	all    bool
	format string
	// filter []string

	// listLong bool
	// host      string
	// tags      restic.TagLists
	// paths     []string
	// recursive bool
}

var listOpts listOptions

var defaultListFormat = "{{.Name}} {{if .Ok}}✓{{end}}" // {{else}}×

var cmdList = &cobra.Command{
	Use:     "list [flags]", //  [snapshotID] [dir...]
	Aliases: []string{"ls"},
	Short:   "List managed files",
	Long:    `The "list" command lists roles and their tasks.`,
	// Example: ``,
	Args:    cobra.NoArgs,
	PreRunE: preRunList,
	RunE:    runList,
	// DisableAutoGenTag: true,
}

func init() {
	cmdRoot.AddCommand(cmdList)

	flags := cmdList.Flags()
	flags.BoolVarP(&listOpts.quiet, "quiet", "q", false, "Only show role names")
	flags.BoolVarP(&listOpts.all, "all", "a", false, "Show all roles (default hides incompatible platforms)")
	flags.StringVarP(&listOpts.format, "format", "", defaultListFormat, "Pretty-print roles using a Go template")
	// flags.StringSliceVarP(&listOpts.filter, "filter", "f", []string{}, "Filter task list")

	// flags.BoolVarP(&listOpts.listLong, "long", "l", false, "use a long listing format showing size and mode")
	// flags.StringVarP(&listOpts.host, "host", "H", "", "only consider snapshots for this `host`, when no snapshot ID is given")
	// flags.Var(&listOpts.tags, "tag", "only consider snapshots which include this `taglist`, when no snapshot ID is given")
	// flags.StringArrayVar(&listOpts.paths, "path", nil, "only consider snapshots which include this (absolute) `path`, when no snapshot ID is given")
	// flags.BoolVar(&listOpts.recursive, "recursive", false, "include files in subfolders of the listed directories")

	addActionFlags(cmdList)
}

func preRunList(cmd *cobra.Command, args []string) error {
	if err := setActionEnv(cmd); err != nil {
		return err
	}
	if listOpts.quiet && listOpts.format != "" && listOpts.format != defaultListFormat {
		return fmt.Errorf("--quiet and --format cannot be specified at the same time")
	}
	if listOpts.format == "" { // && dotOpts.Verbose > 0 {
		listOpts.format = "{{.}}"
	}
	// if len(listOpts.filter) > 0 {
	// 	fmt.Fprintf(os.Stderr, "--filter not implemented\n")
	// }
	if err := preRunAction(cmd, args); err != nil {
		if _, ok := err.(*dot.DiffError); !ok {
			return err
		}
	}
	return nil
}

func runList(cmd *cobra.Command, args []string) error {
	// if len(args) == 0 && opts.Host == "" && len(opts.tags) == 0 && len(opts.paths) == 0 {
	// 	return errors.Fatal("Invalid arguments, either give one or more snapshot IDs or set filters.")
	// }
	// if !listOpts.all {
	// 	dotConfig.Roles.FilterOS()
	// }
	w := dotOpts.stdout // tabwriter.NewWriter(os.Stdout, 8, 8, 8, ' ', 0)
	for _, r := range dotConfig.Roles {
		// fmt.Fprintf(w, "%+v\n", r)
		if listOpts.quiet {
			fmt.Fprintln(w, r.Name)
			continue
		}
		format := listOpts.format
		str, err := templateString(r.Name, format, r)
		if err != nil {
			return err
		}
		fmt.Fprintln(w, str)
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
	// for sn := range FindFilteredSnapshots(ctx, repo, opts.Host, opts.tags, opts.paths, args[:1]) {
	// 	Verbosef("snapshot %s of %v filtered by %v at %s):\n", sn.ID().Str(), sn.paths, dirs, sn.Time)

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
	// 			if opts.recursive {
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

func templateString(name, format string, data interface{}) (string, error) {
	t, err := template.New(name).Parse(format)
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, data); err != nil {
		return "", err
	}
	return tpl.String(), nil
}
