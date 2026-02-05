package nghiep_vu

import (
	"fmt"
	"log"
	"sync"
	"time"

	"app/cau_hinh"
	"app/kho_du_lieu"
	"google.golang.org/api/sheets/v4"
)

// Cấu trúc Hàng Chờ: [SpreadsheetID][SheetName][CellKey] -> Giá trị
// CellKey dạng "row_col"
type CauTrucHangCho struct {
	sync.Mutex
	DuLieu map[string]map[string]map[string]interface{}
}

// Khởi tạo hàng chờ rỗng
var HangCho = &CauTrucHangCho{
	DuLieu: make(map[string]map[string]map[string]interface{}),
}

// Giữ nguyên tên hàm và 5 tham số để KHÔNG LỖI các file cũ
func ThemVaoHangCho(spreadsheetId string, sheetName string, row int, col int, value interface{}) {
	HangCho.Lock()
	defer HangCho.Unlock()

	// Init map nếu chưa có
	if HangCho.DuLieu[spreadsheetId] == nil {
		HangCho.DuLieu[spreadsheetId] = make(map[string]map[string]interface{})
	}
	if HangCho.DuLieu[spreadsheetId][sheetName] == nil {
		HangCho.DuLieu[spreadsheetId][sheetName] = make(map[string]interface{})
	}

	// Lưu giá trị vào RAM
	cellKey := fmt.Sprintf("%d_%d", row, col)
	HangCho.DuLieu[spreadsheetId][sheetName][cellKey] = value
}

// Worker 5 giây
func KhoiTaoWorkerGhiSheet() {
	go func() {
		// Dùng chu kỳ từ config (5s)
		log.Printf("⏳ [WORKER] Kích hoạt ghi Batch theo ô (%v/lần)", cau_hinh.ChuKyGhiSheet)
		ticker := time.NewTicker(cau_hinh.ChuKyGhiSheet)
		for range ticker.C {
			ThucHienGhiSheet()
		}
	}()
}

func ThucHienGhiSheet() {
	HangCho.Lock()
	if len(HangCho.DuLieu) == 0 {
		HangCho.Unlock()
		return
	}

	// Copy dữ liệu ra để giải phóng lock
	dataCopy := HangCho.DuLieu
	HangCho.DuLieu = make(map[string]map[string]map[string]interface{})
	HangCho.Unlock()

	log.Println("💾 [BATCH] Đang ghi dữ liệu xuống Sheet...")

	srv := kho_du_lieu.DichVuSheet
	if srv == nil { return }

	// Duyệt qua từng File ID
	for spreadId, sheetsMap := range dataCopy {
		var valueRanges []*sheets.ValueRange

		// Duyệt qua từng Sheet (KHACH_HANG, SAN_PHAM...)
		for sheetName, cells := range sheetsMap {
			for cellKey, val := range cells {
				var r, c int
				fmt.Sscanf(cellKey, "%d_%d", &r, &c)

				// Chuyển đổi tọa độ (Row 10, Col 0 -> A10)
				// Lưu ý: Row người dùng truyền vào thường là số thực tế (bắt đầu từ 1)
				cotchu := layTenCot(c)
				rangeStr := fmt.Sprintf("%s!%s%d", sheetName, cotchu, r)

				vr := &sheets.ValueRange{
					Range:  rangeStr,
					Values: [][]interface{}{{val}},
				}
				valueRanges = append(valueRanges, vr)
			}
		}

		if len(valueRanges) == 0 { continue }

		// Gửi 1 request duy nhất chứa hàng trăm ô thay đổi
		req := &sheets.BatchUpdateValuesRequest{
			ValueInputOption: "RAW",
			Data:             valueRanges,
		}

		_, err := srv.Spreadsheets.Values.BatchUpdate(spreadId, req).Do()
		if err != nil {
			log.Printf("❌ Lỗi BatchUpdate file %s: %v", spreadId, err)
			// Nếu cần, bạn có thể thêm logic retry hoặc rollback tại đây
		}
	}
	log.Println("✅ [BATCH] Hoàn tất.")
}

// Hàm hỗ trợ đổi số thành chữ (0 -> A, 1 -> B, ... 26 -> AA)
func layTenCot(i int) string {
	if i < 0 { return "" }
	const abc = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if i < 26 {
		return string(abc[i])
	}
	return layTenCot(i/26-1) + string(abc[i%26])
}
