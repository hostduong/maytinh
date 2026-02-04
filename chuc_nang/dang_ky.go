package chuc_nang

import (
	"net/http"
	"strings" // [MỚI] Thêm thư viện xử lý chuỗi
	"time"

	"app/bao_mat"
	"app/cau_hinh"
	"app/mo_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

// GET /register
func TrangDangKy(c *gin.Context) {
	// CHẶN NẾU ĐÃ ĐĂNG NHẬP
	cookie, _ := c.Cookie("session_id")
	if cookie != "" {
		if _, ok := nghiep_vu.TimNhanVienTheoCookie(cookie); ok {
			c.Redirect(http.StatusFound, "/") // Đá về trang chủ
			return
		}
	}
	
	c.HTML(http.StatusOK, "dang_ky", gin.H{})
}

// POST /register
func XuLyDangKy(c *gin.Context) {
	// 1. Nhận dữ liệu & CẮT KHOẢNG TRẮNG [QUAN TRỌNG]
	hoTen := strings.TrimSpace(c.PostForm("ho_ten"))
	user  := strings.TrimSpace(c.PostForm("ten_dang_nhap"))
	pass  := strings.TrimSpace(c.PostForm("mat_khau"))
	email := strings.TrimSpace(c.PostForm("email"))
	maPin := strings.TrimSpace(c.PostForm("ma_pin"))

	// Validate cơ bản: Không cho phép nhập rỗng
	if user == "" || pass == "" || hoTen == "" {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Vui lòng nhập đầy đủ thông tin bắt buộc!"})
		return
	}

	// 2. Kiểm tra trùng lặp
	if nghiep_vu.KiemTraTonTaiUserOrEmail(user, email) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Tên đăng nhập hoặc Email đã tồn tại!"})
		return
	}

	// 3. Mã hóa mật khẩu
	passHash, _ := bao_mat.HashMatKhau(pass)

	// 4. Logic Quyền hạn (Admin vs Khách)
	var maDinhDanh, quyenHan, chucVu string
	if nghiep_vu.DemSoLuongNhanVien() == 0 {
		maDinhDanh = nghiep_vu.TaoMaNhanVienMoi()
		quyenHan = "admin"
		chucVu = "Quản lý cửa hàng"
	} else {
		maDinhDanh = nghiep_vu.TaoMaKhachHangMoi()
		quyenHan = ""
		chucVu = "Khách hàng"
	}

	// 5. Tạo Session Siêu Bảo Mật
	cookie := bao_mat.TaoSessionIDAnToan()
	expiredTime := time.Now().Add(cau_hinh.ThoiGianHetHanCookie).Unix()

	// 6. Tạo Struct
	newNV := &mo_hinh.NhanVien{
		MaNhanVien:      maDinhDanh,
		TenDangNhap:     user,
		Email:           email,
		MatKhauHash:     passHash,
		HoTen:           hoTen,
		ChucVu:          chucVu,
		MaPin:           maPin,
		Cookie:          cookie,
		CookieExpired:   expiredTime,
		VaiTroQuyenHan:  quyenHan,
		TrangThai:       1,
		LanDangNhapCuoi: time.Now().Format("2006-01-02 15:04:05"),
	}

	// 7. Lưu vào hệ thống
	nghiep_vu.ThemNhanVienMoi(newNV)

	// 8. Auto Login
	c.SetCookie("session_id", cookie, int(cau_hinh.ThoiGianHetHanCookie.Seconds()), "/", "", false, true)

	// 9. Điều hướng
	if quyenHan == "admin" {
		c.Redirect(http.StatusFound, "/admin/tong-quan")
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}
