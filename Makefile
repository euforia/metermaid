
NAME = metermaid

COMMIT = $(shell git rev-parse --short HEAD)
COMMIT_COUNT = $(shell git rev-list --count HEAD)
VERSION = $(shell git describe 2> /dev/null || echo "v0.0.0-$(COMMIT_COUNT)-$(COMMIT)")
BUILDTIME = $(shell date +%Y-%m-%dT%T%z)

GOOS ?= $(shell go env GOOS)

SRC_FILES = ./cmd/*.go
BUILD_DIR = ./build
BUILD_OPTS = -a -tags netgo -installsuffix netgo
LD_OPTS = -ldflags="-X main.version=$(VERSION) -X main.buildtime=$(BUILDTIME) -w"

clean:
	rm -rf $(BUILD_DIR)/$(NAME)

deps:
	vgo build

$(BUILD_DIR)/$(NAME):
	GOOS=$(GOOS) go build $(BUILD_OPTS) $(LD_OPTS) -o $(BUILD_DIR)/$(NAME) $(SRC_FILES)

all: $(BUILD_DIR)/$(NAME)