.PHONY: all build clean plugins plugin-example

VERSION:=$(shell git describe --tags --always --dirty --match='v*' 2> /dev/null || echo v0)
LDFLAGS:=-ldflags "-X main.Version=$(VERSION)"

GO:=$(shell which go)

all: build

build: plugins tpl

clean:
	rm -f tpl
	find . -type f -name "*.so" | xargs -I {} rm -f {}

plugins: plugin-example

plugin-example:
	$(GO) build -buildmode=plugin $(LDFLAGS) example

tpl:
	$(GO) build $(LDFLAGS)