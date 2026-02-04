package chuc_nang

import (
	"net/http"
	"strings" // [MỚI] Thêm thư viện
	"time"

	"app/bao_mat"
	"app/cau_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

// Hiển thị trang đăng nhập
func TrangDangNhap(c *gin.Context) {
	// CHẶN NẾU ĐÃ ĐĂNG NHẬP
	cookie, _ := c.Cookie("session_id")
	if cookie != "" {
		if _, ok := nghiep_vu.TimNhanVienTheoCookie(cookie); ok {
			c.Redirect(http.StatusFound, "/") 
			return
		}
	}
	c.HTML(http.StatusOK, "dang_nhap", gin.H{})
}

// Xử lý đăng nhập
func XuLyDangNhap(c *gin.Context) {
	// [MỚI] Cắt khoảng trắng thừa ở đầu đuôi
	inputTaiKhoan := strings.TrimSpace(c.PostForm("ten_dang_nhap"))
	pass          := strings.TrimSpace(c.PostForm("mat_khau"))

	// 1. Kiểm tra tài khoản
	nv, ok := nghiep_vu.TimNhanVienTheoUserHoacEmail(inputTaiKhoan)
	if !ok {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Tài khoản hoặc Email không tồn tại!"})
		return
	}

	// 2. Kiểm tra mật khẩu
	if !bao_mat.KiemTraMatKhau(pass, nv.MatKhauHash) {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Sai mật khẩu!"})
		return
	}

	// 3. Tạo Session Siêu Bảo Mật
	sessionID := bao_mat.TaoSessionIDAnToan()
	expiredTime := time.Now().Add(cau_hinh.ThoiGianHetHanCookie).Unix()

	// 4. Cập nhật RAM & Sheet
	nghiep_vu.CapNhatPhienDangNhap(nv.MaNhanVien, sessionID, expiredTime)

	// 5. Set Cookie
	c.SetCookie("session_id", sessionID, int(cau_hinh.ThoiGianHetHanCookie.Seconds()), "/", "", false, true)

	// 6. Điều hướng
	if nv.VaiTroQuyenHan == "admin" {
		c.Redirect(http.StatusFound, "/admin/tong-quan")
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}

// Đăng xuất
func DangXuat(c *gin.Context) {
	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}
