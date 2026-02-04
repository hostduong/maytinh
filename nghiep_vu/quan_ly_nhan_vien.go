package nghiep_vu

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	// "app/cau_hinh" <-- ĐÃ XÓA DÒNG NÀY
	"app/mo_hinh"
)

var mtxNV sync.Mutex // Khóa an toàn

// 1. Tìm nhân viên theo Cookie
func TimNhanVienTheoCookie(cookie string) (*mo_hinh.NhanVien, bool) {
	for _, nv := range CacheNhanVien.DuLieu {
		if nv.Cookie == cookie {
			return nv, true
		}
	}
	return nil, false
}

// 2. Tìm nhân viên theo Tên Đăng Nhập
func TimNhanVienTheoUsername(username string) (*mo_hinh.NhanVien, bool) {
	for _, nv := range CacheNhanVien.DuLieu {
		if nv.TenDangNhap == username {
			return nv, true
		}
	}
	return nil, false
}

// --- Hàm kiểm tra trùng User hoặc Email ---
func KiemTraTonTaiUserOrEmail(username string, email string) bool {
	for _, nv := range CacheNhanVien.DuLieu {
		if nv.TenDangNhap == username || (email != "" && nv.Email == email) {
			return true
		}
	}
	return false
}

// --- Hàm sinh Mã Nhân Viên tự động (NV_0001 -> NV_0002) ---
func TaoMaNhanVienMoi() string {
	maxID := 0
	
	for _, nv := range CacheNhanVien.DuLieu {
		// Mã có dạng "NV_xxxx". Cắt bỏ "NV_" lấy phần số
		parts := strings.Split(nv.MaNhanVien, "_")
		if len(parts) == 2 {
			id, err := strconv.Atoi(parts[1])
			if err == nil && id > maxID {
				maxID = id
			}
		}
	}
	
	// Tăng lên 1 và format lại thành 4 chữ số (NV_0005)
	return fmt.Sprintf("NV_%04d", maxID+1)
}

// 3. Cập nhật Phiên làm việc
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

// 4. Gia hạn Cookie
func CapNhatHanCookieRAM(maNV string, newExpired int64) {
	mtxNV.Lock()
	defer mtxNV.Unlock()
	nv, ok := CacheNhanVien.DuLieu[maNV]
	if !ok { return }
	nv.CookieExpired = newExpired
	ThemVaoHangCho(CacheNhanVien.SpreadsheetID, "NHAN_VIEN", nv.DongTrongSheet, mo_hinh.CotNV_CookieExpired, newExpired)
}

// 5. Thêm Nhân Viên Mới (Đăng Ký)
func ThemNhanVienMoi(nv *mo_hinh.NhanVien) {
	mtxNV.Lock()
	defer mtxNV.Unlock()

	// Tìm dòng trống tiếp theo
	maxRow := mo_hinh.DongBatDauDuLieu - 1
	for _, item := range CacheNhanVien.DuLieu {
		if item.DongTrongSheet > maxRow {
			maxRow = item.DongTrongSheet
		}
	}
	newRow := maxRow + 1
	nv.DongTrongSheet = newRow

	// Lưu vào RAM
	CacheNhanVien.DuLieu[nv.MaNhanVien] = nv

	// Đẩy vào Hàng Chờ Ghi (Ghi đủ cột A -> L)
	sID := CacheNhanVien.SpreadsheetID
	sName := "NHAN_VIEN"

	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_MaNhanVien, nv.MaNhanVien)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_TenDangNhap, nv.TenDangNhap)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_Email, nv.Email)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_MatKhauHash, nv.MatKhauHash)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_HoTen, nv.HoTen)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_ChucVu, nv.ChucVu)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_MaPin, nv.MaPin)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_Cookie, nv.Cookie)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_CookieExpired, nv.CookieExpired)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_VaiTroQuyenHan, nv.VaiTroQuyenHan)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_TrangThai, nv.TrangThai)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_LanDangNhapCuoi, nv.LanDangNhapCuoi)
}

// 6. Đếm số lượng
func DemSoLuongNhanVien() int {
	return len(CacheNhanVien.DuLieu)
}

// 7. Lấy dòng
func LayDongNhanVien(maNV string) int {
	if nv, ok := CacheNhanVien.DuLieu[maNV]; ok {
		return nv.DongTrongSheet
	}
	return 0
}
