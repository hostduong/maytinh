package bo_nho_dem

import "app/mo_hinh"

func napCauHinhWeb(target *KhoCauHinhWebStore) {
	raw, err := loadSheetData("CAU_HINH_WEB")
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotCH_MaCauHinh || layString(r, mo_hinh.CotCH_MaCauHinh) == "" { continue }
		item := mo_hinh.CauHinhWeb{
			MaCauHinh: layString(r, mo_hinh.CotCH_MaCauHinh),
			GiaTri:    layString(r, mo_hinh.CotCH_GiaTri),
			TrangThai: layInt(r, mo_hinh.CotCH_TrangThai),
		}
		target.DuLieu[item.MaCauHinh] = item
	}
}
