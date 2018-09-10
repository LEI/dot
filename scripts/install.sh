#!/bin/sh

set -e

# https://github.com/golang/dep/blob/master/install.sh

# This install script is intended to download and install the latest available
# release of the dot files manager.
#
# It attempts to identify the current platform and an error will be thrown if
# the platform is not supported.
#
# Environment variables:
# - INSTALL_DIRECTORY (optional): defaults to /usr/local/bin
# - DOT_RELEASE_TAG (optional): defaults to fetching the latest release
# - DOT_OS (optional): use a specific value for OS (mostly for testing)
# - DOT_ARCH (optional): use a specific value for ARCH (mostly for testing)
#
# You can install using this script:
# https://raw.githubusercontent.com/LEI/dot/master/scripts/install.sh
# $ curl -L https://dot.lei.sh | sh

set -e

INSTALL_DIRECTORY="${INSTALL_DIRECTORY:-${PREFIX:-/usr/local}/bin}"
RELEASES_URL="https://github.com/LEI/dot/releases"
test -z "$TMPDIR" && TMPDIR="$(mktemp -d)"

downloadJSON() {
  url="$2"

  echo "Fetching $url..."
  if test -x "$(command -v curl)"; then
    response=$(curl -s -L -w 'HTTPSTATUS:%{http_code}' -H 'Accept: application/json' "$url")
    body=$(echo "$response" | sed -e 's/HTTPSTATUS\:.*//g')
    code=$(echo "$response" | tr -d '\n' | sed -e 's/.*HTTPSTATUS://')
  elif test -x "$(command -v wget)"; then
    temp=$(mktemp)
    body=$(wget -q --header='Accept: application/json' -O - --server-response "$url" 2>"$temp")
    code=$(awk '/^  HTTP/{print $2}' <"$temp" | tail -1)
    rm "$temp"
  else
    echo >&2 "Neither curl nor wget was available to perform http requests."
    exit 1
  fi
  if [ "$code" != 200 ]; then
    echo >&2 "Request failed with code $code"
    exit 1
  fi

  eval "$1='$body'"
}

downloadFile() {
  url="$1"
  destination="$2"

  echo "Fetching $url..."
  if test -x "$(command -v curl)"; then
    echo curl -s -w '%{http_code}' -L "$url" -o "$destination"
    code=$(curl -s -w '%{http_code}' -L "$url" -o "$destination")
  elif test -x "$(command -v wget)"; then
    code=$(wget -q -O "$destination" --server-response "$url" 2>&1 | awk '/^  HTTP/{print $2}' | tail -1)
  else
    echo >&2 "Neither curl nor wget was available to perform http requests."
    exit 1
  fi

  if [ "$code" != 200 ]; then
    echo >&2 "Request failed with code $code"
    exit 1
  fi
}

downloadArchive() {
  url="$1"
  destination="$2"

  echo "Fetching $url..."
  if test -x "$(command -v curl)"; then
    code=$(curl -s -w '%{http_code}' -L "$url" -o "$destination/archive.tar.gz")
    # elif test -x "$(command -v wget)"; then
    #   code=$(wget -q -O "$destination" --server-response "$url" 2>&1 | awk '/^  HTTP/{print $2}' | tail -1)
  else
    echo >&2 "Neither curl nor wget was available to perform http requests."
    exit 1
  fi

  if [ "$code" != 200 ]; then
    echo >&2 "Request failed with code $code"
    exit 1
  fi

  # echo "Extracting archive..."
  tar -xf "$destination/archive.tar.gz" -C "$destination"
  rm "$destination/archive.tar.gz"
}

# TODO: use $(brew --prefix) to find install directory?
findBinDirectory() {
  if [ -z "$INSTALL_DIRECTORY" ]; then
    echo >&2 "Installation could not determine your \$INSTALL_DIRECTORY."
    exit 1
  fi
  if [ ! -d "$INSTALL_DIRECTORY" ]; then
    echo >&2 "Installation requires your INSTALL_DIRECTORY directory $INSTALL_DIRECTORY to exist. Please create it."
    exit 1
  fi
}

