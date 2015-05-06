	curl -sL https://deb.nodesource.com/setup | sudo bash -
	sudo apt-get install nodejs
	sudo npm install -g bower grunt-cli

	make deps
	make test
	make build

	curl -i -F "artifact=@./puppetlabs-stdlib-4.5.1.tar.gz" localhost:8080/api/repo/puppetlabs/stdlib
	curl -i localhost:8080/api/repo/puppetlabs/stdlib
