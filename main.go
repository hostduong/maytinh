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

	// 1. Nạp cấu hình
	cau_hinh.KhoiTaoCauHinh()

	// 2. Kết nối Sheet (Chế độ Public - Tránh lỗi thiếu file JSON)
	kho_du_lieu.KhoiTaoKetNoiGoogle()

	// 3. Khởi tạo bộ nhớ & Worker (Chạy ngầm để Server lên ngay lập tức)
	go func() {
		nghiep_vu.KhoiTaoBoNho()
	}()
	nghiep_vu.KhoiTaoWorkerGhiSheet()
	chuc_nang.KhoiTaoBoDemRateLimit()

	// 4. Cấu hình Web Server
	router := gin.Default()
	
	// [SỬA LỖI QUAN TRỌNG NHẤT]: Dùng *.html thay vì **/* // Vì thư mục giao_dien không có thư mục con, dùng ** sẽ gây Crash.
	router.LoadHTMLGlob("giao_dien/*.html")

	// --- PUBLIC ROUTES ---
	router.GET("/", chuc_nang.TrangChu)
	router.GET("/san-pham/:id", chuc_nang.ChiTietSanPham)
	
	router.GET("/login", chuc_nang.TrangDangNhap)
	router.POST("/login", chuc_nang.XuLyDangNhap)
	router.GET("/register", chuc_nang.TrangDangKy)
	router.POST("/register", chuc_nang.XuLyDangKy)
	router.GET("/logout", chuc_nang.DangXuat)

	router.GET("/forgot-password", chuc_nang.TrangQuenMatKhau)
	router.POST("/api/auth/reset-by-pin", chuc_nang.XuLyQuenPassBangPIN)
	router.POST("/api/auth/send-otp", chuc_nang.XuLyGuiOTPEmail)
	router.POST("/api/auth/reset-by-otp", chuc_nang.XuLyQuenPassBangOTP)

	// --- USER ROUTES ---
	userGroup := router.Group("/api/user")
	{
		userGroup.POST("/update-info", chuc_nang.API_DoiThongTin)
		userGroup.POST("/change-pass", chuc_nang.API_DoiMatKhau)
		userGroup.POST("/change-pin", chuc_nang.API_DoiMaPin)
		userGroup.POST("/send-otp-pin", chuc_nang.API_GuiOTPPin)
		// Đã xóa API_ResetPinBangOTP để tránh lỗi Build
	}

	router.GET("/tai-khoan", func(c *gin.Context) {
		cookie, _ := c.Cookie("session_id")
		if cookie == "" {
			 c.Redirect(http.StatusFound, "/login")
			 return
		}
		if kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
			 c.HTML(http.StatusOK, "ho_so", gin.H{
			 	"TieuDe":       "Hồ sơ của bạn",
			 	"NhanVien":     kh,
			 	"DaDangNhap":   true,
			 	"TenNguoiDung": kh.TenKhachHang,
			 	"QuyenHan":     kh.VaiTroQuyenHan,
			 })
		} else {
			 c.Redirect(http.StatusFound, "/login")
		}
	})

	router.GET("/tool/hash/:pass", func(c *gin.Context) {
		pass := c.Param("pass")
		hash, _ := bao_mat.HashMatKhau(pass)
		c.String(200, "Pass: %s\nHash: %s", pass, hash)
	})

	// --- ADMIN ROUTES ---
	admin := router.Group("/admin")
	admin.Use(chuc_nang.KiemTraQuyenHan)
	{
		admin.GET("/tong-quan", func(c *gin.Context) {
			userID, _ := c.Get("USER_ID")
			kh, _ := nghiep_vu.TimKhachHangTheoCookie(mustGetCookie(c))
			c.HTML(http.StatusOK, "quan_tri", gin.H{
				"TieuDe":       "Quản trị hệ thống",
				"NhanVien":     kh,
				"DaDangNhap":   true,
				"TenNguoiDung": kh.TenKhachHang,
				"QuyenHan":     kh.VaiTroQuyenHan,
				"UserID":       userID,
			})
		})
		admin.GET("/reload", chuc_nang.API_NapLaiDuLieu)
	}

	// [PORT CHO CLOUD RUN]
	port := os.Getenv("PORT")
	if port == "" {
		port = cau_hinh.BienCauHinh.CongChayWeb
	}
	if port == "" {
		port = "8080"
	}
	
	// FIX LỖI 2: Phải nghe 0.0.0.0
	srv := &http.Server{ Addr: "0.0.0.0:" + port, Handler: router }

	go func() {
		log.Printf("✅ Server đang lắng nghe tại 0.0.0.0:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Lỗi server: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("⚠️ Đang tắt Server...")
	nghiep_vu.ThucHienGhiSheet(true)
	log.Println("✅ Server đã tắt an toàn.")
}

func mustGetCookie(c *gin.Context) string {
	cookie, _ := c.Cookie("session_id")
	return cookie
}
