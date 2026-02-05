package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"sync/atomic" // Thêm thư viện này để dùng cờ báo hiệu

	"app/bao_mat"
	"app/cau_hinh"
	"app/chuc_nang"
	"app/kho_du_lieu"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

// Biến cờ đánh dấu trạng thái nạp dữ liệu (0: chưa xong, 1: xong)
var DaNapDuLieuXong int32 = 0

func main() {
	log.Println(">>> ĐANG KHỞI ĐỘNG HỆ THỐNG MAYTINHSHOP...")

	// 1. Cấu hình cơ bản
	cau_hinh.KhoiTaoCauHinh()
	kho_du_lieu.KhoiTaoKetNoiGoogle()

	// 2. [QUAN TRỌNG] Chạy nạp dữ liệu ở luồng riêng (Background)
	// Để không chặn việc mở cổng Server bên dưới
	go func() {
		log.Println("--- [BACKGROUND] Bắt đầu nạp dữ liệu từ Sheet... ---")
		nghiep_vu.KhoiTaoBoNho()
		atomic.StoreInt32(&DaNapDuLieuXong, 1) // Bật cờ báo hiệu đã xong
		log.Println("--- [BACKGROUND] Đã nạp xong toàn bộ dữ liệu! ---")
	}()
	
	nghiep_vu.KhoiTaoWorkerGhiSheet()
	chuc_nang.KhoiTaoBoDemRateLimit()

	// 3. Cấu hình Web Server
	router := gin.Default()
	router.LoadHTMLGlob("giao_dien/*.html") // Đã đúng từ phiên bản trước

	// Middleware kiểm tra trạng thái khởi động
	// Nếu dữ liệu chưa nạp xong, trả về thông báo chờ thay vì để Crash
	router.Use(func(c *gin.Context) {
		if atomic.LoadInt32(&DaNapDuLieuXong) == 0 {
			c.JSON(503, gin.H{
				"status": "starting",
				"msg": "Hệ thống đang khởi động và nạp dữ liệu. Vui lòng thử lại sau 30 giây.",
			})
			c.Abort()
			return
		}
		c.Next()
	})

	// --- CÁC ROUTES GIỮ NGUYÊN ---
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

	userGroup := router.Group("/api/user")
	{
		userGroup.POST("/update-info", chuc_nang.API_DoiThongTin)
		userGroup.POST("/change-pass", chuc_nang.API_DoiMatKhau)
		userGroup.POST("/change-pin", chuc_nang.API_DoiMaPin)
		userGroup.POST("/send-otp-pin", chuc_nang.API_GuiOTPPin)
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

	// 4. MỞ CỔNG NGAY LẬP TỨC
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	
	srv := &http.Server{ Addr: "0.0.0.0:" + port, Handler: router }

	go func() {
		log.Printf("✅ Server đang lắng nghe tại: 0.0.0.0:%s (Chờ dữ liệu nạp...)", port)
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
