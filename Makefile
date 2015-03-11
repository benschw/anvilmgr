SHELL=/bin/bash
VERSION := $(shell cat VERSION)
ITTERATION := $(shell date +%s)

default: build

clean:
	rm puppetlabs-stdlib-4.5.1.tar.gz
	rm bindata.go
	rm -rf build
	rm -rf repo
	rm -rf webapp/node_modules
	rm -rf webapp/dist
	rm -rf webapp/bower_components

deps-web:
	cd webapp; \
	npm install; \
	bower install

deps-server:
	go get -u github.com/jteeuwen/go-bindata/...
	go get -u -t -v ./...
	# just for test
	wget -q https://forgeapi.puppetlabs.com/v3/files/puppetlabs-stdlib-4.5.1.tar.gz -O puppetlabs-stdlib-4.5.1.tar.gz

deps: deps-web deps-server



test:
	go test ./...

build-web:
	cd webapp; \
	grunt build

build-bindata:
	${GOPATH}/bin/go-bindata webapp/dist/...

build-server:
	mkdir -p build
	go build -o build/anvilmgr


build: build-web build-bindata build-server


run:
	./build/anvilmgr -config ./anvilmgr.yaml serve

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

