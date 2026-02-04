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

	// 1. Khởi tạo hệ thống
	cau_hinh.KhoiTaoCauHinh()
	kho_du_lieu.KhoiTaoKetNoiGoogle()
	nghiep_vu.KhoiTaoBoNho()
	nghiep_vu.KhoiTaoWorkerGhiSheet()
	chuc_nang.KhoiTaoBoDemRateLimit()

	// 2. Web Server
	router := gin.Default()
	router.LoadHTMLGlob("giao_dien/**/*")

	// 3. Router - ĐỊNH TUYẾN
	
	// --- NHÓM PUBLIC (Ai cũng vào được) ---
	router.GET("/", chuc_nang.TrangChu)
	router.GET("/san-pham/:id", chuc_nang.ChiTietSanPham)
	
	// Đăng Ký & Đăng Nhập & Đăng Xuất
	router.GET("/login", chuc_nang.TrangDangNhap)
	router.POST("/login", chuc_nang.XuLyDangNhap)
	
	router.GET("/register", chuc_nang.TrangDangKy)
	router.POST("/register", chuc_nang.XuLyDangKy)
	
	router.GET("/logout", chuc_nang.DangXuat)
	
	// --- [BƯỚC 5] TRANG HỒ SƠ TÀI KHOẢN ---
	// (Khi bấm vào tên "Hi, Nguyễn Văn A")
	router.GET("/tai-khoan", func(c *gin.Context) {
		// Kiểm tra xem có đang đăng nhập không
		cookie, _ := c.Cookie("session_id")
		if cookie == "" {
			 c.Redirect(http.StatusFound, "/login")
			 return
		}
		
		// Lấy thông tin từ RAM
		if nv, ok := nghiep_vu.TimNhanVienTheoCookie(cookie); ok {
			 // Tạm thời hiển thị text đơn giản. 
			 // Sau này bạn có thể làm file html riêng như "ho_so.html"
			 c.String(200, "=== THÔNG TIN TÀI KHOẢN ===\n\nHọ tên: %s\nEmail: %s\nChức vụ: %s\nMã thành viên: %s\n\n(Chức năng đổi mật khẩu đang phát triển...)", 
				nv.HoTen, nv.Email, nv.ChucVu, nv.MaNhanVien)
		} else {
			 // Cookie không hợp lệ -> Bắt đăng nhập lại
			 c.Redirect(http.StatusFound, "/login")
		}
	})
	// ---------------------------------------

	// Tool Hash Pass (Dùng cho Admin tạo pass thủ công nếu cần)
	router.GET("/tool/hash/:pass", func(c *gin.Context) {
		pass := c.Param("pass")
		hash, _ := bao_mat.HashMatKhau(pass)
		c.String(200, "Pass: %s\nHash: %s", pass, hash)
	})

	// --- NHÓM ADMIN (Phải Login + Có quyền Admin) ---
	admin := router.Group("/admin")
	admin.Use(chuc_nang.KiemTraQuyenHan)
	{
		admin.GET("/tong-quan", func(c *gin.Context) {
			userID, _ := c.Get("USER_ID")
			userRole, _ := c.Get("USER_ROLE")
			// Sau này sẽ thay bằng c.HTML(200, "admin_dashboard", ...)
			c.String(200, "CHÀO MỪNG QUẢN TRỊ VIÊN!\n\nUser: %v\nRole: %v\n\nTại đây bạn sẽ quản lý kho, đơn hàng...", userID, userRole)
		})
		
		// Nút Reload dữ liệu từ Sheet (quan trọng)
		admin.GET("/reload", chuc_nang.API_NapLaiDuLieu)
	}

	// 4. Chạy Server
	port := cau_hinh.BienCauHinh.CongChayWeb
	if port == "" { port = "8080" }
	
	srv := &http.Server{ Addr: ":" + port, Handler: router }

	// Chạy ngầm
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Lỗi server: %s\n", err)
		}
	}()

	// Graceful Shutdown (Tắt an toàn)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	
	log.Println("⚠️ Đang tắt Server... Xả hàng đợi...")
	nghiep_vu.ThucHienGhiSheet(true) // Ghi nốt dữ liệu chưa lưu
	log.Println("✅ Server đã tắt.")
}
