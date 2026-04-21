BINARY_NAME := assistant-engine
INSTALL_DIR := $(shell go env GOPATH)/bin
CONFIG_DIR := $(HOME)/.assistant-engine
CONFIG_FILE := $(CONFIG_DIR)/config.json
SRC := ./cmd/assistant-engine

.PHONY: build install setup uninstall test clean

build:
	go build -o $(BINARY_NAME) $(SRC)

install: build setup
	cp $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installed $(BINARY_NAME) to $(INSTALL_DIR)"
	@echo "Run 'assistant-engine --help' to verify."

setup:
	@if [ -f "$(CONFIG_FILE)" ]; then \
		echo "Config already exists at $(CONFIG_FILE), skipping setup."; \
	else \
		mkdir -p $(CONFIG_DIR); \
		echo ""; \
		echo "=== Assistant Engine Setup ==="; \
		echo ""; \
		echo "Follow the README for Teams webhook setup instructions."; \
		echo ""; \
		read -p "Webhook URL: " webhook_url; \
		read -p "Webhook type (workflow/classic) [workflow]: " webhook_type; \
		webhook_type=$${webhook_type:-workflow}; \
		read -p "Your Teams email for @mentions: " mention_id; \
		read -p "Your first name for @mentions: " mention_name; \
		read -p "Default reminder time [09:00]: " default_time; \
		default_time=$${default_time:-09:00}; \
		read -p "Default delay in hours [24]: " default_delay; \
		default_delay=$${default_delay:-24}; \
		echo "{" > $(CONFIG_FILE); \
		echo "  \"webhook_url\": \"$$webhook_url\"," >> $(CONFIG_FILE); \
		echo "  \"webhook_type\": \"$$webhook_type\"," >> $(CONFIG_FILE); \
		echo "  \"mention_id\": \"$$mention_id\"," >> $(CONFIG_FILE); \
		echo "  \"mention_name\": \"$$mention_name\"," >> $(CONFIG_FILE); \
		echo "  \"default_time\": \"$$default_time\"," >> $(CONFIG_FILE); \
		echo "  \"default_delay_hours\": $$default_delay" >> $(CONFIG_FILE); \
		echo "}" >> $(CONFIG_FILE); \
		echo ""; \
		echo "Config saved to $(CONFIG_FILE)"; \
	fi

uninstall:
	rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Removed $(BINARY_NAME) from $(INSTALL_DIR)"

test:
	go test ./... -v

clean:
	rm -f $(BINARY_NAME)
