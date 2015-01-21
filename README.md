# Hound

Hound is an extremely fast source code search engine. The core is based on this article (and code) from Russ Cox:
[Regular Expression Matching with a Trigram Index](http://swtch.com/~rsc/regexp/regexp4.html). Hound itself is a static
[React](http://facebook.github.io/react/) frontend that talks to a [Go](http://golang.org/) backend. The backend keeps an up-to-date index for each
repository and and answers searches through a minimal API.

## Why Another Code Search Tool?

We've used many similar tools in the past, and most of them are either too slow, too hard to configure, or require too much software to be installed.
Which brings us to...

## Requirements

### Hard Requirements
* Go 1.3+

### Optional, Recommended Software
* Rake (for building the binaries, not strictly required)
* nodejs (for the command line react-tools)

Yup, that's it. You can proxy requests to the Go service through Apache/nginx/etc., but that's not required.

## Hacking on Hound

### Building

```
rake
```

This will build `./bin/houndd` which is the server binary and `./bin/hound` which is the command line client.

### Running in development

```
./bin/houndd
```

This will start up the combined server and indexer. The first time you start the server, it will take a bit of time to initialize your `data` directory with the repository data.
You can access the web frontend at http://localhost:6080/

### Running in production

```
./bin/houndd --prod --addr=address:port
```

The will start up the combined server/indexer and build all static assets in production mode. The default addr is ":6080", and thus the `--addr` flag can be used to have the server listen on a different port.

## Get in Touch

IRC: #codeascraft on freenode

Created at [Etsy](https://www.etsy.com) by:

* [Kelly Norton](https://github.com/kellegous)
* [Jonathan Klein](https://github.com/jklein)
