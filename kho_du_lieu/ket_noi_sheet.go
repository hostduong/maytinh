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

	// [FIX LỖI 3]: Dùng chế độ Public để không đòi file JSON gây Crash
	srv, err := sheets.NewService(ctx, option.WithoutAuthentication())

	if err != nil {
		log.Fatalf("❌ LỖI KHỞI TẠO SHEET: %v", err)
	}

	DichVuSheet = srv
	log.Println("--- [KẾT NỐI] Đã kết nối Google Sheets (Public Mode) ---")
}
