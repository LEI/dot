FROM debian:jessie

# RUN apt-get update -qy && \
# apt-get install -qy --no-install-suggests --no-install-recommends --force-yes \
# apt-utils locales \
# && rm -rf /var/lib/apt/lists/* \
# && localedef -i en_GB -c -f UTF-8 -A /usr/share/locale/locale.alias en_GB.UTF-8
RUN apt-get update -qy && \
apt-get install -qy --no-install-suggests --no-install-recommends --force-yes \
curl \
git-core \
tmux \
vim

# ENV LANG en_GB.UTF-8
# RUN echo "$LANG UTF-8" > /etc/locale.gen && locale-gen
# RUN echo "LANG=$LANG" > /etc/locale.conf
ENV HOME /root
ENV DOT $HOME/.dot
WORKDIR $HOME

RUN ln -s "$DOT/bin/dot" "/usr/local/bin/dot"
# RUN printf "%s\n" "alias dsrc=\"dot -s $DOT/.dotrc\"" >> $HOME/.bashrc

ENTRYPOINT ["/bin/bash"]
CMD ["-l", "-c", "dot -R "$DOT""] # ; bash -l

ADD . $DOT
