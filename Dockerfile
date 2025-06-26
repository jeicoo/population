ARG GO_IMAGE=golang:1.24
FROM ${GO_IMAGE} AS builder

ARG APP_VERSION=dev
ARG COMMIT_SHA=unknown

ENV CGO_ENABLED=0 GOOS=linux

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-extldflags=-static" -o app

FROM scratch

WORKDIR /

COPY --from=builder /app/app /app

EXPOSE 8080

ENTRYPOINT ["/app"]