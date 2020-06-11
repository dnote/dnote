PACKR2 := $(shell command -v packr2 2> /dev/null)
NPM := $(shell command -v npm 2> /dev/null)
HUB := $(shell command -v hub 2> /dev/null)

currentDir = $(shell pwd)
serverOutputDir = ${currentDir}/build/server
cliOutputDir = ${currentDir}/build/cli
cliHomebrewDir = ${currentDir}/../homebrew-dnote

## installation
install: install-go install-js
.PHONY: install

install-go:
ifndef PACKR2
	@echo "==> installing packr2"
	@go get -u github.com/gobuffalo/packr/v2/packr2
endif

	@echo "==> installing go dependencies"
	@go mod download
.PHONY: install-go

install-js:
ifndef NPM
	$(error npm is not installed)
endif

	@echo "==> installing js dependencies"

ifeq ($(CI), true)
	@(cd ${currentDir} && npm ci --cache $(NPM_CACHE_DIR) --prefer-offline --unsafe-perm=true)
	@(cd ${currentDir}/web && npm ci --cache $(NPM_CACHE_DIR) --prefer-offline --unsafe-perm=true)
	@(cd ${currentDir}/browser && npm ci --cache $(NPM_CACHE_DIR) --prefer-offline --unsafe-perm=true)
	@(cd ${currentDir}/jslib && npm ci --cache $(NPM_CACHE_DIR) --prefer-offline --unsafe-perm=true)
else
	@(cd ${currentDir} && npm install)
	@(cd ${currentDir}/web && npm install)
	@(cd ${currentDir}/browser && npm install)
	@(cd ${currentDir}/jslib && npm install)
endif
.PHONY: install-js

lint:
	@(cd ${currentDir}/web && npm run lint)
	@(cd ${currentDir}/jslib && npm run lint)
	@(cd ${currentDir}/browser && npm run lint)
.PHONY: lint

lint-fix:
	@(cd ${currentDir}/web && npm run lint:fix)
	@(cd ${currentDir}/jslib && npm run lint:fix)
	@(cd ${currentDir}/browser && npm run lint:fix)
.PHONY: lint

## test
test: test-cli test-api test-web test-jslib
.PHONY: test

test-cli:
	@echo "==> running CLI test"
	@(${currentDir}/scripts/cli/test.sh)
.PHONY: test-cli

test-api:
	@echo "==> running API test"
	@(${currentDir}/scripts/server/test-local.sh)
.PHONY: test-api

test-web:
	@echo "==> running web test"

ifeq ($(WATCH), true)
	@(cd ${currentDir}/web && npm run test:watch)
else 
	@(cd ${currentDir}/web && npm run test)
endif
.PHONY: test-web

test-jslib:
	@echo "==> running jslib test"

ifeq ($(WATCH), true)
	@(cd ${currentDir}/jslib && npm run test:watch)
else
	@(cd ${currentDir}/jslib && npm run test)
endif
.PHONY: test-jslib

test-selfhost:
	@echo "==> running a smoke test for self-hosting"

	@${currentDir}/host/smoketest/run_test.sh ${tarballPath}
.PHONY: test-jslib

# development
dev-server:
	@echo "==> running dev environment"
	@VERSION=master ${currentDir}/scripts/web/dev.sh
.PHONY: dev-server

## build
build-web:
ifndef version
	$(error version is required. Usage: make version=0.1.0 build-web)
endif
	@echo "==> building web"
	@VERSION=${version} ${currentDir}/scripts/web/build-prod.sh
.PHONY: build-web

build-server: build-web
ifndef version
	$(error version is required. Usage: make version=0.1.0 build-server)
endif

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
ifndef HUB
	$(error please install hub)
endif

	if [ ! -d ${cliHomebrewDir} ]; then \
		@echo "homebrew-dnote not found locally. did you clone it?"; \
		@exit 1; \
	fi

	@echo "==> releasing cli"
	@${currentDir}/scripts/release.sh cli $(version) ${cliOutputDir}

	@echo "===> releading on Homebrew"
	@(cd "${cliHomebrewDir}" && \
		./release.sh "$(version)" "${cliOutputDir}/dnote_$(version)_darwin_amd64.tar.gz")
.PHONY: release-cli

release-server:
ifndef version
	$(error version is required. Usage: make version=0.1.0 release-server)
endif
ifndef HUB
	$(error please install hub)
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
	@rm -rf web/public
.PHONY: clean

clean-dep:
	@rm -rf ${currentDir}/web/node_modules
	@rm -rf ${currentDir}/jslib/node_modules
	@rm -rf ${currentDir}/browser/node_modules
.PHONY: clean-dep
