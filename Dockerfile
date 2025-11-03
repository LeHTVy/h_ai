# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/h-ai-server ./main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/h-ai-mcp ./cmd/mcp

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binaries
COPY --from=builder /app/bin/h-ai-server .
COPY --from=builder /app/bin/h-ai-mcp .

EXPOSE 8888

CMD ["./h-ai-server"]
