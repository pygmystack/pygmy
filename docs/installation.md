#Installation of Pygmy

## Prerequisites
Make sure you have the following dependencies installed:

* Docker, see [the official guides](https://docs.docker.com/engine/installation/) on how to install docker on your system.
* Go (optional), see [the official guides](https://golang.org/doc/install)

## Installation

### Installing from a precompiled binary

Releases on GitHub accompany binaries [available for download](https://github.com/pygmystack/pygmy/releases).

To install it, put the binary into your system's `$PATH` environment variable and make it executable.

The following is an example of how you would do this, note the URL and location may change depending on your needs.
```console
$ wget https://github.com/pygmystack/pygmy/releases/download/v0.1.0/pygmy-go-darwin
$ mv ./pygmy-go-darwin /usr/local/bin/pygmy
$ chmod u+x /usr/local/bin/pygmy
```

### Build from source

Pygmy comes with a Make file, which you can simply run `make build && make clean` to build binaries for Linux (amd64), Windows (x86) & MacOS (Darwin).

From here you can follow the guidance to install the relevant executable in the `builds/` folder usign the instructions above.

### Installing from source

The installation of `pygmy` is fairly simple and can be accomplished via the go toolchain

```console
$ go get github.com/pygmystack/pygmy
```