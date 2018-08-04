# https://bugs.alpinelinux.org/issues/1792
FROM debian:9-slim AS download

ENV CPU x86_64
ENV DEBIAN_FRONTEND noninteractive
ENV LIBGUESTFS_BACKEND direct

# build-essential
RUN apt-get update -qqy && \
apt-get install --no-install-suggests --no-install-recommends -qqy \
ca-certificates \
curl \
unzip

# ADD doesn't unzip :(

RUN curl -sSL https://dl.google.com/android/repository/sys-img/android/$CPU-24_r07.zip \
-o android.zip && \
unzip android.zip -d android && \
rm android.zip

RUN curl -sSL https://termux.net/bootstrap/bootstrap-$CPU.zip \
-o bootstrap.zip && \
unzip bootstrap.zip -d bootstrap && \
rm bootstrap.zip

FROM download AS builder

# FIXME: https://stackoverflow.com/a/36575437/7796750

RUN apt-get install --no-install-suggests --no-install-recommends -qqy \
linux-image-amd64 \
mount # kmod # busybox
# RUN mkdir system && \

# libguestfs-tools \
# qemu-utils \
# fuse udev # or makedev

# RUN mkdir /dev/fuse && chmod 0777 /dev/fuse
# RUN depmod -F /lib/modules/4.9.0-7-amd64

# RUN echo '#!/bin/sh'"\n"\
# 'if [ $# -eq 1 ] && [ "$1" = "-r" ]; \
# then echo "4.9.0-7-amd64"; \
# else exec uname "$@"; \
# fi' \
# > /usr/local/bin/uname && \
# chmod +x /usr/local/bin/uname

RUN cp -r /lib/modules/4.9.0-7-amd64 /lib/modules/4.9.93-boot2docker

# ENV LIBGUESTFS_DEBUG 1
# ENV LIBGUESTFS_TRACE 1
# RUN mkdir system-readonly && \
# modprobe fuse && \
# guestmount -a android/$CPU/system.img -m /dev/sda system-readonly/ && \
# cp -r system-readonly/ system/ && \
# guestunmount system-readonly/ && \

#fdisk -l android/$CPU/system.img && \
RUN mkdir system && \
modprobe loop && \
ls -la system && \
mount -o loop,offset=$((1*512)) android/$CPU/system.img system && \
echo sys && \
ls -la system
# mcopy -nsv android/$CPU/system.img system/ # mtools

RUN while read -r l; do src="${l%←*}"; dst="${l#*←}"; \
if [ -z "$src" ] || [ -z "$dst" ] || [ "$src" = "$dst" ]; then \
continue; \
fi; \
echo ln -s "bootstrap/$src" "bootstrap/$dst"; \
done < bootstrap/SYMLINKS.txt

RUN ls -la system

FROM scratch
# https://github.com/Rudloff/termux-docker-image

COPY --from=builder system/ /system/
COPY --from=builder bootstrap/ /data/data/com.termux/files/usr/

ENV ANDROID_ROOT /system
ENV TERMUX_ROOT /data/data/com.termux/files/usr
ENV PATH $TERMUX_ROOT/bin/:$TERMUX_ROOT/bin/applets/:$ANDROID_ROOT/bin/:$PATH
ENV LD_LIBRARY_PATH $TERMUX_ROOT/lib/

SHELL ["/data/data/com.termux/files/usr/bin/sh", "-c"]

RUN ls -la /
RUN which sh

# sh -c 'mkdir /bin; ln -s $ANDROID_ROOT/bin/sh /bin/sh'
# RUN chmod +x $TERMUX_ROOT/bin/*
# RUN echo '\n104.18.37.234 termux.net' >> $ANDROID_ROOT/etc/hosts
# RUN chmod +x $TERMUX_ROOT/lib/apt/methods/*
# RUN chmod +x $TERMUX_ROOT/libexec/termux/command-not-found
# RUN echo "mkdir -p /dev/socket/; chmod 777 /dev/socket/; logd &" >> $TERMUX_ROOT/etc/bash.bashrc
# RUN mkdir -p $TERMUX_ROOT
# RUN mkdir -p $TERMUX_ROOT/etc/apt/preferences.d/
# RUN mkdir -p $TERMUX_ROOT/etc/apt/apt.conf.d/
# RUN mkdir -p $TERMUX_ROOT/var/cache/apt/archives/partial/
# RUN mkdir -p $TERMUX_ROOT/var/lib/dpkg/updates/
# RUN mkdir -p /$TERMUX_ROOT/var/log/apt/

