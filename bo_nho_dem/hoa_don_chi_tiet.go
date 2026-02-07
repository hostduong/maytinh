package bo_nho_dem

import "app/mo_hinh"

func napHoaDonChiTiet(target *KhoHoaDonChiTietStore) {
	raw, err := loadSheetData("HOA_DON_CHI_TIET")
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotHDCT_MaHoaDon || layString(r, mo_hinh.CotHDCT_MaHoaDon) == "" { continue }
		item := mo_hinh.HoaDonChiTiet{
			MaHoaDon:   layString(r, mo_hinh.CotHDCT_MaHoaDon),
			MaSanPham:  layString(r, mo_hinh.CotHDCT_MaSanPham),
			SoLuong:    layInt(r, mo_hinh.CotHDCT_SoLuong),
			ThanhTien:  layFloat(r, mo_hinh.CotHDCT_ThanhTien),
		}
		target.DanhSach = append(target.DanhSach, item)
	}
}
