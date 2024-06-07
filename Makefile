# Define the binary name
BINARY_NAME=./bin/qa-challenge-application

# Define the package path (optional)
PACKAGE_PATH=./

# Default target: build the binary
all: run

# Run the application
run: build
	@echo "Running the application from: $(BINARY_NAME)"
	@./$(BINARY_NAME)

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@echo "Dependencies installed"

# Build the binary
build:
	@echo "Building the binary: $(BINARY_NAME)"
	@go build -o $(BINARY_NAME) $(PACKAGE_PATH)
	@echo "Build completed: $(BINARY_NAME)"

# Install
install:
	@echo "Installing application..."
	@go install
	@echo "Application installed"


fmt:
	@echo "Formatting the code..."
	@go fmt ./...
	@echo "Formatting completed"



.PHONY: all build test clean run lint fmt deps install-tools