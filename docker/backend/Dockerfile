# ── Build stage ──────────────────────────────────────────────────────────────
# build context: 仓库根目录 notes/
FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY zw-note-backend/go.mod zw-note-backend/go.sum ./
RUN go mod download

COPY zw-note-backend/ .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w -X main.Version=$(git describe --tags --always --dirty 2>/dev/null || echo dev)" \
    -o /app/bin/server ./cmd/server/main.go


# ── Final stage ──────────────────────────────────────────────────────────────
FROM alpine:3.21

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /app/bin/server   ./server
COPY zw-note-backend/configs/         ./configs/

# 预先创建上传目录，便于挂载持久化卷时继承 appuser 权限（见 docker-compose.yml）
RUN mkdir -p /app/uploads && chown -R appuser:appgroup /app
USER appuser

EXPOSE 8080

ENTRYPOINT ["./server", "start"]
