#!/bin/bash

sep() { printf %${COLUMNS:-100}s |tr " " "${1:-=}"; printf "\n"; }
log() { sep "-"; printf "\n\t%s\n\n" "$@"; sep "-"; }
run() { log "\$ $*"; "$@" || exit $?; }
