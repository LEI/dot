FROM alpine

ENV OS linux
ENV ARCH amd64
ENV USER docker

#RUN echo "ipv6" >> /etc/modules
RUN apk add --update --no-cache --quiet \
bash \
# ca-certificates \
curl \
# coreutils \
git \
# locales \
shadow \
sudo

COPY ./dist/${OS}_${ARCH}/dot /usr/local/bin/dot
COPY ./scripts /tmp/scripts

RUN /tmp/scripts/setup-user.sh $USER --password ''
RUN /tmp/scripts/setup-pacapt.sh

USER $USER

WORKDIR /home/$USER

COPY ./.dot.yml .dot.yml

# RUN touch /home/$USER/.bashrc
RUN touch /$HOME/.bashrc

RUN which bash

#RUN dot install --packages --sudo

ENTRYPOINT ["/bin/bash"]
