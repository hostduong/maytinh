package nghiep_vu

import (
	"app/mo_hinh"
)

// =================================================================================
// NHÓM 1: MASTER DATA (Sản phẩm, Danh mục, Đối tác, Nhân viên...)
// =================================================================================

// LayDanhSachSanPham : Lấy tất cả sản phẩm
func LayDanhSachSanPham() []mo_hinh.SanPham {
	khoa := BoQuanLyKhoa.LayKhoa(CacheSanPham.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()

	ketQua := make([]mo_hinh.SanPham, len(CacheSanPham.DanhSach))
	copy(ketQua, CacheSanPham.DanhSach)
	return ketQua
}

// LayChiTietSanPham : Lấy 1 sản phẩm theo ID
func LayChiTietSanPham(maSP string) (mo_hinh.SanPham, bool) {
	khoa := BoQuanLyKhoa.LayKhoa(CacheSanPham.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()
	sp, tonTai := CacheSanPham.DuLieu[maSP]
	return sp, tonTai
}

func LayDanhSachDanhMuc() map[string]mo_hinh.DanhMuc {
	khoa := BoQuanLyKhoa.LayKhoa(CacheDanhMuc.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()
	
	kq := make(map[string]mo_hinh.DanhMuc)
	for k, v := range CacheDanhMuc.DuLieu { kq[k] = v }
	return kq
}

func LayDanhSachThuongHieu() map[string]mo_hinh.ThuongHieu {
	khoa := BoQuanLyKhoa.LayKhoa(CacheThuongHieu.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()
	
	kq := make(map[string]mo_hinh.ThuongHieu)
	for k, v := range CacheThuongHieu.DuLieu { kq[k] = v }
	return kq
}

func LayDanhSachNhaCungCap() map[string]mo_hinh.NhaCungCap {
	khoa := BoQuanLyKhoa.LayKhoa(CacheNhaCungCap.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()
	
	kq := make(map[string]mo_hinh.NhaCungCap)
	for k, v := range CacheNhaCungCap.DuLieu { kq[k] = v }
	return kq
}

func LayThongTinKhachHang(maKH string) (mo_hinh.KhachHang, bool) {
	khoa := BoQuanLyKhoa.LayKhoa(CacheKhachHang.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()
	kh, tonTai := CacheKhachHang.DuLieu[maKH]
	return kh, tonTai
}

func LayThongTinNhanVien(maNV string) (mo_hinh.NhanVien, bool) {
	khoa := BoQuanLyKhoa.LayKhoa(CacheNhanVien.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()
	nv, tonTai := CacheNhanVien.DuLieu[maNV]
	return nv, tonTai
}

func LayCauHinhWeb() map[string]mo_hinh.CauHinhWeb {
	khoa := BoQuanLyKhoa.LayKhoa(CacheCauHinhWeb.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()
	
	kq := make(map[string]mo_hinh.CauHinhWeb)
	for k, v := range CacheCauHinhWeb.DuLieu { kq[k] = v }
	return kq
}

// =================================================================================
// NHÓM 2: GIAO DỊCH BÁN HÀNG (Phiếu Xuất & Voucher)
// =================================================================================

func LayThongTinDonHang(maPX string) (mo_hinh.PhieuXuat, bool) {
	khoa := BoQuanLyKhoa.LayKhoa(CachePhieuXuat.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()
	px, tonTai := CachePhieuXuat.DuLieu[maPX]
	return px, tonTai
}

// LayChiTietDonHang : Lọc ra các dòng sản phẩm thuộc mã phiếu này
func LayChiTietDonHang(maPX string) []mo_hinh.ChiTietPhieuXuat {
	khoa := BoQuanLyKhoa.LayKhoa(CacheChiTietXuat.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()

	var ketQua []mo_hinh.ChiTietPhieuXuat
	for _, ct := range CacheChiTietXuat.DanhSach {
		if ct.MaPhieuXuat == maPX {
			ketQua = append(ketQua, ct)
		}
	}
	return ketQua
}

func LayThongTinVoucher(maVoucher string) (mo_hinh.KhuyenMai, bool) {
	khoa := BoQuanLyKhoa.LayKhoa(CacheKhuyenMai.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()
	km, tonTai := CacheKhuyenMai.DuLieu[maVoucher]
	return km, tonTai
}

// =================================================================================
// NHÓM 3: GIAO DỊCH NHẬP KHO (Phiếu Nhập)
// =================================================================================

func LayThongTinPhieuNhap(maPN string) (mo_hinh.PhieuNhap, bool) {
	khoa := BoQuanLyKhoa.LayKhoa(CachePhieuNhap.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()
	pn, tonTai := CachePhieuNhap.DuLieu[maPN]
	return pn, tonTai
}

func LayChiTietPhieuNhap(maPN string) []mo_hinh.ChiTietPhieuNhap {
	khoa := BoQuanLyKhoa.LayKhoa(CacheChiTietNhap.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()

	var ketQua []mo_hinh.ChiTietPhieuNhap
	for _, ct := range CacheChiTietNhap.DanhSach {
		if ct.MaPhieuNhap == maPN {
			ketQua = append(ketQua, ct)
		}
	}
	return ketQua
}

// =================================================================================
// NHÓM 4: KHO VẬN & BẢO HÀNH (Serial, Bảo hành)
// =================================================================================

func TraCuuSerial(imei string) (mo_hinh.SerialSanPham, bool) {
	khoa := BoQuanLyKhoa.LayKhoa(CacheSerial.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()
	serial, tonTai := CacheSerial.DuLieu[imei]
	return serial, tonTai
}

func LayThongTinBaoHanh(maPBH string) (mo_hinh.PhieuBaoHanh, bool) {
	khoa := BoQuanLyKhoa.LayKhoa(CachePhieuBaoHanh.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()
	pbh, tonTai := CachePhieuBaoHanh.DuLieu[maPBH]
	return pbh, tonTai
}

// =================================================================================
// NHÓM 5: TÀI CHÍNH (Hóa đơn, Thu Chi)
// =================================================================================

func LayThongTinHoaDon(maHD string) (mo_hinh.HoaDon, bool) {
	khoa := BoQuanLyKhoa.LayKhoa(CacheHoaDon.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()
	hd, tonTai := CacheHoaDon.DuLieu[maHD]
	return hd, tonTai
}

func LayChiTietHoaDon(maHD string) []mo_hinh.HoaDonChiTiet {
	khoa := BoQuanLyKhoa.LayKhoa(CacheHoaDonChiTiet.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()

	var ketQua []mo_hinh.HoaDonChiTiet
	for _, ct := range CacheHoaDonChiTiet.DanhSach {
		if ct.MaHoaDon == maHD {
			ketQua = append(ketQua, ct)
		}
	}
	return ketQua
}

func LayPhieuThuChi(maPTC string) (mo_hinh.PhieuThuChi, bool) {
	khoa := BoQuanLyKhoa.LayKhoa(CachePhieuThuChi.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()
	ptc, tonTai := CachePhieuThuChi.DuLieu[maPTC]
	return ptc, tonTai
}

// LayDanhSachThuChi : Lấy toàn bộ sổ quỹ (Cẩn thận nếu dữ liệu lớn)
func LayDanhSachThuChi() []mo_hinh.PhieuThuChi {
	khoa := BoQuanLyKhoa.LayKhoa(CachePhieuThuChi.TenKey)
	khoa.RLock()
	defer khoa.RUnlock()

	ketQua := make([]mo_hinh.PhieuThuChi, len(CachePhieuThuChi.DanhSach))
	copy(ketQua, CachePhieuThuChi.DanhSach)
	return ketQua
}
