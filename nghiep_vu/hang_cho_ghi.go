package nghiep_vu

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"app/kho_du_lieu"
	"google.golang.org/api/sheets/v4"
)

// Cấu trúc Hàng Chờ: [SpreadsheetID][SheetName][CellKey] -> Giá trị
// CellKey ví dụ: "10_5" (Dòng 10, Cột 5)
type CauTrucHangCho struct {
	sync.Mutex
	DuLieu map[string]map[string]map[string]interface{}
}

// Biến toàn cục
var HangCho = &CauTrucHangCho{
	DuLieu: make(map[string]map[string]map[string]interface{}),
}

// ---------------------------------------------------------
// 1. HÀM THÊM VÀO HÀNG CHỜ (Ghi RAM - Tốc độ ánh sáng)
// ---------------------------------------------------------
func ThemVaoHangCho(spreadsheetId string, sheetName string, row int, col int, value interface{}) {
	HangCho.Lock()
	defer HangCho.Unlock()

	// Khởi tạo Map cấp 1 (Spreadsheet)
	if HangCho.DuLieu[spreadsheetId] == nil {
		HangCho.DuLieu[spreadsheetId] = make(map[string]map[string]interface{})
	}
	// Khởi tạo Map cấp 2 (SheetName)
	if HangCho.DuLieu[spreadsheetId][sheetName] == nil {
		HangCho.DuLieu[spreadsheetId][sheetName] = make(map[string]interface{})
	}

	// Tạo Key duy nhất cho ô: "Dòng_Cột"
	cellKey := fmt.Sprintf("%d_%d", row, col)

	// Ghi đè giá trị (Last-Write-Wins: Cái mới nhất sẽ thắng)
	HangCho.DuLieu[spreadsheetId][sheetName][cellKey] = value
}

// ---------------------------------------------------------
// 2. WORKER (Người công nhân cần mẫn - Chạy ngầm)
// ---------------------------------------------------------
func KhoiTaoWorkerGhiSheet() {
	go func() {
		log.Println(">>> [WRITE-QUEUE] Worker đã khởi động. Chu kỳ: 10 giây.")
		for {
			// Ngủ 10 giây
			time.Sleep(10 * time.Second)
			
			// Thức dậy và làm việc
			ThucHienGhiSheet(false)
		}
	}()
}

// Hàm thực thi ghi (Được gọi bởi Worker hoặc khi SIGTERM)
func ThucHienGhiSheet(isEmergency bool) {
	// BƯỚC 1: SNAPSHOT (Cắt dữ liệu ra biến tạm)
	HangCho.Lock()
	if len(HangCho.DuLieu) == 0 {
		HangCho.Unlock()
		return // Không có gì để ghi
	}

	// Copy dữ liệu sang biến tạm (BatchDangXuLy)
	batchDangXuLy := HangCho.DuLieu
	
	// Reset hàng chờ chính về rỗng để đón request mới
	HangCho.DuLieu = make(map[string]map[string]map[string]interface{})
	HangCho.Unlock()

	if !isEmergency {
		log.Printf(">>> [WRITE-QUEUE] Bắt đầu ghi %d file Spreadsheets...", len(batchDangXuLy))
	}

	// BƯỚC 2: THỰC THI (Gọi Google API)
	for spreadId, sheetsData := range batchDangXuLy {
		err := guiBatchUpdateGoogle(spreadId, sheetsData)
		
		if err != nil {
			log.Printf("❌ LỖI GHI SHEET [%s]: %v. Đang ROLLBACK...", spreadId, err)
			// BƯỚC 3: ROLLBACK (Nếu lỗi -> Merge ngược lại)
			rollbackData(spreadId, sheetsData)
		} else {
			if !isEmergency {
				log.Printf("✅ Đã ghi xong Sheet [%s]", spreadId)
			}
		}
	}
}

// ---------------------------------------------------------
// 3. LOGIC GỌI GOOGLE API (BatchUpdate)
// ---------------------------------------------------------
func guiBatchUpdateGoogle(spreadId string, data map[string]map[string]interface{}) error {
	var valueRanges []*sheets.ValueRange

	// Duyệt qua từng Sheet và từng Ô để đóng gói
	for sheetName, cells := range data {
		for cellKey, val := range cells {
			// Parse lại row, col từ key "row_col"
			var r, c int
			fmt.Sscanf(cellKey, "%d_%d", &r, &c)

			// Chuyển đổi tọa độ (0, 0) thành A1 Notation (VD: "Sheet1!A1")
			rangeStr := fmt.Sprintf("%s!%s%d", sheetName, layTenCot(c), r+1) // Sheet index từ 1

			vr := &sheets.ValueRange{
				Range:  rangeStr,
				Values: [][]interface{}{{val}}, // Mảng 2 chiều 1x1
			}
			valueRanges = append(valueRanges, vr)
		}
	}

	if len(valueRanges) == 0 { return nil }

	// Gọi API
	req := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "RAW", // Ghi thô (quan trọng cho số và ngày tháng)
		Data:             valueRanges,
	}

	_, err := kho_du_lieu.DichVuSheet.Spreadsheets.Values.BatchUpdate(spreadId, req).Do()
	return err
}

// ---------------------------------------------------------
// 4. LOGIC ROLLBACK (Merge ngược thông minh)
// ---------------------------------------------------------
func rollbackData(spreadId string, failedData map[string]map[string]interface{}) {
	HangCho.Lock()
	defer HangCho.Unlock()

	// Logic: Chỉ merge lại những ô mà trong lúc chờ đợi CHƯA CÓ AI SỬA
	// Nếu user mới đã sửa đè lên rồi -> Giữ cái của user mới (Bỏ cái cũ bị lỗi đi)
	
	if HangCho.DuLieu[spreadId] == nil {
		HangCho.DuLieu[spreadId] = make(map[string]map[string]interface{})
	}

	count := 0
	for sheetName, cells := range failedData {
		if HangCho.DuLieu[spreadId][sheetName] == nil {
			HangCho.DuLieu[spreadId][sheetName] = make(map[string]interface{})
		}

		for key, val := range cells {
			// Kiểm tra: Nếu trong hàng chờ chính CHƯA CÓ key này -> Trả lại
			if _, exists := HangCho.DuLieu[spreadId][sheetName][key]; !exists {
				HangCho.DuLieu[spreadId][sheetName][key] = val
				count++
			}
		}
	}
	log.Printf("🔄 Đã khôi phục %d mục vào hàng chờ để thử lại lần sau.", count)
}

// Tiện ích: Đổi số cột thành chữ (0 -> A, 1 -> B, ...)
func layTenCot(i int) string {
	// Đơn giản hóa cho cột A-Z (Hệ thống nhỏ thường không quá cột Z)
	// Nếu cần > Z (AA, AB...) thì cần thuật toán phức tạp hơn chút
	if i >= 0 && i < 26 {
		return string(rune('A' + i))
	}
	return "A" // Mặc định fallback
}