initArch() {
  ARCH=$(uname -m)
  if [ -n "$DOT_ARCH" ]; then
    echo "Using DOT_ARCH"
    ARCH="$DOT_ARCH"
  fi
  case $ARCH in
    amd64) ARCH="amd64" ;;
    x86_64) ARCH="amd64" ;;
    i386) ARCH="386" ;;
    ppc64) ARCH="ppc64" ;;
    ppc64le) ARCH="ppc64le" ;;
    aarch64) ARCH="arm64" ;;
    # ?) ARCH="armv6" ;;
    *)
      echo >&2 "Architecture ${ARCH} is not supported by this installation script"
      exit 1
      ;;
  esac
  echo "ARCH = $ARCH"
}

initOS() {
  OS=$(uname | tr '[:upper:]' '[:lower:]')
  OS_CYGWIN=0
  if [ -n "$DOT_OS" ]; then
    echo "Using DOT_OS"
    OS="$DOT_OS"
  fi
  case "$OS" in
    darwin) OS='darwin' ;;
    linux) OS='linux' ;;
    freebsd) OS='freebsd' ;;
    mingw*) OS='windows' ;;
    msys*) OS='windows' ;;
    cygwin*)
      OS='windows'
      OS_CYGWIN=1
      ;;
    *)
      echo >&2 "OS ${OS} is not supported by this installation script"
      exit 1
      ;;
  esac
  echo "OS = $OS"
}

# identify platform based on uname output
initArch
initOS

# Determine the location if required
findBinDirectory
echo "Will install into $INSTALL_DIRECTORY"

# assemble expected release artifact name
if [ "${OS}" != "linux" ] && { [ "${ARCH}" = "ppc64" ] || [ "${ARCH}" = "ppc64le" ]; }; then
  # ppc64 and ppc64le are only supported on Linux.
  echo "${OS}-${ARCH} is not supported by this instalation script"
else
  BINARY="dot-${OS}-${ARCH}"
fi

# # add .exe if on windows
# if [ "$OS" = "windows" ]; then
#   BINARY="$BINARY.exe"
# fi

# if DOT_RELEASE_TAG was not provided, assume latest
if [ -z "$DOT_RELEASE_TAG" ]; then
  downloadJSON LATEST_RELEASE "$RELEASES_URL/latest"
  DOT_RELEASE_TAG=$(echo "${LATEST_RELEASE}" | tr -s '\n' ' ' | sed 's/.*"tag_name":"//' | sed 's/".*//')
fi
echo "Release Tag = $DOT_RELEASE_TAG"

# fetch the real release data to make sure it exists before we attempt a download
downloadJSON RELEASE_DATA "$RELEASES_URL/tag/$DOT_RELEASE_TAG"

# BINARY_URL="$RELEASES_URL/download/$DOT_RELEASE_TAG/$BINARY"
# DOWNLOAD_FILE=$(mktemp)

# downloadFile "$BINARY_URL" "$DOWNLOAD_FILE"

# echo "Setting executable permissions."
# chmod +x "$DOWNLOAD_FILE"

ARCHIVE_URL="$RELEASES_URL/download/$DOT_RELEASE_TAG/$BINARY.tar.gz"
DOWNLOAD_DIR=$(mktemp -d)

downloadArchive "$ARCHIVE_URL" "$DOWNLOAD_DIR"

echo "Setting executable permissions."
chmod +x "$DOWNLOAD_DIR/dot"

INSTALL_NAME="dot"

if [ "$OS" = "windows" ]; then
  INSTALL_NAME="$INSTALL_NAME.exe"
fi

echo "Moving executable to $INSTALL_DIRECTORY/$INSTALL_NAME"
mv "$DOWNLOAD_DIR/dot" "$INSTALL_DIRECTORY/$INSTALL_NAME"
