package nghiep_vu

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"app/bao_mat"
	"app/mo_hinh"
)

// =============================================================
// CÁC HÀM TRA CỨU & KIỂM TRA (Fix lỗi undefined)
// =============================================================

// 1. Tìm theo Cookie (Session)
func TimKhachHangTheoCookie(cookie string) (*mo_hinh.KhachHang, bool) {
	// Vì map lưu theo User/Email nên phải duyệt (Tuy chậm hơn chút nhưng an toàn)
	// Do số lượng user trong RAM ít nên không đáng kể.
	// Nếu user lớn > 10.000, ta sẽ tối ưu sau.
	for _, kh := range CacheKhachHang.DuLieu {
		if kh.Cookie == cookie && kh.Cookie != "" {
			// Kiểm tra hạn sử dụng cookie
			if time.Now().Unix() > kh.CookieExpired {
				return nil, false
			}
			return kh, true
		}
	}
	return nil, false
}

// 2. Tìm theo User HOẶC Email (Dùng Map nên cực nhanh)
func TimKhachHangTheoUserOrEmail(input string) (*mo_hinh.KhachHang, bool) {
	input = strings.ToLower(strings.TrimSpace(input))
	if kh, ok := CacheKhachHang.DuLieu[input]; ok {
		return kh, true
	}
	return nil, false
}

// 3. Kiểm tra tồn tại (Trả về bool)
func KiemTraTonTaiUserEmail(user, email string) bool {
	user = strings.ToLower(strings.TrimSpace(user))
	email = strings.ToLower(strings.TrimSpace(email))
	
	if _, ok := CacheKhachHang.DuLieu[user]; ok {
		return true
	}
	if email != "" {
		if _, ok := CacheKhachHang.DuLieu[email]; ok {
			return true
		}
	}
	return false
}

// 4. Đếm tổng số khách hàng (Unique)
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

// 5. Lấy dòng trong Sheet
func LayDongKhachHang(maKH string) int {
	if kh, ok := CacheKhachHang.DuLieu[maKH]; ok {
		return kh.DongTrongSheet
	}
	return 0
}

// 6. Cập nhật Phiên đăng nhập (Cookie)
func CapNhatPhienDangNhapKH(kh *mo_hinh.KhachHang) {
	// Cập nhật trong RAM (Vì kh là con trỏ nên nó tự update vào Cache)
	// Chỉ cần đẩy lệnh Update xuống Sheet
	
	// Map struct ra mảng dữ liệu
	row := ConvertKhachHangToRow(*kh) // Hàm này nhận value nên phải *kh
	
	// Đẩy vào hàng chờ (LaGhiMoi = false => Update)
	ThemVaoHangCho("KHACH_HANG", kh.DongTrongSheet, row, false)
}

// =============================================================
// LOGIC NGHIỆP VỤ CHÍNH
// =============================================================

// Hàm xử lý đăng ký tài khoản mới (Input là con trỏ để khớp với code cũ)
func ThemKhachHangMoi(input *mo_hinh.KhachHang) error {
	// 1. Chuẩn hóa
	input.TenDangNhap = strings.ToLower(strings.TrimSpace(input.TenDangNhap))
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	// 2. Kiểm tra trùng
	if KiemTraTonTaiUserEmail(input.TenDangNhap, input.Email) {
		return errors.New("Tên đăng nhập hoặc Email đã tồn tại")
	}

	// 3. Logic Founder
	var chucVu, vaiTro string
	if DemSoLuongKhachHang() == 0 {
		chucVu = "Quản trị viên cấp cao"
		vaiTro = "admin_root"
		log.Println("👑 [FOUNDER] Admin Root khởi tạo")
	} else {
		chucVu = "Khách hàng"
		vaiTro = "customer"
	}

	// 4. Tạo dữ liệu
	maMoi := TaoMaKhachHangMoi()
	now := time.Now().Format("2006-01-02 15:04:05")
	
	hashPass, _ := bao_mat.HashMatKhau(input.MatKhauHash)
	hashPin, _ := bao_mat.HashMatKhau(input.MaPinHash)

	// Cập nhật trực tiếp vào con trỏ input để trả về cho Controller nếu cần
	input.MaKhachHang = maMoi
	input.MatKhauHash = hashPass
	input.MaPinHash = hashPin
	input.ChucVu = chucVu
	input.VaiTroQuyenHan = vaiTro
	input.TrangThai = 1
	input.NgayTao = now
	input.NgayCapNhat = now

	// 5. Lưu vào Cache (RAM)
	// Lưu ý: Phải tạo bản copy hoặc lưu con trỏ cẩn thận. 
	// Ở đây ta lưu con trỏ input vào map.
	CacheKhachHang.DuLieu[maMoi] = input
	CacheKhachHang.DuLieu[input.TenDangNhap] = input
	if input.Email != "" {
		CacheKhachHang.DuLieu[input.Email] = input
	}

	// 6. Đẩy vào Hàng chờ (Worker 5s sẽ ghi)
	row := ConvertKhachHangToRow(*input)
	ThemVaoHangCho("KHACH_HANG", 0, row, true) // Append

	return nil
}

