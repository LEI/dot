# FROM gliderlabs/alpine:4.4

# FROM debian:jessie
FROM debian:testing
RUN apt-get update -qy && apt-get install -qy apt-utils locales \
&& rm -rf /var/lib/apt/lists/* \
&& localedef -i en_GB -c -f UTF-8 -A /usr/share/locale/locale.alias en_GB.UTF-8
RUN apt-get update -qy && apt-get install -qy git-core stow

# FROM base/archlinux
# RUN pacman --noconfirm -Sy archlinux-keyring
# RUN pacman-db-upgrade
# RUN pacman --noconfirm -Syyu ca-certificates ca-certificates-mozilla \
# && trust extract-compat
# RUN pacman --noconfirm -Sy git stow

# ENV ROOT /root/.dotfiles
ENV GIT_AUTHOR_NAME "John Doe"
ENV GIT_AUTHOR_USERNAME "JD"
ENV GIT_AUTHOR_EMAIL "j@d.c"
# ARG CACHEBUST=1
# COPY . /root/.dotfiles
RUN echo "alias dot='source /root/.dotfiles/bootstrap'" >> ~/.bashrc
ENTRYPOINT ["/bin/bash"]
# CMD ["-l" , "-c", "source", "$BOOTSTRAP/bootstrap"]
