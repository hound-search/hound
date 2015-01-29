# Hound

Hound is an extremely fast source code search engine. The core is based on this article (and code) from Russ Cox:
[Regular Expression Matching with a Trigram Index](http://swtch.com/~rsc/regexp/regexp4.html). Hound itself is a static
[React](http://facebook.github.io/react/) frontend that talks to a [Go](http://golang.org/) backend. The backend keeps an up-to-date index for each repository and answers searches through a minimal API. Here it is in action:

![Hound Screen Capture](screen_capture.gif)

## Quick Start Guide

1. Clone the repo: `git clone git@github.com:etsy/Hound.git`
2. Edit [config-example.json](config-example.json) to add the repos you want: `cd Hound && vim config-example.json`
3. Rename the (now edited) config file: `mv config-example.json config.json`
4. Set your GOPATH: ``export GOPATH=`pwd` ``
5. Run the server: `go run src/hound/cmds/houndd/main.go`
6. See Hound in action in your browser at [http://localhost:6080/](http://localhost:6080/)

Have [Rake](http://docs.seattlerb.org/rake/) installed? Steps 4 and 5 change to:

* Run rake to create binaries: `rake`
* Run the binary: `./bin/houndd`

This is the preferred approach, since the binaries are generally easier to work with, and rake will build both the server and the CLI binaries at the same time.

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

## Docker
* docker 1.4+

You should follow the quickstart guide up to step (2) and then run:

    $ docker build . -t houndd
    $ docker run -it --rm -p 0.0.0.0:6080:6080 --name houndd houndd

You should be able to navigate to [http://localhost:6080/](http://localhost:6080/) as usual.

## Support

Currently Hound is only tested on MacOS and CentOS, but it should work on any *nix system. There is no plan to support Windows, and we've heard that it fails to compile on Windows, but we would be happy to accept a PR that fixes this!

Similarly, right now Hound only supports git repositories, although adding SVN and Mercurial wouldn't take too much work. Pull requests for this are welcome.

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
