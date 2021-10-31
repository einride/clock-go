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

include tools/git-verify-nodiff/rules.mk
include tools/golangci-lint/rules.mk
include tools/prettier/rules.mk
include tools/goreview/rules.mk
include tools/mockgen/rules.mk

.PHONY: go-test
go-test:
	go test -v ./...

# go-mod-tidy: update Go module files
.PHONY: go-mod-tidy
go-mod-tidy:
	go mod tidy -v

# mockgen-generate: generate Go mocks
.PHONY: mockgen-generate
mockgen-generate: pkg/mockclock/clock.go

pkg/mockclock/clock.go: pkg/clock/clock.go
	$(mockgen) \
		-destination $@ \
		-package mockclock \
		github.com/einride/clock-go/pkg/clock Clock,Ticker
