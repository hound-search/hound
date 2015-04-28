# Hound

[![Join the chat at https://gitter.im/etsy/Hound](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/etsy/Hound?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

Hound is an extremely fast source code search engine. The core is based on this article (and code) from Russ Cox:
[Regular Expression Matching with a Trigram Index](http://swtch.com/~rsc/regexp/regexp4.html). Hound itself is a static
[React](http://facebook.github.io/react/) frontend that talks to a [Go](http://golang.org/) backend. The backend keeps an up-to-date index for each repository and answers searches through a minimal API. Here it is in action:

![Hound Screen Capture](screen_capture.gif)

## Quick Start Guide

1. Use the Go tools to install Hound. The binaries `houndd` (server) and `hound` (cli) will be installed in your $GOPATH.

    ```
    go get github.com/etsy/hound/cmds/...
    ```

2. Create a [config.json](config-example.json) with your list of repositories.

3. Run the Hound server with `houndd` and you should see output similer to:
````
2015/03/13 09:07:42 Searcher started for statsd
2015/03/13 09:07:42 Searcher started for Hound
2015/03/13 09:07:42 All indexes built!
2015/03/13 09:07:42 running server at http://localhost:6080...
```

## Running in Production

There are no special flags to run Hound in production. You can use the `--addr=:6880` flag to control the port to which the server binds. Currently, Hound does not supports SSL/TLS as most users simply run Hound behind either Apache or nginx. Adding TLS support is pretty straight forward though if anyone wants to add it.

## Why Another Code Search Tool?

We've used many similar tools in the past, and most of them are either too slow, too hard to configure, or require too much software to be installed.
Which brings us to...

## Requirements

### Hard Requirements
* Go 1.3+

### Optional, Recommended Software
* Make (for building the binaries, not strictly required)
* nodejs (for the command line react-tools)

Yup, that's it. You can proxy requests to the Go service through Apache/nginx/etc., but that's not required.

## Docker
* docker 1.4+

You should follow the quickstart guide up to step (3) and then run:

    $ docker build -t houndd .
    $ docker run -it --rm -p 0.0.0.0:6080:6080 --name houndd houndd

You should be able to navigate to [http://localhost:6080/](http://localhost:6080/) as usual.

## Support

Currently Hound is only tested on MacOS and CentOS, but it should work on any *nix system. There is no plan to support Windows, and we've heard that it fails to compile on Windows, but we would be happy to accept a PR that fixes this!

Hound supports the following version control systems: 

* Git - This is the default
* Mercurial - use `"vcs" : "hg"` in the config
* SVN - use `"vcs" : "svn"` in the config
* Bazaar - use `"vcs" : "bzr"` in the config

See [config-example.json](config-example.json) for examples of how to use each VCS.

## Private Repositories

There are a couple of ways to get Hound to index private repositories:

* Use the `file://` protocol. This allows you to index any local folder, so you can clone the repository locally 
and then reference the files directly. The downside here is that the polling to keep the repo up to date will
not work.
* Use SSH style URLs in the config: `"url" : "git@github.com:foo/bar.git"`. As long as you have your 
[SSH keys](https://help.github.com/articles/generating-ssh-keys/) set up on the box where Hound is running this will work. There is currently an [issue](https://github.com/etsy/Hound/issues/19) with URLs in this case that we hope to fix soon.

## Keeping Repos Updated

By default Hound polls the URL in the config for updates every 30 seconds. You can override this value by setting the `ms-between-poll` key on a per repo basis in the config. You can see how this works in the [example config](config-example.json).

### Connection Limiting

During the update phase, Hound will limit the number of concurrent connections open against your vcs system to 1. 
To override this setting, add a "max-connections" value to your top-level config.json file. When going against an 
internal service, you can probably set this value to whatever you want, but against an externally-hosted provider, 
we recommend limiting the connections to 20 or 50. We recommend working with your eternal host to find the right 
value for your organization.

## Editor Integration

Currently the following editors have plugins that support Hound:

* [Sublime Text](https://github.com/bgreenlee/SublimeHound)
* [Vim](https://github.com/urthbound/hound.vim)
* [Emacs](https://github.com/ryoung786/hound.el)

## Hacking on Hound

### Editing & Building

Hound uses the standard Go tools for development, so your favorite Go workflow should work. If you are looking for something that will work, here is one option:

```
git clone https://github.com/etsy/Hound.git hound/src/github.com/etsy/hound
cd hound
GOPATH=`pwd` go install github.com/etsy/hound/cmds/...
```

### Testing

There are an increasing number of tests in each of the packages in Hound. Please make sure these pass before uploading your Pull Request. You can run the tests with the following command.

```
GOPATH=`pwd` go test github.com/etsy/hound/...
```

### Working on the web UI

Hound includes a web UI that is composed of several files (html, css, javascript, etc.). To make sure hound works seamlessly with the standard Go tools, these resources are all built inside of the `houndd` binary. This adds a small burden on developers to re-package the UI files after each change. If you make changes to the UI, please follow these steps:

1. To make development easier, there is a flag that will read the files from the file system (allowing the much-loved edit/refresh cycle).

    ```
    bin/houndd --dev
    ```

2. Before uploading a Pull Request, please run the following command. This should regenerate the file `ui/bindata.go` which should be included in your Pull Request.

    ```
    cd src/github.com/etsy/hound
    make
    ```

## Get in Touch

Created at [Etsy](https://www.etsy.com) by:

* [Kelly Norton](https://github.com/kellegous)
* [Jonathan Klein](https://github.com/jklein)
