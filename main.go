package main

import (
	"log"
	
	"app/cau_hinh"
	"app/chuc_nang"
	"app/kho_du_lieu"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

func main() {
	// 1. Nạp cấu hình (Biến môi trường)
	cau_hinh.KhoiTaoCauHinh()

	// 2. Kết nối Google Sheets
	kho_du_lieu.KhoiTaoKetNoiGoogle()

	// 3. Nạp dữ liệu từ Sheet vào RAM (Chạy 1 lần đầu)
	nghiep_vu.KhoiTaoBoNho()

	// 4. Khởi tạo Web Server (Gin)
	router := gin.Default()

	// Cấu hình Proxy (Quan trọng khi chạy trên Cloud Run)
	router.SetTrustedProxies(nil)

	// --- ĐỊNH NGHĨA CÁC ĐƯỜNG DẪN (ROUTES) ---
	
	// Nhóm API Public (Ai cũng xem được)
	api := router.Group("/api")
	{
		api.GET("/san-pham", chuc_nang.API_LayDanhSachSanPham)
		api.GET("/san-pham/:id", chuc_nang.API_ChiTietSanPham)
		api.GET("/cau-hinh", chuc_nang.API_LayMenu)
	}

	// Trang chủ (Tạm thời trả về text, sau sẽ gắn HTML)
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hệ thống bán hàng Golang + Google Sheet đang chạy!")
	})

	// 5. Bấm nút CHẠY
	port := cau_hinh.BienCauHinh.CongChayWeb
	log.Println(">> Server đang chạy tại: http://localhost:" + port)
	
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Lỗi khởi động Server: %v", err)
	}
}
