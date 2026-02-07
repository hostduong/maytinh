package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"app/bo_nho_dem"
	"app/cau_hinh"
	"app/chuc_nang"
	"app/kho_du_lieu"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

//go:embed giao_dien/*.html
var f embed.FS

func main() {
	log.Println(">>> [SYSTEM] KHỞI ĐỘNG...")

	// 1. Tải cấu hình
	cau_hinh.KhoiTaoCauHinh()
	
	// 2. Kết nối Google Sheet
	func() { defer func() { recover() }(); kho_du_lieu.KhoiTaoKetNoiGoogle() }()

	// 3. Đấu nối Callback
	bo_nho_dem.CallbackGhiSheet = nghiep_vu.ThucHienGhiSheet

	// 4. Khởi tạo RAM rỗng
	bo_nho_dem.KhoiTaoCacStore()
	
	// [CHIẾN LƯỢC KHỞI ĐỘNG ƯU TIÊN - PRIORITY BOOT]
	// Bước 1: Tải dữ liệu cốt lõi (User, Config) -> Bắt buộc chờ xong mới chạy tiếp
	log.Println("🔒 [BOOT] Đang nạp dữ liệu Cốt Lõi (Auth)...")
	bo_nho_dem.NapDuLieuCotLoi()
	
	// Bước 2: Tải dữ liệu nền (Sản phẩm, Đơn hàng) -> Chạy ngầm
	go bo_nho_dem.NapDuLieuNen()
	
	// 5. Khởi động Worker
	nghiep_vu.KhoiTaoWorkerGhiSheet()
	chuc_nang.KhoiTaoBoDemRateLimit()

	// 6. Cấu hình Router
	router := gin.Default()
	templ := template.Must(template.New("").ParseFS(f, "giao_dien/*.html"))
	router.SetHTMLTemplate(templ)

	// --- ĐỊNH NGHĨA ROUTER ---
	
	// Public
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

	// User API (Cần Login)
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
		} else { 
			c.Redirect(http.StatusFound, "/login") 
		}
	})

	// Admin Group
	admin := router.Group("/admin")
	admin.Use(chuc_nang.KiemTraQuyenHan) 
	{
		admin.GET("/tong-quan", chuc_nang.TrangTongQuan)
		admin.GET("/reload", chuc_nang.API_NapLaiDuLieu)
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
	nghiep_vu.ThucHienGhiSheet(true) 
	log.Println("✅ Server tắt an toàn.")
}
