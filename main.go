package main

import (
	"log"
	"net/http" // Quan trọng: Đã thêm thư viện này
	"os"
	"os/signal"
	"syscall"

	"app/bao_mat" // Quan trọng: Để dùng tool hash pass
	"app/cau_hinh"
	"app/chuc_nang"
	"app/kho_du_lieu"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

func main() {
	// ---------------------------------------------------------
	// 1. KHỞI TẠO HỆ THỐNG (THEO THỨ TỰ)
	// ---------------------------------------------------------
	log.Println(">>> ĐANG KHỞI ĐỘNG HỆ THỐNG MAYTINHSHOP...")

	// B1: Nạp cấu hình (ID Sheet, Port...)
	cau_hinh.KhoiTaoCauHinh()

	// B2: Kết nối Google Sheet API
	kho_du_lieu.KhoiTaoKetNoiGoogle()

	// B3: Nạp dữ liệu từ Sheet vào RAM (Cache 17 bảng)
	nghiep_vu.KhoiTaoBoNho()

	// B4: Khởi động các Worker chạy ngầm
	nghiep_vu.KhoiTaoWorkerGhiSheet()   // Worker ghi dữ liệu xuống Sheet (10s/lần)
	chuc_nang.KhoiTaoBoDemRateLimit()   // Worker reset bộ đếm chống Spam (1s/lần)

	// ---------------------------------------------------------
	// 2. KHỞI TẠO WEB SERVER (GIN)
	// ---------------------------------------------------------
	router := gin.Default()
	
	// Nạp giao diện HTML (Lưu ý: Dockerfile phải copy thư mục giao_dien)
	router.LoadHTMLGlob("giao_dien/**/*")
	
	// Phục vụ file tĩnh (CSS, JS, Ảnh) nếu có thư mục static
	// router.Static("/static", "./static")

	// ---------------------------------------------------------
	// 3. ĐỊNH TUYẾN (ROUTER)
	// ---------------------------------------------------------

	// --- NHÓM PUBLIC (AI CŨNG VÀO ĐƯỢC) ---
	router.GET("/", chuc_nang.TrangChu)
	router.GET("/san-pham/:id", chuc_nang.ChiTietSanPham)

	// Login & Logout
	router.GET("/login", chuc_nang.TrangDangNhap)
	router.POST("/login", chuc_nang.XuLyDangNhap)
	router.GET("/logout", chuc_nang.DangXuat)

	// --- CÔNG CỤ TIỆN ÍCH (DEV ONLY) ---
	// Dùng để tạo mật khẩu Hash bỏ vào Sheet (VD: /tool/hash/123456)
	router.GET("/tool/hash/:pass", func(c *gin.Context) {
		pass := c.Param("pass")
		hash, _ := bao_mat.HashMatKhau(pass)
		c.String(200, "Mật khẩu gốc: %s\nHash (Copy vào cột D): %s", pass, hash)
	})

	// --- NHÓM ADMIN (CẦN ĐĂNG NHẬP) ---
	admin := router.Group("/admin")
	admin.Use(chuc_nang.KiemTraQuyenHan) // <--- NGƯỜI GÁC CỔNG (MIDDLEWARE)
	{
		// Trang tổng quan (Dashboard)
		admin.GET("/tong-quan", func(c *gin.Context) {
			// Lấy thông tin user từ Context (do Middleware gán vào)
			userID, _ := c.Get("USER_ID")
			userRole, _ := c.Get("USER_ROLE")
			
			c.HTML(http.StatusOK, "khung_quan_tri", gin.H{ // Giả sử bạn sẽ tạo layout này sau
				"NoiDung": "Chào mừng " + userID.(string) + " (" + userRole.(string) + ")",
			})
			// Hoặc tạm thời trả về text nếu chưa có HTML admin
			// c.String(200, "Chào sếp %v! Quyền hạn: %v", userID, userRole)
		})

		// API nạp lại dữ liệu (Khi sửa Sheet thủ công)
		admin.GET("/reload", chuc_nang.API_NapLaiDuLieu)
	}

	// ---------------------------------------------------------
	// 4. CHẠY SERVER VỚI GRACEFUL SHUTDOWN (TẮT AN TOÀN)
	// ---------------------------------------------------------
	port := cau_hinh.BienCauHinh.CongChayWeb
	if port == "" { port = "8080" }

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Chạy Server trong 1 luồng riêng (Goroutine)
	go func() {
		log.Println(">> Server đang chạy tại cổng: " + port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Lỗi khởi động Server: %s\n", err)
		}
	}()

	// Chờ tín hiệu tắt từ hệ điều hành (Ctrl+C hoặc Cloud Run Stop)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	<-quit // Code sẽ dừng ở đây cho đến khi nhận tín hiệu
	
	log.Println("⚠️  Đang tắt Server... Đang xả hàng đợi lần cuối!")

	// GỌI HÀM XẢ HÀNG KHẨN CẤP (Ghi nốt dữ liệu chưa kịp ghi)
	nghiep_vu.ThucHienGhiSheet(true)

	log.Println("✅ Server đã tắt an toàn. Hẹn gặp lại!")
}
