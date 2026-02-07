package bo_nho_dem

import (
	"strings"
	"app/mo_hinh"
)

// napKhachHang : Hàm chuyên biệt nạp dữ liệu Khách Hàng
func napKhachHang(target *KhoKhachHangStore) {
	raw, err := loadSheetData("KHACH_HANG")
	if err != nil { return }

	for i, r := range raw {
		// Bỏ qua header và dòng trống
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotKH_MaKhachHang || layString(r, mo_hinh.CotKH_MaKhachHang) == "" { continue }
		
		// Map đầy đủ tất cả các cột
		item := &mo_hinh.KhachHang{
			DongTrongSheet: i + 1,
			MaKhachHang:    layString(r, mo_hinh.CotKH_MaKhachHang),
			TenDangNhap:    layString(r, mo_hinh.CotKH_TenDangNhap),
			MatKhauHash:    layString(r, mo_hinh.CotKH_MatKhauHash),
			Cookie:         layString(r, mo_hinh.CotKH_Cookie),
			CookieExpired:  int64(layFloat(r, mo_hinh.CotKH_CookieExpired)),
			MaPinHash:      layString(r, mo_hinh.CotKH_MaPinHash),
			LoaiKhachHang:  layString(r, mo_hinh.CotKH_LoaiKhachHang), // Bạn thiếu cái này
			TenKhachHang:   layString(r, mo_hinh.CotKH_TenKhachHang),
			DienThoai:      layString(r, mo_hinh.CotKH_DienThoai),
			Email:          layString(r, mo_hinh.CotKH_Email),
			UrlFb:          layString(r, mo_hinh.CotKH_UrlFb),
			Zalo:           layString(r, mo_hinh.CotKH_Zalo),
			UrlTele:        layString(r, mo_hinh.CotKH_UrlTele),
			UrlTiktok:      layString(r, mo_hinh.CotKH_UrlTiktok),
			DiaChi:         layString(r, mo_hinh.CotKH_DiaChi),
			NgaySinh:       layString(r, mo_hinh.CotKH_NgaySinh),      // Thiếu
			GioiTinh:       layString(r, mo_hinh.CotKH_GioiTinh),      // Thiếu
			MaSoThue:       layString(r, mo_hinh.CotKH_MaSoThue),
			DangNo:         layFloat(r, mo_hinh.CotKH_DangNo),         // Thiếu
			TongMua:        layFloat(r, mo_hinh.CotKH_TongMua),        // Thiếu
			ChucVu:         layString(r, mo_hinh.CotKH_ChucVu),        // Thiếu
			VaiTroQuyenHan: layString(r, mo_hinh.CotKH_VaiTroQuyenHan),
			TrangThai:      layInt(r, mo_hinh.CotKH_TrangThai),
			GhiChu:         layString(r, mo_hinh.CotKH_GhiChu),        // Thiếu
			NguoiTao:       layString(r, mo_hinh.CotKH_NguoiTao),      // Thiếu
			NgayTao:        layString(r, mo_hinh.CotKH_NgayTao),
			NgayCapNhat:    layString(r, mo_hinh.CotKH_NgayCapNhat),   // Thiếu
		}
		
		// 1. Map ID chính
		target.DuLieu[item.MaKhachHang] = item
		
		// 2. Map User & Email (Chữ thường) để đăng nhập
		if item.TenDangNhap != "" { target.DuLieu[strings.ToLower(item.TenDangNhap)] = item }
		if item.Email != "" { target.DuLieu[strings.ToLower(item.Email)] = item }
		
		// 3. Apped vào danh sách để đếm số lượng thành viên
		target.DanhSach = append(target.DanhSach, item)
	}
}
