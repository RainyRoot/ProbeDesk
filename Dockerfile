FROM golang:1.25-bookworm AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    gcc \
    libgl1-mesa-dev \
    libxrandr-dev \
    libxcursor-dev \
    libxinerama-dev \
    libxi-dev \
    libxxf86vm-dev \
    xorg-dev \
    libgtk-3-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o probedesk .

FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    libgl1 \
    libxrandr2 \
    libxcursor1 \
    libxinerama1 \
    libxi6 \
    libgtk-3-0 \
    ca-certificates \
    xvfb \
    x11-utils \
    && rm -rf /var/lib/apt/lists/*

COPY --from=builder /build/probedesk /usr/local/bin/probedesk
COPY docker-entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
