package bo_nho_dem

import "app/mo_hinh"

func napPhieuBaoHanh(target *KhoPhieuBaoHanhStore) {
	raw, err := loadSheetData("PHIEU_BAO_HANH")
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotPBH_MaPhieuBaoHanh || layString(r, mo_hinh.CotPBH_MaPhieuBaoHanh) == "" { continue }
		item := mo_hinh.PhieuBaoHanh{
			MaPhieuBaoHanh: layString(r, mo_hinh.CotPBH_MaPhieuBaoHanh),
			SerialImei:     layString(r, mo_hinh.CotPBH_SerialImei),
			TrangThai:      layInt(r, mo_hinh.CotPBH_TrangThai),
			TinhTrangLoi:   layString(r, mo_hinh.CotPBH_TinhTrangLoi),
		}
		target.DuLieu[item.MaPhieuBaoHanh] = item
		target.DanhSach = append(target.DanhSach, item)
	}
}
