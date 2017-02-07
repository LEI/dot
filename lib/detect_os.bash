#!/usr/bin/env bash

# export OS PM

detect_os() {
  case "$OSTYPE" in # ${OSTYPE//[0-9.]/}
    darwin*) OS="macos"; PM="brew" ;;
    linux-android) OS="android"; PM="apt-get" ;;
    linux*)
      if has apt-get 2>/dev/null
      then OS="debian"; PM="apt-get"
      elif has apt 2>/dev/null
      then OS="debian"; PM="apt"
      else err "$OSTYPE: unknown package manager"; return 1
      fi
      ;;
    *) err "$OSTYPE: unknown operating system"; return 1 ;;
  esac
  [[ -n "$OS" ]] || [[ -z "$PM" ]] && return 1
}
