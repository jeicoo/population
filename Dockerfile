ARG GO_IMAGE=golang:1.24
FROM ${GO_IMAGE} AS builder

ENV CGO_ENABLED=0 GOOS=linux

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-extldflags=-static" -o app

FROM scratch

WORKDIR /

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=builder --chown=nobody:nogroup /app/app /app

EXPOSE 8080

USER nobody

ENTRYPOINT ["/app"]