# ---

# FROM java:8-alpine
# # FROM java:openjdk-8-jdk-alpine
# # FROM oraclelinux:8-slim

# # https://docs.travis-ci.com/user/languages/android/

# # ENV OS linux
# # ENV ARCH amd64
# # ENV USER docker
# ENV ANDROID_BUILD_TOOLS 27.0.3
# ENV ANDROID_SDK_TOOLS 4333796
# # 92ffee5a1d98d856634e8b71132e8a95d96c83a63fde1099be3d86df3106def9

# ENV ANDROID_HOME /usr/local/android-sdk-linux
# ENV PATH $PATH:$ANDROID_HOME/tools
# ENV PATH $PATH:$ANDROID_HOME/tools/bin
# ENV PATH $PATH:$ANDROID_HOME/platform-tools
# ENV PATH $PATH:$ANDROID_HOME/build-tools/$ANDROID_BUILD_TOOLS
# # ENV PATH ${PATH}:${ANDROID_NDK}

# ENV TERMUX_APP_VERSION 0.65

# RUN apk add --update --no-cache --quiet \
# curl \
# unzip
# # git

# # RUN git clone https://github.com/urho3d/android-ndk.git $HOME/android-ndk
# # ENV ANDROID_NDK_HOME $HOME/android-ndk

# # Download and unzip Android SDK tools
# RUN mkdir -p /usr/local/android-sdk-linux
# RUN curl -sSL https://dl.google.com/android/repository/sdk-tools-linux-${ANDROID_SDK_TOOLS}.zip \
# -o tools.zip && \
# unzip tools.zip -d /usr/local/android-sdk-linux && \
# rm tools.zip

# # RUN mkdir .android && touch .android/repositories.cfg
# RUN mkdir $ANDROID_HOME/licenses && \
# echo 8933bad161af4178b1185d1a37fbf41ea5269c55 > $ANDROID_HOME/licenses/android-sdk-license && \
# echo d56f5187479451eabf01fb78af6dfcb131a6481e >> $ANDROID_HOME/licenses/android-sdk-license && \
# echo 84831b9409646a918e30573bab4c9c91346d8abd > $ANDROID_HOME/licenses/android-sdk-preview-license

# RUN yes | sdkmanager --licenses && \
# sdkmanager --update && \
# sdkmanager \
# "platforms;android-27" \
# "build-tools;${ANDROID_BUILD_TOOLS}" \
# "platform-tools" \
# "tools" \
# "extras;android;m2repository" \
# "ndk-bundle"

# RUN echo $PATH && ls -la $ANDROID_HOME/platforms/android-27

# # gradlew testDebugUnitTest --build-cache --console plain --offline --warning-mode all
# RUN curl -sSL https://github.com/termux/termux-app/archive/v${TERMUX_APP_VERSION}.tar.gz \
# -o termux-app.tar.gz && \
# tar -xzf termux-app.tar.gz && \
# cd termux-app-${TERMUX_APP_VERSION} && \
# echo ./gradlew --quiet testDebugUnitTest && \
# ./gradlew --quiet --offline testDebugUnitTest && \
# echo ./gradlew --quiet clean && \
# ./gradlew --quiet clean && \
# echo ./gradlew --quiet build && \
# ./gradlew --quiet build && \
# ls -la

# ---

#COPY ./dist/${OS}_${ARCH}/dot /usr/local/bin/dot
#COPY ./scripts /tmp/scripts

#RUN /tmp/scripts/setup-user.sh $USER --password ''
#RUN /tmp/scripts/setup-pacapt.sh

#USER $USER

#WORKDIR /home/$USER

#COPY ./.dot.yml .dot.yml

#RUN touch /$HOME/.bashrc

#RUN which bash

##RUN dot install --packages --sudo

#ENTRYPOINT ["/bin/bash"]
