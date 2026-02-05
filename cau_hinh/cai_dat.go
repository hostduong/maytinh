package cau_hinh

import (
	"log"
	"os"
)

// CauHinhHeThong : Cấu trúc chứa các biến môi trường
type CauHinhHeThong struct {
	IdFileSheet string // ID của file Google Sheet
	CongChayWeb string // Port (VD: 8080)
	// Đã xóa FileChungThuc
}

// BienCauHinh : Biến toàn cục để các nơi khác gọi dùng
var BienCauHinh CauHinhHeThong

// KhoiTaoCauHinh : Hàm này chạy đầu tiên để nạp cấu hình
func KhoiTaoCauHinh() {
	// 1. Lấy ID Sheet từ biến môi trường hoặc gán cứng
	idSheet := os.Getenv("SHEET_ID")
	if idSheet == "" {
		idSheet = "17f5js4C9rY7GPd4TOyBidkUPw3vCC6qv6y8KlF3vNs8"
	}

	// 2. Lấy cổng chạy web (Cloud Run sẽ tự điền biến PORT)
	congWeb := os.Getenv("PORT")
	if congWeb == "" {
		congWeb = "8080"
	}

	// 3. Gán vào biến toàn cục
	BienCauHinh = CauHinhHeThong{
		IdFileSheet: idSheet,
		CongChayWeb: congWeb,
	}

	log.Println("--- [CẤU HÌNH] Đã tải xong cài đặt hệ thống ---")
	log.Println("--- ID Sheet:", BienCauHinh.IdFileSheet)
	log.Println("--- Chế độ: Public Sheet (Không cần chứng thực)")
}
