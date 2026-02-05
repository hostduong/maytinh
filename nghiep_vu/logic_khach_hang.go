package nghiep_vu

import (
	"errors"
	"fmt" // <--- ĐÃ THÊM IMPORT NÀY
	"log"
	"strings"
	"time"

	"app/bao_mat"
	"app/cau_hinh"
	"app/mo_hinh"
)

// =============================================================
// CÁC HÀM TRA CỨU & KIỂM TRA
// =============================================================

func TimKhachHangTheoCookie(cookie string) (*mo_hinh.KhachHang, bool) {
	for _, kh := range CacheKhachHang.DuLieu {
		if kh.Cookie == cookie && kh.Cookie != "" {
			if time.Now().Unix() > kh.CookieExpired { return nil, false }
			return kh, true
		}
	}
	return nil, false
}

func TimKhachHangTheoUserOrEmail(input string) (*mo_hinh.KhachHang, bool) {
	input = strings.ToLower(strings.TrimSpace(input))
	if kh, ok := CacheKhachHang.DuLieu[input]; ok { return kh, true }
	return nil, false
}

func KiemTraTonTaiUserEmail(user, email string) bool {
	user = strings.ToLower(strings.TrimSpace(user))
	email = strings.ToLower(strings.TrimSpace(email))
	if _, ok := CacheKhachHang.DuLieu[user]; ok { return true }
	if email != "" {
		if _, ok := CacheKhachHang.DuLieu[email]; ok { return true }
	}
	return false
}

func DemSoLuongKhachHang() int {
	count := 0
	seen := make(map[string]bool)
	for _, v := range CacheKhachHang.DuLieu {
		if !seen[v.MaKhachHang] {
			seen[v.MaKhachHang] = true
			count++
		}
	}
	return count
}

func LayDongKhachHang(maKH string) int {
	if kh, ok := CacheKhachHang.DuLieu[maKH]; ok { return kh.DongTrongSheet }
	return 0
}

func CapNhatPhienDangNhapKH(kh *mo_hinh.KhachHang) {
	// Gọi 5 tham số để tương thích hệ thống
	idFile := cau_hinh.BienCauHinh.IdFileSheet
	ThemVaoHangCho(idFile, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_Cookie, kh.Cookie)
	ThemVaoHangCho(idFile, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_CookieExpired, kh.CookieExpired)
}

// =============================================================
// LOGIC ĐĂNG KÝ CHÍNH
// =============================================================

func ThemKhachHangMoi(input *mo_hinh.KhachHang) error {
	input.TenDangNhap = strings.ToLower(strings.TrimSpace(input.TenDangNhap))
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	if KiemTraTonTaiUserEmail(input.TenDangNhap, input.Email) {
		return errors.New("Tên đăng nhập hoặc Email đã tồn tại")
	}

	// Logic Founder
	var chucVu, vaiTro string
	if DemSoLuongKhachHang() == 0 {
		chucVu = "Quản trị viên cấp cao"; vaiTro = "admin_root"
		log.Println("👑 [FOUNDER] Admin Root khởi tạo")
	} else {
		chucVu = "Khách hàng"; vaiTro = "customer"
	}

	// Tính toán dòng mới: Lấy dòng lớn nhất hiện có + 1
	dongMoi := mo_hinh.DongBatDauDuLieu
	// Nếu có dữ liệu thì tìm dòng max, nếu không thì ghi vào dòng bắt đầu
	if len(CacheKhachHang.DuLieu) > 0 {
		maxRow := 0
		for _, v := range CacheKhachHang.DuLieu {
			if v.DongTrongSheet > maxRow { maxRow = v.DongTrongSheet }
		}
		if maxRow >= mo_hinh.DongBatDauDuLieu {
			dongMoi = maxRow + 1
		}
	}

	maMoi := TaoMaKhachHangMoi()
	now := time.Now().Format("2006-01-02 15:04:05")
	hashPass, _ := bao_mat.HashMatKhau(input.MatKhauHash)
	hashPin, _ := bao_mat.HashMatKhau(input.MaPinHash)

	// Update struct
	input.MaKhachHang = maMoi
	input.MatKhauHash = hashPass
	input.MaPinHash = hashPin
	input.ChucVu = chucVu
	input.VaiTroQuyenHan = vaiTro
	input.TrangThai = 1
	input.NgayTao = now
	input.NgayCapNhat = now
	input.DongTrongSheet = dongMoi // Quan trọng

	// Lưu Cache
	CacheKhachHang.DuLieu[maMoi] = input
	CacheKhachHang.DuLieu[input.TenDangNhap] = input
	if input.Email != "" { CacheKhachHang.DuLieu[input.Email] = input }

	// GHI XUỐNG SHEET (Ghi từng ô - Batch worker sẽ tự gom lại)
	idFile := cau_hinh.BienCauHinh.IdFileSheet
	sheet := "KHACH_HANG"
	
	ThemVaoHangCho(idFile, sheet, dongMoi, mo_hinh.CotKH_MaKhachHang, input.MaKhachHang)
	ThemVaoHangCho(idFile, sheet, dongMoi, mo_hinh.CotKH_TenDangNhap, input.TenDangNhap)
	ThemVaoHangCho(idFile, sheet, dongMoi, mo_hinh.CotKH_MatKhauHash, input.MatKhauHash)
	ThemVaoHangCho(idFile, sheet, dongMoi, mo_hinh.CotKH_MaPinHash, input.MaPinHash)
	ThemVaoHangCho(idFile, sheet, dongMoi, mo_hinh.CotKH_Email, input.Email)
	ThemVaoHangCho(idFile, sheet, dongMoi, mo_hinh.CotKH_DienThoai, input.DienThoai)
	ThemVaoHangCho(idFile, sheet, dongMoi, mo_hinh.CotKH_TenKhachHang, input.TenKhachHang)
	ThemVaoHangCho(idFile, sheet, dongMoi, mo_hinh.CotKH_NgaySinh, input.NgaySinh)
	ThemVaoHangCho(idFile, sheet, dongMoi, mo_hinh.CotKH_GioiTinh, input.GioiTinh)
	ThemVaoHangCho(idFile, sheet, dongMoi, mo_hinh.CotKH_ChucVu, input.ChucVu)
	ThemVaoHangCho(idFile, sheet, dongMoi, mo_hinh.CotKH_VaiTroQuyenHan, input.VaiTroQuyenHan)
	ThemVaoHangCho(idFile, sheet, dongMoi, mo_hinh.CotKH_TrangThai, input.TrangThai)
	ThemVaoHangCho(idFile, sheet, dongMoi, mo_hinh.CotKH_NgayTao, input.NgayTao)
	ThemVaoHangCho(idFile, sheet, dongMoi, mo_hinh.CotKH_NgayCapNhat, input.NgayCapNhat)
	
	// Map thêm các cột rỗng để giữ chỗ nếu cần (Zalo, Fb...) - Tuỳ chọn
	ThemVaoHangCho(idFile, sheet, dongMoi, mo_hinh.CotKH_LoaiKhachHang, "")

	return nil
}

func TaoMaKhachHangMoi() string {
	maxID := 0
	seen := make(map[string]bool)
	for _, kh := range CacheKhachHang.DuLieu {
		if seen[kh.MaKhachHang] { continue }
		seen[kh.MaKhachHang] = true
		parts := strings.Split(kh.MaKhachHang, "_")
		if len(parts) == 2 {
			var id int
			// Fix lỗi cú pháp Sscanf
			fmt.Sscanf(parts[1], "%d", &id)
			if id > maxID { maxID = id }
		}
	}
	return fmt.Sprintf("KH_%04d", maxID+1)
}
