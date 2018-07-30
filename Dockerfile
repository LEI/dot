FROM golang

RUN apt-get update -qq -y && \
apt-get install -qq --no-install-suggests --no-install-recommends -y \
ca-certificates \
curl \
git \
locales # vim tmux
# && rm -rf /var/lib/apt/lists/*

# https://stackoverflow.com/q/28405902/7796750
RUN sed -i -e 's/# en_GB.UTF-8 UTF-8/en_GB.UTF-8 UTF-8/' /etc/locale.gen && \
locale-gen
ENV LANG en_GB.UTF-8
ENV LANGUAGE en_GB:en
ENV LC_ALL en_GB.UTF-8

# ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
ENV DOT $GOPATH/src/github.com/LEI/dot
WORKDIR $DOT

RUN printf "%s\n" \
'PATH="$GOPATH/bin:$PATH"' \
>> ~/.bashrc

ENTRYPOINT ["/bin/bash"]

# Install go dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
# RUN go get gopkg.in/yaml.v2 \
# github.com/imdario/mergo \
# github.com/jessevdk/go-flags

ADD . $DOT

# RUN go vet
# RUN go build
# RUN go clean
# RUN go get github.com/LEI/dot && \
RUN dep ensure
RUN go install

# scripts/realease.sh (goreleaser)
