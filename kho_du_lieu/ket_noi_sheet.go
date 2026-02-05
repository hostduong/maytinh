package kho_du_lieu

import (
	"context"
	"log"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var DichVuSheet *sheets.Service

func KhoiTaoKetNoiGoogle() {
	ctx := context.Background()

	// CÁCH MỚI: Bỏ qua xác thực (Public Mode)
	// Giúp App khởi động được ngay cả khi không có file JSON
	srv, err := sheets.NewService(ctx, option.WithoutAuthentication())

	if err != nil {
		log.Fatalf("❌ LỖI: Không thể khởi tạo dịch vụ Sheets: %v", err)
	}

	DichVuSheet = srv
	log.Println("--- [KẾT NỐI] Đã kết nối Google Sheets (Public Mode) ---")
}
