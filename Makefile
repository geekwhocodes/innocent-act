# Go 
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Versioning
LAST_COMMIT := $(shell git rev-parse --short HEAD)
VERSION := $(shell git describe --tags --always --abbrev=8)
BUILD_DATE := $(shell date '+%Y-%m-%d %H:%M:%S')
BUILD := ${LAST_COMMIT}

# Binary spec 
BIN := innocent-act
LDFLAGS := -s -w -X 'main.buildString=${BUILD}' -X 'main.versionString=${VERSION}' # -X 'main.buildDate=${BUILD_DATE}'

GOOS="linux"
GOARCH="amd64"

# Static files to stuff in 
STATIC := config.sample.yaml \
	web/dist/web:web \
	web/dist/favicon.ico:/web/favicon.ico \

VUE_APP_VERSION=${VERSION}


# Dependencies
.PHONY: deps
deps:
	${GOGET} -u github.com/knadh/stuffbin/...
	cd web && yarn install

# Build the binary ${BIN} for windows
.PHONY: build-win
 build-win: 
	${GOBUILD} -o ${BIN}.exe -ldflags="${LDFLAGS}" ./cmd

# Build the binary ${BIN} for linux
.PHONY: build-linux
 build-linux:
	${GOBUILD} -o ${BIN} -ldflags="${LDFLAGS}" ./cmd

# Build vue app at web/dist
.PHONY: build-frontend
build-frontend:
	set VUE_APP_VERSION=${VUE_APP_VERSION} && cd web && yarn build

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
all: build-win build-frontend stuffin
	stuffbin -a stuff -in ${BIN} -out ${BIN} ${STATIC}

# Bundle all static assets including the JS frontend into the ./listmonk binary
# using stuffbin (installed with make deps).
.PHONY: all-linux
all: build-linux build-frontend stuffin
	stuffbin -a stuff -in ${BIN} -out ${BIN} ${STATIC}
