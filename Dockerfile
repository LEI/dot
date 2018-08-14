FROM golang AS base

RUN apt-get update -qqy && \
apt-get install --no-install-suggests --no-install-recommends -qqy \
# ca-certificates \
# curl \
# git \
locales \
sudo

# # ENV GOPATH /go
ENV PATH $PATH:$GOPATH/bin
# ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
ENV DOT $GOPATH/src/github.com/LEI/dot

RUN printf "%s\n" \
'PATH="$GOPATH/bin:$PATH"' \
>> ~/.profile # ~/.bashrc

ENTRYPOINT ["/bin/bash"]
# ENTRYPOINT ["scripts/install.sh"]

WORKDIR $DOT

COPY . .

RUN ./scripts/setup-lang.sh
ENV LANG en_GB.UTF-8
ENV LANGUAGE en_GB:en
ENV LC_ALL en_GB.UTF-8

RUN if [ -d vendor ]; then make install; else make; fi

# RUN cp .dotrc.yml /root/
