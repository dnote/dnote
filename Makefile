NPM := $(shell command -v npm 2> /dev/null)
GH := $(shell command -v gh 2> /dev/null)

currentDir = $(shell pwd)
serverOutputDir = ${currentDir}/build/server
cliOutputDir = ${currentDir}/build/cli
cliHomebrewDir = ${currentDir}/../homebrew-dnote

## installation
install: install-go install-js
.PHONY: install

install-go:
	@echo "==> installing go dependencies"
	@go mod download
.PHONY: install-go

install-js:
ifndef NPM
	$(error npm is not installed)
endif

	@echo "==> installing js dependencies"

ifeq ($(CI), true)
	@(cd ${currentDir}/pkg/server/assets && npm ci --cache $(NPM_CACHE_DIR) --prefer-offline --unsafe-perm=true)
else
	@(cd ${currentDir}/pkg/server/assets && npm install)
endif
.PHONY: install-js

## test
test: test-cli test-api test-e2e
.PHONY: test

test-cli: generate-cli-schema
	@echo "==> running CLI test"
	@(${currentDir}/scripts/cli/test.sh)
.PHONY: test-cli

test-api:
	@echo "==> running API test"
	@(${currentDir}/scripts/server/test-local.sh)
.PHONY: test-api

test-e2e:
	@echo "==> running E2E test"
	@(${currentDir}/scripts/e2e/test.sh)
.PHONY: test-e2e

# development
dev-server:
	@echo "==> running dev environment"
	@VERSION=master ${currentDir}/scripts/server/dev.sh
.PHONY: dev-server

build-server:
ifndef version
	$(error version is required. Usage: make version=0.1.0 build-server)
endif

	@echo "==> building server assets"
	@(cd "${currentDir}/pkg/server/assets/" && ./styles/build.sh)
	@(cd "${currentDir}/pkg/server/assets/" && ./js/build.sh)

	@echo "==> building server"
	@${currentDir}/scripts/server/build.sh $(version)
.PHONY: build-server

build-server-docker: build-server
ifndef version
	$(error version is required. Usage: make version=0.1.0 [platform=linux/amd64] build-server-docker)
endif

	@echo "==> building Docker image"
	@(cd ${currentDir}/host/docker && ./build.sh $(version) $(platform))
.PHONY: build-server-docker

generate-cli-schema:
	@echo "==> generating CLI database schema"
	@mkdir -p pkg/cli/database
	@touch pkg/cli/database/schema.sql
	@go run -tags fts5 ./pkg/cli/database/schema
.PHONY: generate-cli-schema

build-cli: generate-cli-schema
ifeq ($(debug), true)
	@echo "==> building cli in dev mode"
	@${currentDir}/scripts/cli/dev.sh
else

ifndef version
	$(error version is required. Usage: make version=0.1.0 build-cli)
endif

	@echo "==> building cli"
	@${currentDir}/scripts/cli/build.sh $(version)
endif
.PHONY: build-cli

clean:
	@git clean -f
	@rm -rf build
.PHONY: clean
