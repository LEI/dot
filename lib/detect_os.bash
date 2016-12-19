#!/usr/bin/env bash

LIB="$(dirname "$BASH_SOURCE")"
source $LIB/utils.bash

detect_os() {
  case "$OSTYPE" in # ${OSTYPE//[0-9.]/}
    darwin*) OS="macos"; brew_pkg ;;
    linux-android) OS="android"; apt_pkg ;;
    linux*)
      if has apt-get 2>/dev/null
      then OS="debian"; apt_get_pkg
      elif has apt 2>/dev/null
      then OS="debian"; apt_pkg
      else err "$OSTYPE: unknown package manager"; return 1
      fi
      ;;
    *) err "$OSTYPE: unknown operating system"; return 1 ;;
  esac
}
