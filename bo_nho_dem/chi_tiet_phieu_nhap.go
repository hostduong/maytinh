package bo_nho_dem

import "app/mo_hinh"

func napChiTietPhieuNhap(target *KhoChiTietPhieuNhapStore) {
	raw, err := loadSheetData("CHI_TIET_PHIEU_NHAP")
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotCTPN_MaPhieuNhap || layString(r, mo_hinh.CotCTPN_MaPhieuNhap) == "" { continue }
		item := mo_hinh.ChiTietPhieuNhap{
			MaPhieuNhap:   layString(r, mo_hinh.CotCTPN_MaPhieuNhap),
			MaSanPham:     layString(r, mo_hinh.CotCTPN_MaSanPham),
			SoLuong:       layInt(r, mo_hinh.CotCTPN_SoLuong),
			DonGiaNhap:    layFloat(r, mo_hinh.CotCTPN_DonGiaNhap),
			ThanhTienDong: layFloat(r, mo_hinh.CotCTPN_ThanhTienDong),
		}
		target.DanhSach = append(target.DanhSach, item)
	}
}
