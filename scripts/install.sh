#!/bin/bash

# https://github.com/LEI/dot/blob/master/scripts/install.sh

set -e

REPO="$REPO"

if [[ ! -d "$GOPATH/src/$REPO" ]]; then
  # git clone https://$REPO.git
  # Use --recursive for .gitmodules

  go get "$REPO"
fi

if ! hash dot 2> /dev/null; then
  go install "$REPO"
fi

dot --dry-run --verbose
