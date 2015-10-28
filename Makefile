GO_VERSION = 1.4
BINARY_NAME = wrench
BUILDDIR = ./ARTIFACTS

# Remove prefix since deb, rpm etc don't recognize this as valid version
VERSION = $(shell git describe | sed 's/^v//')

# EXPLICITLY removing artifacts directory to protect from horrible accidents
clean:
	rm -rf ./ARTIFACTS

###############################################################################
## Building
##
## Travis CI Gimme is used to cross-compile wrench
## https://github.com/travis-ci/gimme
###############################################################################

.PHONY: build build_darwin build_linux
build: build_darwin build_linux

compile = bash -c "eval \"$$(GIMME_GO_VERSION=$(GO_VERSION) GIMME_OS=$(1) GIMME_ARCH=$(2) gimme)\"; \
					go build -a \
						-ldflags \"-w -X main.VERSION '$(VERSION)'\" \
						-o $(BUILDDIR)/$(BINARY_NAME)-$(VERSION)-$(1)-$(2)"

build_darwin:
	$(call compile,darwin,amd64)

build_linux:
	$(call compile,linux,amd64)


###############################################################################
## Packaging
##
## Effing Package Management - fpm is used for packaging
## https://github.com/jordansissel/fpm
## gnu-tar, rpmbuild is required
###############################################################################

.PHONY: package package_deb package_rpm
package: package_deb package_rpm

package = fpm -t $(1) \
			-n wrench \
			--force \
			--version $(VERSION) \
			--rpm-os linux \
			--package $(BUILDDIR) \
			-s dir $(BUILDDIR)/wrench-$(VERSION)-linux-amd64=/usr/local/bin/wrench

package_deb:
	$(call package,deb)

package_rpm:
	$(call package,rpm)
