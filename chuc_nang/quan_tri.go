package chuc_nang

import (
	"net/http"
	"app/nghiep_vu"
	"github.com/gin-gonic/gin"
)

// API_NapLaiDuLieu : Hút lại dữ liệu từ Google Sheet vào RAM
func API_NapLaiDuLieu(c *gin.Context) {
	// Gọi lại hàm khởi tạo bộ nhớ (Hàm này sẽ đọc lại toàn bộ 17 sheet)
	nghiep_vu.KhoiTaoBoNho()

	c.JSON(http.StatusOK, gin.H{
		"trang_thai": "thanh_cong",
		"thong_diep": "Đã nạp lại dữ liệu mới nhất từ Google Sheet!",
	})
}
