FROM golang:1.25.0-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o forum ./cmd/server/main.go

# --- Run stage ---
FROM alpine:latest

WORKDIR /app
COPY --from=builder /app/forum .
COPY internal/templates/ ./internal/templates/

RUN mkdir -p /app/data

EXPOSE 8080
CMD ["./forum"]