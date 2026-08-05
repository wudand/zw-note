# ── Build stage ──────────────────────────────────────────────────────────────
# build context: 仓库根目录 notes/
FROM node:22.20.0 AS builder

WORKDIR /app

COPY zw-note/package*.json ./
RUN npm ci

COPY zw-note/ .
RUN npm run build


# ── Runtime stage ────────────────────────────────────────────────────────────
FROM nginx:latest

COPY --from=builder /app/dist /usr/share/nginx/html
COPY docker/frontend/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
