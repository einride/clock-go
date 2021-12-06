# all: run a complete build
all: \
	yaml-format \
	markdown-format \
	mockgen-generate \
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

.PHONY: go-test
go-test:
	go test -v ./...

# go-mod-tidy: update Go module files
.PHONY: go-mod-tidy
go-mod-tidy:
	go mod tidy -v

# mockgen-generate: generate Go mocks
.PHONY: mockgen-generate
mockgen-generate: mockclock/clock.go

mockclock/clock.go: clock.go
	go run -mod=mod github.com/golang/mock/mockgen \
		-destination $@ \
		-package mockclock \
		go.einride.tech/clock Clock,Ticker
