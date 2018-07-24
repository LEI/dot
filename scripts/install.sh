#!/bin/sh

# https://raw.githubusercontent.com/LEI/dot/master/scripts/install.sh
# https://git.io/dot.lei.sh
# curl -sSL https://dot.lei.sh | sh
# wget -qO- https://dot.lei.sh | sh

# https://github.com/rootbeersoup/get.darryl.sh/blob/master/scripts/dotfiles.sh
# https://github.com/golang/dep/blob/master/install.sh
# https://get.docker.com/

set -e

DOT_REPO="github.com/LEI/dot"

has() {
  hash "$1" 2> /dev/null
}

check_go() {
  if ! has go; then
    echo "Go command is unavailable"
    exit 1
  fi
  if [ ! -n "$GOPATH" ]; then
    echo "GOPATH is not set"
    exit 1
  fi
}

check_dot() {
  if [ ! -d "$GOPATH/src/$DOT_REPO" ]; then
    # git clone https://$DOT_REPO.git
    # Use --recursive for .gitmodules
    go get "$DOT_REPO"
  fi
  if ! has dot; then
    go install "$DOT_REPO"
  fi
}

do_install() {
  check_go
  check_dot
  #dot --dry-run --verbose
  echo "Done"
}

do_install
