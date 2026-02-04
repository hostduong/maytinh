package chuc_nang

import (
	"net/http"
	"sync"
	"time"

	"app/cau_hinh"
	"app/mo_hinh" // [MỚI] Import để lấy index cột
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

// Bộ nhớ đếm Request cho Rate Limit
var boDem = make(map[string]int)
var mtx sync.Mutex

// Khởi chạy bộ đếm (Reset mỗi giây)
func KhoiTaoBoDemRateLimit() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			mtx.Lock()
			boDem = make(map[string]int) // Xóa sạch bộ đếm cũ
			mtx.Unlock()
		}
	}()
}

// MIDDLEWARE CHÍNH
func KiemTraQuyenHan(c *gin.Context) {
	// 1. KIỂM TRA RATE LIMIT (CHỐNG SPAM)
	cookie, err := c.Cookie("session_id")
	keyLimit := ""
	
	if err != nil || cookie == "" {
		keyLimit = "LIMIT__IP__" + c.ClientIP()
	} else {
		keyLimit = "LIMIT__COOKIE__" + cookie
	}

	mtx.Lock()
	boDem[keyLimit]++
	soLanGoi := boDem[keyLimit]
	mtx.Unlock()

	// Logic chặn: Nếu người dùng gọi quá 10 req/s
	if soLanGoi > cau_hinh.GioiHanNguoiDung {
		c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"loi": "Thao tác quá nhanh! Vui lòng chậm lại."})
		return
	}

	// 2. KIỂM TRA ĐĂNG NHẬP (AUTH)
	if cookie == "" {
		c.Next()
		return
	}

	// [SỬA] Tìm trong RAM KHÁCH HÀNG
	khachHang, timThay := nghiep_vu.TimKhachHangTheoCookie(cookie)

	if !timThay {
		// Cookie rác -> Xóa cookie trình duyệt
		c.SetCookie("session_id", "", -1, "/", "", false, true)
		c.Next()
		return
	}

	// 3. LOGIC GIA HẠN THÔNG MINH (Auto-Renew)
	thoiGianHetHan := khachHang.CookieExpired // Dạng int64
	now := time.Now().Unix()

	// Nếu đã hết hạn -> Đá ra
	if now > thoiGianHetHan {
		c.SetCookie("session_id", "", -1, "/", "", false, true)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"loi": "Phiên đăng nhập hết hạn"})
		return
	}

	// Nếu còn hạn nhưng sắp hết (trong vùng ân hạn) -> GIA HẠN
	thoiGianConLai := time.Duration(thoiGianHetHan - now) * time.Second
	if thoiGianConLai < cau_hinh.ThoiGianAnHan {
		
		// A. Tính thời gian mới (+30 phút)
		newExp := time.Now().Add(cau_hinh.ThoiGianHetHanCookie).Unix()
		
		// B. Cập nhật vào RAM ngay (Vì khachHang là con trỏ nên cập nhật trực tiếp được)
		khachHang.CookieExpired = newExp

		// C. Đẩy vào Hàng Chờ Ghi (WriteQueue) -> Worker sẽ ghi xuống Sheet sau
		rowID := nghiep_vu.LayDongKhachHang(khachHang.MaKhachHang)
		if rowID > 0 {
			nghiep_vu.ThemVaoHangCho(
				cau_hinh.BienCauHinh.IdFileSheet, // ID file sheet
				"KHACH_HANG",                     // [SỬA] Tên sheet mới
				rowID,                            // Dòng
				mo_hinh.CotKH_CookieExpired,      // [SỬA] Cột E trong struct KhachHang
				newExp,                           // Giá trị mới
			)
		}

		// D. Set lại Cookie mới cho trình duyệt
		c.SetCookie("session_id", cookie, int(cau_hinh.ThoiGianHetHanCookie.Seconds()), "/", "", false, true)
	}

	// Lưu thông tin user vào Context để Controller dùng
	c.Set("USER_ID", khachHang.MaKhachHang)
	
	// Check quyền
	c.Set("USER_ROLE", khachHang.VaiTroQuyenHan)
	
	c.Next()
}
