FROM golang:1.16-buster

ENV GOPATH /go

COPY . /go/src/github.com/hound-search/hound

RUN apt update \
	&& apt upgrade -y \
	&& apt install -y git subversion libc-dev mercurial bzr openssh-client tini \
	&& cd /go/src/github.com/hound-search/hound \
	&& go mod download \
	&& go install github.com/hound-search/hound/cmds/houndd \
	&& apt clean \
	&& rm -rf /go/src /go/pkg

VOLUME ["/data"]

EXPOSE 6080

ENTRYPOINT ["/usr/bin/tini", "--", "/go/bin/houndd", "-conf", "/data/config.json"]
