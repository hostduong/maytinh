# Bước 1: Build
FROM golang:1.20-alpine as builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN go build -o server main.go

# Bước 2: Run
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/server .

# --- ĐÃ XÓA DÒNG COPY FILE JSON ---

EXPOSE 8080
CMD ["./server"]
