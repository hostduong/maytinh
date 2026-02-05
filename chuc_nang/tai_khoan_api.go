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

// --- [HÀM CŨ 1] ---
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
		nghiep_vu.ThemVaoHangCho(sID, "KHACH_HANG", row, mo_hinh.CotKH_TenKhachHang, hoTenMoi)
		nghiep_vu.ThemVaoHangCho(sID, "KHACH_HANG", row, mo_hinh.CotKH_DienThoai, sdtMoi)
		nghiep_vu.ThemVaoHangCho(sID, "KHACH_HANG", row, mo_hinh.CotKH_NgaySinh, ngaySinhMoi)
		nghiep_vu.ThemVaoHangCho(sID, "KHACH_HANG", row, mo_hinh.CotKH_GioiTinh, gioiTinhMoi)

		c.JSON(200, gin.H{"status": "ok", "msg": "Cập nhật thành công!"})
	} else {
		c.JSON(401, gin.H{"status": "error", "msg": "Phiên đăng nhập hết hạn"})
	}
}

// --- [HÀM CŨ 2] ---
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

// --- [HÀM CŨ 3] ---
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

// --- [HÀM MỚI BỔ SUNG - QUÊN MÃ PIN] ---
func API_GuiOTPPin(c *gin.Context) {
	cookie, _ := c.Cookie("session_id")
	kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie)
	if !ok {
		c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên làm việc"})
		return
	}

	// 1. Kiểm tra Rate Limit
	theGui, msgLoi := nghiep_vu.KiemTraRateLimit(kh.Email)
	if !theGui {
		c.JSON(200, gin.H{"status": "error", "msg": msgLoi})
		return
	}

	// 2. Tạo mã PIN mới (8 số)
	newPin := nghiep_vu.TaoMaOTP()

	// 3. Soạn nội dung Email
	body := fmt.Sprintf(`Xin chào,

Chúng tôi đã tạo mã PIN mới cho tài khoản %s theo yêu cầu.

Mã PIN mới của bạn là: %s

Trân trọng,
MayTinhShop`, kh.TenDangNhap, newPin)

	// 4. Gửi Mail
	err := nghiep_vu.GuiMailThongBaoAPI(kh.Email, "Thông báo mã PIN mới", "Hỗ trợ tài khoản", body)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": "Lỗi gửi mail: " + err.Error()})
		return
	}

	// 5. Cập nhật vào DB
	kh.MaPinHash = newPin
	nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MaPinHash, newPin)

	c.JSON(200, gin.H{"status": "ok", "msg": "Mã PIN mới đã được gửi vào Email!"})
}

// Hàm placeholder cho API reset pin (không dùng nữa vì đã gửi thẳng PIN)
func API_ResetPinBangOTP(c *gin.Context) {
	c.JSON(200, gin.H{"status": "error", "msg": "Vui lòng dùng tính năng gửi PIN qua Email."})
}
