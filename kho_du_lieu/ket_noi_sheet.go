package kho_du_lieu

import (
	"context"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Biến toàn cục để kho_chung.go gọi dùng
var DichVuSheet *sheets.Service

func KhoiTaoKetNoiGoogle() {
	ctx := context.Background()

	// [FIX CRASH] Sử dụng chế độ "WithoutAuthentication"
	// Giúp App chạy được ngay cả khi không có file JSON (dành cho Sheet Public)
	srv, err := sheets.NewService(ctx, option.WithoutAuthentication())

	if err != nil {
		log.Fatalf("❌ LỖI KHỞI TẠO SHEET: %v", err)
	}

	DichVuSheet = srv
	log.Println("--- [KẾT NỐI] Đã kết nối Google Sheets (Public Mode) ---")
}

// [QUAN TRỌNG] ĐÃ XÓA HÀM DocToanBoSheet Ở ĐÂY
// VÌ NÓ ĐÃ CÓ TRONG FILE kho_chung.go RỒI.
