package nghiep_vu

import (
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"app/cau_hinh"
	"app/kho_du_lieu"
	"google.golang.org/api/sheets/v4"
)

// =============================================================================
// PHẦN 1: CẤU TRÚC DỮ LIỆU & KÊNH TÍN HIỆU
// =============================================================================

// Cấu trúc lưu trữ thông minh: [SpreadsheetID] -> [SheetName] -> [Row] -> [Col] -> Value
type SmartQueue struct {
	sync.Mutex
	Data map[string]map[string]map[int]map[int]interface{}
}

// Bộ nhớ đệm RAM
var BoNhoGhi = &SmartQueue{
	Data: make(map[string]map[string]map[int]map[int]interface{}),
}

// [HYBRID TRIGGER] Kênh báo thức Worker
// Buffer = 1: Chỉ cần biết "có hàng", không cần đếm bao nhiêu hàng
var KenhBaoThuc = make(chan struct{}, 1)

// =============================================================================
// PHẦN 2: HÀM GIAO TIẾP (GIỮ NGUYÊN 5 THAM SỐ)
// =============================================================================

func ThemVaoHangCho(spreadId string, sheetName string, row int, col int, value interface{}) {
	BoNhoGhi.Lock()
	
	// 1. Init Map 4 cấp (Nếu chưa có)
	if BoNhoGhi.Data[spreadId] == nil {
		BoNhoGhi.Data[spreadId] = make(map[string]map[int]map[int]interface{})
	}
	if BoNhoGhi.Data[spreadId][sheetName] == nil {
		BoNhoGhi.Data[spreadId][sheetName] = make(map[int]map[int]interface{})
	}
	if BoNhoGhi.Data[spreadId][sheetName][row] == nil {
		BoNhoGhi.Data[spreadId][sheetName][row] = make(map[int]interface{})
	}

	// 2. Ghi vào RAM
	BoNhoGhi.Data[spreadId][sheetName][row][col] = value
	BoNhoGhi.Unlock()

	// 3. [HYBRID] Bắn tín hiệu đánh thức Worker
	// Sử dụng select non-blocking: Nếu worker đang thức rồi thì thôi, không cần chặn
	select {
	case KenhBaoThuc <- struct{}{}:
		// Đã gửi tín hiệu thành công
	default:
		// Kênh đã đầy (Worker đang chuẩn bị chạy), không làm gì cả
	}
}

// =============================================================================
// PHẦN 3: WORKER THÔNG MINH (CƠ CHẾ LAI)
// =============================================================================

func KhoiTaoWorkerGhiSheet() {
	go func() {
		log.Printf("🚀 [HYBRID WORKER] Đã khởi động. Chế độ: Ngủ đông -> Chờ %v -> Ghi.", cau_hinh.ChuKyGhiSheet)
		
		// Vòng lặp vô tận xử lý tín hiệu
		for {
			// A. NGỦ ĐÔNG: Chờ tín hiệu từ kênh (Block CPU tại đây)
			<-KenhBaoThuc
			
			// B. TỈNH GIẤC & GOM HÀNG (Debounce)
			// Khi có tín hiệu, chờ thêm 5 giây để gom các request tiếp theo
			time.Sleep(cau_hinh.ChuKyGhiSheet)

			// C. THỰC THI
			XuLyGhiThongMinh()
		}
	}()
}

// =============================================================================
// PHẦN 4: LOGIC TỐI ƯU QUOTA & GHI SHEET (GIỮ NGUYÊN LOGIC GOM CỘT)
// =============================================================================

func XuLyGhiThongMinh() {
	BoNhoGhi.Lock()
	if len(BoNhoGhi.Data) == 0 {
		BoNhoGhi.Unlock()
		return
	}

	// Chép dữ liệu ra biến tạm (Snapshot)
	snapshotData := BoNhoGhi.Data
	BoNhoGhi.Data = make(map[string]map[string]map[int]map[int]interface{}) // Reset sạch
	BoNhoGhi.Unlock()

	log.Println("⚡ [SMART BATCH] Đang xử lý dữ liệu...")
	
	srv := kho_du_lieu.DichVuSheet
	if srv == nil { return }

	// DUYỆT QUA TỪNG WEBSITE (SpreadsheetID)
	for spreadId, sheetsMap := range snapshotData {
		var dataToUpdate []*sheets.ValueRange
		totalCells := 0

		for sheetName, rows := range sheetsMap {
			for r, cols := range rows {
				// --- THUẬT TOÁN GOM CỘT LIỀN KỀ ---
				var colIndexes []int
				for c := range cols { colIndexes = append(colIndexes, c) }
				sort.Ints(colIndexes)

				if len(colIndexes) == 0 { continue }
				
				startCol := colIndexes[0]
				prevCol := colIndexes[0]
				currentValues := []interface{}{}
				currentValues = append(currentValues, cols[startCol])
				totalCells++

				for i := 1; i < len(colIndexes); i++ {
					currCol := colIndexes[i]
					if currCol == prevCol+1 { // Liền kề
						currentValues = append(currentValues, cols[currCol])
						prevCol = currCol
						totalCells++
					} else { // Ngắt quãng -> Đóng gói dải cũ
						rangeStr := fmt.Sprintf("%s!%s%d", sheetName, layTenCot(startCol), r)
						vr := &sheets.ValueRange{Range: rangeStr, Values: [][]interface{}{currentValues}}
						dataToUpdate = append(dataToUpdate, vr)

						startCol = currCol
						prevCol = currCol
						currentValues = []interface{}{cols[currCol]}
						totalCells++
					}
				}
				// Đóng gói dải cuối
				if len(currentValues) > 0 {
					rangeStr := fmt.Sprintf("%s!%s%d", sheetName, layTenCot(startCol), r)
					vr := &sheets.ValueRange{Range: rangeStr, Values: [][]interface{}{currentValues}}
					dataToUpdate = append(dataToUpdate, vr)
				}
			}
		}

		// GỬI 1 REQUEST DUY NHẤT
		if len(dataToUpdate) > 0 {
			req := &sheets.BatchUpdateValuesRequest{
				ValueInputOption: "RAW",
				Data:             dataToUpdate,
			}
			_, err := srv.Spreadsheets.Values.BatchUpdate(spreadId, req).Do()
			if err != nil {
				log.Printf("❌ Lỗi Ghi %s: %v", spreadId[0:5], err)
			} else {
				log.Printf("✅ Ghi xong %d ô (%d dải) vào Sheet.", totalCells, len(dataToUpdate))
			}
		}
	}
}

// Helper đổi số thành chữ (A, B, AA...)
func layTenCot(i int) string {
	if i < 0 { return "A" }
	const abc = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if i < 26 { return string(abc[i]) }
	return layTenCot(i/26-1) + string(abc[i%26])
}
