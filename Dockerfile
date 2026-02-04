# Bước 1: Dùng ảnh Golang chính chủ để Build
FROM golang:1.20-alpine as builder

# Tạo thư mục làm việc trong container
WORKDIR /app

# Copy toàn bộ code vào container TRƯỚC (Để quét được các import)
COPY . .

# --- KHẮC PHỤC LỖI Ở ĐÂY ---
# Chạy lệnh này để nó tự tìm thư viện còn thiếu và tạo go.sum ảo
RUN go mod tidy

# Build ra file chạy tên là "server"
RUN go build -o server main.go

# Bước 2: Dùng ảnh Alpine siêu nhẹ để chạy (Giảm dung lượng)
FROM alpine:latest
WORKDIR /root/

# Copy file chạy từ bước 1 sang
COPY --from=builder /app/server .

# Copy file chứng chỉ Google sang
COPY chung_chi_google.json .

# Mở cổng 8080
EXPOSE 8080

# Chạy lệnh
CMD ["./server"]
