package bo_nho_dem

import "app/mo_hinh"

func napPhieuXuat(target *KhoPhieuXuatStore) {
	raw, err := loadSheetData("PHIEU_XUAT")
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotPX_MaPhieuXuat || layString(r, mo_hinh.CotPX_MaPhieuXuat) == "" { continue }
		item := mo_hinh.PhieuXuat{
			MaPhieuXuat:      layString(r, mo_hinh.CotPX_MaPhieuXuat),
			LoaiXuat:         layString(r, mo_hinh.CotPX_LoaiXuat),
			NgayXuat:         layString(r, mo_hinh.CotPX_NgayXuat),
			MaKho:            layString(r, mo_hinh.CotPX_MaKho),
			MaKhachHang:      layString(r, mo_hinh.CotPX_MaKhachHang),
			TrangThai:        layString(r, mo_hinh.CotPX_TrangThai),
			MaVoucher:        layString(r, mo_hinh.CotPX_MaVoucher),
			TienGiamVoucher:  layFloat(r, mo_hinh.CotPX_TienGiamVoucher),
			TongTienPhieu:    layFloat(r, mo_hinh.CotPX_TongTienPhieu),
			DaThu:            layFloat(r, mo_hinh.CotPX_DaThu),
			ConNo:            layFloat(r, mo_hinh.CotPX_ConNo),
			PhuongThucThanhToan:     layString(r, mo_hinh.CotPX_PhuongThucThanhToan),
			PhiVanChuyen:     layFloat(r, mo_hinh.CotPX_PhiVanChuyen),
			NguonDonHang:     layString(r, mo_hinh.CotPX_NguonDonHang),
			ThongTinGiaoHang: layString(r, mo_hinh.CotPX_ThongTinGiaoHang),
			NguoiTao:         layString(r, mo_hinh.CotPX_NguoiTao),
			NgayTao:          layString(r, mo_hinh.CotPX_NgayTao),
		}
		target.DuLieu[item.MaPhieuXuat] = item
		target.DanhSach = append(target.DanhSach, item)
	}
}
