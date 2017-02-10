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
RUN printf "%s\n" 'PATH="$GOPATH/bin:/usr/local/go/bin:$PATH"' >> ~/.profile

ENV DOT /go/src/github.com/LEI/dot
WORKDIR $DOT

ENTRYPOINT ["/bin/bash"]

ADD . $DOT

# RUN go vet
RUN go install
# , "-s", "$DOT"
