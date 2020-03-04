FROM golang:1.13.8 AS builder

ENV DEBIAN_FRONTEND=noninteractive
RUN apt update && apt install -y binutils

WORKDIR /ops
ADD . .
RUN go build -ldflags='-s -w' -o main && strip -S main && chmod 777 main

############################
# Final container
############################
FROM registry.cto.ai/official_images/base:latest 

ENV DEBIAN_FRONTEND=noninteractive

RUN apt update && apt install -y openssl libssl1.1 ca-certificates && rm -rf /var/lib/apt/lists
COPY --from=builder /ops/main /ops/main
RUN chown -R 9999:9999 /ops 
