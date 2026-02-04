package chuc_nang

import (
	"net/http"
	"strings"
	"time"

	"app/bao_mat"
	"app/cau_hinh"
	"app/mo_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

// GET /register
func TrangDangKy(c *gin.Context) {
	cookie, _ := c.Cookie("session_id")
	if cookie != "" {
		if _, ok := nghiep_vu.TimNhanVienTheoCookie(cookie); ok {
			c.Redirect(http.StatusFound, "/")
			return
		}
	}
	c.HTML(http.StatusOK, "dang_ky", gin.H{})
}

// POST /register
func XuLyDangKy(c *gin.Context) {
	hoTen := strings.TrimSpace(c.PostForm("ho_ten"))
	user  := strings.TrimSpace(c.PostForm("ten_dang_nhap"))
	pass  := strings.TrimSpace(c.PostForm("mat_khau"))
	email := strings.TrimSpace(c.PostForm("email"))
	maPin := strings.TrimSpace(c.PostForm("ma_pin"))

	// VALIDATION
	if !bao_mat.KiemTraHoTen(hoTen) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Họ tên phải từ 6-50 ký tự, chỉ chứa chữ cái!"})
		return
	}
	if !bao_mat.KiemTraTenDangNhap(user) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Tên đăng nhập 6-30 ký tự (chữ, số, . _), không dấu cách/Việt!"})
		return
	}
	if !bao_mat.KiemTraEmail(email) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Email không hợp lệ hoặc chứa ký tự lạ!"})
		return
	}
	if !bao_mat.KiemTraMaPin(maPin) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Mã PIN phải là 8 chữ số!"})
		return
	}
	
	// [ĐÃ SỬA] Gọi hàm KiemTraDinhDangMatKhau
	if !bao_mat.KiemTraDinhDangMatKhau(pass) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Mật khẩu 8-30 ký tự, không chứa dấu cách, ' \" < >"})
		return
	}

	if nghiep_vu.KiemTraTonTaiUserOrEmail(user, email) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Tên đăng nhập hoặc Email đã tồn tại!"})
		return
	}

	passHash, _ := bao_mat.HashMatKhau(pass)

	var maDinhDanh, quyenHan, chucVu string
	if nghiep_vu.DemSoLuongNhanVien() == 0 {
		maDinhDanh = nghiep_vu.TaoMaNhanVienMoi()
		quyenHan = "admin"
		chucVu = "Quản lý cửa hàng"
	} else {
		maDinhDanh = nghiep_vu.TaoMaKhachHangMoi()
		quyenHan = ""
		chucVu = "Khách hàng"
	}

	cookie := bao_mat.TaoSessionIDAnToan()
	expiredTime := time.Now().Add(cau_hinh.ThoiGianHetHanCookie).Unix()

	newNV := &mo_hinh.NhanVien{
		MaNhanVien:      maDinhDanh,
		TenDangNhap:     user,
		Email:           email,
		MatKhauHash:     passHash,
		HoTen:           hoTen,
		ChucVu:          chucVu,
		MaPin:           maPin,
		Cookie:          cookie,
		CookieExpired:   expiredTime,
		VaiTroQuyenHan:  quyenHan,
		TrangThai:       1,
		LanDangNhapCuoi: time.Now().Format("2006-01-02 15:04:05"),
	}

	nghiep_vu.ThemNhanVienMoi(newNV)
	c.SetCookie("session_id", cookie, int(cau_hinh.ThoiGianHetHanCookie.Seconds()), "/", "", false, true)

	if quyenHan == "admin" {
		c.Redirect(http.StatusFound, "/admin/tong-quan")
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}
