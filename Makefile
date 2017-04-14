MAKEFLAGS += --warn-undefined-variables
SHELL := /bin/bash
.SHELLFLAGS := -o pipefail -euc
.DEFAULT_GOAL := build

include Makefile.variables

.PHONY: help
help:
	@echo 'Management commands for cicdtest:'
	@echo
	@echo 'Usage:'
	@echo '  ## Build Commands'
	@echo '    make tag-build       Add git tag for latest build.'
	@echo
	@echo '  ## Generator Commands'
	@echo '    make generate        Run code generator for project.'
	@echo
	@echo '  ## Develop / Test Commands'
	@echo '    make vendor          Install dependencies using glide.'
	@echo '    make format          Run code formatter.'
	@echo '    make check           Run static code analysis (lint).'
	@echo '    make test            Run tests on project.'
	@echo '    make cover           Run tests and capture code coverage metrics on project.'
	@echo '    make clean           Clean the directory tree of produced artifacts.'
	@echo
	@echo '  ## Utility Commands'
	@echo '    make setup           Configures Minishfit/Docker directory mounts.'
	@echo


.PHONY: clean
clean:
	@rm -rf bin cover *.out *.xml

veryclean: clean
	rm -rf tmp
	${DOCKER} rmi -f ${DEV_IMAGE} > /dev/null 2>&1 || true

## prefix before other make targets to run in your local dev environment
local: | quiet
	@$(eval DOCKRUN= )
quiet: # this is silly but shuts up 'Nothing to be done for `local`'
	@:

prepare: tmp/dev_image_id
tmp/dev_image_id:
	@mkdir -p tmp
	@${DOCKER} rmi -f ${DEV_IMAGE} > /dev/null 2>&1 || true
	${DOCKER} build -t ${DEV_IMAGE} -f Dockerfile.dev .
	@${DOCKER} inspect -f "{{ .ID }}" ${DEV_IMAGE} > tmp/dev_image_id

# ----------------------------------------------
# build
.PHONY: build
build: build/dev

.PHONY: build/dev
build/dev: check */*.go *.go
	@mkdir -p bin/
	${DOCKRUN} go build -o bin/swagger --ldflags "$(LDFLAGS)" ./cmd/swagger

.PHONY: install
install: build
	${DOCKRUN} go install ./cmd/swagger

.PHONY: vendor
vendor: tmp/dev_image_id

# ----------------------------------------------
# develop and test

.PHONY: format
format: vendor
	${DOCKRUN} bash ./hack/format.sh

.PHONY: check
check: format
	${DOCKRUN} bash ./hack/check.sh

.PHONY: test
test: check
	${DOCKRUN} bash ./hack/test.sh

.PHONY: cover
cover: check
	@rm -rf cover/
	@mkdir -p cover
	${DOCKRUN} bash ./hack/cover.sh

# does not work
.PHONY: canary
canary: build
	${DOCKRUN} bash -c 'SWAGGER_BIN=$$(pwd)/bin/swagger ./hack/run-canary.sh'

# generate bindata when templates are updated
.PHONY: generate
generate: build/image_build
	${DOCKRUN} go generate ./generator

# ----------------------------------------------
# utilities

.PHONY: setup
setup:
	@bash ./hack/setup.sh
