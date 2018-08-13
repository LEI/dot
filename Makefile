# https://github.com/moby/moby/blob/master/Makefile
# https://github.com/vincentbernat/hellogopher/blob/master/Makefile
# https://sahilm.com/makefiles-for-golang/

# SHELL := /bin/sh
PROJECT := github.com/LEI/dot
PACKAGES := $(shell go list ./... | grep -v /vendor)
EXECUTABLE := dot # BINARY

GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)

GO_TEST_VERBOSE ?= 0
GO_VET_VERBOSE ?= 0
GOLINT_MIN_CONFIDENCE ?= 1

# .DEFAULT_GOAL := default
.PHONY: default
default: ensure test install

# .PHONY: check
# check: test vet lint format

# .PHONY: fix
# fix: fmt

.PHONY: dep
dep:
DEP := $(shell command -v dep 2> /dev/null)
ifndef DEP
	curl -sSL https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif

.PHONY: ensure
ensure:
	make dep
	dep ensure

.PHONY: test
test: format $(PACKAGES)

$(PACKAGES):
ifeq ($(GO_TEST_VERBOSE),1)
	go test -v $@
else
	go test $@
endif
ifeq ($(GO_VET_VERBOSE),1)
	go vet -v $@
else
	go vet $@
endif
GOLINT := $(shell command -v golint 2> /dev/null)
ifndef GOLINT
	go get -u golang.org/x/lint/golint
endif
	golint -set_exit_status $@ # ./...

# .PHONY: lint
# lint:
# 	golint -set_exit_status -min_confidence=$(GOLINT_MIN_CONFIDENCE) $$(go list ./...)

.PHONY: goimports
goimports:
GOIMPORTS := $(shell command -v goimports 2> /dev/null)
ifndef GOIMPORTS
	go get -u golang.org/x/tools/cmd/goimports
endif
	goimports

.PHONY: format
format:
	test -z $(gofmt -s -l $GO_FILES)

# .PHONY: fmt
# fmt:
# 	# gofmt -s -w .
# 	go fmt ./...

# .PHONY: build
# build:
# 	# go build $(PROJECT)
# 	go build -o bin/$(EXECUTABLE) main.go

.PHONY: install
install:
	go install $(PROJECT)

.PHONY: goreleaser
REPO_GORELEASER := github.com/goreleaser/goreleaser
# git clone https://$(REPO_GORELEASER).git "$$GOPATH/src/$(REPO_GORELEASER)"
goreleaser:
GORELEASER := $(shell command -v goreleaser 2> /dev/null)
ifndef GORELEASER
	go get -d $(REPO_GORELEASER)
	cd "$$GOPATH/src/$(REPO_GORELEASER)"; \
		dep ensure -vendor-only; \
		make setup build
	go install $(REPO_GORELEASER)
endif

.PHONY: snapshot
# curl -sL https://git.io/goreleaser | bash --rm-dist --snapshot
snapshot:
	make goreleaser
	goreleaser --rm-dist --snapshot

.phony: release
release:
	make goreleaser
	goreleaser release --help

# .PHONY: docker-test
# docker-test:
# 	make snapshot
# 	docker-compose build test
# 	docker-compose run test

# .PHONY: docker-test-os
# OS := alpine
# docker-test-os:
# 	make snapshot
# 	OS=$(OS) docker-compose build test_os
# 	OS=$(OS) docker-compose run test_os

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

# ci: clean dependencies build test

# .PHONY: test
