
CC_TEST_REPORTER_ID := c09efa7c67c269bfdc6f8a356785d8f7ed55c9dc2b9a1d07b78c384f55c4e527
GO_HTML_COV         := ./coverage.html
GO_TEST_OUTFILE     := ./c.out
CC_PREFIX       	:= github.com/waku-org/go-waku

SHELL := bash # the shell used internally by Make

GOBIN ?= $(shell which go)

.PHONY: all build lint lint-full test coverage build-example static-library dynamic-library test-c test-c-template mobile-android mobile-ios

ifeq ($(OS),Windows_NT)     # is Windows_NT on XP, 2000, 7, Vista, 10...
 detected_OS := Windows
else
 detected_OS := $(strip $(shell uname))
endif

ifeq ($(detected_OS),Darwin)
 GOBIN_SHARED_LIB_EXT := dylib
 TEST_REPORTER_URL := https://codeclimate.com/downloads/test-reporter/test-reporter-latest-darwin-amd64
else ifeq ($(detected_OS),Windows)
 # on Windows need `--export-all-symbols` flag else expected symbols will not be found in libgowaku.dll
 GOBIN_SHARED_LIB_CGO_LDFLAGS := CGO_LDFLAGS="-Wl,--export-all-symbols"
 GOBIN_SHARED_LIB_EXT := dll
else
 TEST_REPORTER_URL := https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64
 GOBIN_SHARED_LIB_EXT := so
 GOBIN_SHARED_LIB_CGO_LDFLAGS := CGO_LDFLAGS="-Wl,-soname,libgowaku.so.0"
endif

GIT_COMMIT = $(shell git rev-parse --short HEAD)
VERSION = $(shell cat ./VERSION)
UID := $(shell id -u)
GID := $(shell id -g)


BUILD_FLAGS ?= $(shell echo "-ldflags='\
	-X github.com/waku-org/go-waku/waku/v2/node.GitCommit=$(GIT_COMMIT) \
	-X github.com/waku-org/go-waku/waku/v2/node.Version=$(VERSION)'")

ANDROID_TARGET ?= 23

# control rln code compilation
ifeq ($(NO_RLN), true)
BUILD_TAGS := gowaku_no_rln
endif

all: build

deps: lint-install

build-with-race:
	${GOBIN} build -race -tags="${BUILD_TAGS}" $(BUILD_FLAGS) -o build/waku ./cmd/waku

build:
	${GOBIN} build -tags="${BUILD_TAGS}" $(BUILD_FLAGS) -o build/waku ./cmd/waku

chat2:
	pushd ./examples/chat2 && \
	${GOBIN} build -o ../../build/chat2 . && \
	popd

vendor:
	${GOBIN} mod tidy

lint-install:
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
		bash -s -- -b $(shell ${GOBIN} env GOPATH)/bin v1.52.2

lint:
	@echo "lint"
	@golangci-lint run ./... --deadline=5m

lint-full:
	@echo "lint"
	@golangci-lint run ./... --config=./.golangci.full.yaml --deadline=5m

test-with-race:
	${GOBIN} test -race -timeout 300s ./waku/... 

test:
	${GOBIN} test -timeout 300s ./waku/... -coverprofile=${GO_TEST_OUTFILE}.tmp
	cat ${GO_TEST_OUTFILE}.tmp | grep -v ".pb.go" > ${GO_TEST_OUTFILE}
	${GOBIN} tool cover -html=${GO_TEST_OUTFILE} -o ${GO_HTML_COV}

COVERAGE_FILE := ./coverage/cc-test-reporter
$(COVERAGE_FILE):
	curl -sfL $(TEST_REPORTER_URL) --output ./coverage/cc-test-reporter #TODO: Support windows
	chmod +x ./coverage/cc-test-reporter

_before-cc: $(COVERAGE_FILE)
   
	CC_TEST_REPORTER_ID=${CC_TEST_REPORTER_ID} ./coverage/cc-test-reporter before-build
	
_after-cc:
	GIT_COMMIT=$(git log | grep -m1 -oE '[^ ]+$') CC_TEST_REPORTER_ID=${CC_TEST_REPORTER_ID} ./coverage/cc-test-reporter after-build --prefix ${CC_PREFIX}

test-ci: _before-cc test _after-cc

generate:
	${GOBIN} generate ./...
	
coverage:
	${GOBIN} test -count 1 -coverprofile=coverage.out ./...
	${GOBIN} tool cover -html=coverage.out -o=coverage.html

# build a docker image for the fleet
docker-image: DOCKER_IMAGE_TAG ?= latest
docker-image: DOCKER_IMAGE_NAME ?= statusteam/go-waku:$(DOCKER_IMAGE_TAG)
docker-image:
	docker build --tag $(DOCKER_IMAGE_NAME) \
		--build-arg="GIT_COMMIT=$(shell git rev-parse HEAD)" .

