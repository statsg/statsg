# statsg

[![codebeat badge](https://codebeat.co/badges/53703c14-bc44-489d-b0ea-7e3d3b2d8d82)](https://codebeat.co/projects/github-com-statsg-statsg) [![circleci` badge](https://circleci.com/gh/statsg/statsg.svg?style=shield&circle-token=456ad54019146fcaa685adb00e91d7bb73f9d58e)](https://circleci.com/gh/statsg/statsg.svg?style=shield&circle-token=456ad54019146fcaa685adb00e91d7bb73f9d58e) [![Go Report Card](https://goreportcard.com/badge/github.com/statsg/statsg)](https://goreportcard.com/report/github.com/statsg/statsg) [![codecov](https://codecov.io/gh/statsg/statsg/branch/master/graph/badge.svg)](https://codecov.io/gh/statsg/statsg) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A reimplementation of statsd to use CoAP as the transport mechanism, and dynamically compressed metric keys.

## Getting started

This project requires Go to be installed. On OS X with Homebrew you can just run `brew install go`, then install [glide](https://github.com/Masterminds/glide).

Running it then should be as simple as:

```console
$ make
$ ./bin/statsg
```

### Testing

``make test``

## Concept

As I always like to consider bandwidth for my services, I thought I'd take a look at a low-hanging
fruit such as statsd, and how I could possibly consume less bandwidth while still receiving relevant
metrics. Two things came to mind. Using a different transport protocol, and creating a easy method
to compress the metric key.

### CoAP for transport

CoAP, or [Constrained Application Protocol](http://coap.technology/) as defined by RFC 7252 is

> The Constrained Application Protocol (CoAP) is a specialized web transfer protocol for use with constrained nodes and constrained networks in the Internet of Things. The protocol is designed for machine-to-machine (M2M) applications such as smart energy and building automation.

The idea is to use a minimalistic transport protocol, which is just UDP on IP with a 4-byte fixed header and compact encoding.

### Compressing the metric key

In general, statsd metric keys are dot-notated plaintext which follows a form such as 
`com.example.api.service_name.some.specific.metric` which is nice for easy `tcpdump`
sessions, but can be a little wasteful in terms of bandwidth. Building the representation
as a tree could simplify it as 

```
| Field | Size    | Representation | Plaintext            |
|----------------------------------------------------------
|     1 | 1 byte  |    [0 - 2^8)   | com.example.api      |
|     2 | 1 byte  |    [0 - 2^8)   | service_name         |
|     3 | 1 byte  |    [0 - 2^8)   | some.specific.metric |
|     4 | n bytes | You metric val | Plaintext metric val |
```

This would yield 255 FQDNs, 255 services per FQDN, and 255 metrics per service yielding
16,777,216 keys in total expressed in 5 bytes.

### Todo

- [ ] Persistance for the keyspace
- [ ] Windows make file
- [ ] Statsd output writer
- [ ] Direct WhisperDB writer

### Presentations

- [Post Hack-a-thon presentation](https://docs.google.com/presentation/d/18qABweNyxOXPynj0iUH6-h6LAdr6r5g-mxJRjKDU_wE/)
