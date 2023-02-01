build:
	cd Server && go build
release: build
	mkdir js-linux
	cp Server/config.json js-linux/.
	cp Server/jServ js-linux/.
	tar -czvf jserv-linux-amd64.tar.gz js-linux/*
	rm -rf js-linux