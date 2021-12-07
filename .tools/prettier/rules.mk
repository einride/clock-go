PRETTIER_VERSION := 2.3.2
PRETTIER_DIR := $(abspath $(dir $(lastword $(MAKEFILE_LIST))))/bin/$(PRETTIER_VERSION)
PRETTIER := $(PRETTIER_DIR)/node_modules/.bin/prettier

$(PRETTIER):
	$(info [$@] installing...)
	@npm install --no-save --no-audit --prefix $(PRETTIER_DIR) prettier@$(PRETTIER_VERSION)
	@chmod +x $@
	@touch $@

.PHONY: markdown-format
markdown-format: $(PRETTIER)
	$(info [$@] formatting markdown files...)
	@$(PRETTIER) --parser markdown --check *.md --loglevel warn -w

.PHONY: yaml-format
yaml-format: $(PRETTIER)
	$(info [$@] linting yaml files...)
	@$(PRETTIER) --parser yaml --check --write ./**/*.y*ml --loglevel warn