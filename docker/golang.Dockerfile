FROM golang

RUN apt-get update -qq -y && \
apt-get install -qq --no-install-suggests --no-install-recommends -y \
ca-certificates \
curl \
git \
locales

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

COPY . $DOT

RUN ./scripts/install.sh
