# statsg

[![codebeat badge](https://codebeat.co/badges/53703c14-bc44-489d-b0ea-7e3d3b2d8d82)](https://codebeat.co/projects/github-com-statsg-statsg) [![circleci` badge](https://circleci.com/gh/statsg/statsg.svg?style=shield&circle-token=456ad54019146fcaa685adb00e91d7bb73f9d58e)](https://circleci.com/gh/statsg/statsg.svg?style=shield&circle-token=456ad54019146fcaa685adb00e91d7bb73f9d58e)

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
