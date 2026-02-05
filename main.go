package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
    // Đã xóa "time" vì không dùng
    // Đã xóa "app/bao_mat" vì không dùng trong main

	"app/cau_hinh"
	"app/chuc_nang"
	"app/kho_du_lieu"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

//go:embed giao_dien/*.html
var f embed.FS

// Middleware để bảo vệ người dùng khi hệ thống đang reload
func MW_KiemTraHeThong(c *gin.Context) {
    // Xin quyền "Đọc" (RLock)
    nghiep_vu.KhoaHeThong.RLock()
    defer nghiep_vu.KhoaHeThong.RUnlock()
    c.Next()
}

func main() {
	log.Println(">>> [SYSTEM] KHỞI ĐỘNG...")

	cau_hinh.KhoiTaoCauHinh()
    // Sử dụng ADC mặc định của Cloud Run (Không JSON)
	func() { defer func() { recover() }(); kho_du_lieu.KhoiTaoKetNoiGoogle() }()

    // Tạo hộp rỗng trước
	nghiep_vu.KhoiTaoCacStore()
    
    // Nạp dữ liệu lần đầu
	go func() {
		log.Println("--- [BOOT] Đang nạp dữ liệu khởi động... ---")
		nghiep_vu.KhoiTaoBoNho() 
	}()
	
	nghiep_vu.KhoiTaoWorkerGhiSheet()
	chuc_nang.KhoiTaoBoDemRateLimit()

	router := gin.Default()
    
    // Áp dụng Middleware "Êm ái" cho toàn bộ web
    router.Use(MW_KiemTraHeThong)

	templ := template.Must(template.New("").ParseFS(f, "giao_dien/*.html"))
	router.SetHTMLTemplate(templ)

	// --- CÁC ROUTE ---
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

    // --- ADMIN & RELOAD ---
	admin := router.Group("/admin")
	admin.Use(chuc_nang.KiemTraQuyenHan)
	{
		admin.GET("/tong-quan", func(c *gin.Context) {
            userID, _ := c.Get("USER_ID"); kh, _ := nghiep_vu.TimKhachHangTheoCookie(mustGetCookie(c))
			c.HTML(http.StatusOK, "quan_tri", gin.H{"TieuDe": "Quản trị", "NhanVien": kh, "DaDangNhap": true, "TenNguoiDung": kh.TenKhachHang, "QuyenHan": kh.VaiTroQuyenHan, "UserID": userID})
		})

        // [LOGIC RELOAD CHUẨN]
		admin.GET("/reload", func(c *gin.Context) {
            log.Println("⚡ [RELOAD] Bắt đầu quy trình nạp lại dữ liệu...")
            
            // B1: Ép ghi toàn bộ hàng chờ xuống Sheet
            nghiep_vu.ThucHienGhiSheet(true) 
            
            // B2: Khóa toàn hệ thống
            nghiep_vu.KhoaHeThong.Lock()
            log.Println("🔒 [LOCK] Đã khóa hệ thống.")
            
            // Chạy ngầm việc nạp để trả về response ngay cho admin đỡ treo
            go func() {
                defer nghiep_vu.KhoaHeThong.Unlock() // B5: Mở khóa khi xong
                
                // B3: Reset RAM
                nghiep_vu.KhoiTaoCacStore()
                
                // B4: Tải lại từ Sheet
                nghiep_vu.KhoiTaoBoNho()
                
                log.Println("🔓 [UNLOCK] Đã mở khóa hệ thống.")
            }()

            c.JSON(200, gin.H{"status": "ok", "msg": "Hệ thống đang nạp lại. Vui lòng đợi 10-20 giây."})
		})
	}

	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	
	srv := &http.Server{ Addr: "0.0.0.0:" + port, Handler: router }

	go func() {
		log.Printf("✅ Server chạy tại 0.0.0.0:%s", port)
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
