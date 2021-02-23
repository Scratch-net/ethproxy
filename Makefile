IMAGE_TAG := ethproxy
SHELL=/bin/bash
OSNAME=$(shell go env GOOS)

.PHONY: all
all: deps lint unit_test build

include linting.mk #linter settings

.PHONY: deps
deps:
	go mod tidy
	go mod download

.PHONY: build
build: deps
	go build -trimpath -a -v -ldflags '-w -s -extldflags "-static"' -tags 'osusergo netgo static_build' -o ./bin/ethproxy

.PHONY: tests
tests:
	go test -v -cover ./... -count=1

.PHONY: dockerise
dockerise:
	docker build -t "${IMAGE_TAG}" .

.PHONY: docker_run
docker_run:
	docker run -p 8080:8080 "${IMAGE_TAG}"
