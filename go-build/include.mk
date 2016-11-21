GO := go
GLIDE := glide
GOFMT := gofmt
GOLINT := golint
GOVET := $(GO) vet

TMP = $(CURDIR)/.tmp
COVERFILE = $(TMP)/coverage.out
TARFILE = challenge.tar.gz
VENDOR = $(CURDIR)/vendor

# all .go files that don't exist in hidden directories
ALL_SRC := $(shell find . -name "*.go" | grep -v -e vendor -e go-build \
        -e ".*/mocks.*")
ALL_PKGS = $(shell $(GO) list $(sort $(dir $(ALL_SRC))))

BUILD_GC_FLAGS ?= -gcflags "-trimpath=$(GOPATH)/src"
BUILD_FLAGS ?=

TEST_FLAGS += $(BUILD_GC_FLAGS)

$(TMP):
	@mkdir .tmp

$(COVERFILE): $(TMP) test

$(TARFILE): package

$(VENDOR):
	$(GLIDE) install

clean:
	@rm -rf $(TMP) $(VENDOR) bin/challenge-executable $(TARFILE)

fmt:
	@$(GOFMT) -s -w $(ALL_SRC)

lint:
	$(foreach pkg, $(ALL_PKGS), \
		$(GOLINT) $(pkg) || true;)

linux: vendor
	@GOOS=linux GOARCH=amd64 $(GO) build -o $(CURDIR)/bin/challenge-executable -i $(BUILD_FLAGS) $(BUILD_GC_FLAGS)

package: linux
	cd $(CURDIR)/.. &&  \
	tar -zcvf challenge/$(TARFILE) --exclude=$(TARFILE) --exclude=.git --exclude=.tmp challenge && \
	echo "TARBALL CREATED: $(CURDIR)/$(TARFILE)"

vendor: $(VENDOR)

vet:
	$(foreach pkg, $(ALL_PKGS), \
		$(GOVET) $(pkg) || true;)

.PHONY: fmt lint vet


include go-build/docker.mk
