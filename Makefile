
NAME = metermaid

COMMIT = $(shell git rev-parse --short HEAD)
COMMIT_COUNT = $(shell git rev-list --count HEAD)
VERSION = $(shell git describe 2> /dev/null || echo "v0.0.0-$(COMMIT_COUNT)-$(COMMIT)")
BUILDTIME = $(shell date +%Y-%m-%dT%T%z)

SRC_FILES = ./cmd/*.go
BUILD_OPTS = -a -tags netgo -installsuffix netgo
LD_OPTS = -ldflags="-X main.version=$(VERSION) -X main.buildtime=$(BUILDTIME) -w"


$(NAME):
	go build $(BUILD_OPTS) $(LD_OPTS) -o $(NAME) $(SRC_FILES)