package bo_nho_dem

import "app/mo_hinh"

func napDanhMuc(target *KhoDanhMucStore) {
	raw, err := loadSheetData("DANH_MUC")
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotDM_MaDanhMuc || layString(r, mo_hinh.CotDM_MaDanhMuc) == "" { continue }
		item := mo_hinh.DanhMuc{
			MaDanhMuc:    layString(r, mo_hinh.CotDM_MaDanhMuc),
			ThuTuHienThi: layInt(r, mo_hinh.CotDM_ThuTuHienThi),
			TenDanhMuc:   layString(r, mo_hinh.CotDM_TenDanhMuc),
			Slug:         layString(r, mo_hinh.CotDM_Slug),
			MaDanhMucCha: layString(r, mo_hinh.CotDM_MaDanhMucCha),
		}
		target.DuLieu[item.MaDanhMuc] = item
	}
}
