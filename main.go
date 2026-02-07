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

	cau_hinh.KhoiTaoCauHinh()
	
	func() { defer func() { recover() }(); kho_du_lieu.KhoiTaoKetNoiGoogle() }()

	bo_nho_dem.CallbackGhiSheet = nghiep_vu.ThucHienGhiSheet
	bo_nho_dem.KhoiTaoCacStore()
	
	// [QUAN TRỌNG] CHIẾN LƯỢC KHỞI ĐỘNG NHANH
	// BƯỚC 1: Tải dữ liệu cốt lõi (User, Config) -> Bắt buộc chờ xong mới chạy tiếp
	log.Println("🔒 [BOOT] Đang nạp dữ liệu Cốt Lõi (Auth)...")
	bo_nho_dem.NapDuLieuCotLoi()
	
	// BƯỚC 2: Tải dữ liệu nền (Sản phẩm, Đơn hàng) -> Chạy ngầm (Non-blocking)
	// Server sẽ mở cổng ngay lập tức sau bước 1
	go bo_nho_dem.NapDuLieuNen()
	
	nghiep_vu.KhoiTaoWorkerGhiSheet()
	chuc_nang.KhoiTaoBoDemRateLimit()

	router := gin.Default()
	templ := template.Must(template.New("").ParseFS(f, "giao_dien/*.html"))
	router.SetHTMLTemplate(templ)

	// ... (Phần Router giữ nguyên như cũ) ...
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

	admin := router.Group("/admin")
	admin.Use(chuc_nang.KiemTraQuyenHan) 
	{
		admin.GET("/tong-quan", chuc_nang.TrangTongQuan)
		admin.GET("/reload", chuc_nang.API_NapLaiDuLieu)
		admin.POST("/api/member/update", chuc_nang.API_Admin_SuaThanhVien)
	}

	port := os.Getenv("PORT")
	if port == "" { port = "8080" }
	
	srv := &http.Server{ Addr: "0.0.0.0:" + port, Handler: router }

	go func() {
