package khodulieu

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/api/sheets/v4"
)

// TaiDuLieuGoc: Kéo toàn bộ dữ liệu từ Sheet về RAM
// Input: tenBang (Ví dụ: "SanPham")
// Output: Mảng 2 chiều chứa dữ liệu
func (k *KhoRAM) TaiDuLieuGoc(tenBang string) ([][]interface{}, error) {
	// 1. Định nghĩa vùng dữ liệu (Lấy từ cột A đến Z)
	vungDuLieu := fmt.Sprintf("%s!A:Z", tenBang)

	// 2. Gọi API Google (Đọc dữ liệu)
	resp, err := k.DichVuSheet.Spreadsheets.Values.Get(k.IDFileSheet, vungDuLieu).Do()
	if err != nil {
		log.Printf("❌ Lỗi khi tải dữ liệu bảng %s: %v", tenBang, err)
		return nil, err
	}

	// 3. Trả về dữ liệu thô
	if len(resp.Values) == 0 {
		fmt.Printf("⚠️ Bảng %s trống, chưa có dữ liệu.\n", tenBang)
		return [][]interface{}{}, nil
	}

	fmt.Printf("📥 Đã tải %d dòng từ bảng %s vào RAM.\n", len(resp.Values), tenBang)
	return resp.Values, nil
}

// GhiMeXuongSheet: Xử lý danh sách các yêu cầu đang chờ (Batch Processing)
// Hàm này được gọi bởi Worker (ThoSan) sau mỗi 5s hoặc khi hàng đợi đầy
func (k *KhoRAM) GhiMeXuongSheet(danhSach []YeuCauGhi) {
	// Bước 1: Phân loại dữ liệu theo từng Bảng (Sheet) để ghi 1 lần
	// Map: Key = Tên Bảng, Value = Danh sách các dòng cần thêm
	duLieuGomNhom := make(map[string][][]interface{})

	for _, yeuCau := range danhSach {
		// Chỉ xử lý thao tác THÊM (Append) theo lô
		if yeuCau.LoaiThaoTac == "THEM" {
			duLieuGomNhom[yeuCau.TenBang] = append(duLieuGomNhom[yeuCau.TenBang], yeuCau.DongDuLieu)
		} else if yeuCau.LoaiThaoTac == "SUA" {
			// Với thao tác SỬA: Cần xử lý riêng (Update từng cell hoặc row)
			// Để đơn giản cho MVP, ta gọi hàm sửa lẻ ở đây (hoặc implement BatchUpdate sau)
			k.suaDongLe(yeuCau)
		}
	}

	// Bước 2: Duyệt qua từng nhóm và bắn API lên Google
	for tenBang, cacDongMoi := range duLieuGomNhom {
		go k.goiApiThemDong(tenBang, cacDongMoi)
	}
}

// goiApiThemDong: Hàm thực thi gọi Google API (Append)
func (k *KhoRAM) goiApiThemDong(tenBang string, cacDong [][]interface{}) {
	vungGhi := fmt.Sprintf("%s!A1", tenBang) // Google tự tìm dòng trống cuối cùng để chèn

	valueRange := &sheets.ValueRange{
		MajorDimension: "ROWS",
		Values:         cacDong,
	}

	// Gọi API Append
	_, err := k.DichVuSheet.Spreadsheets.Values.Append(k.IDFileSheet, vungGhi, valueRange).
		ValueInputOption("USER_ENTERED"). // Tự động parse số, ngày tháng
		Context(context.Background()).
		Do()

	if err != nil {
		log.Printf("❌ LỖI NGHIÊM TRỌNG: Không thể ghi %d dòng vào bảng %s. Error: %v", len(cacDong), tenBang, err)
		// TODO: Logic Retry (Thử lại) nếu mạng lag -> Đẩy lại vào hàng đợi
	} else {
		fmt.Printf("✅ Đã lưu %d dòng mới vào bảng %s.\n", len(cacDong), tenBang)
	}
}

// suaDongLe: Xử lý sửa 1 dòng cụ thể (Update)
func (k *KhoRAM) suaDongLe(yc YeuCauGhi) {
	// Xác định vùng cần sửa (Ví dụ: SanPham!A5:Z5)
	// Lưu ý: ViTriDong trong Sheets bắt đầu từ 1
	vungSua := fmt.Sprintf("%s!A%d:Z%d", yc.TenBang, yc.ViTriDong, yc.ViTriDong)

	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{yc.DongDuLieu},
	}

	_, err := k.DichVuSheet.Spreadsheets.Values.Update(k.IDFileSheet, vungSua, valueRange).
		ValueInputOption("USER_ENTERED").
		Context(context.Background()).
		Do()

	if err != nil {
		log.Printf("❌ Lỗi khi cập nhật dòng %d bảng %s: %v", yc.ViTriDong, yc.TenBang, err)
	} else {
		fmt.Printf("✏️ Đã cập nhật dòng %d bảng %s.\n", yc.ViTriDong, yc.TenBang)
	}
}
