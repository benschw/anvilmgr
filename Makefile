SHELL=/bin/bash

default: build

clean:
	rm -rf build
	rm -rf repo

deps:
	go get github.com/jteeuwen/go-bindata/...
	go get github.com/elazarl/go-bindata-assetfs/...
	go get -t -v ./...

test:
	go test ./...

web:
	/home/ben/go/bin/go-bindata dist/...

build: web
	mkdir -p build
	go build -o build/anvilmgr


run: deps clean test build
	./build/anvilmgr serve


.PHONY: build
