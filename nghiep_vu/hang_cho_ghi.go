package nghiep_vu

import (
	"fmt"
	"log"
	"sync"
	"time"

	"app/kho_du_lieu"
	"google.golang.org/api/sheets/v4"
)

// Cấu trúc Hàng Chờ: [SpreadsheetID][SheetName][CellKey] -> Giá trị
type CauTrucHangCho struct {
	sync.Mutex
	DuLieu map[string]map[string]map[string]interface{}
}

var HangCho = &CauTrucHangCho{
	DuLieu: make(map[string]map[string]map[string]interface{}),
}

func ThemVaoHangCho(spreadsheetId string, sheetName string, row int, col int, value interface{}) {
	HangCho.Lock()
	defer HangCho.Unlock()

	if HangCho.DuLieu[spreadsheetId] == nil {
		HangCho.DuLieu[spreadsheetId] = make(map[string]map[string]interface{})
	}
	if HangCho.DuLieu[spreadsheetId][sheetName] == nil {
		HangCho.DuLieu[spreadsheetId][sheetName] = make(map[string]interface{})
	}

	cellKey := fmt.Sprintf("%d_%d", row, col)
	HangCho.DuLieu[spreadsheetId][sheetName][cellKey] = value
}

func KhoiTaoWorkerGhiSheet() {
	go func() {
		log.Println(">>> [WRITE-QUEUE] Worker đã khởi động. Chu kỳ: 10 giây.")
		for {
			time.Sleep(10 * time.Second)
			ThucHienGhiSheet(false)
		}
	}()
}

func ThucHienGhiSheet(isEmergency bool) {
	HangCho.Lock()
	if len(HangCho.DuLieu) == 0 {
		HangCho.Unlock()
		return 
	}

	batchDangXuLy := HangCho.DuLieu
	HangCho.DuLieu = make(map[string]map[string]map[string]interface{})
	HangCho.Unlock()

	if !isEmergency {
		log.Printf(">>> [WRITE-QUEUE] Bắt đầu ghi %d file Spreadsheets...", len(batchDangXuLy))
	}

	for spreadId, sheetsData := range batchDangXuLy {
		err := guiBatchUpdateGoogle(spreadId, sheetsData)
		if err != nil {
			log.Printf("❌ LỖI GHI SHEET [%s]: %v. Đang ROLLBACK...", spreadId, err)
			rollbackData(spreadId, sheetsData)
		} else {
			if !isEmergency {
				log.Printf("✅ Đã ghi xong Sheet [%s]", spreadId)
			}
		}
	}
}

func guiBatchUpdateGoogle(spreadId string, data map[string]map[string]interface{}) error {
	var valueRanges []*sheets.ValueRange

	for sheetName, cells := range data {
		for cellKey, val := range cells {
			var r, c int
			fmt.Sscanf(cellKey, "%d_%d", &r, &c)
			rangeStr := fmt.Sprintf("%s!%s%d", sheetName, layTenCot(c), r+1) 

			vr := &sheets.ValueRange{
				Range:  rangeStr,
				Values: [][]interface{}{{val}},
			}
			valueRanges = append(valueRanges, vr)
		}
	}

	if len(valueRanges) == 0 { return nil }

	req := &sheets.BatchUpdateValuesRequest{
		ValueInputOption: "RAW",
		Data:             valueRanges,
	}

	_, err := kho_du_lieu.DichVuSheet.Spreadsheets.Values.BatchUpdate(spreadId, req).Do()
	return err
}

func rollbackData(spreadId string, failedData map[string]map[string]interface{}) {
	HangCho.Lock()
	defer HangCho.Unlock()

	if HangCho.DuLieu[spreadId] == nil {
		HangCho.DuLieu[spreadId] = make(map[string]map[string]interface{})
	}

	count := 0
	for sheetName, cells := range failedData {
		if HangCho.DuLieu[spreadId][sheetName] == nil {
			HangCho.DuLieu[spreadId][sheetName] = make(map[string]interface{})
		}
		for key, val := range cells {
			if _, exists := HangCho.DuLieu[spreadId][sheetName][key]; !exists {
				HangCho.DuLieu[spreadId][sheetName][key] = val
				count++
			}
		}
	}
	log.Printf("🔄 Đã khôi phục %d mục vào hàng chờ để thử lại lần sau.", count)
}

func layTenCot(i int) string {
	if i >= 0 && i < 26 {
		return string(rune('A' + i))
	}
	return "A"
}
