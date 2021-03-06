# rtl-gopow
Render tables from rtl_power to a nice heat map. Faster and easier to use than other tools, gopow does not require a scripting enviroment, dependencies or development enviroment to run. Just download the binary and execute. At the same time gopow offers 2-3.6 times the performance compared to script based tools, depending on input file.

## Availability
Since Go is easy to cross compile, this tool can be easily distributed as a binary without any dependencies. You'll find it under [Releases](https://github.com/dhogborg/rtl-gopow/releases) here on github. The following platforms are avalible as a ready to run binary file:

* OS X (x64)
* Linux (x64)
* Linux (arm5)
* Linux (arm7)
* Windows (x64)

https://github.com/dhogborg/rtl-gopow/releases

## Building
(only needed if making changes to the source, otherwise you can download a [release binary](https://github.com/dhogborg/rtl-gopow/releases))

### Install tooling
`make setup` will install [go-bindata](https://github.com/jteeuwen/go-bindata) which is needed to build the resources. Make sure $GOPATH/bin is in your $PATH so your shell can find the tool.

### Make resources
`make resources` will compile rtl-gopow/resources/ to rtl-gopow/internal/resources/resources.go. This file is disposable and will re-generate on every build. 

### Make build
Every platform has it's own make command. 
* `make build_darwin`
* `make build_linux`
* `make build_arm5`
* `make build_arm7`
* `make build_win`

All commands will compile a binary to rtl-gopow/build/gopow[.exe] and create a distributable zip file. For more info on cross compilation see [this document](http://dave.cheney.net/2013/07/09/an-introduction-to-cross-compilation-with-go-1-1).

`make all` will build all platforms and create zip files in rtl-gopow/build/ and then remove the executable files.

### go build/run
Before runnning `go build *.go` or `go run *.go` make sure the resources has been generated by `make resources`. 

## Performance
A render of a 600 MB csv file takes about 2 minutes on a 2,4 GHz Intel Core i5. There is still lots of room for improvement on that though. Memory usage is quite horrid.

Compared to script based tools gopow will run more than 2x faster for smaller files, and more than 3x faster for bigger files.

## Options
```
GLOBAL OPTIONS:
   --input, -i      CSV input file generated by rtl_power [required]
   --output, -o     Output file, default same as input file with new extension
   --format, -f 'png'   Output file format, default png [png,jpeg]
   --verbose        Enable more verbose output
   --no-annotations Disabled annotations such as time and frequency scales
   --help, -h       show help
   --version, -v    print the version

```

## Demo
Here is an render of rtl_power tool scanning 80-90 MHz during 2.5 hours moving in a car. ![80-90 MHz](http://i.imgur.com/knkzLXO.jpg).