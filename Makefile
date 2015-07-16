
.PHONY: setup build resources lint clean

VERSION = $(shell git describe --always --dirty)
TIMESTAMP = $(shell git show -s --format=%ct)

default: build

setup: 
	go get -u github.com/jteeuwen/go-bindata/...

resources:
	go-bindata -pkg resources -o internal/resources/resources.go resources/...

build: resources
	go build -o ./gopow *.go

lint:
	golint .

clean:
	rm -f gopow
	rm -rf internal/resources

