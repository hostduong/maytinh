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
	router.GET("/", chuc_nang.TrangChu)
	router.GET("/san-pham/:id", chuc_nang.ChiTietSanPham)
	
	// Auth
	router.GET("/login", chuc_nang.TrangDangNhap)
	router.POST("/login", chuc_nang.XuLyDangNhap)
	router.GET("/register", chuc_nang.TrangDangKy)
	router.POST("/register", chuc_nang.XuLyDangKy)
	router.GET("/logout", chuc_nang.DangXuat)

	// --- [MỚI] QUÊN MẬT KHẨU ---
	router.GET("/forgot-password", chuc_nang.TrangQuenMatKhau)
	router.POST("/api/auth/reset-by-pin", chuc_nang.XuLyQuenPassBangPIN)
	router.POST("/api/auth/send-otp", chuc_nang.XuLyGuiOTPEmail)
	router.POST("/api/auth/reset-by-otp", chuc_nang.XuLyQuenPassBangOTP)
	// ---------------------------

	// --- [MỚI] API NGƯỜI DÙNG (Yêu cầu Login) ---
	// Nhóm này xử lý trong trang Hồ Sơ
	userGroup := router.Group("/api/user")
	{
		userGroup.POST("/update-info", chuc_nang.API_DoiThongTin)
		userGroup.POST("/change-pass", chuc_nang.API_DoiMatKhau)
		userGroup.POST("/change-pin", chuc_nang.API_DoiMaPin)
	}
	// --------------------------------------------

	// Trang Hồ Sơ
	router.GET("/tai-khoan", func(c *gin.Context) {
		cookie, _ := c.Cookie("session_id")
		if cookie == "" {
			 c.Redirect(http.StatusFound, "/login")
			 return
		}
		
		// [SỬA] Đổi sang TimKhachHang
		if kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
			 c.HTML(http.StatusOK, "ho_so", gin.H{
			 	"TieuDe":       "Hồ sơ của bạn",
			 	"NhanVien":     kh,              // Truyền object KhachHang vào key NhanVien để template dùng
			 	"DaDangNhap":   true,
			 	"TenNguoiDung": kh.TenKhachHang, // [SỬA] Dùng TenKhachHang
			 	"QuyenHan":     kh.VaiTroQuyenHan,
			 })
		} else {
			 c.Redirect(http.StatusFound, "/login")
		}
	})

	// Tool Hash
	router.GET("/tool/hash/:pass", func(c *gin.Context) {
		pass := c.Param("pass")
		hash, _ := bao_mat.HashMatKhau(pass)
		c.String(200, "Pass: %s\nHash: %s", pass, hash)
	})

	// Admin
	admin := router.Group("/admin")
	admin.Use(chuc_nang.KiemTraQuyenHan)
	{
		admin.GET("/tong-quan", func(c *gin.Context) {
			userID, _ := c.Get("USER_ID")
			// [SỬA] Đổi sang TimKhachHang
			kh, _ := nghiep_vu.TimKhachHangTheoCookie(mustGetCookie(c))
			
			c.HTML(http.StatusOK, "quan_tri", gin.H{
				"TieuDe":       "Quản trị hệ thống",
				"NhanVien":     kh,              // Truyền object KhachHang
				"DaDangNhap":   true,
				"TenNguoiDung": kh.TenKhachHang, // [SỬA] Dùng TenKhachHang
				"QuyenHan":     kh.VaiTroQuyenHan,
				"UserID":       userID,
			})
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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("⚠️ Đang tắt Server... Xả hàng đợi...")
	nghiep_vu.ThucHienGhiSheet(true)
	log.Println("✅ Server đã tắt an toàn.")
}

func mustGetCookie(c *gin.Context) string {
	cookie, _ := c.Cookie("session_id")
	return cookie
}
