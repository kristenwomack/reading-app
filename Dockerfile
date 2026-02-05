# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files first for better caching
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source
COPY backend/ ./

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server .

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS calls
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/server .

# Copy frontend files
COPY frontend/ ./frontend/

# Copy books.json for initial data
COPY books.json .

# Create data directory for SQLite
RUN mkdir -p /data

# Environment variables
ENV PORT=8080
ENV DATA_DIR=/data

EXPOSE 8080

CMD ["./server"]
