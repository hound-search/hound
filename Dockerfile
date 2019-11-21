FROM alpine

ARG DELETE_GO=no
ARG INSTALL_CLI=no

ENV GOPATH /go

COPY . /go/src/github.com/it-projects-llc/hound

COPY default-config.json /data/config.json

RUN apk update \
	&& apk add go git subversion libc-dev mercurial bzr openssh \
	go install github.com/it-projects-llc/hound/cmds/houndd

RUN [ "INSTALL_CLI" = "yes" ] \
    && go install github.com/it-projects-llc/hound/cmds/hound

RUN [ "$DELETE_GO" = "yes" ] && apk del go \
    && rm -f /var/cache/apk/* \
    && rm -rf /go/src /go/pkg

VOLUME ["/data"]

EXPOSE 6080

ENTRYPOINT ["/go/bin/houndd", "-conf", "/data/config.json"]
