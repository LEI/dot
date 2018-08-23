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

# go test ./cmd/... ./internal/...
.PHONY: test
test:
	$(RUN) -v test

.PHONY: build
build:
	$(RUN) build

.PHONY: docker
docker:
	$(RUN) docker

# MAKEFLAGS += --silent
ifndef VERBOSE
.SILENT:
endif

# # https://www.gnu.org/software/make/manual/make.html#Last-Resort
# .DEFAULT:

# %::
# 	$(RUN) $@
