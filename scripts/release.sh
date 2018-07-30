#!/bin/bash

set -e

if ! hash dep 2> /dev/null; then
  curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
fi

# Install dependencies
dep ensure

if ! hash goreleaser 2> /dev/null; then
  # curl -sL https://git.io/goreleaser | bash --rm-dist "$@"
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
