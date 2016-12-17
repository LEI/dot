FROM gliderlabs/alpine:3.4
# RUN apk update
RUN apk add --no-cache bash git
RUN git clone https://github.com/LEI/termux-config.git .dotfiles
RUN source ~/.dotfiles/bootstrap
ENTRYPOINT ["/bin/bash"]
