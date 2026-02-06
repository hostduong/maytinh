package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"app/cau_hinh"
	"app/chuc_nang"
	"app/kho_du_lieu"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

//go:embed giao_dien/*.html
var f embed.FS

// Middleware để bảo vệ người dùng khi hệ thống đang reload (Read Lock)
func MW_KiemTraHeThong(c *gin.Context) {
	nghiep_vu.KhoaHeThong.RLock()
	defer nghiep_vu.KhoaHeThong.RUnlock()
	c.Next()
}

func main() {
	log.Println(">>> [SYSTEM] KHỞI ĐỘNG...")

	cau_hinh.KhoiTaoCauHinh()
	
	// Kết nối Google Sheet (ADC)
	func() { defer func() { recover() }(); kho_du_lieu.KhoiTaoKetNoiGoogle() }()

	// Khởi tạo các Store rỗng
	nghiep_vu.KhoiTaoCacStore()
	
	// Nạp dữ liệu lần đầu (Chạy ngầm để server start nhanh)
	go func() {
		log.Println("--- [BOOT] Đang nạp dữ liệu khởi động... ---")
		// Hàm này giờ đã bao gồm cả nạp PHAN_QUYEN (do bạn đã sửa bo_nho_dem.go)
		nghiep_vu.KhoiTaoBoNho() 
	}()
	
	// Khởi động Worker ghi Sheet và Rate Limit
	nghiep_vu.KhoiTaoWorkerGhiSheet()
	chuc_nang.KhoiTaoBoDemRateLimit()

	router := gin.Default()
	
	// Áp dụng Middleware "Êm ái" cho toàn bộ web
	router.Use(MW_KiemTraHeThong)

	templ := template.Must(template.New("").ParseFS(f, "giao_dien/*.html"))
	router.SetHTMLTemplate(templ)

	// --- 1. CÁC ROUTE PUBLIC & USER ---
	router.GET("/", chuc_nang.TrangChu)
	router.GET("/san-pham/:id", chuc_nang.ChiTietSanPham)
	
	// Auth
	router.GET("/login", chuc_nang.TrangDangNhap)
	router.POST("/login", chuc_nang.XuLyDangNhap)
	router.GET("/register", chuc_nang.TrangDangKy)
	router.POST("/register", chuc_nang.XuLyDangKy)
	router.GET("/logout", chuc_nang.DangXuat)
	
	// Forgot Password
	router.GET("/forgot-password", chuc_nang.TrangQuenMatKhau)
	router.POST("/api/auth/reset-by-pin", chuc_nang.XuLyQuenPassBangPIN)
	router.POST("/api/auth/send-otp", chuc_nang.XuLyGuiOTPEmail)
	router.POST("/api/auth/reset-by-otp", chuc_nang.XuLyQuenPassBangOTP)

	// User API
	userGroup := router.Group("/api/user")
	{
		userGroup.POST("/update-info", chuc_nang.API_DoiThongTin)
		userGroup.POST("/change-pass", chuc_nang.API_DoiMatKhau)
		userGroup.POST("/change-pin", chuc_nang.API_DoiMaPin)
		userGroup.POST("/send-otp-pin", chuc_nang.API_GuiOTPPin)
	}

	// Trang cá nhân
	router.GET("/tai-khoan", func(c *gin.Context) {
		cookie, _ := c.Cookie("session_id")
		if cookie == "" { c.Redirect(http.StatusFound, "/login"); return }
		if kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
			c.HTML(http.StatusOK, "ho_so", gin.H{"TieuDe": "Hồ sơ", "NhanVien": kh, "DaDangNhap": true, "TenNguoiDung": kh.TenKhachHang, "QuyenHan": kh.VaiTroQuyenHan})
		} else { c.Redirect(http.StatusFound, "/login") }
	})

	// --- 2. NHÓM ADMIN (CÓ PHÂN QUYỀN RBAC) ---
	admin := router.Group("/admin")
	admin.Use(chuc_nang.KiemTraQuyenHan) // Middleware chặn người không phận sự
	{
		// Dashboard
		admin.GET("/tong-quan", func(c *gin.Context) {
			userID, _ := c.Get("USER_ID")
			// Lấy lại info user để hiển thị avatar/tên
			kh, _ := nghiep_vu.TimKhachHangTheoCookie(mustGetCookie(c))
			
			c.HTML(http.StatusOK, "quan_tri", gin.H{
				"TieuDe": "Quản trị", 
				"NhanVien": kh, 
				"DaDangNhap": true, 
				"TenNguoiDung": kh.TenKhachHang, 
				"QuyenHan": kh.VaiTroQuyenHan, 
				"UserID": userID,
			})
		})

		// [UPDATED] API Nạp lại dữ liệu (Đã được bảo vệ bởi check quyền 'system.reload')
		admin.GET("/reload", chuc_nang.API_NapLaiDuLieu)

		// [NEW] API Sửa thành viên (Đã được bảo vệ bởi check quyền 'member.edit')
		admin.POST("/api/member/update", chuc_nang.API_Admin_SuaThanhVien)
	}

	// --- KHỞI CHẠY SERVER ---
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	
	srv := &http.Server{ Addr: "0.0.0.0:" + port, Handler: router }

	go func() {
		log.Printf("✅ Server chạy tại 0.0.0.0:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ LỖI SERVER: %v", err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("⚠️ Đang tắt Server...")
	nghiep_vu.ThucHienGhiSheet(true) // Ghi nốt dữ liệu còn sót lại
	log.Println("✅ Server tắt an toàn.")
}

func mustGetCookie(c *gin.Context) string { cookie, _ := c.Cookie("session_id"); return cookie }
