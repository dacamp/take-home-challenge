include go-build/docker.mk

GO := go
GLIDE := glide
GOFMT := gofmt
GOLINT := golint
GOVET := $(GO) vet

TMP = $(CURDIR)/.tmp
COVERFILE = $(TMP)/coverage.out
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

$(VENDOR):
	$(GLIDE) install

bins: vendor
	@$(GO) build -o $(CURDIR)/bin/challenge-executable -i $(BUILD_FLAGS) $(BUILD_GC_FLAGS)

clean:
	@rm -rf .tmp vendor bin/challenge-executable

cover:
	@rm -f $(COVERFILE);
	@echo "mode: count" > $(COVERFILE);
	@grep -h -v "mode: " $(TMP)/*.cover >> $(COVERFILE);
	@$(GO) tool cover -html=$(COVERFILE)

fmt:
	@$(GOFMT) -s -w $(ALL_SRC)

lint:
	$(foreach pkg, $(ALL_PKGS), \
		$(GOLINT) $(pkg) || true;)

test: $(TMP) vendor
	@$(foreach pkg, $(ALL_PKGS), \
		 $(GO) test -v -race $(TEST_FLAGS) -coverprofile $(TMP)/$(lastword $(subst /, ,$(pkg))).cover $(pkg) || true;)

vendor: $(VENDOR)

vet:
	$(foreach pkg, $(ALL_PKGS), \
		$(GOVET) $(pkg) || true;)

show-requests:
	open "http://localhost:1234/debug/requests?fam=main.counterHandler&b=0"

.PHONY: fmt lint vet
