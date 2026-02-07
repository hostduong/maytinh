package bo_nho_dem

import "app/mo_hinh"

func napHoaDon(target *KhoHoaDonStore) {
	raw, err := loadSheetData("HOA_DON")
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotHD_MaHoaDon || layString(r, mo_hinh.CotHD_MaHoaDon) == "" { continue }
		item := mo_hinh.HoaDon{
			MaHoaDon:      layString(r, mo_hinh.CotHD_MaHoaDon),
			MaTraCuu:      layString(r, mo_hinh.CotHD_MaTraCuu),
			XmlUrl:        layString(r, mo_hinh.CotHD_XmlUrl),
			MaPhieuXuat:   layString(r, mo_hinh.CotHD_MaPhieuXuat),
			TongTienPhieu: layFloat(r, mo_hinh.CotHD_TongTienPhieu),
			TongVat:       layFloat(r, mo_hinh.CotHD_TongVat),
			TrangThai:     layString(r, mo_hinh.CotHD_TrangThai),
		}
		target.DuLieu[item.MaHoaDon] = item
	}
}
