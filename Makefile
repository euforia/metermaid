
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

clean-$(NAME):
	rm -rf $(BUILD_DIR)

clean-ui:
	rm -rf ui/build
	rm -f ui/ui.go

clean: clean-$(NAME) clean-ui

deps:
	go get -v golang.org/x/vgo

ui/build:
	cd ./ui/ && yarn --verbose build

ui/ui.go:
	cd ./ui/build && go-bindata -pkg ui -o ../ui.go ./...

.PHONY: ui
ui: ui/build ui/ui.go

$(BUILD_DIR)/$(NAME):
	GOOS=$(GOOS) CGO_ENABLED=0 vgo build $(BUILD_OPTS) $(LD_OPTS) -o $(BUILD_DIR)/$(NAME) $(SRC_FILES)

dist: $(BUILD_DIR)/$(NAME)
	cd $(BUILD_DIR) && tar -czf $(NAME)-$(GOOS).tgz $(NAME)

all: $(BUILD_DIR)/$(NAME)
