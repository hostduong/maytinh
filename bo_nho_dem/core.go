package bo_nho_dem

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"app/cau_hinh"
	"app/kho_du_lieu"
	"app/mo_hinh"
)

// [KHÓA HỆ THỐNG]
var KhoaHeThong sync.RWMutex
var HeThongDangBan bool = false

// [CALLBACK] Để main.go gán hàm ghi sheet vào đây (tránh vòng lặp import)
var CallbackGhiSheet func(bool)

// =================================================================================
// ĐỊNH NGHĨA STRUCT STORE (KHO CHỨA)
// =================================================================================

type KhoSanPhamStore struct { DuLieu map[string]mo_hinh.SanPham; DanhSach []mo_hinh.SanPham; TenKey string }
type KhoDanhMucStore struct { DuLieu map[string]mo_hinh.DanhMuc; TenKey string }
type KhoThuongHieuStore struct { DuLieu map[string]mo_hinh.ThuongHieu; TenKey string }
type KhoNhaCungCapStore struct { DuLieu map[string]mo_hinh.NhaCungCap; TenKey string }

// Khách hàng có thêm DanhSach để đếm user
type KhoKhachHangStore struct {
	DuLieu        map[string]*mo_hinh.KhachHang
	DanhSach      []*mo_hinh.KhachHang
	TenKey        string
	SpreadsheetID string
}

type KhoPhieuNhapStore struct { DuLieu map[string]mo_hinh.PhieuNhap; DanhSach []mo_hinh.PhieuNhap; TenKey string }
type KhoChiTietPhieuNhapStore struct { DanhSach []mo_hinh.ChiTietPhieuNhap; TenKey string }
type KhoPhieuXuatStore struct { DuLieu map[string]mo_hinh.PhieuXuat; DanhSach []mo_hinh.PhieuXuat; TenKey string }
type KhoChiTietPhieuXuatStore struct { DanhSach []mo_hinh.ChiTietPhieuXuat; TenKey string }
type KhoSerialStore struct { DuLieu map[string]mo_hinh.SerialSanPham; TenKey string }
type KhoHoaDonStore struct { DuLieu map[string]mo_hinh.HoaDon; TenKey string }
type KhoHoaDonChiTietStore struct { DanhSach []mo_hinh.HoaDonChiTiet; TenKey string }
type KhoPhieuThuChiStore struct { DuLieu map[string]mo_hinh.PhieuThuChi; DanhSach []mo_hinh.PhieuThuChi; TenKey string }
type KhoPhieuBaoHanhStore struct { DuLieu map[string]mo_hinh.PhieuBaoHanh; DanhSach []mo_hinh.PhieuBaoHanh; TenKey string }
type KhoKhuyenMaiStore struct { DuLieu map[string]mo_hinh.KhuyenMai; TenKey string }
type KhoCauHinhWebStore struct { DuLieu map[string]mo_hinh.CauHinhWeb; TenKey string }

// =================================================================================
// HELPER FUNCTIONS (DÙNG CHUNG CHO TẤT CẢ CÁC FILE CON)
// =================================================================================

func TaoKeyCache(tenSheet string) string {
	return fmt.Sprintf("%s__%s", strings.TrimSpace(cau_hinh.BienCauHinh.IdFileSheet), tenSheet)
}

func loadSheetData(sheetName string) ([][]interface{}, error) {
	duLieu, err := kho_du_lieu.DocToanBoSheet(sheetName)
	if err != nil {
		log.Printf("LỖI ĐỌC %s: %v", sheetName, err)
		return nil, err
	}
	return duLieu, nil
}

func layString(dong []interface{}, index int) string {
	if index >= len(dong) || dong[index] == nil { return "" }
	return fmt.Sprintf("%v", dong[index])
}

func layInt(dong []interface{}, index int) int {
	str := layString(dong, index)
	if str == "" { return 0 }
	str = strings.ReplaceAll(str, ".", "")
	str = strings.ReplaceAll(str, ",", "")
	str = strings.ReplaceAll(str, " ", "")
	val, _ := strconv.Atoi(str)
	return val
}

func layFloat(dong []interface{}, index int) float64 {
	str := layString(dong, index)
	if str == "" { return 0 }
	str = strings.ReplaceAll(str, "đ", "")
	str = strings.ReplaceAll(str, "USD", "")
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, ".", "")
	str = strings.ReplaceAll(str, ",", "")
	val, _ := strconv.ParseFloat(str, 64)
	return val
}
