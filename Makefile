MAIN=searchQuery
SHELL=bash

export GOOS?=$(shell go env GOOS)
export GOARCH?=$(shell go env GOARCH)
thisOS:=$(shell uname -s)

HEALTHCHECK_ENDPOINT?=/healthcheck
CMD_DIR?=cmd
BUILD?=build
BUILD_ARCH=$(BUILD)/$(GOOS)-$(GOARCH)
    BIN_DIR?=bin

DATE:=$(shell date '+%Y%m%d-%H%M%S')


build:
	@mkdir -p $(BUILD_ARCH)/$(BIN_DIR)
	go build -i -o $(BUILD_ARCH)/$(BIN_DIR)/$(MAIN) $(CMD_DIR)/$(MAIN)/main.go
	cp -r templates $(BUILD_ARCH)/

package:
	docker build . -f docker/Dockerfile -t guidof/documentindexer-searchquery


deploy: package
	docker push guidof/documentindexer-searchquery:latest


hash:
	@git rev-parse --short HEAD



.PHONY: build package
