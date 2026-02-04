package chuc_nang

import (
	"net/http"
	"time"

	"app/bao_mat"
	"app/cau_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Hiển thị trang đăng nhập HTML
func TrangDangNhap(c *gin.Context) {
	// Nếu đã đăng nhập rồi -> Chuyển vào trong luôn
	cookie, _ := c.Cookie("session_id")
	if cookie != "" {
		if _, ok := nghiep_vu.TimNhanVienTheoCookie(cookie); ok {
			c.Redirect(http.StatusFound, "/admin/tong-quan")
			return
		}
	}
	
	c.HTML(http.StatusOK, "dang_nhap", gin.H{})
}

// Xử lý khi bấm nút "Đăng Nhập"
func XuLyDangNhap(c *gin.Context) {
	user := c.PostForm("ten_dang_nhap")
	pass := c.PostForm("mat_khau")

	// 1. Kiểm tra User có tồn tại không
	nv, ok := nghiep_vu.TimNhanVienTheoUsername(user)
	if !ok {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Tài khoản không tồn tại!"})
		return
	}

	// 2. Kiểm tra Mật khẩu (So sánh Hash)
	if !bao_mat.KiemTraMatKhau(pass, nv.MatKhauHash) {
		c.HTML(http.StatusOK, "dang_nhap", gin.H{"Loi": "Sai mật khẩu!"})
		return
	}

	// 3. Đăng nhập thành công -> Tạo Session Mới
	sessionID := uuid.New().String()
	// Thời gian hết hạn = Hiện tại + 30 phút (Config)
	expiredTime := time.Now().Add(cau_hinh.ThoiGianHetHanCookie).Unix()

	// 4. Lưu vào RAM & Hàng chờ ghi Sheet
	nghiep_vu.CapNhatPhienDangNhap(nv.MaNhanVien, sessionID, expiredTime)

	// 5. Trả Cookie về trình duyệt
	// MaxAge tính bằng giây
	c.SetCookie("session_id", sessionID, int(cau_hinh.ThoiGianHetHanCookie.Seconds()), "/", "", false, true)

	// 6. Chuyển hướng vào trang quản trị
	c.Redirect(http.StatusFound, "/admin/tong-quan")
}

// Xử lý Đăng Xuất
func DangXuat(c *gin.Context) {
	// Xóa cookie trình duyệt
	c.SetCookie("session_id", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/login")
}
