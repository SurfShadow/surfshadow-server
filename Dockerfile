FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o surfshadow-server ./cmd/server/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/surfshadow-server .
COPY --from=builder /app/migrations ./migrations

CMD ["./surfshadow-server"]