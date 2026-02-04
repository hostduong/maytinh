package chuc_nang

import (
	"log"
	"net/http"
	"strings"

	"app/bao_mat"
	"app/mo_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

// Trang HTML Quên Mật Khẩu
func TrangQuenMatKhau(c *gin.Context) {
	c.HTML(http.StatusOK, "quen_mat_khau", gin.H{})
}

// CÁCH 1: Xác thực bằng PIN
func XuLyQuenPassBangPIN(c *gin.Context) {
	user := strings.TrimSpace(c.PostForm("user"))
	pin := strings.TrimSpace(c.PostForm("pin"))
	passMoi := strings.TrimSpace(c.PostForm("pass_moi"))

	nv, ok := nghiep_vu.TimNhanVienTheoUserHoacEmail(user)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Tài khoản không tồn tại!"})
		return
	}

	if nv.MaPin == "" || nv.MaPin != pin {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN không chính xác!"})
		return
	}

	// Đổi pass luôn
	hashMoi, _ := bao_mat.HashMatKhau(passMoi)
	nv.MatKhauHash = hashMoi
	nghiep_vu.ThemVaoHangCho(nghiep_vu.CacheNhanVien.SpreadsheetID, "NHAN_VIEN", nv.DongTrongSheet, mo_hinh.CotNV_MatKhauHash, hashMoi)

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã đặt lại mật khẩu thành công! Hãy đăng nhập."})
}

// CÁCH 2: Gửi OTP qua Email (Giả lập)
func XuLyGuiOTPEmail(c *gin.Context) {
	email := strings.TrimSpace(c.PostForm("email"))
	nv, ok := nghiep_vu.TimNhanVienTheoUserHoacEmail(email)
	
	if !ok {
		// Bảo mật: Không báo lỗi "Email không tồn tại" để tránh hacker dò user
		c.JSON(200, gin.H{"status": "ok", "msg": "Nếu email đúng, mã OTP đã được gửi đi!"})
		return
	}

	// 1. Sinh OTP
	otp := nghiep_vu.TaoMaOTP()
	nghiep_vu.LuuOTP(nv.TenDangNhap, otp) // Lưu theo username cho duy nhất

	// 2. Gọi API Giả Lập (Mock)
	go func() {
		// Giả vờ gọi API bên thứ 3 (Xem trong Log Server của Cloud Run để lấy mã)
		log.Printf(">>> [EMAIL MOCK] Gửi đến: %s | Mã OTP: %s", email, otp)
	}()

	c.JSON(200, gin.H{"status": "ok", "msg": "Mã OTP đã được gửi vào email (Xem Log Server)!"})
}

// Xác nhận OTP và Đổi pass
func XuLyQuenPassBangOTP(c *gin.Context) {
	email := strings.TrimSpace(c.PostForm("email"))
	otpInput := strings.TrimSpace(c.PostForm("otp"))
	passMoi := strings.TrimSpace(c.PostForm("pass_moi"))

	nv, ok := nghiep_vu.TimNhanVienTheoUserHoacEmail(email)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Email không hợp lệ!"})
		return
	}

	// Check OTP
	if !nghiep_vu.KiemTraOTP(nv.TenDangNhap, otpInput) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã OTP sai hoặc đã hết hạn!"})
		return
	}

	// Đổi Pass
	hashMoi, _ := bao_mat.HashMatKhau(passMoi)
	nv.MatKhauHash = hashMoi
	nghiep_vu.ThemVaoHangCho(nghiep_vu.CacheNhanVien.SpreadsheetID, "NHAN_VIEN", nv.DongTrongSheet, mo_hinh.CotNV_MatKhauHash, hashMoi)

	c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
}
