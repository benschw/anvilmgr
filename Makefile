SHELL=/bin/bash
VERSION := $(shell cat VERSION)
ITTERATION := $(shell date +%s)

default: build

clean:
	rm -rf build
	rm -rf repo

deps:
	go get github.com/jteeuwen/go-bindata/...
	go get -t -v ./...

test:
	go test ./...

web-build:
	grunt build
	${GOPATH}/bin/go-bindata dist/...

server-build:
	mkdir -p build
	go build -o build/anvilmgr

build: web-build server-build


run: build
	./build/anvilmgr serve



install:
	install -t /usr/bin build/anvilmgr
	install -t /etc/init.d build/init/anvilmgr

deb:
	mkdir -p build/root/usr/bin
	mkdir -p build/root/etc/init.d
	mkdir -p build/root/etc/anvilmgr
	mkdir -p build/root/var/lib/puppet-anvil/modules
	mkdir -p build/root/var/log/anvilmgr
	
	cp build/anvilmgr build/root/usr/bin/anvilmgr
	cp anvilmgr.init build/root/etc/init.d/anvilmgr
	cp anvilmgr.yaml build/root/etc/anvilmgr/anvilmgr.yaml

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

