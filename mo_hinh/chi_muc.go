package mo_hinh

// DongBatDauDuLieu : Dữ liệu bắt đầu từ dòng 11 trong Excel
const DongBatDauDuLieu = 10

// 1. SAN_PHAM (A-S)
const (
	CotSP_MaSanPham    = 0
	CotSP_TenSanPham   = 1
	CotSP_TenRutGon    = 2
	CotSP_Sku          = 3
	CotSP_MaDanhMuc    = 4
	CotSP_MaThuongHieu = 5
	CotSP_DonVi        = 6
	CotSP_MauSac       = 7
	CotSP_UrlHinhAnh   = 8
	CotSP_ThongSo      = 9
	CotSP_MoTaChiTiet  = 10
	CotSP_BaoHanhThang = 11
	CotSP_TinhTrang    = 12
	CotSP_TrangThai    = 13
	CotSP_GiaBanLe     = 14
	CotSP_GhiChu       = 15
	CotSP_NguoiTao     = 16
	CotSP_NgayTao      = 17
	CotSP_NgayCapNhat  = 18
)

// 2. PHIEU_NHAP (A-Q)
const (
	CotPN_MaPhieuNhap           = 0
	CotPN_MaNhaCungCap          = 1
	CotPN_MaKho                 = 2
	CotPN_NgayNhap              = 3
	CotPN_TrangThai             = 4
	CotPN_SoHoaDon              = 5
	CotPN_NgayHoaDon            = 6
	CotPN_UrlChungTu            = 7
	CotPN_TongTienPhieu         = 8
	CotPN_DaThanhToan           = 9
	CotPN_ConNo                 = 10
	CotPN_PhuongThucThanhToan   = 11
	CotPN_TrangThaiThanhToan    = 12
	CotPN_GhiChu                = 13
	CotPN_NguoiTao              = 14
	CotPN_NgayTao               = 15
	CotPN_NgayCapNhat           = 16
)

// 3. CHI_TIET_PHIEU_NHAP (A-M)
const (
	CotCTPN_MaPhieuNhap    = 0
	CotCTPN_MaSanPham      = 1
	CotCTPN_TenSanPham     = 2
	CotCTPN_DonVi          = 3
	CotCTPN_SoLuong        = 4
	CotCTPN_DonGiaNhap     = 5
	CotCTPN_VatPercent     = 6
	CotCTPN_GiaSauVat      = 7
	CotCTPN_ChietKhauDong  = 8
	CotCTPN_ThanhTienDong  = 9
	CotCTPN_GiaVonThucTe   = 10
	CotCTPN_BaoHanhThang   = 11
	CotCTPN_GhiChuDong     = 12
)

// 4. NHA_CUNG_CAP (A-P)
const (
	CotNCC_MaNhaCungCap  = 0
	CotNCC_TenNhaCungCap = 1
	CotNCC_DienThoai     = 2
	CotNCC_Email         = 3
	CotNCC_DiaChi        = 4
	CotNCC_MaSoThue      = 5
	CotNCC_NguoiLienHe   = 6
	CotNCC_NganHang      = 7
	CotNCC_NoCanTra      = 8
	CotNCC_TongMua       = 9
	CotNCC_HanMucCongNo  = 10
	CotNCC_TrangThai     = 11
	CotNCC_GhiChu        = 12
	CotNCC_NguoiTao      = 13
	CotNCC_NgayTao       = 14
	CotNCC_NgayCapNhat   = 15
)

// 5. SERIAL_SAN_PHAM (A-O)
const (
	CotSerial_SerialImei           = 0
	CotSerial_MaSanPham            = 1
	CotSerial_MaNhaCungCap         = 2
	CotSerial_MaPhieuNhap          = 3
	CotSerial_MaPhieuXuat          = 4
	CotSerial_TrangThai            = 5
	CotSerial_BaoHanhNhaCungCap    = 6
	CotSerial_HanBaoHanhNhaCungCap = 7
	CotSerial_MaKhachHangHienTai   = 8
	CotSerial_NgayXuatKho          = 9
	CotSerial_KichHoatBaoHanhKhach = 10
	CotSerial_HanBaoHanhKhach      = 11
	CotSerial_MaKho                = 12
	CotSerial_GhiChu               = 13
	CotSerial_NgayCapNhat          = 14
)

