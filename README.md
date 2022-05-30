# Pygmy

[![Stability](https://img.shields.io/badge/stability-stable-green.svg)]()
![goreleaser](https://github.com/pygmystack/pygmy/workflows/goreleaser/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/pygmystack/pygmy)](https://goreportcard.com/report/github.com/pygmystack/pygmy)
[![GoDoc](https://godoc.org/github.com/pygmystack/pygmy?status.svg)](https://godoc.org/github.com/pygmystack/pygmy)

This is an application written in Go which is a proposed replacement for [Pygmy](https://pygmy-readthedocs.io/en/master/)
currently written in Ruby. The goal is to provide a better cross-platform experience
for various users running Lagoon, as well as much greater control over configuration
options via YAML.

Please see the existing [Pygmy documentation](https://pygmy.readthedocs.io) for more information
about Pygmy as this is designed to be a drop-in replacement.

## Making Pygmy always use `docker.amazee.io`

To make sure Pygmy always uses the legacy `amazeeio-*` container names and the
`docker.amazee.io` domain, just set the environment variable `PYGMY_DOMAIN` to
`docker.amazee.io`

If the `amazeeio-network' is created and ready for use, Pygmy will
automatically set these variables at runtime. If not, it will now use the
container names `pygmy-*` and the domain `*.pygmy.site`.

```shell
echo 'export PYGMY_DOMAIN=docker.amazee.io' >> ~/.bashrc
```

## Is Pygmy running?

These instructions will currently install the new version as `pygmy` so that the
old version is still available if you have installed it. With no Pygmy running,
you should get "connection refused" when attempting to connect to the local amazee network.

```
curl --HEAD http://myproject.pygmy.site
curl: (7) Failed to connect to myproject.pygmy.site port 80: Connection refused
```

## Installation

These instructions will build Linux, MacOS and Windows binaries of Pygmy on MacOS,
and then test the MacOS version. Apple M1 is not officially supported. Upstream
docker images provided by Amazee are not yet provided. We'll not deviate from these
images as it would deviate from what is officially supported. You can however provide
configuration to switch out the `haproxy` image for a compatible one if you'd like.

### Compile from source

**Works for**: Linux, MacOS & Windows

```shell
git clone https://github.com/pygmystack/pygmy.git && cd pygmy;
make build;
cp ./builds/pygmy-darwin /usr/local/bin/pygmy;
chmod +x /usr/local/bin/pygmy;
```

Pygmy is now an executable as `pygmy`, while any existing Pygmy is still executable
as `pygmy`. Now start Pygmy and use the new `status` command.

### Using Homebrew

**Works for**: Linux & MacOS

```shell
brew tap pygmystack/pygmy;
brew install pygmy;
```

### Using the AUR

**Works for**: [Arch-based Linux Distributions](https://wiki.archlinux.org/title/Arch-based_distributions) (Manjaro, Elementary, ArcoLinux etc)

[pygmy](https://aur.archlinux.org/packages/pygmy/), [pygmy-bin](https://aur.archlinux.org/packages/pygmy-bin/) and
[pygmy-git](https://aur.archlinux.org/packages/pygmy-git/) are available via the Arch User Repository for Arch-based
Linux distributions on the community stream. Unfortunately, Pygmy is not yet available via other distribution methods,
so it is otherwise recommended to use homebrew to install it, download a pre-compiled binary from the releases page, or
to compile from source.  

```shell
# Freshly compile the latest release:
yay -S pygmy;
# Download the latest release precompiled:
yay -S pygmy-bin;
# Download and compile the latest HEAD from GitHub on the main branch:
yay -S pygmy-git;
```

## Usage

If you have an Amazee Lagoon project running, you can test the web address and
expect a `HTTP/1.1 200 OK` response.

```
$ curl --HEAD http://myproject.pygmy.site
HTTP/1.1 200 OK
Server: openresty
Content-Type: text/html; charset=UTF-8
Cache-Control: must-revalidate, no-cache, private
Date: Mon, 11 Nov 2019 11:19:29 GMT
X-UA-Compatible: IE=edge
Content-language: en
X-Content-Type-Options: nosniff
X-Frame-Options: SAMEORIGIN
X-Drupal-Cache-Tags: config:honeypot.settings config:system.site config:user.role.anonymous http_response rendered
X-Drupal-Cache-Contexts: languages:language_interface theme url.path url.query_args user.permissions user.roles:authenticated
Expires: Sun, 19 Nov 1978 05:00:00 GMT
Vary:
X-Frame-Options: SameOrigin
```

If your project is not running you should expect a 503 response:

```
$ curl --HEAD http://FUBARNOTINDAHOUSE.pygmy.site
HTTP/1.0 503 Service Unavailable
Cache-Control: no-cache
Connection: close
Content-Type: text/html
```

Thanks for testing, please post issues and successes in the queue.

## Local development

To run full regression tests locally, you can follow this process if you have `cmake`, `git` and `go` installed. This 
will prevent a significant amount of build failures and problems after committing.

It will use `dind` and your local daemon to walk through several tests which should pass.

1. First clone the project:
   ```
   git clone https://github.com/pygmystack/pygmy-git pygmy && cd pygmy
   ```
2. Perform any updates as required.
3. Clean the environment.
   ```
   go run main.go clean
   ```
4. Build the project.
   ```
   make
   ```
5. Test the project prior to commiting.
   ```
   go test -v
   ```
 
## Releasing
 
We use GitHub Actions for simulating the automated release tagging locally. Using [Act](https://github.com/nektos/act) locally, you can simulate this process and have the same build artifacts in your `dist` folder.
This process will inject the appropriate values into the version logic. To start the process, just run `act`!