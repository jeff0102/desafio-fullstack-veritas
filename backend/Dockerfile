# ---- Build stage (Go >= 1.25.4) ----
FROM golang:1.25-alpine AS build
WORKDIR /app

# Speed up apk if you ever need tools; not strictly required
# RUN apk add --no-cache build-base

# Copy go mod/sum first for better layer caching
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy the rest of the backend
COPY backend ./

# Build API binary
RUN go build -o /out/api ./cmd/api

# ---- Runtime stage ----
FROM alpine:3.20
WORKDIR /srv
COPY --from=build /out/api /usr/local/bin/api
# Persist tasks.json onto a bind-mount (declared in compose)
VOLUME ["/srv"]
EXPOSE 8080
CMD ["api"]
