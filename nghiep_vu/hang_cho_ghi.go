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

// Cấu trúc lưu trữ thông minh:
// [SpreadsheetID] -> [SheetName] -> [Row] -> [Col] -> Value
type SmartQueue struct {
	sync.Mutex
	Data map[string]map[string]map[int]map[int]interface{}
}

// Khởi tạo bộ nhớ đệm
var BoNhoGhi = &SmartQueue{
	Data: make(map[string]map[string]map[int]map[int]interface{}),
}

// Hàm giao tiếp chuẩn (Giữ nguyên 5 tham số để tương thích code cũ)
// Hỗ trợ nhiều Web chạy cùng lúc vì có tham số spreadId
func ThemVaoHangCho(spreadId string, sheetName string, row int, col int, value interface{}) {
	BoNhoGhi.Lock()
	defer BoNhoGhi.Unlock()

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

	// 2. Ghi đè thông minh (Last write wins)
	// Ví dụ: Trong 5s, User đổi tên 3 lần -> Chỉ lưu lần cuối cùng
	BoNhoGhi.Data[spreadId][sheetName][row][col] = value
}

// Worker chạy ngầm (5 giây/lần)
func KhoiTaoWorkerGhiSheet() {
	go func() {
		log.Printf("🚀 [MULTI-TENANT] Kích hoạt Worker ghi đa luồng (%v/lần)", cau_hinh.ChuKyGhiSheet)
		ticker := time.NewTicker(cau_hinh.ChuKyGhiSheet)
		for range ticker.C {
			XuLyGhiThongMinh()
		}
	}()
}

// Hàm xử lý chính: Tối ưu Quota
func XuLyGhiThongMinh() {
	BoNhoGhi.Lock()
	if len(BoNhoGhi.Data) == 0 {
		BoNhoGhi.Unlock()
		return
	}

	// Chép dữ liệu ra biến tạm để giải phóng RAM cho luồng khác ghi tiếp
	// snapshotData chứa toàn bộ dữ liệu của TẤT CẢ các web đang chờ
	snapshotData := BoNhoGhi.Data
	BoNhoGhi.Data = make(map[string]map[string]map[int]map[int]interface{}) // Reset sạch
	BoNhoGhi.Unlock()

	log.Println("⚡ [SMART BATCH] Bắt đầu phân tích và ghi dữ liệu...")
	
	srv := kho_du_lieu.DichVuSheet
	if srv == nil { return }

	// DUYỆT QUA TỪNG WEBSITE (TỪNG SPREADSHEET ID)
	for spreadId, sheetsMap := range snapshotData {
		
		// Danh sách các vùng cần update cho Website này
		var dataToUpdate []*sheets.ValueRange
		totalCells := 0

		// Duyệt từng Sheet (KHACH_HANG, SAN_PHAM...)
		for sheetName, rows := range sheetsMap {
			// Duyệt từng Dòng
			for r, cols := range rows {
				
				// --- THUẬT TOÁN GOM CỘT LIỀN KỀ (CONTIGUOUS RANGE) ---
				
				// B1: Lấy danh sách cột và sắp xếp tăng dần (0, 1, 2, 5, 6...)
				var colIndexes []int
				for c := range cols { colIndexes = append(colIndexes, c) }
				sort.Ints(colIndexes)

				if len(colIndexes) == 0 { continue }
				
				// B2: Gom nhóm
				startCol := colIndexes[0]
				prevCol := colIndexes[0]
				currentValues := []interface{}{}
				currentValues = append(currentValues, cols[startCol])
				totalCells++

				for i := 1; i < len(colIndexes); i++ {
					currCol := colIndexes[i]
					
					// Nếu cột hiện tại liền kề cột trước (VD: 1 tiếp sau 0) -> Gom tiếp
					if currCol == prevCol+1 {
						currentValues = append(currentValues, cols[currCol])
						prevCol = currCol
						totalCells++
					} else {
						// Nếu bị ngắt quãng (VD: đang 2 nhảy sang 5) -> Đóng gói dải trước (A..:C..)
						rangeStr := fmt.Sprintf("%s!%s%d", sheetName, layTenCot(startCol), r)
						vr := &sheets.ValueRange{
							Range: rangeStr,
							Values: [][]interface{}{currentValues},
						}
						dataToUpdate = append(dataToUpdate, vr)

						// Bắt đầu dải mới
						startCol = currCol
						prevCol = currCol
						currentValues = []interface{}{}
						currentValues = append(currentValues, cols[currCol])
						totalCells++
					}
				}
				
				// Đóng gói dải cuối cùng còn sót lại
				if len(currentValues) > 0 {
					rangeStr := fmt.Sprintf("%s!%s%d", sheetName, layTenCot(startCol), r)
					vr := &sheets.ValueRange{
						Range: rangeStr,
						Values: [][]interface{}{currentValues},
					}
					dataToUpdate = append(dataToUpdate, vr)
				}
			}
		}

		// GỬI REQUEST - 1 LẦN DUY NHẤT CHO 1 WEBSITE
		if len(dataToUpdate) > 0 {
			req := &sheets.BatchUpdateValuesRequest{
				ValueInputOption: "RAW",
				Data:             dataToUpdate,
			}
			
			// Gọi API Google
			_, err := srv.Spreadsheets.Values.BatchUpdate(spreadId, req).Do()
			if err != nil {
				log.Printf("❌ [Spreadsheet %s...] Lỗi BatchUpdate: %v", spreadId[0:5], err)
				// Ở đây nếu cần kỹ hơn thì có thể đẩy lại vào hàng chờ (Retry mechanism)
			} else {
				log.Printf("✅ [Spreadsheet %s...] Ghi thành công %d ô dữ liệu (gói trong %d dải).", 
					spreadId[0:5], totalCells, len(dataToUpdate))
			}
		}
	}
}

// Helper: Đổi số thành chữ (0->A, 1->B... 26->AA)
func layTenCot(i int) string {
	if i < 0 { return "A" }
	const abc = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if i < 26 {
		return string(abc[i])
	}
	return layTenCot(i/26-1) + string(abc[i%26])
}
