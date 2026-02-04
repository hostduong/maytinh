# Bước 1: Dùng ảnh Golang chính chủ để Build
FROM golang:1.20-alpine as builder

# Tạo thư mục làm việc trong container
WORKDIR /app

# Copy file quản lý thư viện trước
COPY go.mod ./
# (Nếu sau này bạn chạy go mod tidy nó sẽ sinh ra go.sum, nếu có thì uncomment dòng dưới)
# COPY go.sum ./

# Tải thư viện về
RUN go mod download

# Copy toàn bộ code vào container
COPY . .

# Build ra file chạy tên là "server"
RUN go build -o server main.go

# Bước 2: Dùng ảnh Alpine siêu nhẹ để chạy (Giảm dung lượng)
FROM alpine:latest
WORKDIR /root/

# Copy file chạy từ bước 1 sang
COPY --from=builder /app/server .

# Copy file chứng chỉ Google sang (LƯU Ý: Đọc phần bảo mật bên dưới)
COPY chung_chi_google.json .

# Mở cổng 8080
EXPOSE 8080

# Chạy lệnh
CMD ["./server"]
