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

	// --- [QUAN TRỌNG] CẤU HÌNH GIAO DIỆN HTML ---
	// Lệnh này sẽ tìm tất cả file .html nằm trong thư mục giao_dien và các thư mục con
	// Nếu Cloud Build báo lỗi dòng này, hãy kiểm tra xem bạn đã tạo thư mục "giao_dien" trên GitHub chưa
	router.LoadHTMLGlob("giao_dien/**/*.html")

	// --- ĐỊNH NGHĨA CÁC ĐƯỜNG DẪN (ROUTES) ---
	
	// Nhóm API (Trả về JSON - Dùng cho App hoặc AJAX)
	api := router.Group("/api")
	{
		api.GET("/san-pham", chuc_nang.API_LayDanhSachSanPham)
		api.GET("/san-pham/:id", chuc_nang.API_ChiTietSanPham)
		api.GET("/cau-hinh", chuc_nang.API_LayMenu)
		api.GET("/admin/reload", chuc_nang.API_NapLaiDuLieu)
	}

	// Nhóm WEB (Trả về Giao diện HTML)
	// Khi vào trang chủ, gọi hàm TrangChu để hiển thị giao diện
	router.GET("/", chuc_nang.TrangChu)
    router.GET("/san-pham/:id", chuc_nang.ChiTietSanPham)
	// 5. Bấm nút CHẠY
	port := cau_hinh.BienCauHinh.CongChayWeb
	if port == "" {
		port = "8080"
	}
	log.Println(">> Server đang chạy tại cổng: " + port)
	
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Lỗi khởi động Server: %v", err)
	}
}
