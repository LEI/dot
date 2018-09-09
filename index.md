dot
===

<!--
`dot` is a configuration based file manager. It requires Go 1.8 to compile.
-->

[Documentation](https://dot.lei.sh/dot)
<!-- [Contribution Guide](CONTRIBUTING.md) -->

[![GoDoc](https://godoc.org/github.com/LEI/dot?status.svg)](https://godoc.org/github.com/LEI/dot)
[![Travis](https://travis-ci.org/LEI/dot.svg?branch=master)](https://travis-ci.org/LEI/dot)
[![AppVeyor](https://ci.appveyor.com/api/projects/status/s4qqanrbe62cp1ku?svg=true)](https://ci.appveyor.com/project/LEI/dot)
[![Codecov](https://codecov.io/gh/LEI/dot/branch/master/graph/badge.svg)](https://codecov.io/gh/LEI/dot)
[![Go Report Card](https://goreportcard.com/badge/github.com/LEI/dot)](https://goreportcard.com/report/github.com/LEI/dot)

<!--
## Overview
-->

<!--
## License
-->

Requires [git](https://git-scm.com/) 2.11.0 or greater.

## Installation

Release binaries are available on the
[releases](https://github.com/LEI/dot/releases) page.

### macOS

```sh
brew install lei/dot/dot
```

### Other platforms

#### Archlinux

```sh
sudo pacman -Syu --noconfirm
sudo pacman -S --noconfirm ca-certificates git
curl -sSL https://git.io/dot.lei.sh | sh
```

#### Debian

```sh
sudo apt-get update -yq
sudo apt-get install -yyq ca-certificates curl git
curl -sSL https://git.io/dot.lei.sh | sh
```

### From Source

Requires [go](https://golang.org/dl).

```sh
go get -u github.com/LEI/dot
cd $GOPATH/src/github.com/LEI/dot
go run build.go vendor check install # or `make`
```

![deps](https://dot.lei.sh/deps.png)
[Dependency tree generated with graphviz](https://golang.github.io/dep/docs/daily-dep.html#visualizing-dependencies)

<!--
## Feedback
-->

<!--
## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for more details.
-->
