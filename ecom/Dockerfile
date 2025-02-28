# Build stage: compile the Go application.
FROM golang:1.21 AS builder
WORKDIR /app

# Copy go.mod and go.sum and download dependencies.
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the code.
COPY . .

# Build the Go binary.
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./main.go

# Final stage: use a lightweight image.
FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/server .

# Expose the port the application listens on.
EXPOSE 8080

# Command to run the server.
CMD ["./server"]
