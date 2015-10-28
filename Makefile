.PHONY: build build_darwin build_linux

GO_VERSION = 1.4
BINARY_NAME = wrench
VERSION = $(shell git describe)

# Travis CI Gimme is used to cross-compile wrench
# https://github.com/travis-ci/gimme
compile = bash -c "eval \"$$(GIMME_GO_VERSION=$(GO_VERSION) GIMME_OS=$(1) GIMME_ARCH=$(2) gimme)\"; \
					go build -a \
						-ldflags \"-w -X main.VERSION '$(VERSION)'\" \
						-o $(BINARY_NAME)-$(VERSION)-$(1)-$(2)"

build: build_darwin build_linux

build_darwin:
	$(call compile,darwin,amd64)

build_linux:
	$(call compile,linux,amd64)
