# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install required build tools
RUN apk add --no-cache git

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

# Final stage
FROM alpine:3.18

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

# Install certificates for HTTPS
RUN apk --no-cache add ca-certificates

CMD ["./main"]
