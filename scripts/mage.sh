#!/bin/sh

set -e

# https://magefile.org/zeroinstall
go run mage.go "$@"
