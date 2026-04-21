BINARY_NAME := assistant-engine
INSTALL_DIR := $(shell go env GOPATH)/bin
SRC := ./cmd/assistant-engine

.PHONY: build install uninstall test clean

build:
	go build -o $(BINARY_NAME) $(SRC)

install: build
	cp $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installed $(BINARY_NAME) to $(INSTALL_DIR)"

uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Removed $(BINARY_NAME) from $(INSTALL_DIR)"

test:
	go test ./... -v

clean:
	rm -f $(BINARY_NAME)
