# Define the binary name
BINARY_NAME=./app/qa-challenge-application

# Define the package path (optional)
PACKAGE_PATH=./

# Default target: build the binary
all: build run

# Run the application
run:
	@echo "Running the application from: $(BINARY_NAME)"
	@$(BINARY_NAME)

# Build the binary
build:
	@echo "Building the binary: $(BINARY_NAME)"
	@go build -o $(BINARY_NAME) $(PACKAGE_PATH)
	@echo "Build completed: $(BINARY_NAME)"
.PHONY: all build run