PACKAGE=github.com/nueavv/istioctl-junit/common
CURRENT_DIR=$(shell pwd)
DIST_DIR=${CURRENT_DIR}/dist
CLI_NAME=istioctl-junit
BIN_NAME=istioctl-junit

HOST_OS:=$(shell go env GOOS)
HOST_ARCH:=$(shell go env GOARCH)

TARGET_ARCH?=linux/amd64

VERSION=$(shell cat ${CURRENT_DIR}/VERSION)
BUILD_DATE:=$(if $(BUILD_DATE),$(BUILD_DATE),$(shell date -u +'%Y-%m-%dT%H:%M:%SZ'))
GIT_COMMIT:=$(if $(GIT_COMMIT),$(GIT_COMMIT),$(shell git rev-parse HEAD))
GIT_TAG:=$(if $(GIT_TAG),$(GIT_TAG),$(shell if [ -z "`git status --porcelain`" ]; then git describe --exact-match --tags HEAD 2>/dev/null; fi))
GIT_TREE_STATE:=$(if $(GIT_TREE_STATE),$(GIT_TREE_STATE),$(shell if [ -z "`git status --porcelain`" ]; then echo "clean" ; else echo "dirty"; fi))
VOLUME_MOUNT=$(shell if test "$(go env GOOS)" = "darwin"; then echo ":delegated"; elif test selinuxenabled; then echo ":delegated"; else echo ""; fi)

GOPATH?=$(shell if test -x `which go`; then go env GOPATH; else echo "$(HOME)/go"; fi)
GOCACHE?=$(HOME)/.cache/go-build

override LDFLAGS += \
  -X ${PACKAGE}.version=${VERSION} \
  -X ${PACKAGE}.buildDate=${BUILD_DATE} \
  -X ${PACKAGE}.gitCommit=${GIT_COMMIT} \
  -X ${PACKAGE}.gitTreeState=${GIT_TREE_STATE}\
  -X "${PACKAGE}.extraBuildInfo=${EXTRA_BUILD_INFO}"

ifeq (${STATIC_BUILD}, true)
override LDFLAGS += -extldflags "-static"
endif

ifneq (${GIT_TAG},)
LDFLAGS += -X ${PACKAGE}.gitTag=${GIT_TAG}
endif

# Cleans VSCode debug.test files from sub-dirs to prevent them from being included in by golang embed
.PHONY: clean-debug
clean-debug:
	-find ${CURRENT_DIR} -name debug.test -exec rm -f {} +

.PHONY: cli-local
cli-local: clean-debug
	CGO_ENABLED=0 GODEBUG="tarinsecurepath=0,zipinsecurepath=0" go build -v -ldflags '${LDFLAGS}' -o ${DIST_DIR}/${CLI_NAME} .

.PHONY: release-cli
release-cli: clean-debug 
	make BIN_NAME=${CLI_NAME}-darwin-amd64 GOOS=darwin release-all
	make BIN_NAME=${CLI_NAME}-darwin-arm64 GOOS=darwin GOARCH=arm64 release-all
	make BIN_NAME=${CLI_NAME}-linux-amd64 GOOS=linux release-all
	make BIN_NAME=${CLI_NAME}-linux-arm64 GOOS=linux GOARCH=arm64 release-all
	make BIN_NAME=${CLI_NAME}-linux-ppc64le GOOS=linux GOARCH=ppc64le release-all
	make BIN_NAME=${CLI_NAME}-linux-s390x GOOS=linux GOARCH=s390x release-all
	make BIN_NAME=${CLI_NAME}-windows-amd64.exe GOOS=windows release-all

# consolidated binary for cli, util, server, repo-server, controller
.PHONY: release-all
release-all: clean-debug
	CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} GODEBUG="tarinsecurepath=0,zipinsecurepath=0" go build -v -ldflags '${LDFLAGS}' -o ${DIST_DIR}/${BIN_NAME} .

# Installs all tools required for running codegen (Linux packages)
.PHONY: install-codegen-tools-local
install-codegen-tools-local:
	./hack/install.sh codegen-tools

.PHONY: help
help:
	@echo 'Note: Generally an item w/ (-local) will run inside docker unless you use the -local variant'
	@echo
	@echo 'Common targets'
	@echo
	@echo 'all -- make cli and image'
	@echo
	@echo 'build:'
	@echo '  build(-local)             -- compile go'
	@echo
