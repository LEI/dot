#!/bin/bash

set -e

# sudo apt-get -qq update
# sudo apt-get install -y vim tmux
#GO_FILES=$(find . -iname '*.go' -type f | grep -v /vendor/) # All the .go files, excluding vendor/
# go get ./...
go get github.com/golang/lint/golint
# go get github.com/mattn/goveralls

#test -z $(gofmt -s -l $GO_FILES)
go test -v ./...
go vet ./...
golint -set_exit_status $(go list ./...)
# go install
# yes | dot --https
# # curl -sfLo ~/.vim/autoload/plug.vim --create-dirs 'https://raw.githubusercontent.com/junegunn/vim-plug/master/plug.vim'
# tmux new-session -n test "vim -E -s -u $HOME/.vim/vimrc +PlugInstall +qall; exit"
# test -d $HOME/.tmux/plugins/tpm
# test -d $HOME/.vim/plugged
# go test -v ./.. -covermode=count -coverprofile=coverage.out
# goveralls -coverprofile=coverage.out -service=travis-ci -repotoken $COVERALLS_TOKEN
