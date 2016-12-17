# FROM gliderlabs/alpine:3.4
FROM debian:jessie

RUN apt-get update -y && apt-get install -y locales && rm -rf /var/lib/apt/lists/* \
    && localedef -i en_GB -c -f UTF-8 -A /usr/share/locale/locale.alias en_GB.UTF-8
ENV LANG en_GB.utf8

COPY . /root/.dotfiles
# RUN apt-get install -y git-core
# RUN git clone https://github.com/LEI/termux-config.git "$HOME/.dotfiles"
ENV BOOTSTRAP /root/.dotfiles

ENTRYPOINT ["bash"]
# ARG CACHEBUST=1
CMD ["-c", "cd $BOOTSTRAP; git pull origin master"]
# "-c", "source $HOME/.dotfiles/bootstrap"
