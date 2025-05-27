FROM --platform=linux/amd64 golang:1.24.2-bookworm AS builder

RUN apt-get update && apt-get install -y \
    bash \
    openssl \
    perl \
    libimage-exiftool-perl \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY . .


RUN go mod download


RUN go build -o tawtheeq .


FROM debian:bookworm-slim


RUN apt-get update && apt-get install -y \
    bash \
    openssl \
    perl \
    libimage-exiftool-perl \
 && apt-get clean \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app


COPY --from=builder /app/tawtheeq .
COPY --from=builder /app/.env ./

ENTRYPOINT ["./tawtheeq"]
