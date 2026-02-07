package bo_nho_dem

import "app/mo_hinh"

func napSanPham(target *KhoSanPhamStore) {
	raw, err := loadSheetData("SAN_PHAM")
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotSP_MaSanPham || layString(r, mo_hinh.CotSP_MaSanPham) == "" { continue }
		item := mo_hinh.SanPham{
			MaSanPham:    layString(r, mo_hinh.CotSP_MaSanPham),
			TenSanPham:   layString(r, mo_hinh.CotSP_TenSanPham),
			TenRutGon:    layString(r, mo_hinh.CotSP_TenRutGon),
			Sku:          layString(r, mo_hinh.CotSP_Sku),
			MaDanhMuc:    layString(r, mo_hinh.CotSP_MaDanhMuc),
			MaThuongHieu: layString(r, mo_hinh.CotSP_MaThuongHieu),
			DonVi:        layString(r, mo_hinh.CotSP_DonVi),
			MauSac:       layString(r, mo_hinh.CotSP_MauSac),
			UrlHinhAnh:   layString(r, mo_hinh.CotSP_UrlHinhAnh),
			ThongSo:      layString(r, mo_hinh.CotSP_ThongSo),
			MoTaChiTiet:  layString(r, mo_hinh.CotSP_MoTaChiTiet),
			BaoHanhThang: layInt(r, mo_hinh.CotSP_BaoHanhThang),
			TinhTrang:    layString(r, mo_hinh.CotSP_TinhTrang),
			TrangThai:    layInt(r, mo_hinh.CotSP_TrangThai),
			GiaBanLe:     layFloat(r, mo_hinh.CotSP_GiaBanLe),
			GhiChu:       layString(r, mo_hinh.CotSP_GhiChu),
			NguoiTao:     layString(r, mo_hinh.CotSP_NguoiTao),
			NgayTao:      layString(r, mo_hinh.CotSP_NgayTao),
			NgayCapNhat:  layString(r, mo_hinh.CotSP_NgayCapNhat),
		}
		target.DuLieu[item.MaSanPham] = item
		target.DanhSach = append(target.DanhSach, item)
	}
}
