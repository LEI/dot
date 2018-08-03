FROM alpine

ENV OS linux
ENV ARCH amd64
ENV USER docker

#RUN echo "ipv6" >> /etc/modules
RUN apk add --update --no-cache \
 bash \
curl \
git \
shadow \
sudo
# ca-certificates \
# locales \

COPY ./dist/${OS}_${ARCH}/dot /usr/local/bin/dot
COPY ./scripts /tmp/scripts

RUN /tmp/scripts/setup-user.sh $USER --password ''

# Pre-install pacapt
# ADD https://github.com/icy/pacapt/raw/ng/pacapt /usr/local/bin/pacapt
RUN curl -sSL https://github.com/icy/pacapt/raw/ng/pacapt \
-o /usr/local/bin/pacapt && \
chmod 0755 /usr/local/bin/pacapt

USER $USER

WORKDIR /home/$USER

COPY ./.dot.yml .dot.yml

# RUN touch /home/$USER/.bashrc
RUN touch /$HOME/.bashrc

RUN which bash

#RUN dot install --packages --sudo

ENTRYPOINT ["/bin/bash"]
