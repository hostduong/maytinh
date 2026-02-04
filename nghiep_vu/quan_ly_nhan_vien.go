package nghiep_vu

import (
	"sync"
	"app/cau_hinh"
	"app/mo_hinh"
)

var mtxNV sync.Mutex // Khóa an toàn

// 1. Tìm nhân viên theo Cookie (Dùng cho Middleware kiểm tra quyền)
func TimNhanVienTheoCookie(cookie string) (*mo_hinh.NhanVien, bool) {
	// Duyệt map để tìm (Tốc độ cực nhanh với RAM)
	for _, nv := range CacheNhanVien.DuLieu {
		if nv.Cookie == cookie {
			return nv, true
		}
	}
	return nil, false
}

// 2. Tìm nhân viên theo Tên Đăng Nhập (Dùng cho trang Login)
func TimNhanVienTheoUsername(username string) (*mo_hinh.NhanVien, bool) {
	for _, nv := range CacheNhanVien.DuLieu {
		if nv.TenDangNhap == username {
			return nv, true
		}
	}
	return nil, false
}

// 3. Cập nhật Phiên làm việc mới (Khi đăng nhập thành công)
func CapNhatPhienDangNhap(maNV string, newCookie string, newExpired int64) {
	mtxNV.Lock()
	defer mtxNV.Unlock()

	nv, ok := CacheNhanVien.DuLieu[maNV]
	if !ok { return }

	// A. Cập nhật RAM ngay lập tức
	nv.Cookie = newCookie
	nv.CookieExpired = newExpired

	// B. Đẩy vào Hàng Chờ Ghi (Để lưu xuống Sheet)
	// Ghi cột Cookie (Cột H - Index 7)
	ThemVaoHangCho(CacheNhanVien.SpreadsheetID, "NHAN_VIEN", nv.DongTrongSheet, mo_hinh.CotNV_Cookie, newCookie)
	// Ghi cột Expired (Cột I - Index 8)
	ThemVaoHangCho(CacheNhanVien.SpreadsheetID, "NHAN_VIEN", nv.DongTrongSheet, mo_hinh.CotNV_CookieExpired, newExpired)
}

// 4. Gia hạn thời gian (Khi sắp hết hạn - Auto Renew)
func CapNhatHanCookieRAM(maNV string, newExpired int64) {
	mtxNV.Lock()
	defer mtxNV.Unlock()

	nv, ok := CacheNhanVien.DuLieu[maNV]
	if !ok { return }

	// Chỉ cần update thời gian
	nv.CookieExpired = newExpired
	
	// Đẩy vào hàng chờ ghi
	ThemVaoHangCho(CacheNhanVien.SpreadsheetID, "NHAN_VIEN", nv.DongTrongSheet, mo_hinh.CotNV_CookieExpired, newExpired)
}
