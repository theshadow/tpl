.PHONY: all build plugin-example

VERSION:=$(shell git describe --tags --always --dirty --match='v*' 2> /dev/null || echo v0)
LDFLAGS:=-ldflags "-X main.Version=$(VERSION)"

all: build

build: plugins

plugins: plugin-example

plugin-example:
	cd cmd/plugins/example && go build -buildmode=plugin $(LDFLAGS)