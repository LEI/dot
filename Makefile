# https://github.com/moby/moby/blob/master/Makefile
# https://github.com/vincentbernat/hellogopher/blob/master/Makefile
# https://sahilm.com/makefiles-for-golang/

BINARY = dot
REPO = github.com/LEI/dot

default: install # build

# PACKAGES := \
# 	github.com/eliasson/foo \
# 	github.com/eliasson/bar
# DEPENDENCIES := github.com/eliasson/acme

# all: build silent-test

# build:
# 	go build -o bin/foo main.go

# test:
# 	go test -v $(PACKAGES)

# silent-test:
# 	go test $(PACKAGES)

# format:
# 	go fmt $(PACKAGES)

build:
	go build $(REPO)

install:
	go install $(REPO)/cmd/$(BINARY)

deps:
	dep ensure

test:
	go test -v ./...

# ci: clean dependencies build test

# .PHONY: test
