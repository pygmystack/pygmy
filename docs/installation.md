#Installation of Pygmy

## Prerequisites
Make sure you have the following dependencies installed:

* Docker, see [the official guides](https://docs.docker.com/engine/installation/) on how to install docker on your system.
* Go (optional), see [the official guides](https://golang.org/doc/install)

## Installation

Pygmy now has a growing user-base on Windows, MacOS and Linux and supports a number of growing ways to install `pygmy`, and they're detailed below:

* [Install on Linux](#linux) via PKGBUILD files (Architecture-independent support)
* [Install on MacOS](#macos) & Linux via Homebrew
* [Install on Windows](#windows)
* [Install from pre-compiled binaries](#installing-from-a-precompiled-binary)
* [Install from source](#build-from-source-make)

### Linux

For Linux installations you can use the [PKGBUILD files](https://github.com/fubarhouse/pygmy-go.pkgbuild) to build the application.
   
This process is specifically for Arch Linux, however the repository supports a process to build the Linux artifacts on any architecture via Docker.
   
Please see the [repository](https://github.com/fubarhouse/pygmy-go.pkgbuild) for further details, as this has not yet made it onto the AUR.
   
Other Linux flavours such as `deb` and `rpm` are not yet supported, however you can install it on Linux using `brew`, see the MacOS instructions for `brew` install instructions.

### MacOS

For MacOS installations, you are able to use the namespace to tap and install `pygmy-go`.

This will build the latest tagged release from source.

````console
$ brew tap fubarhouse/pygmy-go && brew install pygmy-go
````

### Windows

Windows-specific package delivery systems aren't supported yet, so it is recommended to either build from source or download a binary from the [releases page](https://github.com/fubarhouse/pygmy-go/releases).

### Installing from a precompiled binary

Releases on GitHub accompany binaries [available for download](https://github.com/fubarhouse/pygmy-go/releases).

To install it, put the binary into your system's `$PATH` environment variable and make it executable.

The following is an example of how you would do this, note the URL and location may change depending on your needs.
```console
$ wget https://github.com/fubarhouse/pygmy-go/releases/download/v0.1.0/pygmy-go-darwin
$ mv ./pygmy-go-darwin /usr/local/bin/pygmy
$ chmod u+x /usr/local/bin/pygmy
```

### Build from source (make)

Pygmy comes with a Make file, which you can simply run `make build && make clean` to build binaries for Linux (amd64), Windows (x86) & MacOS (Darwin).

From here you can follow the guidance to install the relevant executable in the `builds/` folder using the instructions above.

### Installing from source (go)

The installation of `pygmy` is fairly simple and can be accomplished via the go toolchain

```console
$ go get github.com/fubarhouse/pygmy-go
```