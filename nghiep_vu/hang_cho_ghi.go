package nghiep_vu // <--- ĐÃ SỬA CHUẨN

import (
	"fmt"
	"log"
	"sync"
	"time"

	"app/cau_hinh"
	"app/kho_du_lieu"
	"google.golang.org/api/sheets/v4"
)

// Cấu trúc lệnh ghi (Buffer trong RAM)
type LenhGhi struct {
	TenSheet string
	Dong     int
	DuLieu   []interface{}
	LaGhiMoi bool // True = Append, False = Update
}

var (
	HangChoGhi []LenhGhi
	KhoaHangCho sync.Mutex
)

// Hàm 1: Đẩy dữ liệu vào hàng chờ
func ThemVaoHangCho(tenSheet string, dong int, duLieu []interface{}, laGhiMoi bool) {
	KhoaHangCho.Lock()
	defer KhoaHangCho.Unlock()

	HangChoGhi = append(HangChoGhi, LenhGhi{
		TenSheet: tenSheet,
		Dong:     dong,
		DuLieu:   duLieu,
		LaGhiMoi: laGhiMoi,
	})
}

// Hàm 2: Worker chạy ngầm (Trigger mỗi 5 giây)
func KhoiTaoWorkerGhiSheet() {
	go func() {
		log.Printf("⏳ [WORKER] Đã kích hoạt chế độ ghi Batch (%v/lần)", cau_hinh.ChuKyGhiSheet)
		ticker := time.NewTicker(cau_hinh.ChuKyGhiSheet)
		
		for range ticker.C {
			ThucHienGhiSheet(false)
		}
	}()
}

// Hàm 3: Xử lý ghi thực tế
func ThucHienGhiSheet(epBuoc bool) {
	KhoaHangCho.Lock()
	count := len(HangChoGhi)
	if count == 0 {
		KhoaHangCho.Unlock()
		return
	}

	dsCanGhi := make([]LenhGhi, count)
	copy(dsCanGhi, HangChoGhi)
	HangChoGhi = make([]LenhGhi, 0)
	KhoaHangCho.Unlock()

	log.Printf("💾 [BATCH] Worker tỉnh giấc - Đang ghi %d lệnh xuống Sheet...", count)

	srv := kho_du_lieu.DichVuSheet
	if srv == nil {
		log.Println("❌ Lỗi: Mất kết nối Google Sheet API")
		return
	}
	
	spreadId := cau_hinh.BienCauHinh.IdFileSheet

	for _, lenh := range dsCanGhi {
		if lenh.LaGhiMoi {
			// APPEND
			rangeVal := fmt.Sprintf("%s!A1", lenh.TenSheet)
			rb := &sheets.ValueRange{
				Values: [][]interface{}{lenh.DuLieu},
			}
			_, err := srv.Spreadsheets.Values.Append(spreadId, rangeVal, rb).ValueInputOption("RAW").Do()
			if err != nil {
				log.Printf("❌ Lỗi Append %s: %v", lenh.TenSheet, err)
			}
		} else {
			// UPDATE
			rangeVal := fmt.Sprintf("%s!A%d", lenh.TenSheet, lenh.Dong)
			rb := &sheets.ValueRange{
				Values: [][]interface{}{lenh.DuLieu},
			}
			_, err := srv.Spreadsheets.Values.Update(spreadId, rangeVal, rb).ValueInputOption("RAW").Do()
			if err != nil {
				log.Printf("❌ Lỗi Update %s dòng %d: %v", lenh.TenSheet, lenh.Dong, err)
			}
		}
	}
	log.Println("✅ [BATCH] Hoàn tất đợt ghi dữ liệu.")
}
