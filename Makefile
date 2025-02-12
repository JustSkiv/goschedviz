.PHONY: build clean test

# Binary name
BINARY_NAME=goschedviz
# Binary directory
BINARY_DIR=bin

build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BINARY_DIR)
	@go build -o $(BINARY_DIR)/$(BINARY_NAME) ./cmd/goschedviz

clean:
	@echo "Cleaning..."
	@rm -rf $(BINARY_DIR)

test:
	@go test -v ./...