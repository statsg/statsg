# statsg

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
