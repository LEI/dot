#!/bin/sh

set -e

mkdir docs

cp README.md docs/index.md

make docs

command -v dot

ls -la "$GOPATH/bin"
make docs # $GOPATH/bin/dot doc --markdown docs/

# Visualizing dependencies
dep status -dot | /usr/bin/dot -T png > docs/deps.png
