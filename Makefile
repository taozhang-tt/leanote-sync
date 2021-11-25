Bin=leanote-sync

build:
	go build -o ${Bin}
	mv ${Bin} /usr/local/bin/
	sudo cp config.json /etc/leanote-sync.json	
	sudo chmod 0666 /etc/leanote-sync.json
