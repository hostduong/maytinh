package cau_hinh

import (
	"log"
	"os"
)

type CauHinhHeThong struct {
	IdFileSheet string
	CongChayWeb string
	// Đã xóa FileChungThuc
}

var BienCauHinh CauHinhHeThong

func KhoiTaoCauHinh() {
	// Lấy ID Sheet
	idSheet := os.Getenv("SHEET_ID")
	if idSheet == "" {
		idSheet = "17f5js4C9rY7GPd4TOyBidkUPw3vCC6qv6y8KlF3vNs8"
	}

	// Lấy Port
	congWeb := os.Getenv("PORT")
	if congWeb == "" {
		congWeb = "8080"
	}

	BienCauHinh = CauHinhHeThong{
		IdFileSheet: idSheet,
		CongChayWeb: congWeb,
	}

	log.Println("--- [CẤU HÌNH] Đã tải xong (Chế độ Public) ---")
}
