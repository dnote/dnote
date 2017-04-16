release:
	@echo "** Releasing version $(VERSION)..."
	@echo "** Building..."
	@$(MAKE) build
	@echo "** Tagging and pushing..."
	@git tag -a $(VERSION)
	@git push --tags
.PHONY: release

build:
	@gox -os="linux darwin windows openbsd" -output="dnote_{{.OS}}_{{.Arch}}" ./...
.PHONY: build

clean:
	@git clean -f
.PHONY: clean

