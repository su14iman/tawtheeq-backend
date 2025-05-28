FROM mcr.microsoft.com/devcontainers/go:1.21

# Install exiftool and bash
RUN apt-get update && \
    apt-get install -y libimage-exiftool-perl bash && \
    rm -rf /var/lib/apt/lists/*
