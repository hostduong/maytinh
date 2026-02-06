package chuc_nang

import (
	"net/http"
	"app/nghiep_vu"
	"github.com/gin-gonic/gin"
)

// API_NapLaiDuLieu : Hút lại dữ liệu từ Google Sheet vào RAM
// Yêu cầu quyền: system.reload (Chỉ Admin)
func API_NapLaiDuLieu(c *gin.Context) {
	// 1. Lấy vai trò hiện tại (Đã được Middleware set vào Context)
	vaiTro := c.GetString("USER_ROLE")

	// 2. [CHỐT CHẶN AN NINH]
	// Bạn nhớ thêm dòng 'system.reload' vào Sheet PHAN_QUYEN (Cột Admin=1) nhé
	if !nghiep_vu.KiemTraQuyen(vaiTro, "system.reload") {
		c.JSON(http.StatusForbidden, gin.H{
			"trang_thai": "loi",
			"thong_diep": "Bạn không có quyền nạp lại dữ liệu hệ thống!",
		})
		return
	}

	// 3. Logic xử lý (Nếu qua cửa)
	nghiep_vu.KhoiTaoBoNho()

	c.JSON(http.StatusOK, gin.H{
		"trang_thai": "thanh_cong",
		"thong_diep": "Đã nạp lại dữ liệu mới nhất từ Google Sheet!",
	})
}
