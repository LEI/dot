FROM golang

RUN apt-get update -qq --force-yes && \
apt-get install -qq --no-install-suggests --no-install-recommends --force-yes \
ca-certificates \
curl \
git-core \
tmux \
vim
# && rm -rf /var/lib/apt/lists/*

# ENV GOPATH /go
# PATH="$GOPATH/bin:/usr/local/go/bin:$PATH"

RUN printf "%s\n" \
'sep() { printf %${COLUMNS:-100}s |tr " " "${1:-=}"; printf "\n"; }' \
'log() { sep "-"; printf "\n\t%s\n\n" "$@"; sep "-"; }' \
'run() { log "\$ $*"; "$@" || exit $?; }' >> ~/.bashrc

ENV DOT $GOPATH/src/github.com/LEI/dot
WORKDIR $DOT

ENTRYPOINT ["/bin/bash"]

ADD . $DOT

# RUN go version
# RUN go vet
# RUN go build
# RUN go clean
RUN go get github.com/LEI/dot && go install
