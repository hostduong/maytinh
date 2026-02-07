package bo_nho_dem

import "app/mo_hinh"

func napKhuyenMai(target *KhoKhuyenMaiStore) {
	raw, err := loadSheetData("KHUYEN_MAI")
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotKM_MaVoucher || layString(r, mo_hinh.CotKM_MaVoucher) == "" { continue }
		item := mo_hinh.KhuyenMai{
			MaVoucher:      layString(r, mo_hinh.CotKM_MaVoucher),
			LoaiGiam:       layString(r, mo_hinh.CotKM_LoaiGiam),
			GiaTriGiam:     layFloat(r, mo_hinh.CotKM_GiaTriGiam),
			DonToThieu:     layFloat(r, mo_hinh.CotKM_DonToThieu),
			SoLuongConLai:  layInt(r, mo_hinh.CotKM_SoLuongConLai),
			TrangThai:      layInt(r, mo_hinh.CotKM_TrangThai),
		}
		target.DuLieu[item.MaVoucher] = item
	}
}
