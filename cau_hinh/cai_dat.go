package cau_hinh

import (
	"log"
	"os"
	"time"
)

type CauHinhHeThong struct {
	IdFileSheet string
	CongChayWeb string
}

var BienCauHinh CauHinhHeThong

// [CẤU HÌNH] Các tham số hệ thống
const (
	// Chu kỳ ghi dữ liệu xuống Sheet (Đang test: 2s)
	ChuKyGhiSheet = 2 * time.Second
    // Dữ liệu bắt đầu từ dòng 11 trong sheet google
	DongBatDauDuLieu = 11 
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

	log.Println("--- [CẤU HÌNH] Đã tải xong (Mode: Public + Batch 2s) ---")
}
