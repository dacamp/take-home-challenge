include go-build/docker.mk

GO := go
GLIDE := glide
GOFMT := gofmt
GOLINT := golint
GOVET := $(GO) vet
VENDOR = $(CURDIR)/vendor


# all .go files that don't exist in hidden directories
ALL_SRC := $(shell find . -name "*.go" | grep -v -e vendor -e go-build \
        -e ".*/mocks.*")
ALL_PKGS = $(shell $(GO) list $(sort $(dir $(ALL_SRC))))

BUILD_GC_FLAGS ?= -gcflags "-trimpath=$(GOPATH)/src"
BUILD_FLAGS ?=
TEST_FLAGS += $(BUILD_GC_FLAGS)

$(VENDOR):
	$(GLIDE) install

vendor: $(VENDOR)

bins: vendor
	@$(GO) build -o $(CURDIR)/bin/challenge-executable -i $(BUILD_FLAGS) $(BUILD_GC_FLAGS)

test: vendor
	@$(foreach pkg, $(ALL_PKGS), \
		$(GO) test -v -race $(pkg) || true;)

fmt:
	@$(GOFMT) -s -w $(ALL_SRC)

lint:
	$(foreach pkg, $(ALL_PKGS), \
		$(GOLINT) $(pkg) || true;)

vet:
	$(foreach pkg, $(ALL_PKGS), \
		$(GOVET) $(pkg) || true;)

show-requests:
	open "http://localhost:1234/debug/requests?fam=main.counterHandler&b=0"

.PHONY: fmt lint vet
