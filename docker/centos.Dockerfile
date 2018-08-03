# vi: ft=Dockerfile
FROM centos

ENV OS linux
ENV ARCH amd64
ENV USER docker

# RUN printf "errorlevel=0\nrpmverbosity=critical\n" >> /etc/yum.conf

# yum groupinstall --assumeyes --quiet "Development Tools" && \
# yum install --assumeyes --quiet \
# gettext-devel openssl-devel perl-CPAN perl-devel zlib-devel \
# ca-certificates \
# epel-release

# Git requires:
#   curl-devel expat-devel gettext-devel openssl-devel zlib-devel
#   gcc perl-ExtUtils-MakeMaker
# Tmux requires:
#   # libevent file
#   kernel-devel ncurses-devel
# autotools
#   autoconf make
RUN yum update --assumeyes --quiet && \
yum install --assumeyes --quiet \
curl-devel expat-devel gettext-devel openssl-devel zlib-devel \
gcc perl-ExtUtils-MakeMaker \
kernel-devel ncurses-devel file \
autoconf make \
sudo

# Install latest git version instead of 1.8
ENV GIT_VERSION 2.18.0
# ADD https://github.com/git/git/archive/v$GIT_VERSION.tar.gz git.tar.gz
RUN curl -sSL https://github.com/git/git/archive/v$GIT_VERSION.tar.gz \
-o git.tar.gz && \
tar -zxf git.tar.gz && \
cd git-$GIT_VERSION && \
make --quiet configure && \
./configure --prefix=/usr/local --quiet && \
sudo make --quiet install

# Install libevent
ENV LIBEVENT_VERSION 2.1.8
RUN curl -sSL https://github.com/libevent/libevent/releases/download/release-$LIBEVENT_VERSION-stable/libevent-$LIBEVENT_VERSION-stable.tar.gz \
-o libevent.tar.gz && \
tar -xzf libevent.tar.gz && \
cd libevent-$LIBEVENT_VERSION-stable && \
./configure --prefix=/usr/local --quiet && \
make --quiet && \
sudo make --quiet install

# Install latest tmux version instead of 1.8
ENV TMUX_VERSION 2.7
# RUN git clone https://github.com/tmux/tmux.git && \
RUN curl -sSL https://github.com/tmux/tmux/releases/download/$TMUX_VERSION/tmux-$TMUX_VERSION.tar.gz \
-o tmux.tar.gz && \
tar -zxf tmux.tar.gz && \
cd tmux-$TMUX_VERSION && \
LDFLAGS="-L/usr/local/lib -Wl,-rpath=/usr/local/lib" ./configure --prefix=/usr/local --quiet && \
make --quiet && \
sudo make --quiet install

COPY ./dist/${OS}_${ARCH}/dot /usr/local/bin/dot
COPY ./scripts /tmp/scripts

# RUN /tmp/scripts/setup-lang.sh
# ENV LANG en_GB.UTF-8
# ENV LANGUAGE en_GB:en
# ENV LC_ALL en_GB.UTF-8

RUN groupadd sudo
RUN /tmp/scripts/setup-user.sh --groups sudo --password '' $USER
# Add /usr/local/bin to sudo PATH
#sed -e 's#Defaults    secure_path = /sbin:/bin:/usr/sbin:/usr/bin#Defaults    secure_path = /sbin:/bin:/usr/sbin:/usr/bin:/usr/local/bin#' /etc/sudoers
#sed -e 's#Defaults[[:blank:]]+secure_path = /sbin:/bin:/usr/sbin:/usr/bin#Defaults[[:blank:]]+secure_path = /sbin:/bin:/usr/sbin:/usr/bin:/usr/local/bin#' /etc/sudoers
# Deb: Defaults        secure_path="/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"
# sed -e '/secure_path/ s[=.*[&:/usr/local/bin[' /etc/sudoers
# sed -r -e '/^\s*Defaults\s+secure_path/ s[=(.*)[=\1:/usr/local/bin[' /etc/sudoers
# RUN echo 'Defaults secure_path="<default value>:/usr/local/bin"' >> "/etc/sudoers.d/$USER"
RUN sed -i -e 's#Defaults    secure_path = /sbin:/bin:/usr/sbin:/usr/bin#Defaults    secure_path = /sbin:/bin:/usr/sbin:/usr/bin:/usr/local/bin#' /etc/sudoers
RUN /tmp/scripts/setup-pacapt.sh

USER $USER

WORKDIR /home/$USER

COPY ./.dot.yml .dot.yml

#RUN dot install --packages --sudo

ENTRYPOINT ["/bin/bash"]
