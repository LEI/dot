FROM debian:jessie

RUN apt-get update -qy && \
apt-get install -qy --no-install-suggests --no-install-recommends --force-yes \
ca-certificates \
curl \
git-core \
golang \
tmux \
vim
# && rm -rf /var/lib/apt/lists/*

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

ENV DOT /go/src/github.com/LEI/dot
WORKDIR $DOT

ENTRYPOINT ["/bin/bash"]

ADD . $DOT
RUN go install # /go/bin/dot
# , "-s", "$DOT"
