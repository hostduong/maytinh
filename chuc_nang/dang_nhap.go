package chuc_nang

import (
	"net/http"
	"strings"
	"time"

	"app/bao_mat"
	"app/cau_hinh" // Import thêm để dùng các hằng số cấu hình nếu cần
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
	// Ép về chữ thường để so sánh chính xác
	inputTaiKhoan := strings.ToLower(strings.TrimSpace(c.PostForm("ten_dang_nhap")))
	pass          := strings.TrimSpace(c.PostForm("mat_khau"))

	// 1. Tìm user
	kh, ok := nghiep_vu.TimKhachHangTheoUserOrEmail(inputTaiKhoan)
	if !ok {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Tài khoản không tồn tại!"})
		return
	}

	// 2. Kiểm tra mật khẩu
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
	// Thời gian hết hạn: Bây giờ + Thời gian cấu hình (30 phút)
	expTime := time.Now().Add(cau_hinh.ThoiGianHetHanCookie).Unix()

	// Cập nhật vào Struct trong RAM
	kh.Cookie = sessionID
	kh.CookieExpired = expTime

	// [SỬA LỖI TẠI ĐÂY] 
	// Gọi hàm cập nhật phiên chuẩn (Truyền con trỏ struct) thay vì truyền 3 tham số rời
	nghiep_vu.CapNhatPhienDangNhapKH(kh)

	// Set Cookie trình duyệt
	// MaxAge tính bằng giây
	maxAge := int(cau_hinh.ThoiGianHetHanCookie.Seconds())
	c.SetCookie("session_id", sessionID, maxAge, "/", "", false, true)

	c.Redirect(http.StatusFound, "/")
}

func DangXuat(c *gin.Context) {
	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}
