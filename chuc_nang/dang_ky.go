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

func TrangDangKy(c *gin.Context) {
	cookie, _ := c.Cookie("session_id")
	if cookie != "" {
		if _, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
			c.Redirect(http.StatusFound, "/")
			return
		}
	}
	c.HTML(http.StatusOK, "dang_ky", gin.H{})
}

func XuLyDangKy(c *gin.Context) {
	hoTen     := strings.TrimSpace(c.PostForm("ho_ten"))
	user      := strings.TrimSpace(c.PostForm("ten_dang_nhap"))
	pass      := strings.TrimSpace(c.PostForm("mat_khau"))
	email     := strings.TrimSpace(c.PostForm("email"))
	maPin     := strings.TrimSpace(c.PostForm("ma_pin"))
	
	dienThoai := strings.TrimSpace(c.PostForm("dien_thoai_full")) 
	if dienThoai == "" { dienThoai = strings.TrimSpace(c.PostForm("dien_thoai")) }
	ngaySinh  := strings.TrimSpace(c.PostForm("ngay_sinh"))
	gioiTinh  := strings.TrimSpace(c.PostForm("gioi_tinh"))

	// Validation
	if !bao_mat.KiemTraHoTen(hoTen) || !bao_mat.KiemTraTenDangNhap(user) || 
	   !bao_mat.KiemTraEmail(email) || !bao_mat.KiemTraMaPin(maPin) || 
	   !bao_mat.KiemTraDinhDangMatKhau(pass) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Dữ liệu nhập vào không hợp lệ!"})
		return
	}

	// [ĐÃ SỬA] Chỉ kiểm tra User và Email (Bỏ số điện thoại)
	if nghiep_vu.KiemTraTonTaiUserEmail(user, email) {
		c.HTML(http.StatusOK, "dang_ky", gin.H{"Loi": "Tên đăng nhập hoặc Email đã tồn tại!"})
		return
	}

	var maKH, vaiTro, loaiKH string
	if nghiep_vu.DemSoLuongKhachHang() == 0 {
		maKH = "KH_0001"
		vaiTro = "admin"
		loaiKH = "quan_tri_vien"
	} else {
		maKH = nghiep_vu.TaoMaKhachHangMoi()
		vaiTro = "" 
		loaiKH = "khach_le"
	}

	passHash, _ := bao_mat.HashMatKhau(pass)
	cookie := bao_mat.TaoSessionIDAnToan()
	expiredTime := time.Now().Add(cau_hinh.ThoiGianHetHanCookie).Unix()

	newKH := &mo_hinh.KhachHang{
		MaKhachHang:    maKH,
		TenDangNhap:    user,
		Email:          email,
		DienThoai:      dienThoai,
		MatKhauHash:    passHash,
		MaPinHash:      maPin, 
		TenKhachHang:   hoTen,
		NgaySinh:       ngaySinh,
		GioiTinh:       gioiTinh,
		LoaiKhachHang:  loaiKH,
		VaiTroQuyenHan: vaiTro,
		Cookie:         cookie,
		CookieExpired:  expiredTime,
		TrangThai:      1,
		NgayTao:        time.Now().Format("2006-01-02 15:04:05"),
	}

	nghiep_vu.ThemKhachHangMoi(newKH)
	c.SetCookie("session_id", cookie, int(cau_hinh.ThoiGianHetHanCookie.Seconds()), "/", "", false, true)

	if vaiTro == "admin" {
		c.Redirect(http.StatusFound, "/admin/tong-quan")
	} else {
		c.Redirect(http.StatusFound, "/")
	}
}
