package bo_nho_dem

import "app/mo_hinh"

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
