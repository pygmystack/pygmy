#Installation of Pygmy

## Prerequisites
Make sure you have the following dependencies installed:

* Docker, see [the official guides](https://docs.docker.com/engine/installation/) on how to install docker on your system.
* Go (optional), see [the official guides](https://golang.org/doc/install)

## Support
Any platform that supports Docker & Go will be able to run `pygmy-go`.

If it's not readily available, you can [build](#building) `pygmy-go` onto any supported platform.

## Installation
### Homebrew
Windows (WSL), Linux and MacOS all support Homebrew.

You are able to use tap `fubarhouse/pygmy-go` and install `pygmy-go`.

This will build the latest tagged release from source.

````console
$ brew tap fubarhouse/pygmy-go && brew install pygmy-go
````

### Release binaries
Releases on GitHub accompany binaries [available for download](https://github.com/fubarhouse/pygmy-go/releases) for Windows, MacOS and Linux.

To install it, put the binary into a directory in your system's `$PATH` environment variable and make it executable.

The following is an example of how you would do this, note the URL and location may change depending on your needs.
```console
$ wget https://github.com/fubarhouse/pygmy-go/releases/download/v0.2.0/pygmy-go-darwin
$ mv ./pygmy-go-darwin /usr/local/bin/pygmy
$ chmod u+x /usr/local/bin/pygmy
```

### PKGBUILD
If Homebrew on Linux isn't your style, you can use the [PKGBUILD files](https://github.com/fubarhouse/pygmy-go.pkgbuild) to build the application.

This process is specifically for Arch Linux, however the repository supports a process to build the Linux artifacts on any architecture via Docker.

Please see the [repository](https://github.com/fubarhouse/pygmy-go.pkgbuild) for further details, as this has not yet made it onto the AUR.

### Other
* Other Linux flavours such as `deb` and `rpm` are not yet supported.
* Windows distribution platforms such as [chocolatey](https://chocolatey.org/docs/installation) are not yet supported, please use homebrew, install/build from make/source or binary for Windows installations. 

## Building
### Using makefile
Pygmy comes with a Make file, which you can simply run `make build && make clean` to build binaries for Linux (amd64), Windows (x86) & MacOS (Darwin).

From here you can follow the guidance to install the relevant executable in the `builds/` folder using the instructions above.

### From source
The installation of `pygmy` is fairly simple and can be accomplished via the go toolchain

```console
$ go get github.com/fubarhouse/pygmy-go
```