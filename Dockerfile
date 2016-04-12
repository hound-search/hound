FROM alpine

ENV GOPATH /go
COPY . /go/src/github.com/etsy/hound
ONBUILD COPY config.json /hound/
RUN apk update \
	&& apk add go git subversion mercurial bzr openssh \
	&& go install github.com/etsy/hound/cmds/houndd \
	&& apk del go \
	&& rm -f /var/cache/apk/* \
	&& rm -rf /go/src /go/pkg

EXPOSE 6080

ENTRYPOINT ["/go/bin/houndd", "-conf", "/hound/config.json"]
