package bo_nho_dem

import "app/mo_hinh"

func napPhieuNhap(target *KhoPhieuNhapStore) {
	raw, err := loadSheetData("PHIEU_NHAP")
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotPN_MaPhieuNhap || layString(r, mo_hinh.CotPN_MaPhieuNhap) == "" { continue }
		item := mo_hinh.PhieuNhap{
			MaPhieuNhap:   layString(r, mo_hinh.CotPN_MaPhieuNhap),
			MaNhaCungCap:  layString(r, mo_hinh.CotPN_MaNhaCungCap),
			MaKho:         layString(r, mo_hinh.CotPN_MaKho),
			NgayNhap:      layString(r, mo_hinh.CotPN_NgayNhap),
			TrangThai:     layString(r, mo_hinh.CotPN_TrangThai),
			TongTienPhieu: layFloat(r, mo_hinh.CotPN_TongTienPhieu),
			DaThanhToan:   layFloat(r, mo_hinh.CotPN_DaThanhToan),
			ConNo:         layFloat(r, mo_hinh.CotPN_ConNo),
		}
		target.DuLieu[item.MaPhieuNhap] = item
		target.DanhSach = append(target.DanhSach, item)
	}
}
