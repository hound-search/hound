FROM alpine

ENV GOPATH /go

COPY . /go/src/github.com/hound-search/hound

COPY default-config.json /data/config.json

RUN apk update \
	&& apk add go git subversion libc-dev mercurial bzr openssh \
	&& go install github.com/hound-search/hound/cmds/houndd \
	&& apk del go \
	&& rm -f /var/cache/apk/* \
	&& rm -rf /go/src /go/pkg

VOLUME ["/data"]

EXPOSE 6080

ENTRYPOINT ["/go/bin/houndd", "-conf", "/data/config.json"]
