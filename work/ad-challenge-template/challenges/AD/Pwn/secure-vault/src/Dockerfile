FROM ubuntu:24.04

RUN apt-get update && apt-get install -y --no-install-recommends \
    socat \
    gcc \
    libc6-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
COPY . /app/
RUN gcc -O2 -fstack-protector server.c -o /app/server

EXPOSE 9000

# Service dijalankan via daemon/socat
CMD ["socat", "TCP-LISTEN:9000,reuseaddr,fork", "EXEC:/app/server"]
