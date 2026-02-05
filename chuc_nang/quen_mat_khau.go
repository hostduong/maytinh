package chuc_nang

import (
	"net/http"
	"strings"

	"app/bao_mat"
	"app/cau_hinh"
	"app/mo_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

func TrangQuenMatKhau(c *gin.Context) {
	c.HTML(http.StatusOK, "quen_mat_khau", gin.H{})
}

// CÁCH 1: Xác thực bằng mã PIN (Dành cho Quên mật khẩu dùng PIN)
func XuLyQuenPassBangPIN(c *gin.Context) {
	user := strings.ToLower(strings.TrimSpace(c.PostForm("user")))
	pin := strings.TrimSpace(c.PostForm("pin"))
	passMoi := strings.TrimSpace(c.PostForm("pass_moi"))

	kh, ok := nghiep_vu.TimKhachHangTheoUserOrEmail(user)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Tài khoản không tồn tại!"})
		return
	}

	if kh.MaPinHash != pin {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN không chính xác!"})
		return
	}

	if !bao_mat.KiemTraDinhDangMatKhau(passMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu mới không hợp lệ!"})
		return
	}

	hashMoi, _ := bao_mat.HashMatKhau(passMoi)
	kh.MatKhauHash = hashMoi
	nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MatKhauHash, hashMoi)

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã đặt lại mật khẩu thành công!"})
}

// CÁCH 2: Gửi mã OTP qua Email (Dành cho Quên mật khẩu dùng Email)
func XuLyGuiOTPEmail(c *gin.Context) {
	emailInput := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	kh, ok := nghiep_vu.TimKhachHangTheoUserOrEmail(emailInput)
	
	if !ok {
		// Trả về thành công ảo để bảo mật
		c.JSON(200, gin.H{"status": "ok", "msg": "Nếu Email chính xác, mã xác minh đã được gửi!"})
		return
	}

	// 1. Kiểm tra Rate Limit
	theGui, msgLoi := nghiep_vu.KiemTraRateLimit(kh.Email)
	if !theGui {
		c.JSON(200, gin.H{"status": "error", "msg": msgLoi})
		return
	}

	// 2. Tạo mã 6 số và gửi API mail
	code := nghiep_vu.TaoMaOTP6So()
	err := nghiep_vu.GuiMailXacMinhAPI(kh.Email, code)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": "Lỗi gửi thư: " + err.Error()})
		return
	}

	// 3. Lưu vào RAM dùng Username làm key (để đồng bộ KiemTraOTP)
	nghiep_vu.LuuOTP(kh.TenDangNhap, code)

	c.JSON(200, gin.H{"status": "ok", "msg": "Mã xác minh 6 số đã được gửi vào Email của bạn!"})
}

// Bước cuối: Xác nhận OTP và đặt Pass mới
func XuLyQuenPassBangOTP(c *gin.Context) {
	email := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	otpInput := strings.TrimSpace(c.PostForm("otp"))
	passMoi := strings.TrimSpace(c.PostForm("pass_moi"))

	kh, ok := nghiep_vu.TimKhachHangTheoUserOrEmail(email)
	if !ok {
		c.JSON(200, gin.H{"status": "error", "msg": "Email không hợp lệ!"})
		return
	}

	// Xác thực bằng Username làm key
	if !nghiep_vu.KiemTraOTP(kh.TenDangNhap, otpInput) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã xác minh sai hoặc đã hết hạn!"})
		return
	}

	if !bao_mat.KiemTraDinhDangMatKhau(passMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu mới không hợp lệ!"})
		return
	}

	hashMoi, _ := bao_mat.HashMatKhau(passMoi)
	kh.MatKhauHash = hashMoi
	nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MatKhauHash, hashMoi)

	c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công! Hãy đăng nhập lại."})
}
