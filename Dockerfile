# Gunakan image golang
FROM golang:1.23

# Set working directory
WORKDIR /app

# Copy dependency definition
COPY go.mod go.sum ./
RUN go mod download

# Copy semua file project ke container
COPY . .

# Build aplikasi
RUN go build -o tripmate .

# Expose port
EXPOSE 8080

# Jalankan binary
CMD ["./tripmate"]

