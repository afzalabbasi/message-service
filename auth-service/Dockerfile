# auth-service/Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth-service ./cmd/main.go

# Use a small image for the final container
FROM alpine:latest

# Add necessary packages
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/auth-service .

# Expose port
EXPOSE 8081

# Command to run
CMD ["./auth-service"]

# websocket-service/Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o websocket-service ./cmd/main.go

# Use a small image for the final container
FROM alpine:latest

# Add necessary packages
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/websocket-service .

# Expose port
EXPOSE 8082

# Command to run
CMD ["./websocket-service"]

# persistence-service/Dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./
COPY go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o persistence-service ./cmd/main.go

# Use a small image for the final container
FROM alpine:latest

# Add necessary packages
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/persistence-service .

# Command to run
CMD ["./persistence-service"]