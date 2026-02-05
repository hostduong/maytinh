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
	log.Println(">>> [HARDCORE] ĐANG KHỞI ĐỘNG HỆ THỐNG MAYTINHSHOP...")

	// 1. Cấu hình & Kết nối
	cau_hinh.KhoiTaoCauHinh()
	kho_du_lieu.KhoiTaoKetNoiGoogle()

	// 2. Tạo kho rỗng & Chạy ngầm nạp dữ liệu
	nghiep_vu.KhoiTaoCacStore()
	go func() {
		log.Println("--- [BACKGROUND] Đang nạp dữ liệu... ---")
		nghiep_vu.KhoiTaoBoNho()
		log.Println("--- [BACKGROUND] Nạp xong! ---")
	}()
	
	nghiep_vu.KhoiTaoWorkerGhiSheet()
	chuc_nang.KhoiTaoBoDemRateLimit()

	// 3. Web Server
	router := gin.Default()
	router.LoadHTMLGlob("giao_dien/*.html")

	// --- ROUTES ---
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
			c.HTML(http.StatusOK, "ho_so", gin.H{"TieuDe": "Hồ sơ", "NhanVien": kh, "DaDangNhap": true, "TenNguoiDung": kh.TenKhachHang, "QuyenHan": kh.VaiTroQuyenHan})
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

	// [FIX CỨNG CỔNG 8080]
	// Bỏ qua biến môi trường, ép chạy 8080 để khớp với EXPOSE trong Dockerfile
	addr := "0.0.0.0:8080" 
	srv := &http.Server{ Addr: addr, Handler: router }

	go func() {
		log.Printf("✅ Server ĐANG CHẠY CỐ ĐỊNH TẠI: %s", addr)
		// Bỏ qua lỗi server closed để log sạch sẽ hơn
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ LỖI SERVER: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("⚠️ Đang tắt Server...")
	nghiep_vu.ThucHienGhiSheet(true)
	log.Println("✅ Server tắt an toàn.")
}

func mustGetCookie(c *gin.Context) string { cookie, _ := c.Cookie("session_id"); return cookie }
