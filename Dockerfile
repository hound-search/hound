FROM golang

COPY src /go/src
COPY pub /go/pub
COPY config.json /go/pub/
RUN go-wrapper install hound/cmds/houndd

EXPOSE 6080

ENTRYPOINT ["bin/houndd", "-conf", "/go/pub/config.json"]
