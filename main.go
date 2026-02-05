package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"app/cau_hinh"
	"app/chuc_nang"
	"app/kho_du_lieu"
	"app/nghiep_vu"
	"app/bao_mat" // Thêm import này nếu cần dùng trong main

	"github.com/gin-gonic/gin"
)

func main() {
	log.Println(">>> KHỞI ĐỘNG HỆ THỐNG...")

	cau_hinh.KhoiTaoCauHinh()
	kho_du_lieu.KhoiTaoKetNoiGoogle()
	nghiep_vu.KhoiTaoBoNho()
	nghiep_vu.KhoiTaoWorkerGhiSheet()
	chuc_nang.KhoiTaoBoDemRateLimit()

	router := gin.Default()
	router.LoadHTMLGlob("giao_dien/**/*")

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
		// Đã xóa dòng reset-pin-otp
	}

	router.GET("/tai-khoan", func(c *gin.Context) {
		cookie, _ := c.Cookie("session_id")
		if cookie == "" {
			c.Redirect(http.StatusFound, "/login")
			return
		}
		if kh, ok := nghiep_vu.TimKhachHangTheoCookie(cookie); ok {
			c.HTML(http.StatusOK, "ho_so", gin.H{
				"TieuDe": "Hồ sơ", "NhanVien": kh, "DaDangNhap": true, "TenNguoiDung": kh.TenKhachHang,
			})
		} else {
			c.Redirect(http.StatusFound, "/login")
		}
	})

	router.GET("/tool/hash/:pass", func(c *gin.Context) {
		pass := c.Param("pass")
		hash, _ := bao_mat.HashMatKhau(pass)
		c.String(200, "Hash: %s", hash)
	})

	admin := router.Group("/admin")
	admin.Use(chuc_nang.KiemTraQuyenHan)
	{
		admin.GET("/tong-quan", func(c *gin.Context) {
			// Logic admin giữ nguyên
			c.HTML(http.StatusOK, "quan_tri", gin.H{})
		})
		admin.GET("/reload", chuc_nang.API_NapLaiDuLieu)
	}

	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	
	srv := &http.Server{ Addr: "0.0.0.0:" + port, Handler: router }

	go func() {
		log.Printf("Server listening on 0.0.0.0:%s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	nghiep_vu.ThucHienGhiSheet(true)
}
