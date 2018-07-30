#!/bin/bash

set -e

if ! hash dep 2> /dev/null; then
  curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
fi

if ! hash goreleaser 2> /dev/null; then
  REPO_GORELEASER=github.com/goreleaser/goreleaser
  go get -d "$REPO_GORELEASER" && (
    cd "$GOPATH/src/$REPO_GORELEASER"
    dep ensure -vendor-only
    make setup build
    go install "$REPO_GORELEASER"
  )
fi

goreleaser --rm-dist "$@"

ls -la dist
