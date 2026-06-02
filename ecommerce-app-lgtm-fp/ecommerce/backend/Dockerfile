# ─── Stage 1: Build ───────────────────────────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go.mod dan go.sum terlebih dahulu untuk memanfaatkan cache Docker layer
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary dengan optimasi ukuran
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/server \
    ./cmd/server

# ─── Stage 2: Runtime ─────────────────────────────────────────────────────────
FROM gcr.io/distroless/static-debian12

WORKDIR /app

# Copy binary dari stage build
COPY --from=builder /app/server .

# Copy migration files jika ada
# COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

ENTRYPOINT ["/app/server"]

