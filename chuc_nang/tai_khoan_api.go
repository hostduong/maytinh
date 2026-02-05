package chuc_nang

import (
	"strings"

	"app/bao_mat"
	"app/cau_hinh"
	"app/mo_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

// API: Cập nhật Thông tin cá nhân (Tên, SĐT, Ngày sinh, Giới tính)
func API_DoiThongTin(c *gin.Context) {
	// Lấy dữ liệu từ Form
	hoTenMoi    := strings.TrimSpace(c.PostForm("ho_ten"))
	sdtMoi      := strings.TrimSpace(c.PostForm("dien_thoai"))
	ngaySinhMoi := strings.TrimSpace(c.PostForm("ngay_sinh"))
	gioiTinhMoi := strings.TrimSpace(c.PostForm("gioi_tinh"))
	
	cookie, _ := c.Cookie("session_id")

	// 1. Validate Họ Tên
	if !bao_mat.KiemTraHoTen(hoTenMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "Tên không hợp lệ (6-50 ký tự, không số)!"})
		return
	}

	// 2. Validate SĐT (Cơ bản: chỉ số, 8-15 ký tự)
	// Lưu ý: Frontend đã check kỹ bằng lib, Backend check chốt chặn
	if len(sdtMoi) < 8 || len(sdtMoi) > 15 {
		c.JSON(200, gin.H{"status": "error", "msg": "Số điện thoại không hợp lệ!"})
		return
	}

	if kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
		// Cập nhật RAM
		kh.TenKhachHang = hoTenMoi
		kh.DienThoai    = sdtMoi
		kh.NgaySinh     = ngaySinhMoi
		kh.GioiTinh     = gioiTinhMoi

		// Cập nhật Google Sheet (Ghi 4 cột)
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

// API: Đổi Mật Khẩu (Giữ nguyên)
func API_DoiMatKhau(c *gin.Context) {
	passCu := strings.TrimSpace(c.PostForm("pass_cu"))
	passMoi := strings.TrimSpace(c.PostForm("pass_moi"))
	cookie, _ := c.Cookie("session_id")

	if !bao_mat.KiemTraDinhDangMatKhau(passMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu mới không an toàn!"})
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

// API: Đổi Mã PIN (Giữ nguyên)
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
