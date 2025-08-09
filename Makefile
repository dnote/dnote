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

test-cli:
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

test-selfhost:
	@echo "==> running a smoke test for self-hosting"

	@${currentDir}/host/smoketest/run_test.sh ${tarballPath}
.PHONY: test-selfhost

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

build-cli:
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

## release
release-cli: clean build-cli
ifndef version
	$(error version is required. Usage: make version=0.1.0 release-cli)
endif
ifndef GH
	$(error please install github-cli)
endif

	@echo "==> releasing cli"
	@${currentDir}/scripts/release.sh cli $(version) ${cliOutputDir}
.PHONY: release-cli

release-cli-homebrew: clean build-cli
ifndef version
	$(error version is required. Usage: make version=0.1.0 release-cli)
endif

	@echo "==> releasing cli on Homebrew"
	@${currentDir}/scripts/cli/release-homebrew.sh $(version) ${cliOutputDir}
.PHONY: release-cli

release-server:
ifndef version
	$(error version is required. Usage: make version=0.1.0 release-server)
endif
ifndef GH
	$(error please install github-cli)
endif

	@echo "==> releasing server"
	@${currentDir}/scripts/release.sh server $(version) ${serverOutputDir}

	@echo "==> building and releasing docker image"
	@(cd ${currentDir}/host/docker && ./build.sh $(version))
	@(cd ${currentDir}/host/docker && ./release.sh $(version))
.PHONY: release-server

# migrations
create-migration:
ifndef filename
	$(error filename is required. Usage: make filename=your-filename create-migration)
endif

	@(cd ${currentDir}/pkg/server/database && ./scripts/create-migration.sh $(filename))
.PHONY: create-migration

clean:
	@git clean -f
	@rm -rf build
.PHONY: clean
