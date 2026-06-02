# ─── Stage 1: Build ───────────────────────────────────────────────────────────
FROM node:20-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .

# Build Vite app
ARG VITE_API_URL
ARG VITE_FARO_COLLECTOR_URL
ARG VITE_APP_NAME=ecommerce-frontend
ARG VITE_APP_VERSION=1.0.0
ARG VITE_APP_ENV=production

RUN npm run build

# ─── Stage 2: Runtime (Nginx) ─────────────────────────────────────────────────
FROM nginx:alpine

COPY --from=builder /app/dist /usr/share/nginx/html

# Nginx config untuk Vue Router history mode
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
