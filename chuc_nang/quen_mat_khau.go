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

func TrangQuenMatKhau(c *gin.Context) { c.HTML(http.StatusOK, "quen_mat_khau", gin.H{}) }

func XuLyQuenPassBangPIN(c *gin.Context) {
	user := strings.ToLower(strings.TrimSpace(c.PostForm("user")))
	pin := strings.TrimSpace(c.PostForm("pin"))
	passMoi := strings.TrimSpace(c.PostForm("pass_moi"))
	kh, ok := nghiep_vu.TimKhachHangTheoUserOrEmail(user)
	if !ok || kh.MaPinHash != pin { c.JSON(200, gin.H{"status": "error", "msg": "Thông tin sai!"}); return }
	
	hash, _ := bao_mat.HashMatKhau(passMoi)
	kh.MatKhauHash = hash
	nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MatKhauHash, hash)
	c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
}

func XuLyGuiOTPEmail(c *gin.Context) {
	email := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	kh, ok := nghiep_vu.TimKhachHangTheoUserOrEmail(email)
	if !ok { c.JSON(200, gin.H{"status": "ok", "msg": "Đã gửi mã (nếu email đúng)!"}); return }

	okLimit, msg := nghiep_vu.KiemTraRateLimit(kh.Email)
	if !okLimit { c.JSON(200, gin.H{"status": "error", "msg": msg}); return }

	code := nghiep_vu.TaoMaOTP6So()
	if err := nghiep_vu.GuiMailXacMinhAPI(kh.Email, code); err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": "Lỗi gửi mail!"}); return
	}
	nghiep_vu.LuuOTP(kh.TenDangNhap, code)
	c.JSON(200, gin.H{"status": "ok", "msg": "Đã gửi OTP 6 số vào Email!"})
}

func XuLyQuenPassBangOTP(c *gin.Context) {
	email := strings.ToLower(strings.TrimSpace(c.PostForm("email")))
	otp := strings.TrimSpace(c.PostForm("otp"))
	passMoi := strings.TrimSpace(c.PostForm("pass_moi"))
	kh, ok := nghiep_vu.TimKhachHangTheoUserOrEmail(email)
	if !ok || !nghiep_vu.KiemTraOTP(kh.TenDangNhap, otp) { c.JSON(200, gin.H{"status": "error", "msg": "Mã OTP sai!"}); return }

	hash, _ := bao_mat.HashMatKhau(passMoi)
	kh.MatKhauHash = hash
	nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MatKhauHash, hash)
	c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
}
