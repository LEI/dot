# dot

[![GoDoc](https://godoc.org/github.com/LEI/dot?status.svg)](https://godoc.org/github.com/LEI/dot)
[![Build Status](https://travis-ci.org/LEI/dot.svg?branch=master)](https://travis-ci.org/LEI/dot)
<!-- [![Go Report
Card](https://goreportcard.com/badge/github.com/LEI/dot)](https://goreportcard.com/report/github.com/LEI/dot)
-->

- [go-flags](https://github.com/jessevdk/go-flags)
- [pacapt](https://github.com/icy/pacapt)

### Requirements

#### [Install dep](https://golang.github.io/dep/docs/installation.html<Paste>)

##### MacOS

    brew install dep

##### ArchLinux

    pacman -S dep

##### From Source (GOPATH)

    go get -u github.com/golang/dep/cmd/dep

#### Install dependencies

    dep ensure

### Installation

Download a [release](https://github.com/LEI/dot/releases) or install manually:

    got get github.com/LEI/dot

### Configuration

`.dot.yml`

## TODO

- [ ] Copy role, Block in file
- [ ] Advanced conditions e.g. `os: ["!darwin"]`
- [ ] Ansible `with_items` equivalent for templates
- [ ] Chooser (overwrite, backup, remove) instead of (y/n) confirmation, plus prompt for missing variables
- [ ] Handle links in cache and add ability to cleanup backups
- [ ] Split pkg pac command and osx defaults
