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

VOLUME ["/data"]

EXPOSE 6080

RUN chmod +x ./src/entrypoint.sh

CMD ["/sbin/tini", "--", "./src/entrypoint.sh"]
