#!/usr/bin/env make

NAME=luscheduler
BINARY=./bin/${NAME}
SOURCEDIR=./src
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

VERSION := $(shell git describe --abbrev=0 --tags)
SHA := $(shell git rev-parse --short HEAD)

GOPATH ?= /usr/local/go
GOPATH := ${CURDIR}:${GOPATH}
export GOPATH

$(BINARY): $(SOURCES)
	go build -o ${BINARY} -ldflags "-X main.BuildVersion=$(VERSION)-$(SHA)" $(SOURCEDIR)/$(NAME)/cmd/main.go

run: clean $(BINARY)
	${BINARY}

clean:
	rm -f $(BINARY)

.DEFAULT_GOAL: $(BINARY)

include Makefile.git