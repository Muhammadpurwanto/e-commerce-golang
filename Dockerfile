# Tahap 1: Build aplication
# Menggunakan image resmi Go sebagai build environment
FROM golang:1.24.4-alpine AS builder

# Set direktori kerja di dalam container
WORKDIR /app

# Copy go.mod dan go.sum untuk men-download dependensi
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build aplikasi Go.
# CGO_ENABLED=0 untuk static build, -o main untuk nama output file
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

# Tahap 2: Final image
# Menggunakan image dasar yang sangat kecil
FROM alpine:3.19

# Set direktori kerja
WORKDIR /root/

# Copy file .env
COPY .env .

# Copy binary yang sudah di-build dari tahap 'builder'
COPY --from=builder /app/main .

# Expose port yang digunakan aplikasi
EXPOSE 8080

# Command untuk menjalankan aplikasi saat container dimulai
CMD ["./main"]