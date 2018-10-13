FROM golang:alpine as builder
RUN apk --no-cache add git make gcc

WORKDIR /go/luscheduler

RUN git clone http://192.168.23.1:10085/spigell/luscheduler .

RUN make submodule_check
RUN make

FROM alpine:latest
RUN apk --no-cache add ca-certificates bash tzdata
WORKDIR /root/
COPY --from=builder /go/luscheduler/bin/luscheduler .
COPY /docker/entrypoint.sh /tmp/entrypoint.sh

ENTRYPOINT ["/tmp/entrypoint.sh"]
