#!/bin/sh

# https://raw.githubusercontent.com/LEI/dot/master/scripts/install.sh
# https://git.io/dot.lei.sh
# curl -sSL https://dot.lei.sh | sh
# wget -qO- https://dot.lei.sh | sh

# https://github.com/rootbeersoup/get.darryl.sh/blob/master/scripts/dotfiles.sh
# https://github.com/golang/dep/blob/master/install.sh
# https://get.docker.com/

set -e

REPO="$REPO"

has() {
  hash "$1" 2> /dev/null
}

do_install() {
  if ! has go; then
    echo "Go command is unavailable"
    exit 1
  fi
  if [ ! -n "$GOPATH" ]; then
    echo "GOPATH is not set"
    exit 1
  fi
  if [ ! -d "$GOPATH/src/$REPO" ]; then
    # git clone https://$REPO.git
    # Use --recursive for .gitmodules

    go get "$REPO"
  fi

  if ! has dot; then
    go install "$REPO"
  fi

  dot --dry-run --verbose
}

do_install
