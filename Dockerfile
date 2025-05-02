# Stage 1: Build the Go binary for Linux
FROM golang:1.24.2 AS builder

WORKDIR /app
COPY . .

RUN go mod tidy
RUN GOOS=linux GOARCH=amd64 go build -o tcp-server .

# Stage 2: Run the binary in a minimal container
FROM busybox:latest

WORKDIR /root/
COPY --from=builder /app/tcp-server .

EXPOSE 3000
CMD ["./tcp-server"]