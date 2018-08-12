# https://github.com/moby/moby/blob/master/Makefile
# https://github.com/vincentbernat/hellogopher/blob/master/Makefile
# https://sahilm.com/makefiles-for-golang/

BINARY := dot
REPO := github.com/LEI/dot
# PACKAGES := $(shell go list ./...)

GO_TEST_SILENT ?= 1
# GO_TEST_VERBOSE ?= 0
GOLINT_MIN_CONFIDENCE ?= 1

# .DEFAULT_GOAL := default
.PHONY: default
default: dep ensure check install

# .PHONY: all
# all: dep ensure check fix install

.PHONY: check
check: test vet lint format

.PHONY: fix
check: fmt

.PHONY: release
release: goreleaser publish

.PHONY: dep
DEP := $(shell command -v dep 2> /dev/null)
dep:
ifndef DEP
	curl -sSL https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif

.PHONY: ensure
ensure:
	dep ensure

.PHONY: test
test:
ifeq ($(GO_TEST_SILENT),1)
	go test ./...
else
	go test -v ./...
endif

.PHONY: vet
vet:
	go vet ./...

.PHONY: lint
lint:
	golint -set_exit_status -min_confidence=$(GOLINT_MIN_CONFIDENCE) $$(go list ./...)

.PHONY: format
format:
	test -z $(gofmt -s -l $GO_FILES)

.PHONY: fmt
# gofmt -s -w .
fmt:
	go fmt ./...

# .PHONY: build
# build:
# 	# go build $(REPO)
# 	go build -o bin/$(BINARY) main.go

.PHONY: install
install:
	go install $(REPO)

.PHONY: goreleaser
REPO_GORELEASER := github.com/goreleaser/goreleaser
GORELEASER := $(shell command -v goreleaser 2> /dev/null)
# git clone https://$(REPO_GORELEASER).git "$$GOPATH/src/$(REPO_GORELEASER)"
goreleaser:
ifndef GORELEASER
	go get -d $(REPO_GORELEASER)
	cd "$$GOPATH/src/$(REPO_GORELEASER)"; \
		dep ensure -vendor-only; \
		make setup build
	go install $(REPO_GORELEASER)
endif

# .PHONY: tag
# GORELEASER_SNAPSHOT ?= 1
# tag:
# ifeq ($(GORELEASER_SNAPSHOT),1)
# 	goreleaser --snapshot --rm-dist
# else
# 	goreleaser --rm-dist --help
# endif

.PHONY: snapshot
# curl -sL https://git.io/goreleaser | bash --rm-dist --snapshot
snapshot:
	goreleaser --rm-dist --snapshot

# .PHONY: release
# release:
# 	goreleaser release --help

.PHONY: docker-test
docker-test:
	make snapshot
	docker-compose build test
	docker-compose run test

.PHONY: docker-test-os
OS := alpine
docker-test-os:
	make snapshot
	OS=$(OS) docker-compose build test_os
	OS=$(OS) docker-compose run test_os

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
