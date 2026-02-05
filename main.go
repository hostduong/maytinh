package main

import (
	"embed"
	"html/template"
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

// [CÔNG NGHỆ EMBED]
// Dòng này ra lệnh cho Go: "Hãy nhét toàn bộ thư mục giao_dien vào trong file chạy này"
//go:embed giao_dien/*.html
var f embed.FS

func main() {
	log.Println(">>> [EMBED MODE] ĐANG KHỞI ĐỘNG HỆ THỐNG...")

	// 1. Cấu hình & Kết nối (Chạy an toàn)
	cau_hinh.KhoiTaoCauHinh()
	// Gọi kết nối nhưng không để chết chương trình nếu lỗi
	func() {
		defer func() { recover() }() 
		kho_du_lieu.KhoiTaoKetNoiGoogle()
	}()

	// 2. Tạo kho rỗng & Chạy ngầm nạp dữ liệu
	nghiep_vu.KhoiTaoCacStore()
	go func() {
		log.Println("--- [BACKGROUND] Đang nạp dữ liệu... ---")
		nghiep_vu.KhoiTaoBoNho()
		log.Println("--- [BACKGROUND] Nạp xong! ---")
	}()
	
	nghiep_vu.KhoiTaoWorkerGhiSheet()
	chuc_nang.KhoiTaoBoDemRateLimit()

	// 3. Cấu hình Web Server
	router := gin.Default()

	// [QUAN TRỌNG] Nạp HTML từ bộ nhớ Embed (Không phụ thuộc file bên ngoài nữa)
	templ := template.Must(template.New("").ParseFS(f, "giao_dien/*.html"))
	router.SetHTMLTemplate(templ)
	log.Println("✅ Đã nạp giao diện từ Embed (An toàn tuyệt đối)")

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

	// [PORT]
	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	
	srv := &http.Server{ Addr: "0.0.0.0:" + port, Handler: router }

	go func() {
		log.Printf("✅ Server chạy tại 0.0.0.0:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("❌ LỖI SERVER: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	nghiep_vu.ThucHienGhiSheet(true)
}

func mustGetCookie(c *gin.Context) string { cookie, _ := c.Cookie("session_id"); return cookie }
