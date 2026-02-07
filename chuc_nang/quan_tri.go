package chuc_nang

import (
	"net/http"
	"sort"
	"time"

	"app/mo_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

type DuLieuDashboard struct {
	TongDoanhThu    float64
	DonHangHomNay   int
	TongSanPham     int
	TongKhachHang   int
	DonHangMoiNhat  []mo_hinh.PhieuXuat
	ChartNhan       []string
	ChartDoanhThu   []float64
}

func TrangTongQuan(c *gin.Context) {
	userID := c.GetString("USER_ID")
	vaiTro := c.GetString("USER_ROLE")
	kh, _ := nghiep_vu.LayThongTinKhachHang(userID)

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

func tinhToanThongKe() DuLieuDashboard {
	var kq DuLieuDashboard

	nghiep_vu.KhoaHeThong.RLock()
	defer nghiep_vu.KhoaHeThong.RUnlock()

	// 1. Đếm sản phẩm (Dựa trên danh sách)
	kq.TongSanPham = len(nghiep_vu.CacheSanPham.DanhSach)
	
	// 2. Đếm khách hàng (Dựa trên danh sách mới tạo -> ĐẾM ĐÚNG)
	kq.TongKhachHang = len(nghiep_vu.CacheKhachHang.DanhSach)

	now := time.Now().Format("2006-01-02")
	mapDoanhThuNgay := make(map[string]float64)
	var listPX []mo_hinh.PhieuXuat

	for _, px := range nghiep_vu.CachePhieuXuat.DuLieu {
		listPX = append(listPX, px)
		
		if px.TrangThai != "Đã hủy" {
			kq.TongDoanhThu += px.TongTienPhieu
			if len(px.NgayTao) >= 10 {
				ngay := px.NgayTao[:10]
				mapDoanhThuNgay[ngay] += px.TongTienPhieu
			}
		}

		if len(px.NgayTao) >= 10 && px.NgayTao[:10] == now {
			kq.DonHangHomNay++
		}
	}

	sort.Slice(listPX, func(i, j int) bool {
		return listPX[i].NgayTao > listPX[j].NgayTao
	})

	limit := 5
	if len(listPX) < 5 { limit = len(listPX) }
	kq.DonHangMoiNhat = listPX[:limit]

	for i := 6; i >= 0; i-- {
		t := time.Now().AddDate(0, 0, -i)
		key := t.Format("2006-01-02")
		label := t.Format("02/01")
		
		kq.ChartNhan = append(kq.ChartNhan, label)
		kq.ChartDoanhThu = append(kq.ChartDoanhThu, mapDoanhThuNgay[key])
	}

	return kq
}

func API_NapLaiDuLieu(c *gin.Context) {
	vaiTro := c.GetString("USER_ROLE")
	if !nghiep_vu.KiemTraQuyen(vaiTro, "system.reload") {
		c.JSON(http.StatusForbidden, gin.H{"trang_thai": "loi", "thong_diep": "Không có quyền!"})
		return
	}

	nghiep_vu.LamMoiHeThong()

	c.JSON(http.StatusOK, gin.H{
		"trang_thai": "thanh_cong", 
		"thong_diep": "Đã đồng bộ dữ liệu mới nhất!",
	})
}
