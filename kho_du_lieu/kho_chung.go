package kho_du_lieu

import (
	"fmt"
	// Thay đường dẫn này bằng tên module thật của bạn
	"app/cau_hinh" 
)

// DocToanBoSheet : Đọc tất cả dữ liệu text từ 1 sheet
// Input: Tên Sheet (VD: "SAN_PHAM")
// Output: Mảng 2 chiều chứa dữ liệu
func DocToanBoSheet(tenSheet string) ([][]interface{}, error) {
	// Lấy ID Sheet từ cấu hình
	idFile := cau_hinh.BienCauHinh.IdFileSheet
	
	// Phạm vi đọc: Từ A1 đến cột cuối cùng (Z)
	phamViDoc := tenSheet + "!A:Z" 

	// Gọi API Google
	ket_qua, loi := DichVuSheet.Spreadsheets.Values.Get(idFile, phamViDoc).Do()
	if loi != nil {
		return nil, fmt.Errorf("không thể đọc dữ liệu từ sheet %s: %v", tenSheet, loi)
	}

	// Trả về mảng dữ liệu thô
	return ket_qua.Values, nil
}

// GhiDongMoi : Thêm 1 dòng dữ liệu vào cuối Sheet
// Input: Tên Sheet, Mảng dữ liệu 1 chiều (1 dòng)
func GhiDongMoi(tenSheet string, dongDuLieu []interface{}) error {
	idFile := cau_hinh.BienCauHinh.IdFileSheet
	phamViGhi := tenSheet + "!A1" // Google tự tìm dòng cuối để append

	// Tạo đối tượng giá trị
	doiTuongGiaTri := &sheets.ValueRange{
		Values: [][]interface{}{dongDuLieu}, // API yêu cầu mảng 2 chiều
	}

	// Gọi API Append
	_, loi := DichVuSheet.Spreadsheets.Values.Append(idFile, phamViGhi, doiTuongGiaTri).ValueInputOption("USER_ENTERED").Do()
	
	if loi != nil {
		return fmt.Errorf("lỗi khi ghi dòng mới vào %s: %v", tenSheet, loi)
	}

	return nil
}

// CapNhatDong : Sửa dữ liệu tại 1 dòng cụ thể
// Input: Tên Sheet, Số dòng (bắt đầu từ 1), Dữ liệu mới
func CapNhatDong(tenSheet string, soDong int, duLieuMoi []interface{}) error {
	idFile := cau_hinh.BienCauHinh.IdFileSheet
	
	// Tạo phạm vi ghi đè (VD: SAN_PHAM!A11:Z11)
	phamViGhi := fmt.Sprintf("%s!A%d:Z%d", tenSheet, soDong, soDong)

	doiTuongGiaTri := &sheets.ValueRange{
		Values: [][]interface{}{duLieuMoi},
	}

	_, loi := DichVuSheet.Spreadsheets.Values.Update(idFile, phamViGhi, doiTuongGiaTri).ValueInputOption("USER_ENTERED").Do()
	
	if loi != nil {
		return fmt.Errorf("lỗi khi cập nhật dòng %d sheet %s: %v", soDong, tenSheet, loi)
	}

	return nil
}