// 6. PHIEU_XUAT (A-U)
const (
	CotPX_MaPhieuXuat           = 0
	CotPX_LoaiXuat              = 1
	CotPX_NgayXuat              = 2
	CotPX_MaKho                 = 3
	CotPX_MaKhachHang           = 4
	CotPX_TrangThai             = 5
	CotPX_MaVoucher             = 6
	CotPX_TienGiamVoucher       = 7
	CotPX_TongTienPhieu         = 8
	CotPX_LinkChungTu           = 9
	CotPX_DaThu                 = 10
	CotPX_ConNo                 = 11
	CotPX_PhuongThucThanhToan   = 12
	CotPX_TrangThaiThanhToan    = 13
	CotPX_PhiVanChuyen          = 14
	CotPX_NguonDonHang          = 15
	CotPX_ThongTinGiaoHang      = 16
	CotPX_GhiChu                = 17
	CotPX_NguoiTao              = 18
	CotPX_NgayTao               = 19
	CotPX_NgayCapNhat           = 20
)

// 7. CHI_TIET_PHIEU_XUAT (A-M)
const (
	CotCTPX_MaPhieuXuat    = 0
	CotCTPX_MaSanPham      = 1
	CotCTPX_TenSanPham     = 2
	CotCTPX_DonVi          = 3
	CotCTPX_SoLuong        = 4
	CotCTPX_DonGiaBan      = 5
	CotCTPX_VatPercent     = 6
	CotCTPX_GiaSauVat      = 7
	CotCTPX_ChietKhauDong  = 8
	CotCTPX_ThanhTienDong  = 9
	CotCTPX_GiaVon         = 10
	CotCTPX_BaoHanhThang   = 11
	CotCTPX_GhiChuDong     = 12
)

// 8. HOA_DON (A-T)
const (
	CotHD_MaHoaDon            = 0
	CotHD_MaTraCuu            = 1
	CotHD_XmlUrl              = 2
	CotHD_LoaiHoaDon          = 3
	CotHD_MaPhieuXuat         = 4
	CotHD_MaKhachHang         = 5
	CotHD_NgayHoaDon          = 6
	CotHD_KyHieu              = 7
	CotHD_SoHoaDon            = 8
	CotHD_MauSo               = 9
	CotHD_LinkChungTu         = 10
	CotHD_TongTienPhieu       = 11
	CotHD_TongVat             = 12
	CotHD_TongTienSauVat      = 13
	CotHD_TrangThai           = 14
	CotHD_TrangThaiThanhToan  = 15
	CotHD_GhiChu              = 16
	CotHD_NguoiTao            = 17
	CotHD_NgayTao             = 18
	CotHD_NgayCapNhat         = 19
)

// 9. HOA_DON_CHI_TIET (A-I)
const (
	CotHDCT_MaHoaDon   = 0
	CotHDCT_MaSanPham  = 1
	CotHDCT_TenSanPham = 2
	CotHDCT_DonVi      = 3
	CotHDCT_SoLuong    = 4
	CotHDCT_DonGiaBan  = 5
	CotHDCT_VatPercent = 6
	CotHDCT_TienVat    = 7
	CotHDCT_ThanhTien  = 8
)

