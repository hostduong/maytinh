package chuc_nang

import (
	"net/http"
	"sort"
	"time"

	"app/mo_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

// Cấu trúc dữ liệu báo cáo gửi ra HTML
type DuLieuDashboard struct {
	TongDoanhThu    float64
	DonHangHomNay   int
	TongSanPham     int
	TongKhachHang   int
	DonHangMoiNhat  []mo_hinh.PhieuXuat
	ChartNhan       []string  // Label ngày (VD: "01/10", "02/10")
	ChartDoanhThu   []float64 // Dữ liệu tiền tương ứng
}

// Handler hiển thị trang Dashboard
func TrangTongQuan(c *gin.Context) {
	// Lấy thông tin Admin đang login
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")
	kh, _ := nghiep_vu.LayThongTinKhachHang(userID)

	// --- TÍNH TOÁN SỐ LIỆU TỪ RAM ---
	stats := tinhToanThongKe()

	c.HTML(http.StatusOK, "quan_tri", gin.H{
		"TieuDe":       "Tổng quan hệ thống",
		"NhanVien":     kh,
		"DaDangNhap":   true,
		"TenNguoiDung": kh.TenKhachHang,
		"QuyenHan":     vaiTro,
		"ThongKe":      stats,
	})
}

// Hàm nội bộ tính toán logic thống kê
func tinhToanThongKe() DuLieuDashboard {
	var kq DuLieuDashboard

	// 1. Lấy dữ liệu thô từ Cache (Có Read Lock)
	listSP := nghiep_vu.LayDanhSachSanPham()
	
	nghiep_vu.KhoaHeThong.RLock()
	defer nghiep_vu.KhoaHeThong.RUnlock()

	kq.TongSanPham = len(listSP)
	kq.TongKhachHang = len(nghiep_vu.CacheKhachHang.DuLieu)

	// 2. Duyệt Phiếu Xuất để tính Doanh thu & Đơn hàng
	now := time.Now().Format("2006-01-02")
	mapDoanhThuNgay := make(map[string]float64)
	var listPX []mo_hinh.PhieuXuat

	// Copy map ra slice để sort
	for _, px := range nghiep_vu.CachePhieuXuat.DuLieu {
		listPX = append(listPX, px)
		
		// Tính tổng doanh thu (Chỉ tính đơn đã hoàn thành hoặc đang xử lý, trừ đơn hủy)
		if px.TrangThai != "Đã hủy" {
			kq.TongDoanhThu += px.TongTienPhieu
			
			// Gom nhóm theo ngày cho biểu đồ (Lấy 10 ký tự đầu yyyy-mm-dd)
			if len(px.NgayTao) >= 10 {
				ngay := px.NgayTao[:10]
				mapDoanhThuNgay[ngay] += px.TongTienPhieu
			}
		}

		// Đếm đơn hôm nay
		if len(px.NgayTao) >= 10 && px.NgayTao[:10] == now {
			kq.DonHangHomNay++
		}
	}

	// 3. Sắp xếp đơn hàng mới nhất lên đầu
	sort.Slice(listPX, func(i, j int) bool {
		return listPX[i].NgayTao > listPX[j].NgayTao
	})

	// Lấy 5 đơn mới nhất
	limit := 5
	if len(listPX) < 5 { limit = len(listPX) }
	kq.DonHangMoiNhat = listPX[:limit]

	// 4. Chuẩn bị dữ liệu biểu đồ (7 ngày gần nhất)
	// Tạo slice 7 ngày qua
	for i := 6; i >= 0; i-- {
		t := time.Now().AddDate(0, 0, -i)
		key := t.Format("2006-01-02")
		label := t.Format("02/01") // dd/mm
		
		kq.ChartNhan = append(kq.ChartNhan, label)
		kq.ChartDoanhThu = append(kq.ChartDoanhThu, mapDoanhThuNgay[key])
	}

	return kq
}

// API Nạp lại dữ liệu (Giữ nguyên code cũ)
func API_NapLaiDuLieu(c *gin.Context) {
	vaiTro := c.GetString("USER_ROLE")
	if !nghiep_vu.KiemTraQuyen(vaiTro, "system.reload") {
		c.JSON(http.StatusForbidden, gin.H{"trang_thai": "loi", "thong_diep": "Không có quyền!"})
		return
	}
	nghiep_vu.KhoiTaoBoNho()
	c.JSON(http.StatusOK, gin.H{"trang_thai": "thanh_cong", "thong_diep": "Đã nạp lại dữ liệu mới nhất từ Google Sheet!"})
}