// Helper: Tạo mã KH
func TaoMaKhachHangMoi() string {
	maxID := 0
	seen := make(map[string]bool)
	for _, kh := range CacheKhachHang.DuLieu {
		if seen[kh.MaKhachHang] { continue }
		seen[kh.MaKhachHang] = true
		parts := strings.Split(kh.MaKhachHang, "_")
		if len(parts) == 2 {
			id, _ := fmt.Sscanf(parts[1], "%d", &maxID)
			if id > maxID { maxID = id }
		}
	}
	return fmt.Sprintf("KH_%04d", maxID+1)
}

// [UPDATED] Helper: Map ĐẦY ĐỦ các cột để Admin có thể sửa trên Sheet
func ConvertKhachHangToRow(kh mo_hinh.KhachHang) []interface{} {
	row := make([]interface{}, 35)
	
	row[mo_hinh.CotKH_MaKhachHang] = kh.MaKhachHang
	row[mo_hinh.CotKH_TenDangNhap] = kh.TenDangNhap
	row[mo_hinh.CotKH_MatKhauHash] = kh.MatKhauHash
	row[mo_hinh.CotKH_Cookie] = kh.Cookie
	row[mo_hinh.CotKH_CookieExpired] = kh.CookieExpired
	row[mo_hinh.CotKH_MaPinHash] = kh.MaPinHash

	row[mo_hinh.CotKH_LoaiKhachHang] = kh.LoaiKhachHang
	row[mo_hinh.CotKH_TenKhachHang] = kh.TenKhachHang
	row[mo_hinh.CotKH_DienThoai] = kh.DienThoai
	row[mo_hinh.CotKH_Email] = kh.Email
	
	row[mo_hinh.CotKH_UrlFb] = kh.UrlFb
	row[mo_hinh.CotKH_Zalo] = kh.Zalo
	row[mo_hinh.CotKH_UrlTele] = kh.UrlTele
	row[mo_hinh.CotKH_UrlTiktok] = kh.UrlTiktok
	row[mo_hinh.CotKH_DiaChi] = kh.DiaChi
	
	row[mo_hinh.CotKH_NgaySinh] = kh.NgaySinh
	row[mo_hinh.CotKH_GioiTinh] = kh.GioiTinh
	
	row[mo_hinh.CotKH_MaSoThue] = kh.MaSoThue
	row[mo_hinh.CotKH_DangNo] = kh.DangNo
	row[mo_hinh.CotKH_TongMua] = kh.TongMua

	row[mo_hinh.CotKH_ChucVu] = kh.ChucVu
	row[mo_hinh.CotKH_VaiTroQuyenHan] = kh.VaiTroQuyenHan
	row[mo_hinh.CotKH_TrangThai] = kh.TrangThai
	row[mo_hinh.CotKH_GhiChu] = kh.GhiChu
	row[mo_hinh.CotKH_NguoiTao] = kh.NguoiTao
	
	row[mo_hinh.CotKH_NgayTao] = kh.NgayTao
	row[mo_hinh.CotKH_NgayCapNhat] = kh.NgayCapNhat

	return row
}
