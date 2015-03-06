SHELL=/bin/bash
VERSION := $(shell cat VERSION)
ITTERATION := $(shell date +%s)

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
	grunt build
	${GOPATH}/bin/go-bindata dist/...

build: deps web
	mkdir -p build/init
	cp anvilmgr.init build/init/anvilmgr
	mkdir -p build
	go build -o build/anvilmgr


run: deps clean test build
	./build/anvilmgr serve



install:
	install -t /usr/bin build/anvilmgr
	install -t /etc/init.d build/init/anvilmgr

deb: build
	fpm -s dir -t deb -n anvilmgr -v $(VERSION) -p build/anvilmgr-amd64.deb \
		--deb-priority optional \
		--category util \
		--force \
		--iteration $(ITTERATION) \
		--deb-compression bzip2 \
		--url http://git.bvops.net/projects/AUTO/repos/anvilmgr/browse \
		--description "Anvil Manager web app and api for anvil-puppet" \
		-m "Ben Schwartz <benschw@gmail.com>" \
		--license "Apache License 2.0" \
		--vendor "bancvue.com" -a amd64 \
		build/root/=/

