# Bước 1: Build
FROM golang:1.20-alpine as builder
WORKDIR /app
COPY . .
# Tắt CGO để file chạy độc lập hoàn toàn
ENV CGO_ENABLED=0 
RUN go mod tidy
RUN go build -o server main.go

# Bước 2: Run
FROM alpine:latest
WORKDIR /root/
# Cài thêm chứng chỉ bảo mật để gọi HTTPS Google không bị lỗi
RUN apk add --no-cache ca-certificates

COPY --from=builder /app/server .

# Mở cổng
EXPOSE 8080
CMD ["./server"]
