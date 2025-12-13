FROM mcr.microsoft.com/devcontainers/go:1.24-bullseye

ENV LANG=C.UTF-8 \
    TZ=Asia/Tokyo

RUN apt-get update && apt-get install -y \
    curl \
    git \
    && rm -rf /var/lib/apt/lists/*

RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.5.1
