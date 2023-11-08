# Define the output binary name
BINARY_NAME := tfc-collect

# Specify the platforms you want to support
PLATFORMS := windows darwin linux

# Define the architectures for each platform
ARCH_windows := amd64
ARCH_darwin := amd64
ARCH_linux := amd64 386

# Calculate the OS-ARCH pairs for cross-compilation
OS_ARCH_PAIRS := $(foreach platform, $(PLATFORMS), $(addprefix $(platform)/, $(ARCH_$(platform))))

# Output directory for the binary files
BIN_DIR := bin

# Default target: build for all specified platforms
all: $(OS_ARCH_PAIRS)

# Create a target for each platform-architecture pair
$(OS_ARCH_PAIRS):
	GOOS=$(firstword $(subst /, ,$@)) GOARCH=$(lastword $(subst /, ,$@)) go build -buildvcs=false -o $(BIN_DIR)/$(BINARY_NAME)_$(firstword $(subst /, ,$@))_$(lastword $(subst /, ,$@))

# Clean up generated binaries
clean:
	rm -rf $(BIN_DIR)

# Build for the current platform (useful for local development)
build:
	go build -o $(BIN_DIR)/$(BINARY_NAME)

# Phony targets
.PHONY: all $(OS_ARCH_PAIRS) clean build

test:
	go test ./... -coverprofile=./cov.out -covermode=atomic -coverpkg=./... 

run: 
	go run main.go

lint:
	golangci-lint run
