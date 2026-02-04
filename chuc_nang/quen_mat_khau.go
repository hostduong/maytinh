package chuc_nang

import (
	"log"
	"net/http"
	"strings"

	"app/bao_mat"
	"app/cau_hinh" // [MỚI] Import cau_hinh
	"app/mo_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

func TrangQuenMatKhau(c *gin.Context) {
	c.HTML(http.StatusOK, "quen_mat_khau", gin.H{})
}

// CÁCH 1: Xác thực bằng PIN
func XuLyQuenPassBangPIN(c *gin.Context) {
	user := strings.TrimSpace(c.PostForm("user"))
	pin := strings.TrimSpace(c.PostForm("pin"))
	passMoi := strings.TrimSpace(c.PostForm("pass_moi"))

	// [SỬA] TimKhachHangTheoUserOrEmail
	kh, ok := nghiep_vu.TimKhachHangTheoUserOrEmail(user)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Tài khoản không tồn tại!"})
		return
	}

	// [SỬA] Dùng MaPinHash (nếu chưa hash thì so sánh thẳng)
	if kh.MaPinHash != pin {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN không chính xác!"})
		return
	}

	hashMoi, _ := bao_mat.HashMatKhau(passMoi)
	kh.MatKhauHash = hashMoi
	// [SỬA] Ghi vào KHACH_HANG
	nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MatKhauHash, hashMoi)

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã đặt lại mật khẩu thành công! Hãy đăng nhập."})
}

// CÁCH 2: Gửi OTP qua Email
func XuLyGuiOTPEmail(c *gin.Context) {
	email := strings.TrimSpace(c.PostForm("email"))
	// [SỬA] TimKhachHang
	kh, ok := nghiep_vu.TimKhachHangTheoUserOrEmail(email)
	
	if !ok {
		c.JSON(200, gin.H{"status": "ok", "msg": "Nếu email đúng, mã OTP đã được gửi đi!"})
		return
	}

	otp := nghiep_vu.TaoMaOTP()
	nghiep_vu.LuuOTP(kh.TenDangNhap, otp)

	go func() {
		log.Printf(">>> [EMAIL MOCK] Gửi đến: %s | Mã OTP: %s", email, otp)
	}()

	c.JSON(200, gin.H{"status": "ok", "msg": "Mã OTP đã được gửi vào email (Xem Log Server)!"})
}

// Xác nhận OTP
func XuLyQuenPassBangOTP(c *gin.Context) {
	email := strings.TrimSpace(c.PostForm("email"))
	otpInput := strings.TrimSpace(c.PostForm("otp"))
	passMoi := strings.TrimSpace(c.PostForm("pass_moi"))

	// [SỬA] TimKhachHang
	kh, ok := nghiep_vu.TimKhachHangTheoUserOrEmail(email)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Email không hợp lệ!"})
		return
	}

	if !nghiep_vu.KiemTraOTP(kh.TenDangNhap, otpInput) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã OTP sai hoặc đã hết hạn!"})
		return
	}

	hashMoi, _ := bao_mat.HashMatKhau(passMoi)
	kh.MatKhauHash = hashMoi
	// [SỬA] Ghi vào KHACH_HANG
	nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MatKhauHash, hashMoi)

	c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
}
