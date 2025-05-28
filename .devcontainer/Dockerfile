FROM mcr.microsoft.com/devcontainers/go:1.21

RUN apt-get update && apt-get install -y \
    bash \
    openssl \
    perl \
    libimage-exiftool-perl \
    mupdf-tools \
    libmupdf-dev \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*
