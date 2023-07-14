.PHONY: fmt vet test tidy vendor

GOCMD := go
GOMOD := $(GOCMD) mod
GOFMT := $(GOCMD) fmt
GOVET := $(GOCMD) vet
GOTEST := $(GOCMD) test
GOTIDY := $(GOMOD) tidy
GOVENDOR := $(GOMOD) vendor

fmt:
	$(GOFMT) ./...

vet:
	$(GOVET) ./...

test:
	$(GOTEST) ./...

tidy:
	$(GOTIDY)

vendor:
	$(GOVENDOR)