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
