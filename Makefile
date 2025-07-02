# Variables
PACKAGE := github.com/ffais/yaml-sort
VERSION := $(shell git describe --tags --always --abbrev=0 --match='v[0-9]*.[0-9]*.[0-9]*' 2> /dev/null | sed 's/^.//')
COMMIT_HASH := $(shell git rev-parse --short HEAD)
BUILD_TIMESTAMP := $(shell date '+%Y-%m-%dT%H:%M:%S')

# LDFLAGS construction
LDFLAGS := -X '$(PACKAGE)/cmd.version=$(VERSION)' \
           -X '$(PACKAGE)/cmd.commitHash=$(COMMIT_HASH)' \
           -X '$(PACKAGE)/cmd.buildTimestamp=$(BUILD_TIMESTAMP)'

.PHONY: build
build:
	@echo "Building yaml-sort..."
	go build -ldflags="$(LDFLAGS)"

.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Commit Hash: $(COMMIT_HASH)"
	@echo "Build Time: $(BUILD_TIMESTAMP)"

.PHONY: clean
clean:
	@echo "Cleaning..."
	rm -f yaml-sort

# Default target
.DEFAULT_GOAL := build
