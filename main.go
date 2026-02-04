package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	
	"app/cau_hinh"
	"app/chuc_nang"
	"app/kho_du_lieu"
	"app/nghiep_vu" // Import gói nghiệp vụ

	"github.com/gin-gonic/gin"
)

func main() {
	// ... (Các bước khởi tạo cũ: 1, 2, 3 giữ nguyên) ...
	cau_hinh.KhoiTaoCauHinh()
	kho_du_lieu.KhoiTaoKetNoiGoogle()
	nghiep_vu.KhoiTaoBoNho()

	// --- [MỚI] KHỞI ĐỘNG WORKER GHI SHEET ---
	nghiep_vu.KhoiTaoWorkerGhiSheet()
    chuc_nang.KhoiTaoBoDemRateLimit() // Worker reset bộ đếm request
	// 4. Khởi tạo Web Server
	router := gin.Default()
	// ... (Các config router giữ nguyên) ...
	
func main() {
    // ... (Phần khởi tạo bên trên giữ nguyên) ...

    router := gin.Default()
    router.LoadHTMLGlob("giao_dien/**/*")

    // ------------------------------------------
    // 1. NHÓM PUBLIC (Ai cũng vào được)
    // ------------------------------------------
    router.GET("/", chuc_nang.TrangChu)
    router.GET("/san-pham/:id", chuc_nang.ChiTietSanPham)
    
    // Login & Logout
    router.GET("/login", chuc_nang.TrangDangNhap)
    router.POST("/login", chuc_nang.XuLyDangNhap)
    router.GET("/logout", chuc_nang.DangXuat)

    // Công cụ tạo Hash mật khẩu (Dùng tạm để lấy chuỗi Hash bỏ vào Sheet)
    router.GET("/tool/hash/:pass", func(c *gin.Context) {
        pass := c.Param("pass")
        hash, _ := bao_mat.HashMatKhau(pass)
        c.String(200, "Mật khẩu: %s\nHash: %s", pass, hash)
    })

    // ------------------------------------------
    // 2. NHÓM ADMIN (Phải đăng nhập mới vào được)
    // ------------------------------------------
    admin := router.Group("/admin")
    admin.Use(chuc_nang.KiemTraQuyenHan) // <--- CÀI NGƯỜI GÁC CỔNG Ở ĐÂY
    {
        // Sau này ta sẽ thêm các trang quản trị vào đây
        admin.GET("/tong-quan", func(c *gin.Context) {
            c.String(200, "Chào mừng Sếp! Đây là trang quản trị (Đã bảo mật).")
        })
        
        // Nút Reload dữ liệu (Chỉ admin mới được bấm)
        admin.GET("/reload", chuc_nang.API_NapLaiDuLieu)
    }

    // ... (Phần chạy server bên dưới giữ nguyên) ...
}
	
	// --- [MỚI] CHẠY SERVER VỚI CƠ CHẾ GRACEFUL SHUTDOWN ---
	port := cau_hinh.BienCauHinh.CongChayWeb
	if port == "" { port = "8080" }

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Chạy server trong 1 Goroutine riêng
	go func() {
		log.Println(">> Server đang chạy tại cổng: " + port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Lỗi khởi động: %s\n", err)
		}
	}()

	// --- LẮNG NGHE TÍN HIỆU TẮT (SIGTERM) ---
	quit := make(chan os.Signal, 1)
	// Cloud Run gửi SIGTERM khi tắt container
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	<-quit // Code sẽ dừng ở đây chờ tín hiệu...
	log.Println("⚠️  Đang tắt Server... Đang xả hàng đợi lần cuối!")

	// GỌI HÀM XẢ HÀNG KHẨN CẤP
	nghiep_vu.ThucHienGhiSheet(true)

	log.Println("✅ Server đã tắt an toàn.")
}
