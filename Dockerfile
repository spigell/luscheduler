FROM golang:1.12-alpine3.9 as builder
RUN apk --no-cache add git make gcc musl-dev

WORKDIR /go/luscheduler
COPY . ${WORKDIR}

RUN make

FROM alpine:3.9.2
RUN apk --no-cache add ca-certificates bash tzdata
WORKDIR /root/
COPY --from=builder /go/luscheduler/bin/luscheduler .
COPY /docker/entrypoint.sh /tmp/entrypoint.sh

ENTRYPOINT ["/tmp/entrypoint.sh"]
