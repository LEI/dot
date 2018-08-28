#!/bin/sh

set -e

mkdir docs

cp README.md docs/index.md

make docs

# Visualizing dependencies
# https://golang.github.io/dep/docs/daily-dep.html#visualizing-dependencies
dep status -dot | /usr/bin/dot -T png > docs/deps.png
