
VERSION = $(shell git describe --always --dirty)
TIMESTAMP = $(shell git show -s --format=%ct)

default: build

.PHONY: setup
setup: 
	go get -u github.com/jteeuwen/go-bindata/...

.PHONY: resources
resources:
	go-bindata -pkg resources -o internal/resources/resources.go resources/...

.PHONY: build
build:
	go build -o ./build/gopow *.go

.PHONY: all
all: build_darwin_x86 build_darwin_arm64 build_linux build_arm5 build_arm7 build_win64 build_win32
	rm ./build/gopow
	rm ./build/gopow.exe

.PHONY: build_darwin_x86
build_darwin_x86:
	GOOS=darwin GOARCH=amd64 go build -a -o ./build/gopow *.go
	zip ./build/gopow_darwin_amd64.zip ./build/gopow

.PHONY: build_darwin_arm64
build_darwin_arm64:
	GOOS=darwin GOARCH=arm64 go build -a -o ./build/gopow *.go
	zip ./build/gopow_darwin_arm64.zip ./build/gopow

.PHONY: build_linux
build_linux:
	GOOS=linux GOARCH=amd64 go build -a -o ./build/gopow *.go
	zip ./build/gopow_linux_amd64.zip ./build/gopow

.PHONY: build_arm5
build_arm5:
	GOOS=linux GOARM=5 GOARCH=arm go build -a -o ./build/gopow *.go
	zip ./build/gopow_linux_arm5.zip ./build/gopow

.PHONY: build_arm7
build_arm7:
	GOOS=linux GOARM=7 GOARCH=arm go build -a -o ./build/gopow *.go
	zip ./build/gopow_linux_arm7.zip ./build/gopow

.PHONY: build_win64
build_win64:
	GOOS=windows GOARCH=amd64 go build -a -o ./build/gopow.exe *.go
	zip ./build/gopow_win64.zip ./build/gopow.exe

.PHONY: build_win32
build_win32:
	GOOS=windows GOARCH=386 go build -a -o ./build/gopow.exe *.go
	zip ./build/gopow_win32.zip ./build/gopow.exe

.PHONY: clean
clean:
	- rm -r build
	- rm -rf internal/resources

