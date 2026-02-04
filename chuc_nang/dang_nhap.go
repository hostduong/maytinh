package chuc_nang

import (
	"net/http"
	"strings"
	"time"
	"app/bao_mat"
	"app/cau_hinh"
	"app/nghiep_vu"
	"github.com/gin-gonic/gin"
)

func TrangDangNhap(c *gin.Context) {
	cookie, _ := c.Cookie("session_id")
	if cookie != "" {
		if _, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
			c.Redirect(http.StatusFound, "/") 
			return
		}
	}
	c.HTML(http.StatusOK, "dang_nhap", gin.H{})
}

func XuLyDangNhap(c *gin.Context) {
	inputTaiKhoan := strings.TrimSpace(c.PostForm("ten_dang_nhap"))
	pass          := strings.TrimSpace(c.PostForm("mat_khau"))

	// Tìm Khách Hàng thay vì Nhân Viên
	kh, ok := nghiep_vu.TimKhachHangTheoUserOrEmail(inputTaiKhoan)
	if !ok {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Tài khoản không tồn tại!"})
		return
	}

	if !bao_mat.KiemTraMatKhau(pass, kh.MatKhauHash) { // Kiểm tra với MatKhauHash của KH
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Sai mật khẩu!"})
		return
	}

	sessionID := bao_mat.TaoSessionIDAnToan()
	expiredTime := time.Now().Add(cau_hinh.ThoiGianHetHanCookie).Unix()

	// Cập nhật phiên cho Khách Hàng
	nghiep_vu.CapNhatPhienDangNhapKH(kh.MaKhachHang, sessionID, expiredTime)

	c.SetCookie("session_id", sessionID, int(cau_hinh.ThoiGianHetHanCookie.Seconds()), "/", "", false, true)

	if kh.VaiTroQuyenHan == "admin" {
		c.Redirect(http.StatusFound, "/admin/tong-quan")
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}
