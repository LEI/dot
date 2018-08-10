#!/bin/sh

set -e

# https://wiki.archlinux.org/index.php/locale
# LANG, LANGUAGE, LC_ALL...
LOCALE="${1:-en_GB.UTF-8 UTF-8}"

if [ -z "$LOCALE" ]; then
  exit 1
fi

sed -i -e "s/#[[:blank:]]*$LOCALE/$LOCALE/" /etc/locale.gen

locale-gen
