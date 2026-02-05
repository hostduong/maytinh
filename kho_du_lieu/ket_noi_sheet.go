package kho_du_lieu

import (
	"context"
	"log"

	"google.golang.org/api/sheets/v4"
    // Bỏ dòng option.WithoutAuthentication
)

var DichVuSheet *sheets.Service

func KhoiTaoKetNoiGoogle() {
	ctx := context.Background()

	// [NÂNG CẤP] Không truyền option gì cả.
	// Go sẽ tự động tìm Service Account của Cloud Run (ADC) để lấy quyền Ghi.
	srv, err := sheets.NewService(ctx)

	if err != nil {
		log.Printf("❌ LỖI KẾT NỐI GOOGLE SHEET: %v", err)
		return
	}

	DichVuSheet = srv
	log.Println("✅ [ADC] Đã kết nối Google Sheets bằng quyền Cloud Run (Đọc/Ghi)")
}
