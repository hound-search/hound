FROM golang

COPY . /go/src/github.com/etsy/hound
ONBUILD COPY config.json /hound/
RUN adduser --uid 999 --gecos ,,, --disabled-password --home /hound hound
RUN go-wrapper install github.com/etsy/hound/cmds/houndd

USER hound
EXPOSE 6080

ENTRYPOINT ["/go/bin/houndd", "-conf", "/hound/config.json"]
