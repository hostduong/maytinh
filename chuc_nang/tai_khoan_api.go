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

// API_DoiThongTin : (Giữ nguyên)
func API_DoiThongTin(c *gin.Context) {
	hoTenMoi    := strings.TrimSpace(c.PostForm("ho_ten"))
	sdtMoi      := strings.TrimSpace(c.PostForm("dien_thoai"))
	ngaySinhMoi := strings.TrimSpace(c.PostForm("ngay_sinh"))
	gioiTinhMoi := strings.TrimSpace(c.PostForm("gioi_tinh"))
	cookie, _ := c.Cookie("session_id")

	if !bao_mat.KiemTraHoTen(hoTenMoi) { c.JSON(200, gin.H{"status": "error", "msg": "Tên không hợp lệ!"}); return }
	if len(sdtMoi) > 0 && (len(sdtMoi) < 8 || len(sdtMoi) > 15) { c.JSON(200, gin.H{"status": "error", "msg": "SĐT lỗi!"}); return }

	if kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
		kh.TenKhachHang = hoTenMoi; kh.DienThoai = sdtMoi; kh.NgaySinh = ngaySinhMoi; kh.GioiTinh = gioiTinhMoi
		sID := cau_hinh.BienCauHinh.IdFileSheet; row := kh.DongTrongSheet
		nghiep_vu.ThemVaoHangCho(sID, "KHACH_HANG", row, mo_hinh.CotKH_TenKhachHang, hoTenMoi)
		nghiep_vu.ThemVaoHangCho(sID, "KHACH_HANG", row, mo_hinh.CotKH_DienThoai, sdtMoi)
		nghiep_vu.ThemVaoHangCho(sID, "KHACH_HANG", row, mo_hinh.CotKH_NgaySinh, ngaySinhMoi)
		nghiep_vu.ThemVaoHangCho(sID, "KHACH_HANG", row, mo_hinh.CotKH_GioiTinh, gioiTinhMoi)
		c.JSON(200, gin.H{"status": "ok", "msg": "Cập nhật thành công!"})
	} else { c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên"}) }
}

// API_DoiMatKhau : (Giữ nguyên)
func API_DoiMatKhau(c *gin.Context) {
	passCu := strings.TrimSpace(c.PostForm("pass_cu"))
	passMoi := strings.TrimSpace(c.PostForm("pass_moi"))
	cookie, _ := c.Cookie("session_id")
	if !bao_mat.KiemTraDinhDangMatKhau(passMoi) { c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu mới lỗi!"}); return }
	if kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
		if !bao_mat.KiemTraMatKhau(passCu, kh.MatKhauHash) { c.JSON(200, gin.H{"status": "error", "msg": "Mật khẩu cũ sai!"}); return }
		hash, _ := bao_mat.HashMatKhau(passMoi)
		kh.MatKhauHash = hash
		nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MatKhauHash, hash)
		c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mật khẩu thành công!"})
	} else { c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên"}) }
}

// [CẬP NHẬT CHUẨN]: Đổi mã PIN có băm (Hash)
func API_DoiMaPin(c *gin.Context) {
	pinCu := strings.TrimSpace(c.PostForm("pin_cu"))
	pinMoi := strings.TrimSpace(c.PostForm("pin_moi"))
	cookie, _ := c.Cookie("session_id")
	
	if !bao_mat.KiemTraMaPin(pinMoi) { c.JSON(200, gin.H{"status": "error", "msg": "PIN phải đủ 8 số!"}); return }
	
	if kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
		// Dùng KiemTraMatKhau để so sánh mã PIN cũ đã băm
		if !bao_mat.KiemTraMatKhau(pinCu, kh.MaPinHash) {
			c.JSON(200, gin.H{"status": "error", "msg": "Mã PIN hiện tại không đúng!"})
			return
		}
		
		// Băm mã PIN mới trước khi lưu
		hashMoi, _ := bao_mat.HashMatKhau(pinMoi)
		kh.MaPinHash = hashMoi
		nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MaPinHash, hashMoi)
		
		c.JSON(200, gin.H{"status": "ok", "msg": "Đổi mã PIN thành công!"})
	} else { c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên làm việc"}) }
}

// [CẬP NHẬT CHUẨN]: Gửi PIN mới qua Email có băm (Hash)
func API_GuiOTPPin(c *gin.Context) {
	cookie, _ := c.Cookie("session_id")
	kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie)
	if !ok { c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên làm việc"}); return }

	theGui, msg := nghiep_vu.KiemTraRateLimit(kh.Email)
	if !theGui { c.JSON(200, gin.H{"status": "error", "msg": msg}); return }

	// 1. Tạo PIN ngẫu nhiên 8 số (mã thô để gửi mail)
	newPinRaw := nghiep_vu.TaoMaOTP()

	// 2. Nội dung Email chuẩn
	body := fmt.Sprintf(`Xin chào,

Chúng tôi đã tạo mã PIN mới cho tài khoản %s theo yêu cầu của bạn trên hệ thống.

Mã PIN mới của bạn là: %s

Vì lý do bảo mật, vui lòng đổi mã PIN này ngay sau khi đăng nhập.

Nếu bạn không yêu cầu thay đổi mã PIN, bạn có thể bỏ qua email này.

Trân trọng,
Đội ngũ hỗ trợ`, kh.Email, newPinRaw)

	err := nghiep_vu.GuiMailThongBaoAPI(kh.Email, "Thông báo thay đổi mã PIN", "Hỗ trợ tài khoản", body)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": "Lỗi gửi mail: " + err.Error()})
		return
	}

	// 3. Băm mã PIN mới trước khi lưu vào RAM và Sheet
	hashNewPin, _ := bao_mat.HashMatKhau(newPinRaw)
	kh.MaPinHash = hashNewPin
	nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MaPinHash, hashNewPin)

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã gửi mã PIN mới vào Email!"})
}
