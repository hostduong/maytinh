package chuc_nang

import (
	"strings"

	"app/bao_mat"
	"app/cau_hinh"
	"app/mo_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

// API: Cập nhật Thông tin cá nhân (Họ tên, SĐT, Ngày sinh, Giới tính)
func API_DoiThongTin(c *gin.Context) {
	hoTenMoi    := strings.TrimSpace(c.PostForm("ho_ten"))
	sdtMoi      := strings.TrimSpace(c.PostForm("dien_thoai"))
	ngaySinhMoi := strings.TrimSpace(c.PostForm("ngay_sinh"))
	gioiTinhMoi := strings.TrimSpace(c.PostForm("gioi_tinh"))
	
	cookie, _ := c.Cookie("session_id")

	if !bao_mat.KiemTraHoTen(hoTenMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "Tên không hợp lệ (6-50 ký tự, không số)!"})
		return
	}

	if len(sdtMoi) > 0 && (len(sdtMoi) < 8 || len(sdtMoi) > 15) {
		c.JSON(200, gin.H{"status": "error", "msg": "Số điện thoại không hợp lệ!"})
		return
	}

	if kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
		kh.TenKhachHang = hoTenMoi
		kh.DienThoai    = sdtMoi
		kh.NgaySinh     = ngaySinhMoi
		kh.GioiTinh     = gioiTinhMoi

		sID := cau_hinh.BienCauHinh.IdFileSheet
		row := kh.DongTrongSheet
		sheet := "KHACH_HANG"

		nghiep_vu.ThemVaoHangCho(sID, sheet, row, mo_hinh.CotKH_TenKhachHang, hoTenMoi)
		nghiep_vu.ThemVaoHangCho(sID, sheet, row, mo_hinh.CotKH_DienThoai, sdtMoi)
		nghiep_vu.ThemVaoHangCho(sID, sheet, row, mo_hinh.CotKH_NgaySinh, ngaySinhMoi)
		nghiep_vu.ThemVaoHangCho(sID, sheet, row, mo_hinh.CotKH_GioiTinh, gioiTinhMoi)

		c.JSON(200, gin.H{"status": "ok", "msg": "Cập nhật hồ sơ thành công!"})
	} else {
		c.JSON(401, gin.H{"status": "error", "msg": "Phiên đăng nhập hết hạn"})
	}
}

// API: Đổi Mật Khẩu (Giữ nguyên logic cũ)
func API_DoiMatKhau(c *gin.Context) {
	passCu := strings.TrimSpace(c.PostForm("pass_cu"))
	passMoi := strings.TrimSpace(c.PostForm("pass_moi"))
	cookie, _ := c.Cookie("session_id")

	if !bao_mat.KiemTraDinhDangMatKhau(passMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu mới không hợp lệ!"})
		return
	}

	if kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
		if !bao_mat.KiemTraMatKhau(passCu, kh.MatKhauHash) {
			c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu cũ không đúng!"})
			return
		}
		hashMoi, _ := bao_mat.HashMatKhau(passMoi)
		kh.MatKhauHash = hashMoi
		nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MatKhauHash, hashMoi)
		c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
	} else {
		c.JSON(401, gin.H{"status": "error", "msg": "Phiên đăng nhập hết hạn"})
	}
}

// API: Đổi Mã PIN (Giữ nguyên logic cũ)
func API_DoiMaPin(c *gin.Context) {
	pinCu := strings.TrimSpace(c.PostForm("pin_cu"))
	pinMoi := strings.TrimSpace(c.PostForm("pin_moi"))
	cookie, _ := c.Cookie("session_id")

	if !bao_mat.KiemTraMaPin(pinMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN mới phải đúng 8 số!"})
		return
	}

	if kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
		if kh.MaPinHash != pinCu {
			c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN cũ không đúng!"})
			return
		}
		kh.MaPinHash = pinMoi
		nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MaPinHash, pinMoi)
		c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mã PIN thành công!"})
	} else {
		c.JSON(401, gin.H{"status": "error", "msg": "Phiên đăng nhập hết hạn"})
	}
}

// --- API BỔ SUNG CHO QUÊN MÃ PIN ---

func API_GuiOTPPin(c *gin.Context) {
	cookie, _ := c.Cookie("session_id")
	kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie)
	if !ok { c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên làm việc"}); return }

	theGui, msgLoi := nghiep_vu.KiemTraRateLimit(kh.Email)
	if !theGui { c.JSON(200, gin.H{"status": "error", "msg": msgLoi}); return }

	code := nghiep_vu.TaoMaOTP6So() // API mail cần 6 số
	if err := nghiep_vu.GuiMailXacMinhAPI(kh.Email, code); err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": "Lỗi API gửi thư: " + err.Error()}); return
	}

	nghiep_vu.LuuOTPVaUpdateRate(kh.Email, code)
	c.JSON(200, gin.H{"status": "ok", "msg": "Mã xác minh đã được gửi vào Email của bạn!"})
}

func API_ResetPinBangOTP(c *gin.Context) {
	otpInput := strings.TrimSpace(c.PostForm("otp"))
	pinMoi   := strings.TrimSpace(c.PostForm("pin_moi"))
	cookie, _ := c.Cookie("session_id")

	if !bao_mat.KiemTraMaPin(pinMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN mới phải đúng 8 số!"}); return
	}

	kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie)
	if !ok { c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên làm việc"}); return }

	if !nghiep_vu.KiemTraOTP(kh.Email, otpInput) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mã xác minh sai hoặc đã hết hạn!"}); return
	}

	kh.MaPinHash = pinMoi
	nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MaPinHash, pinMoi)
	c.JSON(200, gin.H{"status": "ok", "msg": "Đã khôi phục mã PIN thành công!"})
}
