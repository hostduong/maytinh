package bo_nho_dem

import "app/mo_hinh"

func napNhaCungCap(target *KhoNhaCungCapStore) {
	raw, err := loadSheetData("NHA_CUNG_CAP")
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotNCC_MaNhaCungCap || layString(r, mo_hinh.CotNCC_MaNhaCungCap) == "" { continue }
		item := mo_hinh.NhaCungCap{
			MaNhaCungCap:  layString(r, mo_hinh.CotNCC_MaNhaCungCap),
			TenNhaCungCap: layString(r, mo_hinh.CotNCC_TenNhaCungCap),
			DienThoai:     layString(r, mo_hinh.CotNCC_DienThoai),
			Email:         layString(r, mo_hinh.CotNCC_Email),
			DiaChi:        layString(r, mo_hinh.CotNCC_DiaChi),
			NoCanTra:      layFloat(r, mo_hinh.CotNCC_NoCanTra),
			TrangThai:     layInt(r, mo_hinh.CotNCC_TrangThai),
		}
		target.DuLieu[item.MaNhaCungCap] = item
	}
}

func napSerial(target *KhoSerialStore) {
	raw, err := loadSheetData("SERIAL_SAN_PHAM")
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotSerial_SerialImei || layString(r, mo_hinh.CotSerial_SerialImei) == "" { continue }
		item := mo_hinh.SerialSanPham{
			SerialImei:         layString(r, mo_hinh.CotSerial_SerialImei),
			MaSanPham:          layString(r, mo_hinh.CotSerial_MaSanPham),
			MaPhieuNhap:        layString(r, mo_hinh.CotSerial_MaPhieuNhap),
			MaPhieuXuat:        layString(r, mo_hinh.CotSerial_MaPhieuXuat),
			TrangThai:          layInt(r, mo_hinh.CotSerial_TrangThai),
			MaKhachHangHienTai: layString(r, mo_hinh.CotSerial_MaKhachHangHienTai),
		}
		target.DuLieu[item.SerialImei] = item
	}
}
