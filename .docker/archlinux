# vi: ft=Dockerfile
FROM base/archlinux

ENV OS linux
ENV ARCH amd64
ENV USER docker

RUN pacman -Syu --needed --noconfirm --noprogressbar --quiet \
base-devel \
ca-certificates \
git \
sudo

COPY ./dist/${OS}_${ARCH}/dot /usr/local/bin/dot
COPY ./scripts /tmp/scripts

RUN /tmp/scripts/setup-lang.sh
ENV LANG en_GB.UTF-8
ENV LANGUAGE en_GB:en
ENV LC_ALL en_GB.UTF-8

RUN groupadd sudo
RUN /tmp/scripts/setup-user.sh --groups sudo --password '' $USER

USER $USER

WORKDIR /home/$USER

COPY ./.dot.yml .dot.yml

#RUN dot install --packages --sudo

ENTRYPOINT ["/bin/bash"]
