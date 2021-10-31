mockgen_cwd := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))
mockgen := $(mockgen_cwd)/bin/mockgen
mockgen_version := 1.5.0

go_version := $(shell go version | sed -E "s/^go version go(.\...).*/\1/")

$(mockgen): $(mockgen_cwd)/go.mod
	$(info [mockgen] building...)
ifeq ($(go_version),1.15)
	cd $(mockgen_cwd) && go build -o $@ github.com/golang/mock/mockgen
else
	GOBIN=$(dir $(mockgen)) go install github.com/golang/mock/mockgen@v$(mockgen_version)
endif
