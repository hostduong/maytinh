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

	// 1. KHỞI TẠO CÁC THÀNH PHẦN HỆ THỐNG
	cau_hinh.KhoiTaoCauHinh()
	kho_du_lieu.KhoiTaoKetNoiGoogle()
	nghiep_vu.KhoiTaoBoNho()
	nghiep_vu.KhoiTaoWorkerGhiSheet()
	chuc_nang.KhoiTaoBoDemRateLimit()

	// 2. KHỞI TẠO WEB SERVER
	router := gin.Default()
	
	// Nạp toàn bộ file HTML trong thư mục giao_dien (bao gồm cả thư mục con)
	router.LoadHTMLGlob("giao_dien/**/*")

	// 3. ĐỊNH TUYẾN (ROUTING)

	// --- NHÓM PUBLIC (Ai cũng truy cập được) ---
	router.GET("/", chuc_nang.TrangChu)
	router.GET("/san-pham/:id", chuc_nang.ChiTietSanPham)

	// Đăng Ký - Đăng Nhập - Đăng Xuất
	router.GET("/login", chuc_nang.TrangDangNhap)
	router.POST("/login", chuc_nang.XuLyDangNhap)
	
	router.GET("/register", chuc_nang.TrangDangKy)
	router.POST("/register", chuc_nang.XuLyDangKy)
	
	router.GET("/logout", chuc_nang.DangXuat)

	// --- TRANG HỒ SƠ TÀI KHOẢN (Cần đăng nhập) ---
	router.GET("/tai-khoan", func(c *gin.Context) {
		cookie, _ := c.Cookie("session_id")
		if cookie == "" {
			 c.Redirect(http.StatusFound, "/login")
			 return
		}
		
		// Tìm thông tin người dùng từ RAM
		if nv, ok := nghiep_vu.TimNhanVienTheoCookie(cookie); ok {
			 // Render giao diện HTML "ho_so" đẹp mắt
			 c.HTML(http.StatusOK, "ho_so", gin.H{
			 	"TieuDe":       "Hồ sơ của bạn",
			 	"NhanVien":     nv,                 // Truyền object NV để lấy Email, Chức vụ...
			 	"DaDangNhap":   true,               // Báo cho Header biết đã login
			 	"TenNguoiDung": nv.HoTen,           // Hiển thị tên trên Header
			 	"QuyenHan":     nv.VaiTroQuyenHan,  // Để hiện menu Admin nếu có quyền
			 })
		} else {
			 // Cookie không hợp lệ (Session ảo/hết hạn) -> Bắt đăng nhập lại
			 c.Redirect(http.StatusFound, "/login")
		}
	})

	// --- CÔNG CỤ TIỆN ÍCH (Dùng cho Dev/Admin tạo pass thủ công) ---
	router.GET("/tool/hash/:pass", func(c *gin.Context) {
		pass := c.Param("pass")
		hash, _ := bao_mat.HashMatKhau(pass)
		c.String(200, "Pass: %s\nHash: %s", pass, hash)
	})

	// --- NHÓM ADMIN (Phải Login + Có quyền Admin) ---
	admin := router.Group("/admin")
	admin.Use(chuc_nang.KiemTraQuyenHan) // Middleware bảo vệ
	{
		// Trang Dashboard Quản Trị
		admin.GET("/tong-quan", func(c *gin.Context) {
			userID, _ := c.Get("USER_ID")
			
			// Lấy lại thông tin nhân viên để hiển thị tên đầy đủ
			nv, _ := nghiep_vu.TimNhanVienTheoCookie(mustGetCookie(c))

			// Render giao diện HTML "quan_tri"
			c.HTML(http.StatusOK, "quan_tri", gin.H{
				"TieuDe":       "Quản trị hệ thống",
				"NhanVien":     nv,
				"DaDangNhap":   true,
				"TenNguoiDung": nv.HoTen,
				"QuyenHan":     nv.VaiTroQuyenHan,
				"UserID":       userID,
			})
		})
		
		// API Nạp lại dữ liệu từ Sheet (Khi sửa tay trên Excel)
		admin.GET("/reload", chuc_nang.API_NapLaiDuLieu)
	}

	// 4. CHẠY SERVER
	port := cau_hinh.BienCauHinh.CongChayWeb
	if port == "" { port = "8080" } // Mặc định 8080 nếu không có biến môi trường
	
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Chạy server trong luồng riêng (Goroutine)
	go func() {
		log.Println(">> Server đang chạy tại cổng: " + port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Lỗi khởi động Server: %s\n", err)
		}
	}()

	// 5. GRACEFUL SHUTDOWN (Tắt server an toàn)
	// Lắng nghe tín hiệu Ctrl+C hoặc Kill từ hệ điều hành
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Chờ tín hiệu...
	
	log.Println("⚠️ Đang tắt Server... Đang xả hàng đợi lần cuối!")
	nghiep_vu.ThucHienGhiSheet(true) // Ghi nốt dữ liệu chưa kịp lưu
	log.Println("✅ Server đã tắt an toàn.")
}

// Hàm phụ trợ: Lấy cookie nhanh (tránh lặp lại code)
func mustGetCookie(c *gin.Context) string {
	cookie, _ := c.Cookie("session_id")
	return cookie
}
