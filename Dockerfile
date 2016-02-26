FROM golang:alpine

COPY . /go/src/github.com/etsy/hound
ONBUILD COPY config.json /hound/
RUN adduser -u 999 -g ,,, -D -h /hound hound
RUN go install github.com/etsy/hound/cmds/houndd

USER hound
EXPOSE 6080

ENTRYPOINT ["/go/bin/houndd", "-conf", "/hound/config.json"]
