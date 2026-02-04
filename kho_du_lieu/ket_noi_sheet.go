package kho_du_lieu

import (
	"context"
	"log"
	
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var DichVuSheet *sheets.Service 

func KhoiTaoKetNoiGoogle() {
	ctx := context.Background()

	// CÁCH MỚI: Tự động tìm Credential (File, Env Var, hoặc Cloud Run Identity)
	// Scope: Full quyền đọc/ghi
	creds, err := google.FindDefaultCredentials(ctx, "https://www.googleapis.com/auth/spreadsheets")
	
	var srv *sheets.Service
	var errService error

	if err == nil {
		// Trường hợp 1: Tìm thấy Credential (JSON hoặc Cloud Run)
		log.Println("--- [AUTH] Tìm thấy Credential hợp lệ. Đang kết nối... ---")
		srv, errService = sheets.NewService(ctx, option.WithCredentials(creds))
	} else {
		// Trường hợp 2: Không tìm thấy (Dùng cho Public Sheet - Không xác thực)
		// Lưu ý: Cách này chỉ ĐỌC được, GHI có thể bị lỗi 401/403 tùy config Sheet
		log.Println("--- [AUTH] Cảnh báo: Không tìm thấy Credential. Chạy chế độ No-Auth... ---")
		srv, errService = sheets.NewService(ctx, option.WithHTTPClient(nil)) // Client rỗng
	}

	if errService != nil {
		log.Fatalf("LỖI: Không thể khởi tạo dịch vụ Sheets: %v", errService)
	}

	DichVuSheet = srv
	log.Println("--- [KẾT NỐI] Đã kết nối thành công với Google Sheets API ---")
}
