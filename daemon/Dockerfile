FROM golang:1.18.3 AS builder

COPY . /go/src/github.com/ubiq/ubiq-burn-stats/daemon
WORKDIR /go/src/github.com/ubiq/ubiq-burn-stats/daemon
RUN go build

FROM debian:buster-slim

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/ubiq/ubiq-burn-stats/daemon/daemon daemon

ENTRYPOINT ["/daemon"]
