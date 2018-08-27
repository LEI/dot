## dot install

Install tasks

### Synopsis

The "install" command installs roles by executing their tasks.

```
dot install [flags]
```

### Options

```
  -c, --config-file DOT_FILE   main configuration DOT_FILE (default ".dotrc.yml")
  -d, --dry-run                do not execute tasks
  -F, --force                  force execution
  -h, --help                   help for install
  -P, --packages               manage system packages
  -r, --role-filter strings    filter roles by name
  -s, --source DOT_SOURCE      DOT_SOURCE directory (default "$HOME")
  -S, --sync                   synchronize repositories
  -t, --target DOT_TARGET      DOT_TARGET directory (default "$HOME")
```

### Options inherited from parent commands

```
  -q, --quiet       do not output
  -v, --verbose n   be verbose (specify --verbose multiple times or level n)
```

### SEE ALSO

* [dot](dot.md)	 - Manage files
* [dot install copy](dot_install_copy.md)	 - File tasks
* [dot install directory](dot_install_directory.md)	 - Directory tasks
* [dot install line](dot_install_line.md)	 - Line in file tasks
* [dot install link](dot_install_link.md)	 - Symbolic link tasks
* [dot install package](dot_install_package.md)	 - Package tasks
* [dot install template](dot_install_template.md)	 - Template tasks