build-example-basic2:
	cd examples/basic2 && $(MAKE)

build-example-chat-2:
	cd examples/chat2 && $(MAKE)

build-example-filter2:
	cd examples/filter2 && $(MAKE)

build-example-c-bindings:
	cd examples/c-bindings && $(MAKE)

build-example-rln:
	cd examples/rln && $(MAKE)

build-example: build-example-basic2 build-example-chat-2 build-example-filter2 build-example-c-bindings build-example-rln

static-library:
	@echo "Building static library..."
	${GOBIN} build \
		-buildmode=c-archive \
		-tags="${BUILD_TAGS} gowaku_no_rln" \
		-o ./build/lib/libgowaku.a \
		./library/c/
	@echo "Static library built:"
ifeq ($(detected_OS),Darwin)
	sed -i '' -e "s/#include <cgo_utils.h>//gi" ./build/lib/libgowaku.h
else
	sed -i "s/#include <cgo_utils.h>//gi" ./build/lib/libgowaku.h
endif
	@ls -la ./build/lib/libgowaku.*

dynamic-library:
	@echo "Building shared library..."
	rm -f ./build/lib/libgowaku.$(GOBIN_SHARED_LIB_EXT)*
	$(GOBIN_SHARED_LIB_CFLAGS) $(GOBIN_SHARED_LIB_CGO_LDFLAGS) ${GOBIN} build \
		-buildmode=c-shared \
		-tags="${BUILD_TAGS} gowaku_no_rln" \
		-o ./build/lib/libgowaku.$(GOBIN_SHARED_LIB_EXT) \
		./library/c/
ifeq ($(detected_OS),Darwin)
	sed -i '' -e "s/#include <cgo_utils.h>//gi" ./build/lib/libgowaku.h
else
	sed -i "s/#include <cgo_utils.h>//gi" ./build/lib/libgowaku.h
endif
ifeq ($(detected_OS),Linux)
	cd ./build/lib && \
	mv ./libgowaku.$(GOBIN_SHARED_LIB_EXT) ./libgowaku.$(GOBIN_SHARED_LIB_EXT).0 && \
	ln -s ./libgowaku.$(GOBIN_SHARED_LIB_EXT).0 ./libgowaku.$(GOBIN_SHARED_LIB_EXT)
endif
	@echo "Shared library built:"
	@ls -la ./build/lib/libgowaku.*

mobile-android:
	@echo "Android target: ${ANDROID_TARGET} (override with ANDROID_TARGET var)"
	gomobile init && \
	${GOBIN} get -d golang.org/x/mobile/cmd/gomobile && \
	CGO=1 gomobile bind -v -target=android -androidapi=${ANDROID_TARGET} -ldflags="-s -w" -tags="${BUILD_TAGS} gowaku_no_rln" $(BUILD_FLAGS) -o ./build/lib/gowaku.aar ./library/mobile
	@echo "Android library built:"
	@ls -la ./build/lib/*.aar ./build/lib/*.jar

mobile-ios:
	gomobile init && \
	${GOBIN} get -d golang.org/x/mobile/cmd/gomobile && \
	gomobile bind -target=ios -ldflags="-s -w" -tags="nowatchdog ${BUILD_TAGS} gowaku_no_rln" $(BUILD_FLAGS) -o ./build/lib/Gowaku.xcframework ./library/mobile
	@echo "IOS library built:"
	@ls -la ./build/lib/*.xcframework

install-xtools:
	${GOBIN} install golang.org/x/tools/...@v0.1.10

install-bindata:
	${GOBIN} install github.com/kevinburke/go-bindata/go-bindata@v3.13.0

install-gomobile: install-xtools
	${GOBIN} install golang.org/x/mobile/cmd/gomobile@v0.0.0-20220518205345-8578da9835fd
	${GOBIN} install golang.org/x/mobile/cmd/gobind@v0.0.0-20220518205345-8578da9835fd

build-linux-pkg:
	docker build --build-arg UID=${UID} --build-arg GID=${GID} -f ./scripts/linux/Dockerfile -t statusteam/gowaku-linux-pkgs:latest .
	./scripts/linux/docker-run.sh
	ls -la ./build/*.rpm ./build/*.deb

TEST_MNEMONIC="swim relax risk shy chimney please usual search industry board music segment"

start-ganache:
	docker run -p 8545:8545 --name ganache-cli --rm -d trufflesuite/ganache-cli:latest -m ${TEST_MNEMONIC}

stop-ganache:
	docker stop ganache-cli

test-onchain: BUILD_TAGS += include_onchain_tests
test-onchain:
	${GOBIN} test -v -count 1 -tags="${BUILD_TAGS}" github.com/waku-org/go-waku/waku/v2/protocol/rln

test-onchain-with-race:
	${GOBIN} test -race -v -count 1 -tags="${BUILD_TAGS}" github.com/waku-org/go-waku/waku/v2/protocol/rln