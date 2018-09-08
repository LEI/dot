#!/bin/sh

set -e

mkdir docs

cp README.md docs/index.md

make docs

# Visualizing dependencies
dep status -dot | /usr/bin/dot -T png >docs/deps.png
