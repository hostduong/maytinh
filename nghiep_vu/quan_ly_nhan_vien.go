package nghiep_vu

import (
	"app/mo_hinh"
	"sync"
)

var mtxNV sync.Mutex // Khóa để an toàn khi nhiều người cùng login

// 1. Tìm nhân viên theo Cookie (Dùng cho Middleware)
func TimNhanVienTheoCookie(cookie string) (*mo_hinh.NhanVien, bool) {
	// Vì map key là MaNV, nên ta phải duyệt qua map để tìm Cookie
	// (Với số lượng nhân viên < 1000 thì việc này siêu nhanh, không lo chậm)
	
	for _, nv := range CacheNhanVien.DuLieu {
		if nv.Cookie == cookie {
			return nv, true
		}
	}
	return nil, false
}

// 2. Tìm nhân viên theo Tên Đăng Nhập (Dùng cho Login)
func TimNhanVienTheoUsername(username string) (*mo_hinh.NhanVien, bool) {
	for _, nv := range CacheNhanVien.DuLieu {
		if nv.TenDangNhap == username {
			return nv, true
		}
	}
	return nil, false
}

// 3. Cập nhật hạn Cookie trong RAM (Auto-Renew)
func CapNhatHanCookieRAM(maNV string, newExpired int64) {
	mtxNV.Lock()
	defer mtxNV.Unlock()

	if nv, ok := CacheNhanVien.DuLieu[maNV]; ok {
		nv.CookieExpired = newExpired
	}
}

// 4. Lấy số dòng trong Sheet (Để WriteQueue biết ghi vào đâu)
func LayDongNhanVien(maNV string) int {
	if nv, ok := CacheNhanVien.DuLieu[maNV]; ok {
		return nv.DongTrongSheet
	}
	return 0
}
