package bo_nho_dem

import "app/mo_hinh"

func napChiTietPhieuXuat(target *KhoChiTietPhieuXuatStore) {
	raw, err := loadSheetData("CHI_TIET_PHIEU_XUAT")
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotCTPX_MaPhieuXuat || layString(r, mo_hinh.CotCTPX_MaPhieuXuat) == "" { continue }
		item := mo_hinh.ChiTietPhieuXuat{
			MaPhieuXuat:   layString(r, mo_hinh.CotCTPX_MaPhieuXuat),
			MaSanPham:     layString(r, mo_hinh.CotCTPX_MaSanPham),
			TenSanPham:    layString(r, mo_hinh.CotCTPX_TenSanPham),
			SoLuong:       layInt(r, mo_hinh.CotCTPX_SoLuong),
			DonGiaBan:     layFloat(r, mo_hinh.CotCTPX_DonGiaBan),
			ThanhTienDong: layFloat(r, mo_hinh.CotCTPX_ThanhTienDong),
			GiaVon:        layFloat(r, mo_hinh.CotCTPX_GiaVon),
		}
		target.DanhSach = append(target.DanhSach, item)
	}
}
