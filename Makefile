SHELL := /usr/bin/env bash

# all: run a complete build
all: \
	yaml-format \
	markdown-format \
	go-mod-tidy \
	go-test \
	go-lint \
	go-review \
	go-mod-tidy \
	git-verify-nodiff

include .tools/git-verify-nodiff/rules.mk
include .tools/golangci-lint/rules.mk
include .tools/prettier/rules.mk
include .tools/goreview/rules.mk
include .tools/semantic-release/rules.mk

.PHONY: go-test
go-test:
	go test -v ./...

# go-mod-tidy: update Go module files
.PHONY: go-mod-tidy
go-mod-tidy:
	go mod tidy -v
