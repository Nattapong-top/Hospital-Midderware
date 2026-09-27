# Stage 1: Build binary
FROM golang:1.27.1-alpine AS builder

WORKDIR /app

# Copy dependency modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build application
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Stage 2: Minimal runtime image
FROM alpine:3.19

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]