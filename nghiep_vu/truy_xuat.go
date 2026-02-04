package nghiep_vu

import (
	"app/mo_hinh"
)

// ========================================================
// CÁC HÀM ĐỌC DỮ LIỆU AN TOÀN (SAFE READ)
// ========================================================

// LayDanhSachSanPham : Lấy toàn bộ SP đang có trong RAM
func LayDanhSachSanPham() []mo_hinh.SanPham {
	// 1. Lấy khóa
	khoa := BoQuanLyKhoa.LayKhoa(CacheSanPham.TenKey)
	
	// 2. Khóa ĐỌC (RLock) - Nhiều người xem cùng lúc được
	khoa.RLock()
	defer khoa.RUnlock()

	// 3. Copy ra mảng mới để trả về (Tránh lỗi tham chiếu vùng nhớ)
	ketQua := make([]mo_hinh.SanPham, len(CacheSanPham.DanhSach))
	copy(ketQua, CacheSanPham.DanhSach)

	return ketQua
}

// LayChiTietSanPham : Tìm SP theo ID
func LayChiTietSanPham(maSP string) (mo_hinh.SanPham, bool) {
	khoa := BoQuanLyKhoa.LayKhoa(CacheSanPham.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()

	sp, tonTai := CacheSanPham.DuLieu[maSP]
	return sp, tonTai
}

// LayDanhSachDanhMuc : Lấy menu danh mục
func LayDanhSachDanhMuc() map[string]mo_hinh.DanhMuc {
	khoa := BoQuanLyKhoa.LayKhoa(CacheDanhMuc.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()

	// Copy map
	ketQua := make(map[string]mo_hinh.DanhMuc)
	for k, v := range CacheDanhMuc.DuLieu {
		ketQua[k] = v
	}
	return ketQua
}

// LayCauHinhWeb : Lấy banner, cấu hình
func LayCauHinhWeb() map[string]mo_hinh.CauHinhWeb {
	khoa := BoQuanLyKhoa.LayKhoa(CacheCauHinhWeb.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()

	ketQua := make(map[string]mo_hinh.CauHinhWeb)
	for k, v := range CacheCauHinhWeb.DuLieu {
		ketQua[k] = v
	}
	return ketQua
}

// (Bạn có thể viết thêm các hàm lấy Khách hàng, Đơn hàng tương tự...)
