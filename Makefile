release:
	@echo "** Releasing version $(VERSION)..."
	@echo "** Tagging and pushing..."
	@git tag -a $(VERSION) -m "$(VERSION)"
	@git push --tags
	@goreleaser --rm-dist
.PHONY: release

build-snapshot:
	@goreleaser --snapshot --rm-dist
.PHONY: build-snapshot

clean:
	@git clean -f
.PHONY: clean
