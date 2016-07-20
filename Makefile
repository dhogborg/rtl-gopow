
.PHONY: setup build resources lint clean

VERSION = $(shell git describe --always --dirty)
TIMESTAMP = $(shell git show -s --format=%ct)

default: build

setup: 
	go get -u github.com/jteeuwen/go-bindata/...
	go get github.com/tools/godep

resources:
	go-bindata -pkg resources -o internal/resources/resources.go resources/...

build: resources
	godep go build -o ./build/gopow *.go

all: build_darwin build_linux build_arm5 build_arm7 build_win64 build_win32
	rm ./build/gopow
	rm ./build/gopow.exe

build_darwin: resources
	GOOS=darwin GOARCH=amd64 godep go build -a -o ./build/gopow *.go
	zip ./build/gopow_darwin64.zip ./build/gopow

build_linux: resources
	GOOS=linux GOARCH=amd64 godep go build -a -o ./build/gopow *.go
	zip ./build/gopow_linux64.zip ./build/gopow

build_arm5: resources
	GOOS=linux GOARM=5 GOARCH=arm godep go build -a -o ./build/gopow *.go
	zip ./build/gopow_linux_arm5.zip ./build/gopow

build_arm7: resources
	GOOS=linux GOARM=7 GOARCH=arm godep go build -a -o ./build/gopow *.go
	zip ./build/gopow_linux_arm7.zip ./build/gopow

build_win64: resources
	GOOS=windows GOARCH=amd64 godep go build -a -o ./build/gopow.exe *.go
	zip ./build/gopow_win64.zip ./build/gopow.exe

build_win32: resources
	GOOS=windows GOARCH=386 godep go build -a -o ./build/gopow.exe *.go
	zip ./build/gopow_win32.zip ./build/gopow.exe

lint:
	golint .

# Save dependencies to vendor folder
deps:
	- rm -r vendor Godeps
	godep save ./...

deps_restore:
	godep restore ./...
	- rm -r vendor

clean:
	- rm -r build
	- rm -rf internal/resources

