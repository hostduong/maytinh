package chuc_nang

import (
	"net/http"
	"strings"
	"time"

	"app/bao_mat"
	"app/cau_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

func TrangDangNhap(c *gin.Context) {
	cookie, _ := c.Cookie("session_id")
	if cookie != "" {
		if _, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
			c.Redirect(http.StatusFound, "/") 
			return
		}
	}
	c.HTML(http.StatusOK, "dang_nhap", gin.H{})
}

func XuLyDangNhap(c *gin.Context) {
	// [MỚI] Nhận input đa năng (Mã KH / User / Email)
	inputDinhDanh := strings.ToLower(strings.TrimSpace(c.PostForm("input_dinh_danh")))
	pass          := strings.TrimSpace(c.PostForm("mat_khau"))

	// 1. Tìm user (Hàm này đã được update ở bo_nho_dem.go để tìm cả Mã KH)
	kh, ok := nghiep_vu.TimKhachHangTheoUserOrEmail(inputDinhDanh)
	if !ok {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Tài khoản không tồn tại!"})
		return
	}

	// 2. Kiểm tra mật khẩu (So sánh hash)
	if !bao_mat.KiemTraMatKhau(pass, kh.MatKhauHash) {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Mật khẩu không đúng!"})
		return
	}

	// 3. Kiểm tra trạng thái
	if kh.TrangThai == 0 {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Tài khoản đã bị khóa vĩnh viễn!"})
		return
	}
	if kh.TrangThai == 2 {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Tài khoản đang bị tạm khóa!"})
		return
	}

	// 4. Tạo Session & Cookie mới
	sessionID := bao_mat.TaoSessionIDAnToan()
	expTime := time.Now().Add(cau_hinh.ThoiGianHetHanCookie).Unix()

	// Cập nhật vào Struct trong RAM
	kh.Cookie = sessionID
	kh.CookieExpired = expTime

	// Gọi hàm cập nhật phiên chuẩn
	nghiep_vu.CapNhatPhienDangNhapKH(kh)

	// Set Cookie trình duyệt
	maxAge := int(cau_hinh.ThoiGianHetHanCookie.Seconds())
	c.SetCookie("session_id", sessionID, maxAge, "/", "", false, true)

	c.Redirect(http.StatusFound, "/")
}

func DangXuat(c *gin.Context) {
	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}
