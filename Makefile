
.PHONY: setup build resources lint clean

VERSION = $(shell git describe --always --dirty)
TIMESTAMP = $(shell git show -s --format=%ct)

default: build_darwin

setup: 
	go get -u github.com/jteeuwen/go-bindata/...

resources:
	go-bindata -pkg resources -o internal/resources/resources.go resources/...

build_darwin: resources
	GOOS=darwin GOARCH=amd64 go build -a -o ./build/gopow *.go
	zip ./build/gopow_darwin64.zip ./build/gopow

build_linux: resources
	GOOS=linux GOARCH=amd64 go build -a -o ./build/gopow *.go
	zip ./build/gopow_linux64.zip ./build/gopow

build_arm5: resources
	GOOS=linux GOARM=5 GOARCH=arm go build -a -o ./build/gopow *.go
	zip ./build/gopow_linux_arm5.zip ./build/gopow

build_arm7: resources
	GOOS=linux GOARM=7 GOARCH=arm go build -a -o ./build/gopow *.go
	zip ./build/gopow_linux_arm7.zip ./build/gopow

build_win: resources
	GOOS=windows GOARCH=amd64 go build -a -o ./build/gopow.exe *.go
	zip ./build/gopow_win64.zip ./build/gopow.exe

all: build_darwin build_linux build_arm5 build_arm7 build_win
	rm ./build/gopow
	rm ./build/gopow.exe

lint:
	golint .

clean:
	rm -f build
	rm -rf internal/resources

