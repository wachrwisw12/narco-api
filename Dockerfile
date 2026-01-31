# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o server main.go

# Run stage
FROM alpine:latest

WORKDIR /app
COPY .env .env
COPY keys /run/keys
COPY --from=builder /app/server .
# Copy .env เข้า container


EXPOSE 8000

CMD ["./server"]
