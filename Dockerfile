FROM golang

COPY . /go/src/github.com/etsy/Hound
COPY config.json /hound/
RUN go-wrapper install github.com/etsy/Hound/cmds/houndd

EXPOSE 6080

ENTRYPOINT ["/go/bin/houndd", "-conf", "/hound/config.json"]
