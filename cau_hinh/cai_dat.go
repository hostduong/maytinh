package cau_hinh

import (
	"log"
	"os"
)

type CauHinhHeThong struct {
	IdFileSheet string // ID của file Google Sheet
	CongChayWeb string // Port (VD: 8080)
	// [ĐÃ XÓA] FileChungThuc -> Để không gây lỗi tìm file
}

var BienCauHinh CauHinhHeThong

func KhoiTaoCauHinh() {
	// 1. Lấy ID Sheet
	idSheet := os.Getenv("SHEET_ID")
	if idSheet == "" {
		// ID mặc định của bạn
		idSheet = "17f5js4C9rY7GPd4TOyBidkUPw3vCC6qv6y8KlF3vNs8"
	}

	// 2. Lấy Port
	congWeb := os.Getenv("PORT")
	if congWeb == "" {
		congWeb = "8080"
	}

	// 3. Gán vào biến
	BienCauHinh = CauHinhHeThong{
		IdFileSheet: idSheet,
		CongChayWeb: congWeb,
	}

	log.Println("--- [CẤU HÌNH] Đã tải xong cài đặt (Public Mode) ---")
	log.Println("--- ID Sheet:", BienCauHinh.IdFileSheet)
}
