package chuc_nang

import (
	"net/http"
	"time"

	"app/bao_mat"
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
	// 1. Nhận dữ liệu từ Form HTML
	hoTen := c.PostForm("ho_ten")
	user := c.PostForm("ten_dang_nhap")
	pass := c.PostForm("mat_khau")
	email := c.PostForm("email")   // Mới
	maPin := c.PostForm("ma_pin")  // Mới

	// 2. Kiểm tra trùng lặp (User hoặc Email)
	if nghiep_vu.KiemTraTonTaiUserOrEmail(user, email) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Tên đăng nhập hoặc Email đã tồn tại!"})
		return
	}

	// 3. Mã hóa mật khẩu
	passHash, _ := bao_mat.HashMatKhau(pass)

	// 4. Logic Quyền hạn & Chức vụ (Người đầu tiên là Admin)
	role := "nhan_vien"
	chucVu := "Nhân viên"
	soLuong := nghiep_vu.DemSoLuongNhanVien()
	if soLuong == 0 {
		role = "admin"
		chucVu = "Quản lý cửa hàng"
	}

	// 5. Sinh dữ liệu tự động
	maNV := nghiep_vu.TaoMaNhanVienMoi() // NV_0001
	
	// Tạo Cookie ngẫu nhiên (UUID) - Kiểm tra trùng cookie là thừa vì xác suất UUID trùng là 0
	cookie := uuid.New().String() 

	// 6. Tạo struct nhân viên mới
	newNV := &mo_hinh.NhanVien{
		MaNhanVien:      maNV,
		TenDangNhap:     user,
		Email:           email,
		MatKhauHash:     passHash,
		HoTen:           hoTen,
		ChucVu:          chucVu,
		MaPin:           maPin,
		Cookie:          cookie,
		CookieExpired:   0, // Mới tạo chưa login nên chưa có hạn
		VaiTroQuyenHan:  role,
		TrangThai:       1, // Active
		LanDangNhapCuoi: time.Now().Format("2006-01-02 15:04:05"),
	}

	// 7. Lưu vào hệ thống
	nghiep_vu.ThemNhanVienMoi(newNV)

	// 8. Thông báo thành công
	c.HTML(http.StatusOK, "dang_nhap", gin.H{
		"Loi": "Đăng ký thành công! Mã NV của bạn là: " + maNV, 
	})
}
