#!/bin/bash

set -e

goreleaser --rm-dist --snapshot

docker-compose build

docker-compose run test

ls -la dist
