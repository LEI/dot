# FROM gliderlabs/alpine:3.4
FROM debian:jessie
RUN apt-get update -y \
    && apt-get install -y apt-utils dialog git-core stow locales \
    && rm -rf /var/lib/apt/lists/* \
    && localedef -i en_GB -c -f UTF-8 -A /usr/share/locale/locale.alias en_GB.UTF-8
# RUN git clone https://github.com/LEI/termux-config.git "$HOME/.dotfiles"
ENV BOOTSTRAP /root/.dotfiles
ENV GIT_AUTHOR_NAME "John Doe"
ENV GIT_AUTHOR_USERNAME "JD"
ENV GIT_AUTHOR_EMAIL "j@d.c"
# ARG CACHEBUST=1
COPY . $BOOTSTRAP
# RUN echo "alias dot='source $BOOTSTRAP/bootstrap'" >> ~/.bashrc
ENTRYPOINT ["/bin/bash"]
CMD ["-l", "-c", "source $BOOTSTRAP/bootstrap; bash -l"]
# CMD ["-c", "cd $BOOTSTRAP; git pull origin master; bash; dot"]
