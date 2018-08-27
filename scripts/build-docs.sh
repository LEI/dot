#!/bin/sh

set -e

mkdir docs

cp README.md docs/index.md

# make install && "$GOPATH/bin/dot" doc --markdown docs

make docs

# Visualizing dependencies
# https://golang.github.io/dep/docs/daily-dep.html#visualizing-dependencies
if hash /usr/bin/dot 2>/dev/null; then
  dep status -dot | /usr/bin/dot -T png > docs/deps.png
fi
