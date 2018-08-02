#!/bin/sh

set -e

echo "SETUP LOCALE: $LANG"

# LANG="${LANG:-$1}"

sed -i -e 's/#[[:blank:]]*en_GB.UTF-8 UTF-8/en_GB.UTF-8 UTF-8/' /etc/locale.gen

locale-gen

LANG=en_GB.UTF-8
LANGUAGE=en_GB:en
LC_ALL=en_GB.UTF-8
