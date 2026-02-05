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

// [GIỮ NGUYÊN 3 HÀM CŨ: API_DoiThongTin, API_DoiMatKhau, API_DoiMaPin]
// (Bạn hãy giữ nguyên code cũ của 3 hàm này, không thay đổi gì)

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

func API_DoiMaPin(c *gin.Context) {
	pinCu := strings.TrimSpace(c.PostForm("pin_cu"))
	pinMoi := strings.TrimSpace(c.PostForm("pin_moi"))
	cookie, _ := c.Cookie("session_id")
	if !bao_mat.KiemTraMaPin(pinMoi) { c.JSON(200, gin.H{"status": "error", "msg": "PIN phải 8 số!"}); return }
	if kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
		if kh.MaPinHash != pinCu { c.JSON(200, gin.H{"status": "error", "msg": "PIN cũ sai!"}); return }
		kh.MaPinHash = pinMoi
		nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MaPinHash, pinMoi)
		c.JSON(200, gin.H{"status": "ok", "msg": "Đổi PIN thành công!"})
	} else { c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên"}) }
}

// --- [HÀM BỔ SUNG - GỬI PIN MỚI] ---
func API_GuiOTPPin(c *gin.Context) {
	cookie, _ := c.Cookie("session_id")
	kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie)
	if !ok { c.JSON(401, gin.H{"status": "error", "msg": "Hết phiên làm việc"}); return }

	// 1. Rate Limit
	theGui, msg := nghiep_vu.KiemTraRateLimit(kh.Email)
	if !theGui { c.JSON(200, gin.H{"status": "error", "msg": msg}); return }

	// 2. Tạo PIN mới 8 số (Dùng hàm cũ TaoMaOTP 8 số)
	newPin := nghiep_vu.TaoMaOTP()

	// 3. Gửi Mail
	body := fmt.Sprintf("Chào %s,\n\nMã PIN mới của bạn là: %s\n\nTrân trọng!", kh.TenDangNhap, newPin)
	err := nghiep_vu.GuiMailThongBaoAPI(kh.Email, "Cấp lại Mã PIN", "Hỗ trợ", body)
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "msg": "Lỗi gửi mail: " + err.Error()})
		return
	}

	// 4. Cập nhật
	kh.MaPinHash = newPin
	nghiep_vu.ThemVaoHangCho(cau_hinh.BienCauHinh.IdFileSheet, "KHACH_HANG", kh.DongTrongSheet, mo_hinh.CotKH_MaPinHash, newPin)

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã gửi mã PIN mới vào Email!"})
}
