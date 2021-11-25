NAME = wrench
BUILDDIR = ./ARTIFACTS

# Remove prefix since deb, rpm etc don't recognize this as valid version
VERSION = $(shell git describe --tags --match 'v[0-9]*\.[0-9]*\.[0-9]*' | sed 's/^v//')

###############################################################################
## Building
###############################################################################

.PHONY: build build_darwin build_linux
build: build_darwin build_linux

compile = bash -c "env GOOS=$(1) GOARCH=$(2) go build -a \
						-ldflags \"-w -X main.VERSION='$(VERSION)'\" \
						-o $(BUILDDIR)/$(NAME)-$(VERSION)-$(1)-$(2)"

build_darwin:
	$(call compile,darwin,amd64)
	$(call compile,darwin,arm64)

build_linux:
	$(call compile,linux,amd64)

###############################################################################
## Clean
##
## EXPLICITLY removing artifacts directory to protect from horrible accidents
###############################################################################
clean:
	rm -rf ./ARTIFACTS
