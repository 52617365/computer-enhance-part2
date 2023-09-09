
# Output binary name
BINARY_NAME = gen

# Build directory
BUILD_DIR = build

.PHONY: all build clean

# Build the project and place the binary in the build directory
all: build

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME)

# Clean up the build directory
clean:
	@rm -rf $(BUILD_DIR)
