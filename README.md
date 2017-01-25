# statsg

[![codebeat badge](https://codebeat.co/badges/53703c14-bc44-489d-b0ea-7e3d3b2d8d82)](https://codebeat.co/projects/github-com-statsg-statsg) [![circleci` badge](https://circleci.com/gh/statsg/statsg.svg?style=shield&circle-token=456ad54019146fcaa685adb00e91d7bb73f9d58e)](https://circleci.com/gh/statsg/statsg.svg?style=shield&circle-token=456ad54019146fcaa685adb00e91d7bb73f9d58e) [![Go Report Card](https://goreportcard.com/badge/github.com/statsg/statsg)](https://goreportcard.com/report/github.com/statsg/statsg) [![Code Coverage](https://codecov.io/gh/statsg/statsg)](https://codecov.io/gh/statsg/statsg/branch/master/graph/badge.svg)

A reimplementation of statsd to use CoAP as the transport mechanism, and dynamically compressed metric keys.

## Getting started

This project requires Go to be installed. On OS X with Homebrew you can just run `brew install go`.

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
|     1 | 2 bytes |    [0 - 2^16)  | com.example.api      |
|     2 | 1 byte  |    [0 - 2^8)   | service_name         |
|     3 | 2 bytes |    [0 - 2^16)  | some.specific.metric |
|     4 | n bytes | You metric val | Plaintext metric val |
```

This would yield 65,535 FQDNs, 255 services per FQDN, and 65,535 metrics per service yielding
1,099,511,627,776 keys in total expressed in 5 bytes.
