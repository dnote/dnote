DEP := $(shell command -v dep 2> /dev/null)
PACKR2 := $(shell command -v packr2 2> /dev/null)
NPM := $(shell command -v npm 2> /dev/null)
HUB := $(shell command -v hub 2> /dev/null)
COMPILEDAEMON := $(shell command -v CompileDaemon 2> /dev/null)

serverOutputDir = ${GOPATH}/src/github.com/dnote/dnote/build/server
cliOutputDir = ${GOPATH}/src/github.com/dnote/dnote/build/cli
cliHomebrewDir = ${GOPATH}/src/github.com/dnote/homebrew-dnote

## installation
install: install-go install-js
.PHONY: install

install-go:
ifndef DEP
	@echo "==> installing dep"
	@curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif

ifndef PACKR2
	@echo "==> installing packr2"
	@go get -u github.com/gobuffalo/packr/v2/packr2
endif

ifndef COMPILEDAEMON
	@echo "==> installing CompileDaemon"
	@go get -u github.com/githubnemo/CompileDaemon
endif

	@echo "==> installing go dependencies"
	@dep ensure
.PHONY: install-go

install-js:
ifndef NPM
	$(error npm is not installed)
endif

	@echo "==> installing js dependencies"

ifeq ($(CI), true)
	@(cd ${GOPATH}/src/github.com/dnote/dnote/web && npm install --unsafe-perm=true)
	@(cd ${GOPATH}/src/github.com/dnote/dnote/browser && npm install --unsafe-perm=true)
	@(cd ${GOPATH}/src/github.com/dnote/dnote/jslib && npm install --unsafe-perm=true)
else
	@(cd ${GOPATH}/src/github.com/dnote/dnote/web && npm install)
	@(cd ${GOPATH}/src/github.com/dnote/dnote/browser && npm install)
	@(cd ${GOPATH}/src/github.com/dnote/dnote/jslib && npm install)
endif
.PHONY: install-js

## test
test: test-cli test-api test-web test-jslib
.PHONY: test

test-cli:
	@echo "==> running CLI test"
	@${GOPATH}/src/github.com/dnote/dnote/pkg/cli/scripts/test.sh
.PHONY: test-cli

test-api:
	@echo "==> running API test"
	@${GOPATH}/src/github.com/dnote/dnote/pkg/server/api/scripts/test-local.sh
.PHONY: test-api

test-web:
	@echo "==> running web test"

ifeq ($(WATCH), true)
	@(cd ${GOPATH}/src/github.com/dnote/dnote/web && npm run test:watch)
else 
	@(cd ${GOPATH}/src/github.com/dnote/dnote/web && npm run test)
endif
.PHONY: test-web

test-jslib:
	@echo "==> running jslib test"

ifeq ($(WATCH), true)
	@(cd ${GOPATH}/src/github.com/dnote/dnote/jslib && npm run test:watch)
else
	@(cd ${GOPATH}/src/github.com/dnote/dnote/jslib && npm run test)
endif
.PHONY: test-jslib

# development
dev-server:
	@echo "==> running dev environment"
	@(cd ${GOPATH}/src/github.com/dnote/dnote/web && ./scripts/dev.sh)
.PHONY: dev-server

## build
build-web:
	@echo "==> building web"
	@(cd ${GOPATH}/src/github.com/dnote/dnote/web && ./scripts/build-prod.sh)
.PHONY: build-web

build-server: build-web
ifndef version
	$(error version is required. Usage: make version=0.1.0 build-server)
endif

	@echo "==> building server"
	@(cd ${GOPATH}/src/github.com/dnote/dnote/pkg/server && ./scripts/build.sh $(version))
.PHONY: build-server

build-cli:
ifeq ($(debug), true)
	@echo "==> building cli in dev mode"
	@${GOPATH}/src/github.com/dnote/dnote/pkg/cli/scripts/dev.sh
else

ifndef version
	$(error version is required. Usage: make version=0.1.0 build-cli)
endif

	@echo "==> building cli"
	@${GOPATH}/src/github.com/dnote/dnote/pkg/cli/scripts/build.sh $(version)
endif
.PHONY: build-cli

## release
release-cli: build-cli
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
	@${GOPATH}/src/github.com/dnote/dnote/scripts/release.sh cli $(version) ${cliOutputDir}

	@echo "===> releading on Homebrew"
	@(cd "${cliHomebrewDir}" && \
		./release.sh \
			"$(version)" \
			"${shasum -a 256 "${cliOutputDir}/dnote_$(version)_darwin_amd64.tar.gz" | cut -d ' ' -f 1}" \
	)
.PHONY: release-cli

release-server: build-server
ifndef version
	$(error version is required. Usage: make version=0.1.0 release-server)
endif
ifndef HUB
	$(error please install hub)
endif

	@echo "==> releasing server"
	@${GOPATH}/src/github.com/dnote/dnote/scripts/release.sh server $(version) ${serverOutputDir}
.PHONY: release-server

# migrations
create-migration:
ifndef filename
	$(error filename is required. Usage: make filename=your-filename create-migration)
endif

	@(cd ${GOPATH}/src/github.com/dnote/dnote/pkg/server/database && ./scripts/create-migration.sh $(filename))
.PHONY: create-migration

clean:
	@git clean -f
	@rm -rf build
	@rm -rf web/public
.PHONY: clean
