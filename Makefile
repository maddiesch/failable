ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: clean
clean:
	cd ${ROOT_DIR} && go clean

.PHONY: test
test: clean
	cd ${ROOT_DIR} && go test -v ./...
