SHELL := /bin/bash

OUTPUT_BIN ?= bin/kubectl-catalog

.PHONY: bin/create
bin/create:
	@go build -o $(OUTPUT_BIN) main.go

plugin: bin/create
	@sudo cp $(OUTPUT_BIN) /usr/local/bin
	@kubectl catalog --help > /dev/null || (echo "failed to find the custom plugin in kubectl plugin path"; exit 1)
