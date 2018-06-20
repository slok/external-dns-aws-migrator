
# Name of this service/application
SERVICE_NAME := external-dns-aws-migrator

# Path of the go service inside docker
DOCKER_GO_SERVICE_PATH := /go/src/github.com/slok/external-dns-aws-migrator

# Shell to use for running scripts
SHELL := $(shell which bash)

# Get OS
OSTYPE := $(shell uname)

# Get docker path or an empty string
DOCKER := $(shell command -v docker)

# Get the main unix group for the user running make (to be used by docker-compose later)
GID := $(shell id -g)

# Get the unix user id for the user running make (to be used by docker-compose later)
UID := $(shell id -u)

# Commit version from git
VERSION=$(shell git describe --tags --always)

# cmds
UNIT_TEST_CMD := ./hack/scripts/unit-test.sh
INTEGRATION_TEST_CMD := ./hack/scripts/integration-test.sh
MOCKS_CMD := ./hack/scripts/mockgen.sh
DOCKER_RUN_CMD := docker run --env ostype=$(OSTYPE) -v ${PWD}:$(DOCKER_GO_SERVICE_PATH) --rm -it $(SERVICE_NAME)
BUILD_BINARY_CMD := VERSION=${VERSION} ./hack/scripts/build.sh
CI_RELEASE_CMD := ./hack/scripts/travis-release.sh

# environment dirs
DEV_DIR := docker/dev

# The default action of this Makefile is to build the development docker image
.PHONY: default
default: build

# Test if the dependencies we need to run this Makefile are installed
.PHONY: deps-development
deps-development:
ifndef DOCKER
	@echo "Docker is not available. Please install docker"
	@exit 1
endif

# Build the development docker image
.PHONY: build
build:
	docker build -t $(SERVICE_NAME) --build-arg uid=$(UID) --build-arg  gid=$(GID) -f ./docker/dev/Dockerfile .

# Shell the development docker image
.PHONY: build
shell: build
	$(DOCKER_RUN_CMD) /bin/bash

# Build production stuff.
.PHONY: build-binary
build-binary:
	$(DOCKER_RUN_CMD) /bin/sh -c '$(BUILD_BINARY_CMD)'

# Test stuff in dev
.PHONY: unit-test
unit-test: build
	$(DOCKER_RUN_CMD) /bin/sh -c '$(UNIT_TEST_CMD)'
.PHONY: integration-test
integration-test: build
	$(DOCKER_RUN_CMD) /bin/sh -c '$(INTEGRATION_TEST_CMD)'
.PHONY: test
test: integration-test

# Test stuff in ci
.PHONY: ci-unit-test
ci-unit-test:
	$(UNIT_TEST_CMD)
.PHONY: ci-integration-test
ci-integration-test:
	$(INTEGRATION_TEST_CMD)
.PHONY: ci
ci: ci-integration-test

.PHONY: ci-release
ci-release:
	$(CI_RELEASE_CMD)

# Mocks stuff in dev
.PHONY: mocks
mocks: build
	$(DOCKER_RUN_CMD) /bin/sh -c '$(MOCKS_CMD)'