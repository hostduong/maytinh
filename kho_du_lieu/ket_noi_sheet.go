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

	// Chế độ Public (Không cần file JSON) -> Tránh lỗi Crash
	srv, err := sheets.NewService(ctx, option.WithoutAuthentication())

	if err != nil {
		log.Fatalf("Lỗi khởi tạo Sheet: %v", err)
	}

	DichVuSheet = srv
	log.Println("--- Đã kết nối Google Sheets (Public Mode) ---")
}

// Hàm đọc dữ liệu (Độc lập, không phụ thuộc nghiep_vu)
func DocToanBoSheet(spreadsheetId string, rangeName string) ([][]interface{}, error) {
	if DichVuSheet == nil {
		return nil, nil
	}
	resp, err := DichVuSheet.Spreadsheets.Values.Get(spreadsheetId, rangeName).Do()
	if err != nil {
		return nil, err
	}
	return resp.Values, nil
}
