package nghiep_vu

import (
	"fmt"
	"log"
	"strconv"
	"app/kho_du_lieu" // Import kho dữ liệu
	"app/mo_hinh"     // Import model struct
)

// ========================================================
// BIẾN TOÀN CỤC (CACHE) - LƯU TRỮ DỮ LIỆU TRÊN RAM
// ========================================================
var (
	// Map để tra cứu nhanh theo ID (Key = Mã, Value = Struct)
	KhoSanPham    map[string]mo_hinh.SanPham
	KhoDanhMuc    map[string]mo_hinh.DanhMuc
	KhoKhachHang  map[string]mo_hinh.KhachHang
	KhoDonHang    map[string]mo_hinh.PhieuXuat
	
	// Slice để duyệt danh sách (Dùng khi cần loop)
	DanhSachSanPham []mo_hinh.SanPham
)

// KhoiTaoBoNho : Hàm này chạy 1 lần khi khởi động để nạp dữ liệu
func KhoiTaoBoNho() {
	log.Println("--- [RAM] Đang nạp dữ liệu vào bộ nhớ đệm... ---")
	
	// Khởi tạo Map
	KhoSanPham = make(map[string]mo_hinh.SanPham)
	KhoDanhMuc = make(map[string]mo_hinh.DanhMuc)
	KhoKhachHang = make(map[string]mo_hinh.KhachHang)
	KhoDonHang = make(map[string]mo_hinh.PhieuXuat)

	// 1. Nạp Danh Mục
	napDanhMuc()
	
	// 2. Nạp Sản Phẩm
	napSanPham()

	// ... (Sẽ nạp tiếp các bảng khác sau)

	log.Println("--- [RAM] Hoàn tất nạp dữ liệu ---")
}

// --------------------------------------------------------
// HÀM HỖ TRỢ NẠP DỮ LIỆU (PRIVATE)
// --------------------------------------------------------

func napSanPham() {
	duLieuTho, loi := kho_du_lieu.DocToanBoSheet("SAN_PHAM")
	if loi != nil {
		log.Println("Lỗi đọc SAN_PHAM:", loi)
		return
	}

	count := 0
	// Bắt đầu đọc từ dòng quy định trong chi_muc.go
	startRow := mo_hinh.DongBatDauDuLieu 

	for i, dong := range duLieuTho {
		if i < startRow {
			continue // Bỏ qua tiêu đề
		}
		
		// Kiểm tra dòng có dữ liệu không (Cột Mã SP phải có)
		if len(dong) <= mo_hinh.CotSP_MaSanPham || dong[mo_hinh.CotSP_MaSanPham] == "" {
			continue
		}

		// Map dữ liệu từ mảng sang Struct
		sp := mo_hinh.SanPham{
			MaSanPham:    layString(dong, mo_hinh.CotSP_MaSanPham),
			TenSanPham:   layString(dong, mo_hinh.CotSP_TenSanPham),
			TenRutGon:    layString(dong, mo_hinh.CotSP_TenRutGon),
			Sku:          layString(dong, mo_hinh.CotSP_Sku),
			MaDanhMuc:    layString(dong, mo_hinh.CotSP_MaDanhMuc),
			DonVi:        layString(dong, mo_hinh.CotSP_DonVi),
			UrlHinhAnh:   layString(dong, mo_hinh.CotSP_UrlHinhAnh),
			GiaBanLe:     layFloat(dong, mo_hinh.CotSP_GiaBanLe),
			TrangThai:    layInt(dong, mo_hinh.CotSP_TrangThai),
			// ... (Bạn có thể map thêm các trường khác nếu cần dùng ngay)
		}

		// Lưu vào RAM
		KhoSanPham[sp.MaSanPham] = sp
		DanhSachSanPham = append(DanhSachSanPham, sp)
		count++
	}
	log.Printf("-> Đã nạp %d Sản Phẩm", count)
}

func napDanhMuc() {
	duLieuTho, loi := kho_du_lieu.DocToanBoSheet("DANH_MUC")
	if loi != nil {
		log.Println("Lỗi đọc DANH_MUC:", loi)
		return
	}

	count := 0
	startRow := mo_hinh.DongBatDauDuLieu

	for i, dong := range duLieuTho {
		if i < startRow { continue }
		if len(dong) <= mo_hinh.CotDM_MaDanhMuc || dong[mo_hinh.CotDM_MaDanhMuc] == "" { continue }

		dm := mo_hinh.DanhMuc{
			MaDanhMuc:    layString(dong, mo_hinh.CotDM_MaDanhMuc),
			TenDanhMuc:   layString(dong, mo_hinh.CotDM_TenDanhMuc),
			Slug:         layString(dong, mo_hinh.CotDM_Slug),
			MaDanhMucCha: layString(dong, mo_hinh.CotDM_MaDanhMucCha),
		}
		
		KhoDanhMuc[dm.MaDanhMuc] = dm
		count++
	}
	log.Printf("-> Đã nạp %d Danh Mục", count)
}

// --------------------------------------------------------
// CÁC HÀM TIỆN ÍCH CHUYỂN ĐỔI KIỂU DỮ LIỆU (Helper)
// --------------------------------------------------------

func layString(dong []interface{}, index int) string {
	if index >= len(dong) || dong[index] == nil {
		return ""
	}
	return fmt.Sprintf("%v", dong[index])
}

func layFloat(dong []interface{}, index int) float64 {
	str := layString(dong, index)
	if str == "" { return 0 }
	val, _ := strconv.ParseFloat(str, 64)
	return val
}

func layInt(dong []interface{}, index int) int {
	str := layString(dong, index)
	if str == "" { return 0 }
	val, _ := strconv.Atoi(str)
	return val
}
