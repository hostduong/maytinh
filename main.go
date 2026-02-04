package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"app/bao_mat"
	"app/cau_hinh"
	"app/chuc_nang"
	"app/kho_du_lieu"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println(">>> ĐANG KHỞI ĐỘNG HỆ THỐNG MAYTINHSHOP...")

	// 1. Khởi tạo
	cau_hinh.KhoiTaoCauHinh()
	kho_du_lieu.KhoiTaoKetNoiGoogle()
	nghiep_vu.KhoiTaoBoNho()
	nghiep_vu.KhoiTaoWorkerGhiSheet()
	chuc_nang.KhoiTaoBoDemRateLimit()

	// 2. Web Server
	router := gin.Default()
	router.LoadHTMLGlob("giao_dien/**/*")

	// 3. Router
	// --- NHÓM PUBLIC (Ai cũng vào được) ---
	router.GET("/", chuc_nang.TrangChu)
	router.GET("/san-pham/:id", chuc_nang.ChiTietSanPham)
	
	// Login & Logout
	router.GET("/login", chuc_nang.TrangDangNhap)
	router.POST("/login", chuc_nang.XuLyDangNhap)
	router.GET("/logout", chuc_nang.DangXuat)
	
	// --- [MỚI] ĐĂNG KÝ TÀI KHOẢN ---
	router.GET("/register", chuc_nang.TrangDangKy)  // Hiển thị form
	router.POST("/register", chuc_nang.XuLyDangKy)  // Xử lý đăng ký
	// -------------------------------
	
	// Tool Hash Pass (Tiện ích)
	router.GET("/tool/hash/:pass", func(c *gin.Context) {
		pass := c.Param("pass")
		hash, _ := bao_mat.HashMatKhau(pass)
		c.String(200, "Pass: %s\nHash: %s", pass, hash)
	})

	// --- NHÓM ADMIN (Cần đăng nhập) ---
	admin := router.Group("/admin")
	admin.Use(chuc_nang.KiemTraQuyenHan)
	{
		admin.GET("/tong-quan", func(c *gin.Context) {
			userID, _ := c.Get("USER_ID")
			userRole, _ := c.Get("USER_ROLE")
			c.String(200, "Chào sếp %v! Quyền hạn của sếp là: %v", userID, userRole)
		})
		admin.GET("/reload", chuc_nang.API_NapLaiDuLieu)
	}

	// 4. Chạy Server
	port := cau_hinh.BienCauHinh.CongChayWeb
	if port == "" { port = "8080" }
	
	srv := &http.Server{ Addr: ":" + port, Handler: router }

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Lỗi server: %s\n", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("⚠️ Đang tắt Server... Xả hàng đợi...")
	nghiep_vu.ThucHienGhiSheet(true)
	log.Println("✅ Server đã tắt.")
}
