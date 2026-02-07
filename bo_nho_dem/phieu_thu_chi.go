package bo_nho_dem

import "app/mo_hinh"

func napPhieuThuChi(target *KhoPhieuThuChiStore) {
	raw, err := loadSheetData("PHIEU_THU_CHI")
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotPTC_MaPhieuThuChi || layString(r, mo_hinh.CotPTC_MaPhieuThuChi) == "" { continue }
		item := mo_hinh.PhieuThuChi{
			MaPhieuThuChi: layString(r, mo_hinh.CotPTC_MaPhieuThuChi),
			LoaiPhieu:     layString(r, mo_hinh.CotPTC_LoaiPhieu),
			SoTien:        layFloat(r, mo_hinh.CotPTC_SoTien),
			HangMucThuChi: layString(r, mo_hinh.CotPTC_HangMucThuChi),
			TrangThaiDuyet: layInt(r, mo_hinh.CotPTC_TrangThaiDuyet),
		}
		target.DuLieu[item.MaPhieuThuChi] = item
		target.DanhSach = append(target.DanhSach, item)
	}
}
