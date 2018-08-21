.PHONY: all
all: dot

RUN := go run build.go

# go build -o bin/$(BINARY) cmd/dot
.PHONY: dot
dot:
	$(RUN) vendor check install

# go test ./cmd/... ./internal/...
.PHONY: test
test:
	$(RUN) -v test
