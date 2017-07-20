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
			-n $(NAME) \
			--force \
			--version $(VERSION) \
			--rpm-os linux \
			--package $(BUILDDIR) \
			-s dir $(BUILDDIR)/$(NAME)-$(VERSION)-linux-amd64=/usr/local/bin/$(NAME)

package_deb:
	$(call package,deb)

package_rpm:
	$(call package,rpm)


###############################################################################
## Clean
##
## EXPLICITLY removing artifacts directory to protect from horrible accidents
###############################################################################
clean:
	rm -rf ./ARTIFACTS
