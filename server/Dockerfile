FROM golang:1.18 AS builder
WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker caching
COPY go.mod go.sum ./
RUN go mod download

# Copy only the server-specific directories and shared packages
COPY server ./server
COPY pkg ./pkg
COPY config ./config

RUN go build -o ./server/cmd/server ./server/cmd

FROM debian:buster
WORKDIR /app

COPY --from=builder /app/server/cmd /app/server/cmd
COPY --from=builder /app/config /app/config

# Expose the application port
EXPOSE 8080

WORKDIR /app/server/cmd

CMD ["./server"]