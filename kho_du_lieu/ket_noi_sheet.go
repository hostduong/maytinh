package kho_du_lieu

import (
	"context"
	"log"
	"os"

	"app/cau_hinh" // Chú ý: import "app/..."

	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// --- ĐÂY LÀ BIẾN MÀ TRÌNH BIÊN DỊCH ĐANG TÌM ---
var DichVuSheet *sheets.Service 

// KhoiTaoKetNoiGoogle : Hàm kết nối API
func KhoiTaoKetNoiGoogle() {
	ngu_canh := context.Background()

	// 1. Đọc file JSON Key
	du_lieu_file, loi := os.ReadFile(cau_hinh.BienCauHinh.FileChungThuc)
	if loi != nil {
		log.Fatalf("LỖI CHÊT NGƯỜI: Không đọc được file %s. Hãy tải file JSON từ Google Cloud về và đổi tên lại.\nChi tiết: %v", cau_hinh.BienCauHinh.FileChungThuc, loi)
	}

	// 2. Tạo cấu hình JWT từ JSON
	cau_hinh_jwt, loi := google.JWTConfigFromJSON(du_lieu_file, "https://www.googleapis.com/auth/spreadsheets")
	if loi != nil {
		log.Fatalf("LỖI: File JSON không đúng định dạng Google: %v", loi)
	}

	// 3. Tạo Client HTTP
	client := cau_hinh_jwt.Client(ngu_canh)

	// 4. Khởi tạo Dịch vụ Google Sheets
	dich_vu, loi := sheets.NewService(ngu_canh, option.WithHTTPClient(client))
	if loi != nil {
		log.Fatalf("LỖI: Không thể khởi tạo dịch vụ Sheets: %v", loi)
	}

	DichVuSheet = dich_vu
	log.Println("--- [KẾT NỐI] Đã kết nối thành công với Google Sheets API ---")
}
