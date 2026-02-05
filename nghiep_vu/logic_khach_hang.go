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

// Hàm xử lý nghiệp vụ Đăng ký mới
func ThemKhachHangMoi(input mo_hinh.KhachHang) error {
	// 1. Chuẩn hóa dữ liệu
	input.TenDangNhap = strings.ToLower(strings.TrimSpace(input.TenDangNhap))
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))

	// 2. Kiểm tra trùng lặp trong Cache
	if _, ok := CacheKhachHang.DuLieu[input.TenDangNhap]; ok {
		return errors.New("Tên đăng nhập đã tồn tại")
	}
	if input.Email != "" {
		if _, ok := CacheKhachHang.DuLieu[input.Email]; ok {
			return errors.New("Email này đã được sử dụng")
		}
	}

	// 3. Logic Founder & Phân quyền
	var chucVu, vaiTro string
	countUsers := 0
	
	seen := make(map[string]bool)
	for _, v := range CacheKhachHang.DuLieu {
		if !seen[v.MaKhachHang] {
			seen[v.MaKhachHang] = true
			countUsers++
		}
	}

	if countUsers == 0 {
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

	newKH := mo_hinh.KhachHang{
		MaKhachHang:    maMoi,
		TenDangNhap:    input.TenDangNhap,
		MatKhauHash:    hashPass,
		MaPinHash:      hashPin,
		
		TenKhachHang:   input.TenKhachHang, 
		Email:          input.Email,
		DienThoai:      input.DienThoai,
		GioiTinh:       input.GioiTinh,
		NgaySinh:       input.NgaySinh,
		
		ChucVu:         chucVu,
		VaiTroQuyenHan: vaiTro,
		LoaiKhachHang:  "", // Để trống để Admin tự điền sau
		TrangThai:      1,
		
		NgayTao:        now,
		NgayCapNhat:    now,
	}

	// 5. Lưu Cache
	CacheKhachHang.DuLieu[maMoi] = &newKH
	CacheKhachHang.DuLieu[newKH.TenDangNhap] = &newKH
	if newKH.Email != "" {
		CacheKhachHang.DuLieu[newKH.Email] = &newKH
	}

	// 6. Đẩy vào Hàng chờ
	row := ConvertKhachHangToRow(newKH)
	ThemVaoHangCho("KHACH_HANG", 0, row, true) 

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
	// Khởi tạo mảng có kích thước đủ lớn (ví dụ 35 cột) để chứa hết các trường
	// Điều này đảm bảo vị trí cột luôn đúng chuẩn
	row := make([]interface{}, 35)
	
	// Nhóm 1: Định danh & Bảo mật
	row[mo_hinh.CotKH_MaKhachHang] = kh.MaKhachHang
	row[mo_hinh.CotKH_TenDangNhap] = kh.TenDangNhap
	row[mo_hinh.CotKH_MatKhauHash] = kh.MatKhauHash
	row[mo_hinh.CotKH_Cookie] = kh.Cookie
	row[mo_hinh.CotKH_CookieExpired] = kh.CookieExpired
	row[mo_hinh.CotKH_MaPinHash] = kh.MaPinHash

	// Nhóm 2: Thông tin cá nhân
	row[mo_hinh.CotKH_LoaiKhachHang] = kh.LoaiKhachHang
	row[mo_hinh.CotKH_TenKhachHang] = kh.TenKhachHang
	row[mo_hinh.CotKH_DienThoai] = kh.DienThoai
	row[mo_hinh.CotKH_Email] = kh.Email
	
	// [MỚI] Map thêm các cột Mạng xã hội & Liên hệ (Admin sẽ điền sau)
	row[mo_hinh.CotKH_UrlFb] = kh.UrlFb
	row[mo_hinh.CotKH_Zalo] = kh.Zalo
	row[mo_hinh.CotKH_UrlTele] = kh.UrlTele
	row[mo_hinh.CotKH_UrlTiktok] = kh.UrlTiktok
	row[mo_hinh.CotKH_DiaChi] = kh.DiaChi
	
	row[mo_hinh.CotKH_NgaySinh] = kh.NgaySinh
	row[mo_hinh.CotKH_GioiTinh] = kh.GioiTinh
	
	// Nhóm 3: Tài chính & Thuế
	row[mo_hinh.CotKH_MaSoThue] = kh.MaSoThue
	row[mo_hinh.CotKH_DangNo] = kh.DangNo
	row[mo_hinh.CotKH_TongMua] = kh.TongMua

	// Nhóm 4: Phân quyền & Quản trị
	row[mo_hinh.CotKH_ChucVu] = kh.ChucVu
	row[mo_hinh.CotKH_VaiTroQuyenHan] = kh.VaiTroQuyenHan
	row[mo_hinh.CotKH_TrangThai] = kh.TrangThai
	row[mo_hinh.CotKH_GhiChu] = kh.GhiChu
	row[mo_hinh.CotKH_NguoiTao] = kh.NguoiTao
	
	// Nhóm 5: Thời gian
	row[mo_hinh.CotKH_NgayTao] = kh.NgayTao
	row[mo_hinh.CotKH_NgayCapNhat] = kh.NgayCapNhat

	return row
}
