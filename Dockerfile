FROM golang

RUN apt-get update -qq --force-yes && \
apt-get install -qq --no-install-suggests --no-install-recommends --force-yes \
ca-certificates \
curl \
git-core \
tmux \
vim-nox
# && rm -rf /var/lib/apt/lists/*

# ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
ENV DOT $GOPATH/src/github.com/LEI/dot
WORKDIR $DOT

RUN printf "%s\n" \
'PATH="$GOPATH/bin:$PATH"' \
>> ~/.bashrc

ENTRYPOINT ["/bin/bash"]

RUN go get github.com/spf13/cobra
RUN go get github.com/spf13/viper
RUN go get github.com/boltdb/bolt
# RUN go get -u gopkg.in/src-d/go-git.v4/...

ADD . $DOT

# RUN go vet
# RUN go build
# RUN go clean
RUN go get github.com/LEI/dot \
&& go install
