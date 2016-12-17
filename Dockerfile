# FROM gliderlabs/alpine:3.4
FROM debian:jessie
RUN apt-get update && apt-get install -y locales && rm -rf /var/lib/apt/lists/* \
    && localedef -i en_GB -c -f UTF-8 -A /usr/share/locale/locale.alias en_GB.UTF-8
ENV LANG en_GB.utf8
RUN apt-get install git-core
RUN git clone https://github.com/LEI/termux-config.git .dotfiles
# RUN source ~/.dotfiles/bootstrap
ENTRYPOINT ["/bin/bash"]
