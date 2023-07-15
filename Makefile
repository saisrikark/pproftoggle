.PHONY: fmt vet test bench tidy vendor

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

bench:
	$(GOTEST) -bench=. -benchmem -run BenchmarkToggle

tidy:
	$(GOTIDY)

vendor:
	$(GOVENDOR)