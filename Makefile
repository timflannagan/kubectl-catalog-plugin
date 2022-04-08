SHELL := /bin/bash

OUTPUT_BIN ?= bin/kubectl-catalog

.PHONY: bin/catalog
bin/catalog:
	@go build -o $(OUTPUT_BIN) cmd/main.go

plugin: bin/catalog
	@sudo cp $(OUTPUT_BIN) /usr/local/bin
	@kubectl catalog --help > /dev/null || (echo "failed to find the custom plugin in kubectl plugin path"; exit 1)
