package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath" // Thêm thư viện này để liệt kê file
	"sync/atomic"
	"syscall"

	"app/bao_mat"
	"app/cau_hinh"
	"app/chuc_nang"
	"app/kho_du_lieu"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

var DaNapDuLieuXong int32 = 0

func main() {
	log.Println(">>> [STARTUP] ĐANG KHỞI ĐỘNG HỆ THỐNG...")

	// 1. Cấu hình & Kết nối
	cau_hinh.KhoiTaoCauHinh()
	kho_du_lieu.KhoiTaoKetNoiGoogle()

	// 2. Chạy ngầm việc nạp dữ liệu (Không chặn Server khởi động)
	go func() {
		log.Println("--- [DATA] Bắt đầu tải dữ liệu từ Google Sheet... ---")
		// Dùng defer recover để tránh việc nạp dữ liệu làm sập cả web
		defer func() {
			if r := recover(); r != nil {
				log.Println("❌ [DATA ERROR] Lỗi nghiêm trọng khi nạp dữ liệu:", r)
			}
		}()
		nghiep_vu.KhoiTaoBoNho()
		atomic.StoreInt32(&DaNapDuLieuXong, 1)
		log.Println("✅ [DATA] Đã nạp xong dữ liệu!")
	}()
	
	nghiep_vu.KhoiTaoWorkerGhiSheet()
	chuc_nang.KhoiTaoBoDemRateLimit()

	// 3. Web Server
	router := gin.Default()

	// --- [ĐOẠN CODE DÒ LỖI QUAN TRỌNG] ---
	// Kiểm tra xem thực sự có file nào trong thư mục giao_dien không
	files, _ := filepath.Glob("giao_dien/*")
	log.Println("📂 [DEBUG] Danh sách file trong thư mục 'giao_dien':", files)

	// Thử nạp HTML, nếu lỗi thì BỎ QUA để Server vẫn chạy được (không bị Crash)
	func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println("⚠️ [HTML ERROR] Không nạp được giao diện (Web sẽ chạy API only). Lỗi:", r)
			}
		}()
		// Load file html phẳng
		router.LoadHTMLGlob("giao_dien/*.html")
		log.Println("✅ [HTML] Đã nạp giao diện thành công.")
	}()
	// --------------------------------------

	// Middleware chặn truy cập khi chưa nạp xong data
	router.Use(func(c *gin.Context) {
		if atomic.LoadInt32(&DaNapDuLieuXong) == 0 {
			c.JSON(503, gin.H{"status": "loading", "msg": "Hệ thống đang khởi động, vui lòng đợi..."})
			c.Abort()
			return
		}
		c.Next()
	})

	// Routes
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
		if cookie == "" { c.Redirect(http.StatusFound, "/login"); return }
		if kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
			// Nếu HTML chưa load được thì trả về JSON để debug
			if len(router.Routes()) > 0 { 
				c.HTML(http.StatusOK, "ho_so", gin.H{"TieuDe": "Hồ sơ", "NhanVien": kh, "DaDangNhap": true, "TenNguoiDung": kh.TenKhachHang, "QuyenHan": kh.VaiTroQuyenHan})
			} else {
				c.JSON(200, kh)
			}
		} else { c.Redirect(http.StatusFound, "/login") }
	})

	router.GET("/tool/hash/:pass", func(c *gin.Context) {
		pass := c.Param("pass"); hash, _ := bao_mat.HashMatKhau(pass)
		c.String(200, "Hash: %s", hash)
	})

	admin := router.Group("/admin")
	admin.Use(chuc_nang.KiemTraQuyenHan)
	{
		admin.GET("/tong-quan", func(c *gin.Context) {
			userID, _ := c.Get("USER_ID"); kh, _ := nghiep_vu.TimKhachHangTheoCookie(mustGetCookie(c))
			c.HTML(http.StatusOK, "quan_tri", gin.H{"TieuDe": "Quản trị", "NhanVien": kh, "DaDangNhap": true, "TenNguoiDung": kh.TenKhachHang, "QuyenHan": kh.VaiTroQuyenHan, "UserID": userID})
		})
		admin.GET("/reload", chuc_nang.API_NapLaiDuLieu)
	}

	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	
	srv := &http.Server{ Addr: "0.0.0.0:" + port, Handler: router }

	go func() {
		log.Printf("✅ Server lắng nghe tại 0.0.0.0:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("❌ Lỗi Server: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	nghiep_vu.ThucHienGhiSheet(true)
}

func mustGetCookie(c *gin.Context) string { cookie, _ := c.Cookie("session_id"); return cookie }
