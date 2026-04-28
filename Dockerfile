FROM golang:1.25-alpine AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd/server

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /usr/src/app/app /usr/local/bin/app
COPY --from=builder /usr/src/app/internal/templates ./internal/templates

CMD ["/usr/local/bin/app"]