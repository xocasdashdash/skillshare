# CI Dockerfile — minimal image for skill validation in pipelines
# Usage: docker run ghcr.io/<org>/skillshare-ci audit /path/to/skills
# Build: docker build -f docker/ci/Dockerfile -t skillshare-ci .

# Stage 1: Build frontend assets
FROM node:22-bookworm-slim AS ui-builder
RUN npm install -g pnpm@latest
WORKDIR /ui
COPY ui/package.json ui/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
COPY ui/ ./
RUN pnpm run build

# Stage 2: Build Go binary
FROM golang:1.25.5-bookworm AS go-builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=ui-builder /ui/dist internal/server/dist/
ARG TARGETARCH=amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o /skillshare ./cmd/skillshare

# Stage 3: Minimal CI runtime
FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y --no-install-recommends \
    git ca-certificates && \
    rm -rf /var/lib/apt/lists/*
COPY --from=go-builder /skillshare /usr/local/bin/skillshare
ENTRYPOINT ["skillshare"]
