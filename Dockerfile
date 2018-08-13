FROM golang AS base

# RUN apt-get update -qqy && \
# apt-get install --no-install-suggests --no-install-recommends -qqy \
# ca-certificates \
# curl \
# git \
# locales \
# sudo

# # https://stackoverflow.com/q/28405902/7796750
# RUN sed -i -e 's/# en_GB.UTF-8 UTF-8/en_GB.UTF-8 UTF-8/' /etc/locale.gen && \
# locale-gen
# ENV LANG en_GB.UTF-8
# ENV LANGUAGE en_GB:en
# ENV LC_ALL en_GB.UTF-8

# # ENV GOPATH /go
# ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
ENV DOT $GOPATH/src/github.com/LEI/dot

WORKDIR $DOT

# RUN printf "%s\n" \
# 'PATH="$GOPATH/bin:$PATH"' \
# >> ~/.bashrc

ENTRYPOINT ["/bin/bash"]
# ENTRYPOINT ["scripts/install"]

COPY . $DOT

RUN if [ -d vendor ]; then make install; else make; fi

# RUN cp .dotrc.yml /root/
