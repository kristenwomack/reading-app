# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Copy Go module files and download dependencies
COPY backend/go.mod backend/go.sum ./backend/
RUN cd backend && go mod download

# Copy source code
COPY backend/ ./backend/
COPY books.json ./

# Build the binary
RUN cd backend && CGO_ENABLED=0 go build -o /reading-tracker main.go

# Runtime stage
FROM alpine:3.21

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy binary and static files
COPY --from=builder /reading-tracker .
COPY --from=builder /build/books.json .
COPY frontend/ ./frontend/

# Create data directory for SQLite persistent volume
RUN mkdir -p /data

EXPOSE 3000

CMD ["./reading-tracker"]
