#!bin/sh

set -e

DOWNLOAD_URL="https://github.com/LEI/dot/releases/download"
PREFIX="${PREFIX:-/usr/local}"
# test -z "$TMPDIR" && TMPDIR="$(mktemp -d)"

last_version() {
  # # Get version from Homebrew formula
  # curl -s https://raw.githubusercontent.com/LEI/homebrew-dot/master/dot.rb |
  #   grep url |
  #   cut -f8 -d'/'
  curl -s https://raw.githubusercontent.com/LEI/dot/master/VERSION
}

download() {
  version="v$(last_version)" # || true
  test -z "$version" && {
    echo "Unable to get dot version."
    exit 1
  }
  echo "Downloading dot $version for $(uname -s)_$(uname -m)..."
  rm -f /tmp/dot.tar.gz
  curl -s -L -o /tmp/dot.tar.gz \
    "$DOWNLOAD_URL/$version/dot-$(uname -s)_$(uname -m).tar.gz"
}

extract() {
  tar -xf /tmp/dot.tar.gz -C "$TMPDIR"
}

download
extract # || true
# sudo mv -f "$TMPDIR"/dot /usr/local/bin/dot
mv -f "$TMPDIR"/dot "$PREFIX"/bin/dot
command -v dot
