package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath" // Thư viện để tìm file
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
	log.Println(">>> [BOOT] BẮT ĐẦU KHỞI ĐỘNG HỆ THỐNG...")

	// 1. In thông tin môi trường để kiểm tra (Debug)
	dir, _ := os.Getwd()
	log.Println("--- [DEBUG] Thư mục hiện tại:", dir)
	
	// Kiểm tra xem có file HTML nào không
	matches, _ := filepath.Glob("giao_dien/*.html")
	log.Printf("--- [DEBUG] Tìm thấy %d file HTML trong thư mục 'giao_dien'", len(matches))
	for _, f := range matches {
		log.Println("    Found:", f)
	}

	// 2. Cấu hình & Kết nối
	cau_hinh.KhoiTaoCauHinh()
	kho_du_lieu.KhoiTaoKetNoiGoogle()

	// 3. Tạo kho rỗng & Chạy ngầm nạp dữ liệu
	nghiep_vu.KhoiTaoCacStore()
	go func() {
		log.Println("--- [DATA] Đang tải dữ liệu ngầm... ---")
		nghiep_vu.KhoiTaoBoNho()
		atomic.StoreInt32(&DaNapDuLieuXong, 1)
		log.Println("--- [DATA] Đã nạp xong! ---")
	}()
	
	nghiep_vu.KhoiTaoWorkerGhiSheet()
	chuc_nang.KhoiTaoBoDemRateLimit()

	// 4. Cấu hình Web Server
	router := gin.Default()

	// [CHỐNG SẬP] Chỉ nạp HTML nếu thực sự tìm thấy file
	if len(matches) > 0 {
		router.LoadHTMLGlob("giao_dien/*.html")
		log.Println("✅ [HTML] Đã nạp giao diện thành công.")
	} else {
		log.Println("⚠️ [HTML WARNING] KHÔNG tìm thấy file HTML nào! Web sẽ chạy ở chế độ API Only.")
		// Không gọi LoadHTMLGlob để tránh Panic
	}

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
			// Nếu HTML chưa load được thì trả JSON để không lỗi
			if len(matches) > 0 {
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
			if len(matches) > 0 {
				c.HTML(http.StatusOK, "quan_tri", gin.H{"TieuDe": "Quản trị", "NhanVien": kh, "DaDangNhap": true, "TenNguoiDung": kh.TenKhachHang, "QuyenHan": kh.VaiTroQuyenHan, "UserID": userID})
			} else {
				c.JSON(200, gin.H{"msg": "Admin Panel (No HTML)", "data": kh})
			}
		})
		admin.GET("/reload", chuc_nang.API_NapLaiDuLieu)
	}

	// ======== [PORT CHO CLOUD RUN] ========
port := os.Getenv("PORT")
if port == "" {
    port = cau_hinh.BienCauHinh.CongChayWeb
}
if port == "" {
    port = "8080"
}

// ======== KHỞI TẠO SERVER NGAY (QUAN TRỌNG) ========
srv := &http.Server{
    Addr:    "0.0.0.0:" + port,
    Handler: router,
}

go func() {
    log.Printf("✅ Server đang lắng nghe tại 0.0.0.0:%s", port)
    if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("❌ Lỗi server: %s\n", err)
    }
}()

// ======== CHỈ LÀM VIỆC NẶNG SAU KHI SERVER ĐÃ CHẠY ========
go func() {
    log.Println("🔄 Khởi tạo Google Sheet (chạy nền)...")
    kho_du_lieu.KhoiTaoKetNoiGoogle()

    log.Println("🔄 Khởi tạo bộ nhớ & worker (chạy nền)...")
    nghiep_vu.KhoiTaoBoNho()
    nghiep_vu.KhoiTaoWorkerGhiSheet()
    chuc_nang.KhoiTaoBoDemRateLimit()
}()

// ======== GRACEFUL SHUTDOWN ========
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

log.Println("⚠️ Đang tắt Server...")
nghiep_vu.ThucHienGhiSheet(true)
log.Println("✅ Server đã tắt an toàn.")


func mustGetCookie(c *gin.Context) string { cookie, _ := c.Cookie("session_id"); return cookie }
