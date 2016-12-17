# FROM gliderlabs/alpine:3.4
FROM debian:jessie
RUN apt-get update -y && apt-get install -y dialog git-core stow
# RUN apt-get update -y && apt-get install -y locales && rm -rf /var/lib/apt/lists/* \
# && localedef -i en_GB -c -f UTF-8 -A /usr/share/locale/locale.alias en_GB.UTF-8
# RUN git clone https://github.com/LEI/termux-config.git "$HOME/.dotfiles"
ENV BOOTSTRAP /root/.dotfiles
COPY . $BOOTSTRAP
# RUN echo "alias dot='source $BOOTSTRAP/bootstrap'" >> ~/.bashrc
ENTRYPOINT ["/bin/bash"]
CMD ["-c", "source $BOOTSTRAP/bootstrap; bash -l"]
# ARG CACHEBUST=0
# CMD ["-c", "cd $BOOTSTRAP; git pull origin master; bash; dot"]
