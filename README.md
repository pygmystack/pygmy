# Pygmy

[![Stability](https://img.shields.io/badge/stability-stable-green.svg)]()
[![Travis CI](https://travis-ci.org/fubarhouse/pygmy-go.svg?branch=master)]()
[![Go Report Card](https://goreportcard.com/badge/github.com/fubarhouse/pygmy-go)](https://goreportcard.com/report/github.com/fubarhouse/pygmy-go)
[![GoDoc](https://godoc.org/github.com/fubarhouse/pygmy-go?status.svg)](https://godoc.org/github.com/fubarhouse/pygmy-go)

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

These instructions will currently install the new version as `pygmy-go` so that the
old version is still available if you have installed it. With no Pygmy running,
you should get "connection refused" when attempting to connect to the local amazee network.

```
curl --HEAD http://myproject.docker.amazee.io
curl: (7) Failed to connect to myproject.docker.amazee.io port 80: Connection refused
```

## Installation (OSX specific)

These instructions will build Linux, OSX and Windows binaries of Pygmy on OSX,
and then test the OSX version.

1. `git clone https://github.com/fubarhouse/pygmy-go.git && cd pygmy-go`
2. `make build`
3. `cp ./builds/pygmy-go-darwin /usr/local/bin/pygmy-go && chmod +x /usr/local/bin/pygmy-go`

Pygmy is now an executable as `pygmy-go`, while any existing Pygmy is still executable
as `pygmy`. Now start Pygmy and use the new `status` command.

4. `pygmy-go up`
5. `pygmy-go status`

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
 