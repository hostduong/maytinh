package nghiep_vu

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"app/cau_hinh"
	"app/mo_hinh"
)

var mtxKH sync.Mutex

// 1. TÌM KIẾM
func TimKhachHangTheoCookie(cookie string) (*mo_hinh.KhachHang, bool) {
	khoa := BoQuanLyKhoa.LayKhoa(CacheKhachHang.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()

	for _, kh := range CacheKhachHang.DuLieu {
		if kh.Cookie == cookie {
			return kh, true 
		}
	}
	return nil, false
}

func TimKhachHangTheoUserOrEmail(input string) (*mo_hinh.KhachHang, bool) {
	khoa := BoQuanLyKhoa.LayKhoa(CacheKhachHang.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()

	for _, kh := range CacheKhachHang.DuLieu {
		if kh.TenDangNhap == input || kh.Email == input {
			return kh, true
		}
	}
	return nil, false
}

// [ĐÃ SỬA] Chỉ kiểm tra User và Email, bỏ tham số phone
func KiemTraTonTaiUserEmail(user, email string) bool {
	khoa := BoQuanLyKhoa.LayKhoa(CacheKhachHang.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()

	for _, kh := range CacheKhachHang.DuLieu {
		if kh.TenDangNhap == user { return true }
		if email != "" && kh.Email == email { return true }
		// Đã xóa dòng check DienThoai
	}
	return false
}

// 2. SINH MÃ & ĐẾM (Giữ nguyên)
func DemSoLuongKhachHang() int {
	khoa := BoQuanLyKhoa.LayKhoa(CacheKhachHang.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()
	return len(CacheKhachHang.DuLieu)
}

func TaoMaKhachHangMoi() string {
	khoa := BoQuanLyKhoa.LayKhoa(CacheKhachHang.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()

	maxID := 0
	for _, kh := range CacheKhachHang.DuLieu {
		if strings.HasPrefix(kh.MaKhachHang, "KH_") {
			parts := strings.Split(kh.MaKhachHang, "_")
			if len(parts) == 2 {
				id, err := strconv.Atoi(parts[1])
				if err == nil && id > maxID {
					maxID = id
				}
			}
		}
	}
	return fmt.Sprintf("KH_%04d", maxID+1)
}

func LayDongKhachHang(maKH string) int {
	khoa := BoQuanLyKhoa.LayKhoa(CacheKhachHang.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()
	
	if kh, ok := CacheKhachHang.DuLieu[maKH]; ok {
		return kh.DongTrongSheet
	}
	return 0
}

// 3. GHI & CẬP NHẬT (Giữ nguyên)
func CapNhatPhienDangNhapKH(maKH string, newCookie string, newExpired int64) {
	mtxKH.Lock()
	defer mtxKH.Unlock()

	khoa := BoQuanLyKhoa.LayKhoa(CacheKhachHang.TenKey)
	khoa.Lock() 
	kh, ok := CacheKhachHang.DuLieu[maKH]
	if ok {
		kh.Cookie = newCookie
		kh.CookieExpired = newExpired
	}
	khoa.Unlock()

	if ok {
		sID := cau_hinh.BienCauHinh.IdFileSheet
		ThemVaoHangCho(sID, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_Cookie, newCookie)
		ThemVaoHangCho(sID, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_CookieExpired, newExpired)
	}
}

func ThemKhachHangMoi(kh *mo_hinh.KhachHang) {
	mtxKH.Lock()
	defer mtxKH.Unlock()

	khoa := BoQuanLyKhoa.LayKhoa(CacheKhachHang.TenKey)
	khoa.Lock()
	
	maxRow := mo_hinh.DongBatDauDuLieu - 1
	for _, item := range CacheKhachHang.DuLieu {
		if item.DongTrongSheet > maxRow {
			maxRow = item.DongTrongSheet
		}
	}
	newRow := maxRow + 1
	kh.DongTrongSheet = newRow

	CacheKhachHang.DuLieu[kh.MaKhachHang] = kh
	khoa.Unlock()

	sID := cau_hinh.BienCauHinh.IdFileSheet
	sName := "KHACH_HANG"

	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotKH_MaKhachHang, kh.MaKhachHang)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotKH_TenDangNhap, kh.TenDangNhap)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotKH_MatKhauHash, kh.MatKhauHash) 
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotKH_Cookie, kh.Cookie)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotKH_CookieExpired, kh.CookieExpired)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotKH_MaPinHash, kh.MaPinHash)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotKH_LoaiKhachHang, kh.LoaiKhachHang)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotKH_TenKhachHang, kh.TenKhachHang)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotKH_DienThoai, kh.DienThoai)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotKH_Email, kh.Email)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotKH_NgaySinh, kh.NgaySinh)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotKH_GioiTinh, kh.GioiTinh)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotKH_VaiTroQuyenHan, kh.VaiTroQuyenHan)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotKH_TrangThai, kh.TrangThai)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotKH_NgayTao, time.Now().Format("2006-01-02 15:04:05"))
}

func CapNhatHanCookieRAM(maKH string, newExpired int64) {
	khoa := BoQuanLyKhoa.LayKhoa(CacheKhachHang.TenKey)
	khoa.Lock()
	defer khoa.Unlock()
	
	if kh, ok := CacheKhachHang.DuLieu[maKH]; ok {
		kh.CookieExpired = newExpired
	}
}
