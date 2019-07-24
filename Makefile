# all: run a complete build
all: \
	markdown-lint \
	mocks \
	go-lint \
	go-review \
	go-mod-tidy \
	git-verify-nodiff \
	git-verify-submodules

export GO111MODULE := on

# clean: remove all generated build files
.PHONY: clean
clean:
	rm -rf build

.PHONY: build
build:
	@git submodule update --init --recursive $@

include build/rules.mk
build/rules.mk: build
	@# included in submodule: build

# markdown-lint: lint Markdown files
.PHONY: markdown-lint
markdown-lint: $(PRETTIER)
	$(PRETTIER) --check **/*.md  --parser markdown

# go-mod-tidy: update Go module files
.PHONY: go-mod-tidy
go-mod-tidy:
	go mod tidy -v

# go-lint: lint Go files
.PHONY: go-lint
go-lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run --enable-all

# go-review: review Go files
.PHONY: go-review
go-review: $(GOREVIEW)
	$(GOREVIEW) -c 1 ./...

# mocks: generate Go mocks
.PHONY: mocks
mocks: pkg/mockclock/clock.go

pkg/mockclock/clock.go: pkg/clock/clock.go $(GOBIN)
	$(GOBIN) -m -run github.com/golang/mock/mockgen \
		-destination $@ \
		-package mockclock \
		github.com/einride/clock-go/pkg/clock Clock,Ticker
