# Go 
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Versioning
LAST_COMMIT := $(shell git rev-parse --short HEAD)
VERSION := $(shell git describe --tags --always --dirty)
BUILD_DATE := $(shell date '+%Y-%m-%d %H:%M:%S')
BUILD := ${LAST_COMMIT}

# Binary spec 
BIN := innocent-act.exe
LDFLAGS := -s -w -X 'main.buildString=${BUILD}' -X 'main.versionString=${VERSION}' # -X 'main.buildDate=${BUILD_DATE}'

# Static files to stuff in 
STATIC := config.sample.yaml

# Dependencies
.PHONY: deps
deps:
	${GOGET} -u github.com/knadh/stuffbin/...

# Build the binary ${BIN} 
.PHONY: build
build: 
	${GOBUILD} -o ${BIN} -ldflags="${LDFLAGS}" ./cmd

# Run tests.
.PHONY: test
test:
	${GOTEST} ./...

# Run the backend.
.PHONY: run
run: build
	./${BIN}

# Clean 
.PHONY: clean
clean: 
	rm -f ${BIN}
	rm -rf ./bin

# stuff in static files in ${BIN}
.PHONY: stuffin
stuffin: 
	stuffbin -a stuff -in ${BIN} -out ${BIN} ${STATIC}

# Bundle all static assets including the JS frontend into the ./listmonk binary
# using stuffbin (installed with make deps).
.PHONY: all
all: build stuffin
	stuffbin -a stuff -in ${BIN} -out ${BIN} ${STATIC}
