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
  log "$PKG_CMD: update packages"
  run $PKG_CMD $PKG_UPD
}

pkg_install() {
  log "$PKG_CMD: install $# packages"
  if [[ $# -ne 0 ]]
  then run $PKG_CMD $PKG_ADD "$@"
  fi
}

pkg_delete() {
  log "$PKG_CMD: remove $# packages"
  if [[ $# -ne 0 ]]
  then run $PKG_CMD $PKG_DEL "$@"
  fi
}
