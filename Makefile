release:
	@echo "** Releasing version $(VERSION)..."
	@echo "** Tagging and pushing..."
	@git tag -a $(VERSION) -m "$(VERSION)"
	@git push --tags
	@API_ENDPOINT=https://api.dnote.io goreleaser --rm-dist
.PHONY: release

build-snapshot:
	@API_ENDPOINT=http://127.0.0.1:5000 goreleaser --snapshot --rm-dist
.PHONY: build-snapshot

clean:
	@git clean -f
.PHONY: clean
