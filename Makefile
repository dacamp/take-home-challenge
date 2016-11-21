include go-build/docker.mk

VENDOR = $(CURDIR)/vendor
GO := go
GOFMT := gofmt
GOVET := $(GO) vet
GOLINT := golint
GLIDE := glide


# all .go files that don't exist in hidden directories
ALL_SRC := $(shell find . -name "*.go" | grep -v -e vendor -e go-build \
	-e ".*/\..*" \
	-e ".*/_.*" \
        -e ".*/mocks.*")

$(VENDOR):
	$(GLIDE) install

vendor: $(VENDOR)

bins: vendor
	$(GO) build -o $(CURDIR)/bin/challenge-executable

test: vendor
	$(GO) test -i -race

fmt:
	@$(GOFMT) -s -w $(ALL_SRC)

lint:
	@$(GOLINT) ./... 2>&1 | grep -v vendor # this is lame, it still vets vendors packages

vet:
	@$(GOVET) ./... 2>&1  | grep -v vendor # this is lame, it still vets vendors packages

show-requests:
	open "http://localhost:1234/debug/requests?fam=main.counterHandler&b=0"
