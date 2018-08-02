# vi: ft=Dockerfile
FROM debian

ENV OS linux
ENV ARCH amd64
ENV USER docker
ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update -qq -y && \
apt-get install -qq --no-install-suggests --no-install-recommends -y \
apt-utils \
ca-certificates \
git \
locales \
sudo

COPY ./dist/${OS}_${ARCH}/dot /usr/local/bin/dot
COPY ./scripts /tmp/scripts

RUN /tmp/scripts/setup-lang.sh
ENV LANG en_GB.UTF-8
ENV LANGUAGE en_GB:en
ENV LC_ALL en_GB.UTF-8
# LC_ALL=$LANG

RUN /tmp/scripts/setup-user.sh $USER --password '' --groups staff


USER $USER

WORKDIR /home/$USER

COPY ./.dot.yml .dot.yml

# Restore interactive
ENV DEBIAN_FRONTEND newt

#RUN dot install --packages --sudo

ENTRYPOINT ["/bin/bash"]
