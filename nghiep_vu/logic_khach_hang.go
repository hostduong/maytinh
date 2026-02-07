package nghiep_vu

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"app/bao_mat"
	"app/bo_nho_dem" // [MỚI]
	"app/cau_hinh"
	"app/mo_hinh"
)

func TimKhachHangTheoCookie(cookie string) (*mo_hinh.KhachHang, bool) {
	// [SỬA] bo_nho_dem.CacheKhachHang
	for _, kh := range bo_nho_dem.CacheKhachHang.DuLieu {
		if kh.Cookie == cookie && kh.Cookie != "" {
			if time.Now().Unix() > kh.CookieExpired { return nil, false }
			return kh, true
		}
	}
	return nil, false
}

func TimKhachHangTheoUserOrEmail(input string) (*mo_hinh.KhachHang, bool) {
	input = strings.ToLower(strings.TrimSpace(input))
	// [SỬA] bo_nho_dem.CacheKhachHang
	if kh, ok := bo_nho_dem.CacheKhachHang.DuLieu[input]; ok { return kh, true }
	return nil, false
}

func KiemTraTonTaiUserEmail(user, email string) bool {
	user = strings.ToLower(strings.TrimSpace(user))
	email = strings.ToLower(strings.TrimSpace(email))
	if _, ok := bo_nho_dem.CacheKhachHang.DuLieu[user]; ok { return true }
	if email != "" {
		if _, ok := bo_nho_dem.CacheKhachHang.DuLieu[email]; ok { return true }
	}
	return false
}

func DemSoLuongKhachHang() int {
	// [SỬA] Đếm trực tiếp từ DanhSach (chuẩn hơn)
	return len(bo_nho_dem.CacheKhachHang.DanhSach)
}

func LayDongKhachHang(maKH string) int {
	if kh, ok := bo_nho_dem.CacheKhachHang.DuLieu[maKH]; ok { return kh.DongTrongSheet }
	return 0
}

func CapNhatPhienDangNhapKH(kh *mo_hinh.KhachHang) {
	idFile := cau_hinh.BienCauHinh.IdFileSheet
	ThemVaoHangCho(idFile, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_Cookie, kh.Cookie)
	ThemVaoHangCho(idFile, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_CookieExpired, kh.CookieExpired)
}

func ThemKhachHangMoi(input *mo_hinh.KhachHang) error {
	input.TenDangNhap = strings.ToLower(strings.TrimSpace(input.TenDangNhap))
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	if KiemTraTonTaiUserEmail(input.TenDangNhap, input.Email) {
		return errors.New("Tên đăng nhập hoặc Email đã tồn tại")
	}

	var chucVu, vaiTro string
	if DemSoLuongKhachHang() == 0 {
		chucVu = "Quản trị viên cấp cao"; vaiTro = "admin_root"
		log.Println("👑 [FOUNDER] Admin Root khởi tạo")
	} else {
		chucVu = "Khách hàng"; vaiTro = "customer"
	}

	dongMoi := mo_hinh.DongBatDauDuLieu
	if len(bo_nho_dem.CacheKhachHang.DuLieu) > 0 {
		maxRow := 0
		for _, v := range bo_nho_dem.CacheKhachHang.DuLieu {
			if v.DongTrongSheet > maxRow { maxRow = v.DongTrongSheet }
		}
		if maxRow >= mo_hinh.DongBatDauDuLieu {
			dongMoi = maxRow + 1
		}
	}

	maMoi := TaoMaKhachHangMoi()
	now := time.Now().Format("2006-01-02 15:04:05")
	
	hashPin, _ := bao_mat.HashMatKhau(input.MaPinHash)

	input.MaKhachHang = maMoi
	input.MaPinHash = hashPin
	input.ChucVu = chucVu
	input.VaiTroQuyenHan = vaiTro
	input.TrangThai = 1
	input.NgayTao = now
	input.NgayCapNhat = now
	input.DongTrongSheet = dongMoi

	// [CẦN KHÓA] Ghi vào RAM cần lock để an toàn
	bo_nho_dem.KhoaHeThong.Lock()
	bo_nho_dem.CacheKhachHang.DuLieu[maMoi] = input
	bo_nho_dem.CacheKhachHang.DuLieu[input.TenDangNhap] = input
	if input.Email != "" { bo_nho_dem.CacheKhachHang.DuLieu[input.Email] = input }
	// Đừng quên thêm vào danh sách
	bo_nho_dem.CacheKhachHang.DanhSach = append(bo_nho_dem.CacheKhachHang.DanhSach, input)
	bo_nho_dem.KhoaHeThong.Unlock()

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
	ThemVaoHangCho(idFile, sheet, dongMoi, mo_hinh.CotKH_LoaiKhachHang, "")

	return nil
}

func TaoMaKhachHangMoi() string {
	maxID := 0
	// [SỬA] Duyệt danh sách thay vì map
	for _, kh := range bo_nho_dem.CacheKhachHang.DanhSach {
		parts := strings.Split(kh.MaKhachHang, "_")
		if len(parts) == 2 {
			var id int
			fmt.Sscanf(parts[1], "%d", &id)
			if id > maxID { maxID = id }
		}
	}
	return fmt.Sprintf("KH_%04d", maxID+1)
}
