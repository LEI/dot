# $(MAKECMDGOALS)

.PHONY: all
all: dot

RUN := go run build.go

# go build -o bin/$(BINARY) cmd/dot
.PHONY: dot
dot:
	$(RUN) vendor check install

.PHONY: vendor
vendor:
	$(RUN) vendor -only

.PHONY: check
check:
	$(RUN) check # -v

# go test ./cmd/... ./internal/...
.PHONY: test
test:
	$(RUN) -v test

.PHONY: build
build:
	$(RUN) build

.PHONY: install
install:
	$(RUN) install

.PHONY: docker
docker:
	$(RUN) docker

.PHONY: docs
docs:
	rm -fr ./docs
	$(RUN) docs

.PHONY: release
release:
	$(RUN) release

.PHONY: snapshot
snapshot:
	$(RUN) snapshot

# MAKEFLAGS += --silent
ifndef VERBOSE
.SILENT:
endif

# # https://www.gnu.org/software/make/manual/make.html#Last-Resort
# .DEFAULT:

# %::
# 	$(RUN) $@
