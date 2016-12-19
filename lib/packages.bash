#!/usr/bin/env bash

LIB="$(dirname "$BASH_SOURCE")"
source $LIB/utils.bash

brew_pkg() {
  # if ! has brew
  # then /usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"
  # fi
  PKG_CMD="brew"
  PKG_ADD="install"
  PKG_DEL="remove"
  PKG_UPD="update"
}

apt_pkg() {
  PKG_CMD="apt"
  PKG_ADD="install -qqy"
  PKG_DEL="remove -qqy"
  PKG_UPD="update -qqy"
}

apt_get_pkg() {
  PKG_CMD="apt-get"
  PKG_ADD="install -qqy"
  PKG_DEL="remove -qqy"
  PKG_UPD="update -qqy"
}

pkg_update() {
  if [[ $# -ne 0 ]]
  then log "$PKG_CMD: update"
    dry_run $PKG_CMD $PKG_UPD
  fi
}

pkg_install() {
  # declare -u p="$PKG_CMD"
  if [[ $# -ne 0 ]]
  then log "$PKG_CMD: install $*"
    dry_run $PKG_CMD $PKG_ADD "$@"
  fi
}

pkg_delete() {
  # declare -u p="$PKG_CMD"
  if [[ $# -ne 0 ]]
  then log "$PKG_CMD: delete $*"
    dry_run $PKG_CMD $PKG_DEL "$@"
  fi
}
