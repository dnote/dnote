all:
	./scripts/build.sh $(VERSION)
.PHONY: all

release:
	./scripts/build.sh $(VERSION)
	./scripts/release.sh $(VERSION)
.PHONY: release

clean:
	@git clean -f
.PHONY: clean
