FROM golang:alpine as builder
RUN apk --no-cache add git make gcc

WORKDIR /go/luscheduler

RUN make submodule_check
RUN make

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/teleport/bin/luscheduler .
COPY entrypoint.sh /tmp/entrypoint.sh

ENTRYPOINT ["/tmp/entrypoint.sh"]
