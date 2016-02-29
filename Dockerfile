FROM golang:alpine

COPY . /go/src/github.com/etsy/hound
ONBUILD COPY config.json /hound/
RUN apk update && apk add git subversion mercurial bzr
RUN go install github.com/etsy/hound/cmds/houndd

EXPOSE 6080

ENTRYPOINT ["/go/bin/houndd", "-conf", "/hound/config.json"]
