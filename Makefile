# $(MAKECMDGOALS)

.PHONY: all
all: dot

RUN := go run build.go

# go build -o bin/$(BINARY) cmd/dot
.PHONY: dot
dot:
	$(RUN) check install

# .PHONY: vendor
# vendor:
# 	$(RUN) vendor -only

.PHONY: test
test:
	# go test ./cmd/... ./internal/...
	$(RUN) -v test

.PHONY: coverage
coverage:
	$(RUN) test:coverage

.PHONY: integration
integration:
	$(RUN) test:integration

.PHONY: docker
docker:
	OS=$(OS) $(RUN) $@

COMMANDS := build check docs install release snapshot

.PHONY: $(COMMANDS)
$(COMMANDS):
	$(RUN) $@

# MAKEFLAGS += --silent
ifndef VERBOSE
.SILENT:
endif

# # https://www.gnu.org/software/make/manual/make.html#Last-Resort
# .DEFAULT:

# %::
# 	$(RUN) $@
