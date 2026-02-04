package nghiep_vu

import (
	"sync"
	"app/mo_hinh"
)

var mtxNV sync.Mutex

func TimNhanVienTheoCookie(cookie string) (*mo_hinh.NhanVien, bool) {
	for _, nv := range CacheNhanVien.DuLieu {
		if nv.Cookie == cookie {
			return nv, true
		}
	}
	return nil, false
}

func TimNhanVienTheoUsername(username string) (*mo_hinh.NhanVien, bool) {
	for _, nv := range CacheNhanVien.DuLieu {
		if nv.TenDangNhap == username {
			return nv, true
		}
	}
	return nil, false
}

func CapNhatPhienDangNhap(maNV string, newCookie string, newExpired int64) {
	mtxNV.Lock()
	defer mtxNV.Unlock()

	nv, ok := CacheNhanVien.DuLieu[maNV]
	if !ok { return }

	nv.Cookie = newCookie
	nv.CookieExpired = newExpired

	ThemVaoHangCho(CacheNhanVien.SpreadsheetID, "NHAN_VIEN", nv.DongTrongSheet, mo_hinh.CotNV_Cookie, newCookie)
	ThemVaoHangCho(CacheNhanVien.SpreadsheetID, "NHAN_VIEN", nv.DongTrongSheet, mo_hinh.CotNV_CookieExpired, newExpired)
}

func CapNhatHanCookieRAM(maNV string, newExpired int64) {
	mtxNV.Lock()
	defer mtxNV.Unlock()

	nv, ok := CacheNhanVien.DuLieu[maNV]
	if !ok { return }

	nv.CookieExpired = newExpired
	
	ThemVaoHangCho(CacheNhanVien.SpreadsheetID, "NHAN_VIEN", nv.DongTrongSheet, mo_hinh.CotNV_CookieExpired, newExpired)
}

func LayDongNhanVien(maNV string) int {
	if nv, ok := CacheNhanVien.DuLieu[maNV]; ok {
		return nv.DongTrongSheet
	}
	return 0
}

// ... (Các hàm cũ giữ nguyên)

// 5. Thêm Nhân Viên Mới (Đăng Ký)
func ThemNhanVienMoi(nv *mo_hinh.NhanVien) {
	mtxNV.Lock()
	defer mtxNV.Unlock()

	// 1. Tìm dòng trống tiếp theo trong Sheet
	// (Duyệt qua Cache để tìm dòng lớn nhất đang có dữ liệu)
	maxRow := mo_hinh.DongBatDauDuLieu - 1
	for _, item := range CacheNhanVien.DuLieu {
		if item.DongTrongSheet > maxRow {
			maxRow = item.DongTrongSheet
		}
	}
	newRow := maxRow + 1
	nv.DongTrongSheet = newRow

	// 2. Lưu vào RAM
	CacheNhanVien.DuLieu[nv.MaNhanVien] = nv

	// 3. Đẩy vào Hàng Chờ Ghi (Ghi từng cột)
	// ID Sheet, Tên Sheet, Dòng, Cột, Giá trị
	sheetID := CacheNhanVien.SpreadsheetID
	sheetName := "NHAN_VIEN"

	ThemVaoHangCho(sheetID, sheetName, newRow, mo_hinh.CotNV_MaNhanVien, nv.MaNhanVien)
	ThemVaoHangCho(sheetID, sheetName, newRow, mo_hinh.CotNV_TenDangNhap, nv.TenDangNhap)
	ThemVaoHangCho(sheetID, sheetName, newRow, mo_hinh.CotNV_MatKhauHash, nv.MatKhauHash)
	ThemVaoHangCho(sheetID, sheetName, newRow, mo_hinh.CotNV_HoTen, nv.HoTen)
	ThemVaoHangCho(sheetID, sheetName, newRow, mo_hinh.CotNV_VaiTroQuyenHan, nv.VaiTroQuyenHan)
	ThemVaoHangCho(sheetID, sheetName, newRow, mo_hinh.CotNV_TrangThai, nv.TrangThai)
    // Mặc định ngày tạo là ngày login cuối để đỡ trống
    ThemVaoHangCho(sheetID, sheetName, newRow, mo_hinh.CotNV_LanDangNhapCuoi, nv.LanDangNhapCuoi) 
}

// 6. Đếm số lượng nhân viên (Để biết ai là người đầu tiên)
func DemSoLuongNhanVien() int {
	return len(CacheNhanVien.DuLieu)
}
