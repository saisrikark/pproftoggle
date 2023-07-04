.PHONY: fmt vet test

GOCMD := go
GOFMT := $(GOCMD) fmt
GOVET := $(GOCMD) vet
GOTEST := $(GOCMD) test

fmt:
	$(GOFMT) ./...

vet:
	$(GOVET) ./...

test:
	$(GOTEST) -v ./...
