.PHONY: build test clean

PLUGIN_NAME := example-deployer

build:
	go build -buildmode=plugin -o $(PLUGIN_NAME).so .

test:
	go test -v ./...

clean:
	rm -rf $(DIST_DIR)
