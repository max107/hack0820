FROM golang:1.16.0-alpine3.13 AS base
WORKDIR $GOPATH/src/github.com/max107/hack0820
ENV USER=appuser
ENV UID=10001
RUN adduser \
  --disabled-password \
  --gecos "" \
  --home "/nonexistent" \
  --shell "/sbin/nologin" \
  --no-create-home \
  --uid "${UID}" \
  "${USER}"
RUN apk add --update --no-cache git ca-certificates gcc libc-dev make
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN make build-linux

FROM alpine:3.13.2
WORKDIR /app/
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /etc/passwd /etc/passwd
COPY --from=base /etc/group /etc/group
COPY --from=base /go/src/github.com/max107/hack0820/bin/linux/app .
USER appuser:appuser
CMD ["/app/app"]
