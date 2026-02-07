package bo_nho_dem

import "app/mo_hinh"

func napThuongHieu(target *KhoThuongHieuStore) {
	raw, err := loadSheetData("THUONG_HIEU")
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotTH_MaThuongHieu || layString(r, mo_hinh.CotTH_MaThuongHieu) == "" { continue }
		item := mo_hinh.ThuongHieu{
			MaThuongHieu:  layString(r, mo_hinh.CotTH_MaThuongHieu),
			TenThuongHieu: layString(r, mo_hinh.CotTH_TenThuongHieu),
			LogoUrl:       layString(r, mo_hinh.CotTH_LogoUrl),
		}
		target.DuLieu[item.MaThuongHieu] = item
	}
}
