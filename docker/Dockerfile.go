# ============================================
# STAGE 1: Builder
# ============================================
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build API server
FROM builder AS api-builder
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/api ./cmd/api

# Build WS server
FROM builder AS ws-builder
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/ws ./cmd/ws

# Build Worker
FROM builder AS worker-builder
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/bin/worker ./cmd/worker

# ============================================
# STAGE 2: Runtime
# ============================================
FROM alpine:3.20 AS base
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app

# API runtime
FROM base AS api
COPY --from=api-builder /app/bin/api .
EXPOSE 8080
CMD ["./api"]

# WS runtime
FROM base AS ws
COPY --from=ws-builder /app/bin/ws .
EXPOSE 8081
CMD ["./ws"]

# Worker runtime
FROM base AS worker
COPY --from=worker-builder /app/bin/worker .
CMD ["./worker"]
