package chuc_nang

import (
	"net/http"
	"time"

	"app/bao_mat"
	"app/cau_hinh"
	"app/mo_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GET /register
func TrangDangKy(c *gin.Context) {
	c.HTML(http.StatusOK, "dang_ky", gin.H{})
}

// POST /register
func XuLyDangKy(c *gin.Context) {
	// 1. Nhận dữ liệu
	hoTen := c.PostForm("ho_ten")
	user := c.PostForm("ten_dang_nhap")
	pass := c.PostForm("mat_khau")
	email := c.PostForm("email")
	maPin := c.PostForm("ma_pin")

	// 2. Kiểm tra trùng
	if nghiep_vu.KiemTraTonTaiUserOrEmail(user, email) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Tên đăng nhập hoặc Email đã tồn tại!"})
		return
	}

	// 3. Mã hóa mật khẩu
	passHash, _ := bao_mat.HashMatKhau(pass)

	// 4. LOGIC QUYỀN HẠN & MÃ SỐ (QUAN TRỌNG)
	var maDinhDanh string
	var quyenHan string
	var chucVu string

	if nghiep_vu.DemSoLuongNhanVien() == 0 {
		// --- NGƯỜI ĐẦU TIÊN (ADMIN) ---
		maDinhDanh = nghiep_vu.TaoMaNhanVienMoi() // NV_0001
		quyenHan = "admin"
		chucVu = "Quản lý cửa hàng"
	} else {
		// --- NGƯỜI THỨ 2 TRỞ ĐI (KHÁCH HÀNG) ---
		maDinhDanh = nghiep_vu.TaoMaKhachHangMoi() // KH_xxxx
		quyenHan = "" // Để rỗng theo yêu cầu (Khách vãng lai)
		chucVu = "Khách hàng"
	}

	// 5. Tạo Session cho Auto-Login
	cookie := uuid.New().String()
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
	// Nếu là Admin -> Vào trang quản trị
	// Nếu là Khách -> Có thể vào trang chủ hoặc trang cá nhân (Tạm thời cứ vào admin/tong-quan để test)
	c.Redirect(http.StatusFound, "/admin/tong-quan")
}
