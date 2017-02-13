FROM debian:jessie

RUN apt-get update -qq --force-yes && \
apt-get install -qq --no-install-suggests --no-install-recommends --force-yes \
ca-certificates \
curl \
git-core \
golang \
tmux \
vim
# && rm -rf /var/lib/apt/lists/*

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
RUN printf "%s\n" \
'PATH="$GOPATH/bin:/usr/local/go/bin:$PATH"' \
'sep() { printf %${COLUMNS:-100}s |tr " " "${1:-=}"; printf "\n"; }' \
'run() { sep "-"; printf "\n\t%s\n\n" "\$ $*"; sep "-"; "$@" || exit 1; }' >> ~/.bashrc

ENV DOT /go/src/github.com/LEI/dot
WORKDIR $DOT

ENTRYPOINT ["/bin/bash"]

ADD . $DOT

# RUN go vet
# RUN go build
RUN go get github.com/LEI/dot && go install
# , "-s", "$DOT"
