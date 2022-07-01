FROM alpine:3.11.7

ENV GOPATH /go

COPY . /go/src/github.com/hound-search/hound

RUN apk update \
	&& apk add go git subversion libc-dev mercurial bzr openssh tini \
	&& cd /go/src/github.com/hound-search/hound \
	&& go mod download \
	&& go install github.com/hound-search/hound/cmds/houndd \
	&& apk del go \
	&& rm -f /var/cache/apk/* \
	&& rm -rf /go/src /go/pkg

# Hardcode gid and uid so that it never changes. This changing will break
# users running this as nonroot in production as you run it with the uid directly,
# not the user name.
RUN addgroup -g 1000 hound && adduser -u 1000 -G hound -D hound

VOLUME ["/data"]

EXPOSE 6080

ENTRYPOINT ["/sbin/tini", "--", "/go/bin/houndd", "-conf", "/data/config.json"]
