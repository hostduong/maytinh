package chuc_nang

import (
	"strings"

	"app/bao_mat"
	"app/mo_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

// API: Cập nhật Họ Tên
func API_DoiThongTin(c *gin.Context) {
	hoTenMoi := strings.TrimSpace(c.PostForm("ho_ten_moi"))
	cookie, _ := c.Cookie("session_id")

	if nv, ok := nghiep_vu.TimNhanVienTheoCookie(cookie); ok {
		// Update RAM
		nv.HoTen = hoTenMoi
		
		// Update Sheet
		nghiep_vu.ThemVaoHangCho(nghiep_vu.CacheNhanVien.SpreadsheetID, "NHAN_VIEN", nv.DongTrongSheet, mo_hinh.CotNV_HoTen, hoTenMoi)
		
		c.JSON(200, gin.H{"status": "ok", "msg": "Cập nhật tên thành công!"})
	} else {
		c.JSON(401, gin.H{"status": "error", "msg": "Phiên đăng nhập hết hạn"})
	}
}

// API: Đổi Mật Khẩu
func API_DoiMatKhau(c *gin.Context) {
	passCu := c.PostForm("pass_cu")
	passMoi := c.PostForm("pass_moi")
	cookie, _ := c.Cookie("session_id")

	if nv, ok := nghiep_vu.TimNhanVienTheoCookie(cookie); ok {
		// 1. Check Pass Cũ
		if !bao_mat.KiemTraMatKhau(passCu, nv.MatKhauHash) {
			c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu cũ không đúng!"})
			return
		}
		// 2. Hash Pass Mới & Lưu
		hashMoi, _ := bao_mat.HashMatKhau(passMoi)
		nv.MatKhauHash = hashMoi
		nghiep_vu.ThemVaoHangCho(nghiep_vu.CacheNhanVien.SpreadsheetID, "NHAN_VIEN", nv.DongTrongSheet, mo_hinh.CotNV_MatKhauHash, hashMoi)
		
		c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
	} else {
		c.JSON(401, gin.H{"status": "error", "msg": "Phiên đăng nhập hết hạn"})
	}
}

// API: Đổi Mã PIN
func API_DoiMaPin(c *gin.Context) {
	pinCu := c.PostForm("pin_cu")
	pinMoi := c.PostForm("pin_moi")
	cookie, _ := c.Cookie("session_id")

	if nv, ok := nghiep_vu.TimNhanVienTheoCookie(cookie); ok {
		// 1. Check PIN Cũ
		if nv.MaPin != pinCu {
			c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN cũ không đúng!"})
			return
		}
		// 2. Lưu PIN Mới
		nv.MaPin = pinMoi
		nghiep_vu.ThemVaoHangCho(nghiep_vu.CacheNhanVien.SpreadsheetID, "NHAN_VIEN", nv.DongTrongSheet, mo_hinh.CotNV_MaPin, pinMoi)
		
		c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mã PIN thành công!"})
	} else {
		c.JSON(401, gin.H{"status": "error", "msg": "Phiên đăng nhập hết hạn"})
	}
}
