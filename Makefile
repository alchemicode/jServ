build:
	cd Server && go build
release: build
	mkdir js-linux
	cp Server/config.json js-linux/.
	cp Server/jServ js-linux/.
	tar -czvf jserv-linux-amd64.tar.gz js-linux/*
	rm -rf js-linux
deploy: build
	[[ ! -d /etc/jserv ]] && mkdir /etc/jserv
	[[ -f /usr/bin/jServ ]] && rm /usr/bin/jServ
	cp Server/jServ /usr/bin/.
	cp Server/config.json /etc/jserv/.
	# systemd unit needed
