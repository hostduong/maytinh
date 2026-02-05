package chuc_nang

import (
	"fmt"
	"strings"

	"app/bao_mat"
	"app/cau_hinh"
	"app/mo_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

// API_DoiThongTin : Cập nhật Tên, SĐT, Ngày sinh, Giới tính
func API_DoiThongTin(c *gin.Context) {
	hoTenMoi    := strings.TrimSpace(c.PostForm("ho_ten"))
	sdtMoi      := strings.TrimSpace(c.PostForm("dien_thoai"))
	ngaySinhMoi := strings.TrimSpace(c.PostForm("ngay_sinh"))
	gioiTinhMoi := strings.TrimSpace(c.PostForm("gioi_tinh"))
	
	cookie, _ := c.Cookie("session_id")

	if !bao_mat.KiemTraHoTen(hoTenMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "Tên không hợp lệ!"})
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

		c.JSON(200, gin.H{"status": "ok", "msg": "Cập nhật thành công!"})
	} else {
		c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên đăng nhập"})
	}
}

// API_DoiMatKhau : Đổi pass bằng pass cũ
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
	}
}

// API_DoiMaPin : Đổi PIN bằng PIN cũ
func API_DoiMaPin(c *gin.Context) {
	pinCu := strings.TrimSpace(c.PostForm("pin_cu"))
	pinMoi := strings.TrimSpace(c.PostForm("pin_moi"))
	cookie, _ := c.Cookie("session_id")

	if !bao_mat.KiemTraMaPin(pinMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "PIN mới phải 8 số!"})
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
	}
}

// API_GuiOTPPin : Quên PIN - Tự cấp PIN mới và gửi Email
func API_GuiOTPPin(c *gin.Context) {
	cookie, _ := c.Cookie("session_id")
	kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie)
	if !ok {
		c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên làm việc"})
		return
	}

	theGui, msgLoi := nghiep_vu.KiemTraRateLimit(kh.Email)
	if !theGui {
		c.JSON(200, gin.H{"status": "error", "msg": msgLoi})
		return
	}

	newPin := nghiep_vu.TaoMaOTP() // Hàm cũ 8 số
	body := fmt.Sprintf("Xin chào,\n\nMã PIN mới cho tài khoản %s là: %s\n\nTrân trọng!", kh.TenDangNhap, newPin)

	err := nghiep_vu.GuiMailThongBaoAPI(kh.Email, "Thay đổi mã PIN", "Hỗ trợ", body)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": "Lỗi gửi mail: " + err.Error()})
		return
	}

	kh.MaPinHash = newPin
	nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MaPinHash, newPin)
	c.JSON(200, gin.H{"status": "ok", "msg": "Mã PIN mới đã được gửi vào Email!"})
}

// API_ResetPinBangOTP : Để trống do đã dùng chiến thuật gửi PIN trực tiếp
func API_ResetPinBangOTP(c *gin.Context) {
	c.JSON(200, gin.H{"status": "error", "msg": "Vui lòng dùng tính năng gửi PIN qua Email."})
}
