.PHONY: test

GOCMD := go
GOTEST := $(GOCMD) test

test:
	$(GOTEST) -v ./...