// 10. KHACH_HANG (A-V)
const (
	CotKH_MaKhachHang      = 0  // A
	CotKH_TenDangNhap      = 1  // B
	CotKH_MatKhauHash      = 2  // C
	CotKH_Cookie           = 3  // D
	CotKH_CookieExpired    = 4  // E
	CotKH_MaPinHash        = 5  // F
	CotKH_LoaiKhachHang    = 6  // G
	CotKH_TenKhachHang     = 7  // H (Họ tên)
	CotKH_DienThoai        = 8  // I
	CotKH_Email            = 9  // J
	CotKH_UrlFb            = 10 // K
	CotKH_Zalo             = 11 // L
	CotKH_UrlTele          = 12 // M
	CotKH_UrlTiktok        = 13 // N
	CotKH_DiaChi           = 14 // O
	CotKH_NgaySinh         = 15 // P
	CotKH_GioiTinh         = 16 // Q
	CotKH_MaSoThue         = 17 // R
	CotKH_DangNo           = 18 // S
	CotKH_TongMua          = 19 // T
	CotKH_ChucVu           = 20 // U
	CotKH_VaiTroQuyenHan   = 21 // V
	CotKH_TrangThai        = 22 // W
	CotKH_GhiChu           = 23 // X
	CotKH_NguoiTao         = 24 // Y
	CotKH_NgayTao          = 25 // Z
	CotKH_NgayCapNhat      = 26 // AA
)

// 11. PHIEU_THU_CHI (A-P)
const (
	CotPTC_MaPhieuThuChi       = 0
	CotPTC_NgayTaoPhieu        = 1
	CotPTC_LoaiPhieu           = 2
	CotPTC_DoiTuongLoai        = 3
	CotPTC_DoiTuongId          = 4
	CotPTC_HangMucThuChi       = 5
	CotPTC_CoHoaDonDo          = 6
	CotPTC_MaChungTuThamChieu  = 7
	CotPTC_SoTien              = 8
	CotPTC_PhuongThucThanhToan = 9
	CotPTC_TrangThaiDuyet      = 10
	CotPTC_NguoiDuyet          = 11
	CotPTC_GhiChu              = 12
	CotPTC_NguoiTao            = 13
	CotPTC_NgayTao             = 14
	CotPTC_NgayCapNhat         = 15
)

// 12. PHIEU_BAO_HANH (A-V)
const (
	CotPBH_MaPhieuBaoHanh    = 0
	CotPBH_LoaiPhieu         = 1
	CotPBH_SerialImei        = 2
	CotPBH_MaSanPham         = 3
	CotPBH_MaKhachHang       = 4
	CotPBH_TenNguoiGui       = 5
	CotPBH_SdtNguoiGui       = 6
	CotPBH_NgayNhan          = 7
	CotPBH_TinhTrangLoi      = 8
	CotPBH_HinhThuc          = 9
	CotPBH_TrangThai         = 10
	CotPBH_NgayTraDuKien     = 11
	CotPBH_NgayTraThucTe     = 12
	CotPBH_ChiPhiSua         = 13
	CotPBH_PhiThuKhach       = 14
	CotPBH_KetQuaSuaChua     = 15
	CotPBH_LinhKienThayThe   = 16
	CotPBH_MaNhanVienKyThuat = 17
	CotPBH_GhiChu            = 18
	CotPBH_NguoiTao          = 19
	CotPBH_NgayTao           = 20
	CotPBH_NgayCapNhat       = 21
)

// 13. DANH_MUC (A-E)
const (
	CotDM_MaDanhMuc    = 0
	CotDM_ThuTuHienThi = 1
	CotDM_TenDanhMuc   = 2
	CotDM_Slug         = 3
	CotDM_MaDanhMucCha = 4
)

// 14. THUONG_HIEU (A-C)
const (
	CotTH_MaThuongHieu  = 0
	CotTH_TenThuongHieu = 1
	CotTH_LogoUrl       = 2
)

// 15. KHUYEN_MAI (A-I)
const (
	CotKM_MaVoucher      = 0
	CotKM_TenChuongTrinh = 1
	CotKM_LoaiGiam       = 2
	CotKM_GiaTriGiam     = 3
	CotKM_DonToThieu     = 4
	CotKM_NgayBatDau     = 5
	CotKM_NgayKetThuc    = 6
	CotKM_SoLuongConLai  = 7
	CotKM_TrangThai      = 8
)

// 16. CAU_HINH_WEB (A-D)
const (
	CotCH_MaCauHinh = 0
	CotCH_GiaTri    = 1
	CotCH_MoTa      = 2
	CotCH_TrangThai = 3
)
