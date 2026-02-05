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

// API: Cập nhật Thông tin cá nhân (Tên, SĐT, Ngày sinh, Giới tính)
func API_DoiThongTin(c *gin.Context) {
	hoTenMoi    := strings.TrimSpace(c.PostForm("ho_ten"))
	sdtMoi      := strings.TrimSpace(c.PostForm("dien_thoai"))
	ngaySinhMoi := strings.TrimSpace(c.PostForm("ngay_sinh"))
	gioiTinhMoi := strings.TrimSpace(c.PostForm("gioi_tinh"))
	
	cookie, _ := c.Cookie("session_id")

	// 1. Validate Họ Tên theo quy tắc Whitelist
	if !bao_mat.KiemTraHoTen(hoTenMoi) {
		c.JSON(200, gin.H{"status": "error", "msg": "Tên không hợp lệ (6-50 ký tự, không số)!"})
		return
	}

	// 2. Validate SĐT (8-15 ký tự)
	if len(sdtMoi) > 0 && (len(sdtMoi) < 8 || len(sdtMoi) > 15) {
		c.JSON(200, gin.H{"status": "error", "msg": "Số điện thoại không hợp lệ!"})
		return
	}

	if kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
		// Cập nhật RAM
		kh.TenKhachHang = hoTenMoi
		kh.DienThoai    = sdtMoi
		kh.NgaySinh     = ngaySinhMoi
		kh.GioiTinh     = gioiTinhMoi

		// Cập nhật Google Sheet
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

// API: Đổi Mật Khẩu (Xác thực bằng pass cũ)
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

// API: Đổi Mã PIN (Xác thực bằng PIN cũ)
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

// API: Quên mã PIN (Chiến thuật mới: Tạo PIN mới và gửi thẳng qua Email)
func API_GuiOTPPin(c *gin.Context) {
	cookie, _ := c.Cookie("session_id")
	kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie)
	if !ok {
		c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên làm việc"})
		return
	}

	// 1. Kiểm tra giới hạn gửi (Rate Limit 1 phút/lần, 10 lần/6h)
	theGui, msgLoi := nghiep_vu.KiemTraRateLimit(kh.Email)
	if !theGui {
		c.JSON(200, gin.H{"status": "error", "msg": msgLoi})
		return
	}

	// 2. Tạo mã PIN mới ngẫu nhiên (8 số)
	newPin := nghiep_vu.TaoMaOTP() 

	// 3. Soạn nội dung Email thông báo
	body := fmt.Sprintf(`Xin chào,

Chúng tôi đã tạo mã PIN mới cho tài khoản %s theo yêu cầu của bạn trên hệ thống.

Mã PIN mới của bạn là: %s

Vì lý do bảo mật, vui lòng đổi mã PIN này ngay sau khi đăng nhập.

Nếu bạn không yêu cầu thay đổi mã PIN, bạn có thể bỏ qua email này.

Trân trọng,
Đội ngũ hỗ trợ`, kh.TenDangNhap, newPin)

	// 4. Gọi API gửi thư (type: sender)
	err := nghiep_vu.GuiMailThongBaoAPI(kh.Email, "Thông báo thay đổi mã PIN", "Hỗ trợ tài khoản", body)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": "Lỗi gửi thư: " + err.Error()})
		return
	}

	// 5. Cập nhật mã PIN mới vào Database (RAM & Sheet)
	kh.MaPinHash = newPin
	nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MaPinHash, newPin)

	c.JSON(200, gin.H{"status": "ok", "msg": "Mã PIN mới đã được gửi vào Email của bạn!"})
}

// API: Reset Pin bằng OTP (Hàm này hiện tại để trống vì đã dùng chiến thuật gửi PIN trực tiếp)
func API_ResetPinBangOTP(c *gin.Context) {
	c.JSON(200, gin.H{"status": "error", "msg": "Chức năng này đã được thay thế bằng gửi PIN trực tiếp."})
}
