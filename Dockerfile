# syntax=docker/dockerfile:1
# ══════════════════════════════════════════════════════════════════════════════
# Stage 1 · Build
# Uses BuildKit cache mounts so Go modules and the build cache persist
# between CI runs, dramatically reducing build time.
# ══════════════════════════════════════════════════════════════════════════════
FROM golang:1.26-bookworm AS builder

ARG VERSION=dev
ARG BUILD_DATE=unknown
ARG GIT_COMMIT=none

WORKDIR /app

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux \
    go build \
      -ldflags "-s -w \
        -X github.com/masudur-rahman/expense-tracker-bot/cmd.Version=${VERSION} \
        -X github.com/masudur-rahman/expense-tracker-bot/cmd.BuildDate=${BUILD_DATE} \
        -X github.com/masudur-rahman/expense-tracker-bot/cmd.GitCommit=${GIT_COMMIT}" \
      -o /bin/expense-tracker .

# ══════════════════════════════════════════════════════════════════════════════
# Stage 2 · Runtime base
# Shared base with fonts and common dependencies.
# ══════════════════════════════════════════════════════════════════════════════
FROM debian:bookworm-slim AS runtime-base

ARG DEBIAN_RELEASE_NAME=bookworm

RUN set -x \
 && apt-get update \
 && apt-get upgrade -y \
 && apt-get install -y --no-install-recommends ca-certificates wget \
 && echo 'Etc/UTC' > /etc/timezone

RUN echo "deb http://deb.debian.org/debian ${DEBIAN_RELEASE_NAME} contrib" >> /etc/apt/sources.list \
 && apt-get update \
 && DEBIAN_FRONTEND=noninteractive apt-get install -y \
      fonts-lohit-beng-bengali \
      fonts-dejavu \
      fontconfig \
      ttf-mscorefonts-installer \
 && rm -rf /var/lib/apt/lists/* /usr/share/doc /usr/share/man /tmp/*

WORKDIR /app
COPY --from=builder /bin/expense-tracker /app/expense-tracker

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
  CMD ["wget", "-q", "--spider", "http://localhost:8080/healthz"]

ENTRYPOINT ["/app/expense-tracker"]
CMD ["serve"]

# ══════════════════════════════════════════════════════════════════════════════
# Stage 3a · wkhtmltopdf edition
# docker build --target wkhtmltopdf .
# ══════════════════════════════════════════════════════════════════════════════
FROM runtime-base AS wkhtmltopdf

ARG TARGETARCH=amd64
ARG WKHTMLTOPDF_VERSION=0.12.6.1-3
ARG DEBIAN_RELEASE_NAME=bookworm

RUN set -x \
 && apt-get update \
 && wget -q https://github.com/wkhtmltopdf/packaging/releases/download/${WKHTMLTOPDF_VERSION}/wkhtmltox_${WKHTMLTOPDF_VERSION}.${DEBIAN_RELEASE_NAME}_${TARGETARCH}.deb \
 && dpkg -i wkhtmltox_${WKHTMLTOPDF_VERSION}.${DEBIAN_RELEASE_NAME}_${TARGETARCH}.deb || true \
 && apt-get install -f -y \
 && ldconfig \
 && rm wkhtmltox_${WKHTMLTOPDF_VERSION}.${DEBIAN_RELEASE_NAME}_${TARGETARCH}.deb \
 && rm -rf /var/lib/apt/lists/* /usr/share/doc /usr/share/man /tmp/*

RUN mkdir -p /app/configs /app/.sqlite \
 && chown -R 65535:65535 /app

USER 65535:65535

# ══════════════════════════════════════════════════════════════════════════════
# Stage 3b · chromedp edition
# docker build --target chromedp .
# ══════════════════════════════════════════════════════════════════════════════
FROM runtime-base AS chromedp

RUN apt-get update \
 && apt-get install -y --no-install-recommends chromium \
 && rm -rf /var/lib/apt/lists/* /usr/share/doc /usr/share/man /tmp/*

ENV CHROME_PATH=/usr/bin/chromium
RUN mkdir -p /app/configs /app/.sqlite \
 && chown -R 65535:65535 /app

USER 65535:65535
