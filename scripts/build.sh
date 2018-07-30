#!/bin/bash

set -e

if ! hash dep 2> /dev/null; then
  curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
fi

if ! hash goreleaser 2> /dev/null; then
  go get -d github.com/goreleaser/goreleaser && (
    cd "$GOPATH/src/github.com/goreleaser/goreleaser"
    dep ensure -vendor-only
    make setup build
  )
fi

goreleaser --rm-dist --snapshot

ls -la dist
