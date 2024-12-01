# Gunakan base image untuk Golang
FROM golang:1.23 AS builder

# Set working directory di dalam container
WORKDIR /app

# Copy semua file plugin ke container
COPY . .

# Build binary plugin
RUN go mod tidy && go build -o plugin-ratelimit main.go

# Gunakan image minimal untuk hasil akhir
FROM alpine:latest

# Copy binary dari stage builder ke stage ini
COPY --from=builder /app/plugin-ratelimit /usr/local/bin/plugin-ratelimit

# Jalankan binary
ENTRYPOINT ["/usr/local/bin/plugin-ratelimit"]
