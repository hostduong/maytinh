package cau_hinh

import (
	"log"
	"os"
	"time"
)

// Giữ nguyên struct cũ của bạn
type CauHinhHeThong struct {
	IdFileSheet string
	CongChayWeb string
}

var BienCauHinh CauHinhHeThong

// Thêm cấu hình Batch Write (Hằng số)
const (
	ChuKyGhiSheet = 5 * time.Second
)

func KhoiTaoCauHinh() {
	idSheet := os.Getenv("SHEET_ID")
	if idSheet == "" {
		idSheet = "17f5js4C9rY7GPd4TOyBidkUPw3vCC6qv6y8KlF3vNs8"
	}

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
