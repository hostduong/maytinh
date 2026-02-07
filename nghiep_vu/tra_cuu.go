package nghiep_vu

import (
	"sort"
	"app/bo_nho_dem"
	"app/mo_hinh"
)

// --- CÁC HÀM CŨ (Giữ nguyên) ---
// func LayDanhSachSanPham() ...

// --- THÊM MỚI 2 HÀM NÀY ---

func LayDanhSachDanhMuc() []mo_hinh.DanhMuc {
	// 1. Lock đọc để an toàn
	bo_nho_dem.KhoaHeThong.RLock()
	defer bo_nho_dem.KhoaHeThong.RUnlock()

	var danhSach []mo_hinh.DanhMuc
	
	// 2. Quét Map trong RAM
	for _, dm := range bo_nho_dem.CacheDanhMuc.DuLieu {
		danhSach = append(danhSach, dm)
	}

	// 3. Sắp xếp theo thứ tự hiển thị (Nếu có) hoặc tên
	sort.Slice(danhSach, func(i, j int) bool {
		// Ưu tiên ThuTuHienThi nhỏ lên trước
		if danhSach[i].ThuTuHienThi != danhSach[j].ThuTuHienThi {
			return danhSach[i].ThuTuHienThi < danhSach[j].ThuTuHienThi
		}
		return danhSach[i].TenDanhMuc < danhSach[j].TenDanhMuc
	})

	return danhSach
}

func LayDanhSachThuongHieu() []mo_hinh.ThuongHieu {
	bo_nho_dem.KhoaHeThong.RLock()
	defer bo_nho_dem.KhoaHeThong.RUnlock()

	var danhSach []mo_hinh.ThuongHieu
	for _, th := range bo_nho_dem.CacheThuongHieu.DuLieu {
		danhSach = append(danhSach, th)
	}

	sort.Slice(danhSach, func(i, j int) bool {
		return danhSach[i].TenThuongHieu < danhSach[j].TenThuongHieu
	})

	return danhSach
}
