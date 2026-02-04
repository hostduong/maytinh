package chuc_nang

import (
	"net/http"
	"app/nghiep_vu" // Nhớ sửa thành tên module của bạn nếu khác "app"
	"github.com/gin-gonic/gin"
)

// TrangChu : Hiển thị trang chủ HTML
func TrangChu(c *gin.Context) {
	// 1. Lấy dữ liệu từ RAM
	danhSachSP := nghiep_vu.LayDanhSachSanPham()
	// (Tạm thời lấy hết, sau này sẽ có logic lấy SP mới nhất/bán chạy)

	// 2. Trả về HTML
	c.HTML(http.StatusOK, "khung_giao_dien", gin.H{
		"TieuDe":          "Trang Chủ",
		"DanhSachSanPham": danhSachSP,
	})
}
