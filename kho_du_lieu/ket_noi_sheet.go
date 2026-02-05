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

	// Kết nối Public (Không xác thực)
	srv, err := sheets.NewService(ctx, option.WithoutAuthentication())

	if err != nil {
		// [QUAN TRỌNG] Chỉ in lỗi, KHÔNG ĐƯỢC DÙNG log.Fatalf (Sẽ làm sập Web)
		log.Printf("❌ CẢNH BÁO: Không thể kết nối Google Sheet! Lỗi: %v", err)
		return
	}

	DichVuSheet = srv
	log.Println("✅ [KẾT NỐI] Đã kết nối Google Sheets (Public Mode)")
}
