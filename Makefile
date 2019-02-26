.PHONY: all build

all: build

build:
	go build -ldflags "-X main.Version=$(shell git describe --tags --always --dirty --match='v*' 2> /dev/null || echo v0)"
