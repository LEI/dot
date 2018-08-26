## dot remove

Remove managed files

### Synopsis

The "remove" command removes roles by executing their tasks.

```
dot remove [flags]
```

### Options

```
  -c, --config-file DOT_FILE   main configuration DOT_FILE (default ".dotrc.yml")
  -d, --dry-run                do not execute tasks
  -F, --force                  force execution
  -h, --help                   help for remove
  -P, --packages               manage system packages
  -r, --role-filter strings    filter roles by name
  -s, --source DOT_SOURCE      DOT_SOURCE directory (default "$HOME")
  -t, --target DOT_TARGET      DOT_TARGET directory (default "$HOME")
```

### Options inherited from parent commands

```
  -q, --quiet       do not output
  -v, --verbose n   be verbose (specify --verbose multiple times or level n)
```

### SEE ALSO

* [dot](dot.md)	 - Manage files
* [dot remove copy](dot_remove_copy.md)	 - Copy managed files
* [dot remove directory](dot_remove_directory.md)	 - Dir managed files
* [dot remove line](dot_remove_line.md)	 - Line managed files
* [dot remove link](dot_remove_link.md)	 - Link managed files
* [dot remove pkg](dot_remove_pkg.md)	 - Package manager
* [dot remove template](dot_remove_template.md)	 - Template managed files

