#!/bin/bash

set -e

DIR="${BASH_SOURCE%/*}"

source "$DIR/install.sh"

# curl -sL https://git.io/goreleaser | bash --rm-dist "$@"

if ! hash goreleaser 2> /dev/null; then
  REPO_GORELEASER=github.com/goreleaser/goreleaser
  go get -d "$REPO_GORELEASER" && (
    cd "$GOPATH/src/$REPO_GORELEASER"
    dep ensure -vendor-only
    make setup build
    go install "$REPO_GORELEASER"
  )
fi

# Build binaries
goreleaser --rm-dist "$@"
