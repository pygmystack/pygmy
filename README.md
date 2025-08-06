# Pygmy

[![Stability](https://img.shields.io/badge/stability-stable-green.svg)]()
![goreleaser](https://github.com/pygmystack/pygmy/workflows/goreleaser/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/pygmystack/pygmy)](https://goreportcard.com/report/github.com/pygmystack/pygmy)
[![GoDoc](https://godoc.org/github.com/pygmystack/pygmy?status.svg)](https://godoc.org/github.com/pygmystack/pygmy)

This is an application written in Go which is a proposed replacement for [Pygmy](https://pygmy.readthedocs.io/en/master/)
currently written in Ruby. The goal is to provide a better cross-platform experience
for various users running Lagoon, as well as much greater control over configuration
options via YAML.

Please see the existing [Pygmy documentation](https://pygmy.readthedocs.io) for more information
about Pygmy as this is designed to be a drop-in replacement.

## Early testing

We welcome testers of this tool. You will probably be an existing user of Pygmy who
can verify the same functionality, or perhaps who has had trouble installing Pygmy in the
past on Windows.

## Is Pygmy running?

These instructions will currently install the new version as `pygmy` so that the
old version is still available if you have installed it. With no Pygmy running,
you should get "connection refused" when attempting to connect to the local amazee network.

```
curl --HEAD http://myproject.docker.amazee.io
curl: (7) Failed to connect to myproject.docker.amazee.io port 80: Connection refused
```

## Installation

These instructions will build Linux, MacOS and Windows binaries of Pygmy on MacOS,
and then test the MacOS version. M1 and arm64 images are available and supported

### Using Homebrew

Homebrew is the recommended way to install pygmy and keep it up to date on compatible systems.

**Works for**: Linux & MacOS (also WSL-based systems)

```shell
brew tap pygmystack/pygmy;
brew install pygmy;
```

### Compile from source

Ensure to select the correct build for your OS and architecture in the `cp` command.

**Works for**: Linux, MacOS & Windows

```shell
git clone https://github.com/pygmystack/pygmy.git && cd pygmy;
make build;
cp ./builds/pygmy-darwin /usr/local/bin/pygmy;
chmod +x /usr/local/bin/pygmy;
```

Pygmy is now executable as `pygmy`. Now start Pygmy and use the new `status` command.
If you still need to use the previous `pygmy`, cp the binary to a different name (e.g. pygmy-go)

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
$ curl --HEAD http://myproject.docker.amazee.io
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
$ curl --HEAD http://FUBARNOTINDAHOUSE.docker.amazee.io
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
   git clone https://github.com/pygmystack/pygmy.git pygmy && cd pygmy
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

## Generating local certificates with mkcert

To generate a certificate pair for use with HAProxy, you typically use the default domain `*.docker.amazee.io`, but you can substitute any domain as needed for your local development environment:

1. **Install mkcert**  
   Follow instructions at [mkcert GitHub](https://github.com/FiloSottile/mkcert).  
   Example for macOS:  
   ```shell
   brew install mkcert
   mkcert -install
   ```

2. **Generate a wildcard certificate and key**  
   Run:  
   ```shell
   mkcert "*.docker.amazee.io"
   ```
   This creates `_wildcard.docker.amazee.io-key.pem` (certificate) and `_wildcard.docker.amazee.io-key.pem` (private key).

3. **Combine certificate and key for HAProxy**  
   Run:  
   ```shell
   mkdir -p ~/pygmy/
   cat _wildcard.docker.amazee.io.pem _wildcard.docker.amazee.io-key.pem > ~/pygmy/server.pem
   ```
   This combined certificate can either be passed as an argument when starting pygmy
   
   ```
   pygmy up --tls-cert=/location_to/haproxy.pem
   ```
   or can be stored in the default location of `~/pygmy/server.pem`.


**Notes:**  
- mkcert certificates are for local development and testing only.

## Releasing
 
We use GitHub Actions for simulating the automated release tagging locally. Using [Act](https://github.com/nektos/act) locally, you can simulate this process and have the same build artifacts in your `dist` folder.
This process will inject the appropriate values into the version logic. To start the process, just run `act`!
