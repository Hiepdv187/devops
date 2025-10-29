# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

# Download dependencies first to leverage Docker layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy application source
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./

# Runtime stage
FROM gcr.io/distroless/base-debian12:latest

WORKDIR /app

COPY --from=builder /app/server /app/server

EXPOSE 3003

# Fiber reads PORT from env; default already 3003
ENV PORT=3003

ENTRYPOINT ["/app/server"]
