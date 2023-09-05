# syntax=docker/dockerfile:1

FROM golang:alpine as builder

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/var/cache/apk \
    apk update \
	&& apk add git subversion mercurial breezy openssh tini npm rsync build-base

COPY . /src

RUN --mount=type=cache,target=/go/pkg/mod \
	cd /src \
	&& make \

FROM alpine:latest
RUN   --mount=type=cache,target=/var/cache/apk \
      apk add git subversion mercurial breezy openssh tini

COPY --from=builder /src/.build/bin/houndd /bin/

VOLUME ["/data"]

EXPOSE 6080

ENTRYPOINT ["/sbin/tini", "--", "/bin/houndd", "-conf", "/data/config.json"]
