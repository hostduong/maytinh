package chuc_nang

import (
	"net/http"
	"time"

	"app/bao_mat"
	"app/cau_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
	// "github.com/google/uuid" <--- Xóa dòng này
)

// Hiển thị trang đăng nhập
func TrangDangNhap(c *gin.Context) {
	// --- [MỚI] CHẶN NẾU ĐÃ ĐĂNG NHẬP ---
	cookie, _ := c.Cookie("session_id")
	if cookie != "" {
		if _, ok := nghiep_vu.TimNhanVienTheoCookie(cookie); ok {
			c.Redirect(http.StatusFound, "/") // Đá về trang chủ ngay
			return
		}
	}
	// ------------------------------------

	c.HTML(http.StatusOK, "dang_nhap", gin.H{})
}

// Xử lý đăng nhập
func XuLyDangNhap(c *gin.Context) {
	inputTaiKhoan := c.PostForm("ten_dang_nhap")
	pass := c.PostForm("mat_khau")

	nv, ok := nghiep_vu.TimNhanVienTheoUserHoacEmail(inputTaiKhoan)
	if !ok {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Tài khoản hoặc Email không tồn tại!"})
		return
	}

	if !bao_mat.KiemTraMatKhau(pass, nv.MatKhauHash) {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Sai mật khẩu!"})
		return
	}

	// 3. Tạo Session Siêu Bảo Mật
	// --- [MỚI] ---
	sessionID := bao_mat.TaoSessionIDAnToan()
	// -------------
	
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

// Đăng xuất (Giữ nguyên)
func DangXuat(c *gin.Context) {
	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}
