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

// ChiTietSanPham : Hiển thị trang chi tiết
func ChiTietSanPham(c *gin.Context) {
	// 1. Lấy ID từ đường dẫn (VD: /san-pham/SP001 -> id = SP001)
	id := c.Param("id")

	// 2. Tìm trong RAM
	sp, tonTai := nghiep_vu.LayChiTietSanPham(id)

	if !tonTai {
		// Nếu không thấy thì báo lỗi 404 (Sau này làm trang 404 đẹp sau)
		c.String(http.StatusNotFound, "Không tìm thấy sản phẩm này!")
		return
	}

	// 3. Trả về HTML
	c.HTML(http.StatusOK, "khung_giao_dien", gin.H{
		"TieuDe":  sp.TenSanPham,
		"SanPham": sp,
	})
}
