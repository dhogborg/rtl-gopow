
.PHONY: setup build resources lint clean

VERSION = $(shell git describe --always --dirty)
TIMESTAMP = $(shell git show -s --format=%ct)

default: build

setup: 
	go get -u github.com/jteeuwen/go-bindata/...

resources:
	go-bindata -pkg resources -o internal/resources/resources.go resources/...

build:
	go build -o ./build/gopow *.go

all: build_darwin build_linux build_arm5 build_arm7 build_win64 build_win32
	rm ./build/gopow
	rm ./build/gopow.exe

build_darwin:
	GOOS=darwin GOARCH=amd64 go build -a -o ./build/gopow *.go
	zip ./build/gopow_darwin64.zip ./build/gopow

build_linux:
	GOOS=linux GOARCH=amd64 go build -a -o ./build/gopow *.go
	zip ./build/gopow_linux64.zip ./build/gopow

build_arm5:
	GOOS=linux GOARM=5 GOARCH=arm go build -a -o ./build/gopow *.go
	zip ./build/gopow_linux_arm5.zip ./build/gopow

build_arm7:
	GOOS=linux GOARM=7 GOARCH=arm go build -a -o ./build/gopow *.go
	zip ./build/gopow_linux_arm7.zip ./build/gopow

build_win64:
	GOOS=windows GOARCH=amd64 go build -a -o ./build/gopow.exe *.go
	zip ./build/gopow_win64.zip ./build/gopow.exe

build_win32:
	GOOS=windows GOARCH=386 go build -a -o ./build/gopow.exe *.go
	zip ./build/gopow_win32.zip ./build/gopow.exe

clean:
	- rm -r build
	- rm -rf internal/resources

