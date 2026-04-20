# Stage 1: build the binary
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Download dependencies first (cached layer)
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build a static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /api ./cmd/server

# Stage 2: minimal runtime image
FROM gcr.io/distroless/static-debian12

COPY --from=builder /api /api

EXPOSE 8080

ENTRYPOINT ["/api"]
