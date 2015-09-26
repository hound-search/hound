
FROM    alpine:edge
RUN     apk -U add git
COPY    dist/bin/houndd /houndd
COPY    config.json /hound/config.json
VOLUME  /hound/data
EXPOSE  6080
ENTRYPOINT ["/houndd", "-conf", "/hound/config.json"]
