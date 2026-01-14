.PHONY: build test

PLUGIN_NAME := example-deployer

build:
	go build -buildmode=plugin -o $(PLUGIN_NAME).so .

test:
	go test -v ./...

