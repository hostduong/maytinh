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

// CÁCH 1: Xác thực bằng PIN (Giữ nguyên)
func XuLyQuenPassBangPIN(c *gin.Context) {
	user := strings.TrimSpace(c.PostForm("user"))
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

	hashMoi, _ := bao_mat.HashMatKhau(passMoi)
	kh.MatKhauHash = hashMoi
	nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MatKhauHash, hashMoi)

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã đặt lại mật khẩu thành công! Hãy đăng nhập."})
}

// CÁCH 2: Gửi OTP qua Email (CẬP NHẬT DÙNG API THẬT)
func XuLyGuiOTPEmail(c *gin.Context) {
	email := strings.TrimSpace(c.PostForm("email"))
	kh, ok := nghiep_vu.TimKhachHangTheoUserOrEmail(email)
	
	if !ok {
		c.JSON(200, gin.H{"status": "ok", "msg": "Nếu email đúng, mã OTP đã được gửi đi!"})
		return
	}

	// 1. Chặn Spam
	okLimit, msg := nghiep_vu.KiemTraRateLimit(kh.Email)
	if !okLimit {
		c.JSON(200, gin.H{"status": "error", "msg": msg})
		return
	}

	// 2. Tạo mã OTP 6 số (Dùng hàm mới 6 số)
	otp := nghiep_vu.TaoMaOTP6So()

	// 3. Gửi Mail thật
	err := nghiep_vu.GuiMailXacMinhAPI(kh.Email, otp)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": "Lỗi gửi mail: " + err.Error()})
		return
	}

	// 4. Lưu vào RAM (Dùng tên đăng nhập làm key)
	nghiep_vu.LuuOTP(kh.TenDangNhap, otp)

	c.JSON(200, gin.H{"status": "ok", "msg": "Mã OTP (6 số) đã được gửi vào email!"})
}

// Xác nhận OTP (Giữ nguyên)
func XuLyQuenPassBangOTP(c *gin.Context) {
	email := strings.TrimSpace(c.PostForm("email"))
	otpInput := strings.TrimSpace(c.PostForm("otp"))
	passMoi := strings.TrimSpace(c.PostForm("pass_moi"))

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
	nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MatKhauHash, hashMoi)

	c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
}
