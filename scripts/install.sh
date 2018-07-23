#!/bin/sh

# https://raw.githubusercontent.com/LEI/dot/master/scripts/install.sh
# https://git.io/dot.lei.sh
# curl -sSL https://dot.lei.sh | sh
# wget -qO- https://dot.lei.sh | sh

# https://github.com/rootbeersoup/get.darryl.sh/blob/master/scripts/dotfiles.sh
# https://github.com/golang/dep/blob/master/install.sh
# https://get.docker.com/

set -e

DOT_REPO="github.com/LEI/dot"
JOBBER_REPO="github.com/LEI/jobber"
# JOBBER_REPO="github.com/dshearer/jobber"
# JOBBER_BRANCH="v1.3.2" # undefined sudo_cmd
JOBBER_VER="v1.3.1"
JOBBER_SRC_HASH="8d8cdeb941710e168f8f63abbfc06aab2aadfdfc22b3f6de7108f56403860476"

has() {
  hash "$1" 2> /dev/null
}

check_go() {
  if ! has go; then
    echo "Go command is unavailable"
    exit 1
  fi
  if [ ! -n "$GOPATH" ]; then
    echo "GOPATH is not set"
    exit 1
  fi
}

check_jobber() {
  # if ! has systemd && ! has daemonize; then
  #   if has brew; then
  #     brew install daemonize
  #   else
  #     echo "Please haz services"
  #   fi
  # fi
  if has jobber; then
    return 0
  fi
  # https://dshearer.github.io/jobber/download/
  if [ ! -d "$GOPATH/src/$JOBBER_REPO" ]; then
    # go get -u "$JOBBER_REPO"
    git clone "https://$JOBBER_REPO.git" "$GOPATH/src/$JOBBER_REPO"
  fi
  cd "$GOPATH/src/$JOBBER_REPO"
  make check
  make install # DESTDIR=/usr/local...
  # sudo make install

  # echo "Starting daemon..."
  # daemonize \
  #   -c /usr/local/ \
  #   -e /usr/local/var/log/jobber/error.log \
  #   -o /usr/local/var/log/jobber/output.log \
  #   -a -v "$GOPATH/bin/jobbermaster"

  # # https://github.com/dshearer/jobber-docker/blob/master/alpine3.7/Dockerfile
  # wget "https://api.github.com/repos/dshearer/jobber/tarball/${JOBBER_VER}" -O jobber.tar.gz
  # if has sha256sum; then
  #   echo "${SRC_HASH}  jobber.tar.gz" | sha256sum -cw
  # else
  #   echo "${SRC_HASH}  jobber.tar.gz" | shasum -a 256 -cw
  # fi
  # tar xzf *.tar.gz && rm *.tar.gz && mv dshearer-* jobber
  # cd jobber
  # make check
  # make install DESTDIR=/jobber-dist/

  # mkdir /jobber
  # cp /jobber-dist/usr/local/libexec/jobberrunner /jobber/jobberrunner
  # cp /jobber-dist/usr/local/bin/jobber /jobber/jobber
  # PATH=/jobber:${PATH}

  # USERID=100
  # addgroup jobberuser
  # adduser -S -u "${USERID}" -G jobberuser jobberuser
  # mkdir -p "/var/jobber/${USERID}"
  # chown -R jobberuser:jobberuser "/var/jobber/${USERID}"

  # cp jobfile /home/jobberuser/.jobber
  # chown jobberuser:jobberuser /home/jobberuser/.jobber
  # chmod 0600 /home/jobberuser/.jobber
}

check_dot() {
  if [ ! -d "$GOPATH/src/$DOT_REPO" ]; then
    # git clone https://$DOT_REPO.git
    # Use --recursive for .gitmodules
    go get "$DOT_REPO"
  fi
  if ! has dot; then
    go install "$DOT_REPO"
  fi
}

do_install() {
  check_go
  # check_jobber
  check_dot
  #dot --dry-run --verbose
  echo "Done"
}

do_install
