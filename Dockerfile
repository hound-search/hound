FROM alpine:3.16

ENV GOPATH /go

COPY . /src

RUN apk update \
	&& apk add go git subversion libc-dev mercurial breezy openssh tini build-base npm rsync \
	&& cd /src \
	&& make \
	&& cp .build/bin/houndd /bin \
	&& rm -r .build \
	&& apk del go build-base rsync npm \
	&& rm -f /var/cache/apk/*

# Hardcode gid and uid so that it never changes. This changing will break
# users running this as nonroot in production as you run it with the uid directly,
# not the user name.
RUN addgroup -g 1000 hound && adduser -u 1000 -G hound -D hound

VOLUME ["/data"]

EXPOSE 6080

ENTRYPOINT ["/sbin/tini", "--", "/bin/houndd", "-conf", "/data/config.json"]
