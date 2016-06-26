APPNAME = dops
VERSION = 0.1.0-dev

setup:
	glide install

build-all: build-mac build-linux

build:
	go build -o ${APPNAME} .

build-linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-s -X main.Version=${VERSION}" -v -o ${APPNAME}-linux-amd64 .

build-mac:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-s -X main.Version=${VERSION}" -v -o ${APPNAME}-darwin-amd64 .

ci:
	APPNAME=${APPNAME} bin/ci-run.sh

clean:
	rm -f ${APPNAME}
	rm -f ${APPNAME}-linux-amd64
	rm -f ${APPNAME}-darwin-amd64

all:
	setup
	build
	install

test:
	go test -v github.com/ashwanthkumar/dops
	go test -v github.com/ashwanthkumar/dops/config
	go test -v github.com/ashwanthkumar/dops/server
	go test -v github.com/ashwanthkumar/dops/server/torrent
	go test -v github.com/ashwanthkumar/dops/server/storage
	go test -v github.com/ashwanthkumar/dops/server/engine
	go test -v github.com/ashwanthkumar/dops/server/engine/docker

test-only:
	go test -v github.com/ashwanthkumar/dops/${name}

install: build
	sudo install -d /usr/local/bin
	sudo install -c ${APPNAME} /usr/local/bin/${APPNAME}

uninstall:
	sudo rm /usr/local/bin/${APPNAME}
