# Hound

[![Build Status](https://travis-ci.org/etsy/hound.svg?branch=master)](https://travis-ci.org/etsy/hound) [![Join the chat at https://gitter.im/etsy/Hound](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/etsy/Hound?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Hound is an extremely fast source code search engine. The core is based on this article (and code) from Russ Cox:
[Regular Expression Matching with a Trigram Index](http://swtch.com/~rsc/regexp/regexp4.html). Hound itself is a static
[React](http://facebook.github.io/react/) frontend that talks to a [Go](http://golang.org/) backend. The backend keeps an up-to-date index for each repository and answers searches through a minimal API. Here it is in action:

![Hound Screen Capture](screen_capture.gif)

## Quick Start Guide

### Using Go Tools

1. Use the Go tools to install Hound. The binaries `houndd` (server) and `hound` (cli) will be installed in your $GOPATH.

```
go get github.com/etsy/hound/cmds/...
```

2. Create a [config.json](config-example.json) in a directory with your list of repositories.

3. Run the Hound server with `houndd` and you should see output similar to:
```
2015/03/13 09:07:42 Searcher started for statsd
2015/03/13 09:07:42 Searcher started for Hound
2015/03/13 09:07:42 All indexes built!
2015/03/13 09:07:42 running server at http://localhost:6080...
```

### Using Docker (1.4+)

1. Create a [config.json](config-example.json) in a directory with your list of repositories.

2. Run 
```
docker run -d -p 6080:6080 --name hound -v $(pwd):/data etsy/hound
```

You should be able to navigate to [http://localhost:6080/](http://localhost:6080/) as usual.


## Running in Production

There are no special flags to run Hound in production. You can use the `--addr=:6880` flag to control the port to which the server binds. Currently, Hound does not supports SSL/TLS as most users simply run Hound behind either Apache or nginx. Adding TLS support is pretty straight forward though if anyone wants to add it.

## Why Another Code Search Tool?

We've used many similar tools in the past, and most of them are either too slow, too hard to configure, or require too much software to be installed.
Which brings us to...

## Requirements
* Go 1.4+

Yup, that's it. You can proxy requests to the Go service through Apache/nginx/etc., but that's not required.


## Support

Currently Hound is only tested on MacOS and CentOS, but it should work on any *nix system. Hound on Windows is not supported but we've heard it compiles and runs just fine.

Hound supports the following version control systems: 

* Git - This is the default
* Mercurial - use `"vcs" : "hg"` in the config
* SVN - use `"vcs" : "svn"` in the config
* Bazaar - use `"vcs" : "bzr"` in the config

See [config-example.json](config-example.json) for examples of how to use each VCS.

## Private Repositories

There are a couple of ways to get Hound to index private repositories:

* Use the `file://` protocol. This allows you to index a local clone of a repository. The downside here is that the polling to keep the repo up to date will
not work. (This also doesn't work on local folders that are not of a supported repository type.)
* Use SSH style URLs in the config: `"url" : "git@github.com:foo/bar.git"`. As long as you have your 
[SSH keys](https://help.github.com/articles/generating-ssh-keys/) set up on the box where Hound is running this will work.

## Keeping Repos Updated

By default Hound polls the URL in the config for updates every 30 seconds. You can override this value by setting the `ms-between-poll` key on a per repo basis in the config. If you are indexing a large number of repositories, you may also be interested in tweaking the `max-concurrent-indexers` property. You can see how these work in the [example config](config-example.json). 

## Editor Integration

Currently the following editors have plugins that support Hound:

* [Sublime Text](https://github.com/bgreenlee/SublimeHound)
* [Vim](https://github.com/urthbound/hound.vim)
* [Emacs](https://github.com/ryoung786/hound.el)

## Hacking on Hound

### Editing & Building

#### Requirements:
 * make
 * Node.js ([Installation Instructions](https://github.com/joyent/node/wiki/Installing-Node.js-via-package-manager))
 * React-tools (install w/ `npm -g install react-tools`)

Hound includes tools to make building locally easy. It is recommended that you use these tools if you are working on Hound. To get setup and build, just run the following commands:

```
git clone https://github.com/etsy/hound.git hound/src/github.com/etsy/hound
cd hound
src/github.com/etsy/hound/tools/setup
make
```

### Testing

There are an increasing number of tests in each of the packages in Hound. Please make sure these pass before uploading your Pull Request. You can run the tests with the following command.

```
make test
```

### Working on the web UI

Hound includes a web UI that is composed of several files (html, css, javascript, etc.). To make sure hound works seamlessly with the standard Go tools, these resources are all bundled inside of the `houndd` binary. Note that changes to the UI will result in local changes to the `ui/bindata.go` file. You must include these changes in your Pull Request.

To make development easier, there is a flag that will read the files from the file system (allowing the much-loved edit/refresh cycle).

```
bin/houndd --dev
```

## Get in Touch

Created at [Etsy](https://www.etsy.com) by:

* [Kelly Norton](https://github.com/kellegous)
* [Jonathan Klein](https://github.com/jklein)
