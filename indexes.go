package models

// Hằng số quan trọng: Dòng bắt đầu chứa dữ liệu (0-based index)
// Dòng 11 trong Sheet => Index = 10
const DataStartRow = 10

// --------------------------------------------------------
// INDEX SHEET: SAN_PHAM
// --------------------------------------------------------
const (
	IdxSP_MaSanPham    = 0 // A
	IdxSP_TenSanPham   = 1 // B
	IdxSP_TenRutGon    = 2 // C
	IdxSP_Sku          = 3 // D
	IdxSP_MaDanhMuc    = 4 // E
	IdxSP_MaThuongHieu = 5 // F
	IdxSP_DonVi        = 6 // G
	IdxSP_MauSac       = 7 // H
	IdxSP_UrlHinhAnh   = 8 // I
	IdxSP_ThongSo      = 9 // J
	IdxSP_MoTaChiTiet  = 10 // K
	IdxSP_BaoHanhThang = 11 // L
	IdxSP_TinhTrang    = 12 // M
	IdxSP_TrangThai    = 13 // N
	IdxSP_GiaBanLe     = 14 // O
	IdxSP_GhiChu       = 15 // P
	IdxSP_NguoiTao     = 16 // Q
	IdxSP_NgayTao      = 17 // R
	IdxSP_NgayCapNhat  = 18 // S
)

// --------------------------------------------------------
// INDEX SHEET: PHIEU_NHAP
// --------------------------------------------------------
const (
	IdxPN_MaPhieuNhap    = 0 // A
	IdxPN_MaNCC          = 1 // B
	IdxPN_MaKho          = 2 // C
	IdxPN_NgayNhap       = 3 // D
	IdxPN_TrangThai      = 4 // E
	IdxPN_SoHoaDon       = 5 // F
	IdxPN_NgayHoaDon     = 6 // G
	IdxPN_UrlChungTu     = 7 // H
	IdxPN_TongTienPhieu  = 8 // I
	IdxPN_DaThanhToan    = 9 // J
	IdxPN_ConNo          = 10 // K
	IdxPN_PhuongThucTT   = 11 // L
	IdxPN_TrangThaiTT    = 12 // M
	IdxPN_GhiChu         = 13 // N
	IdxPN_NguoiTao       = 14 // O
	IdxPN_NgayTao        = 15 // P
	IdxPN_NgayCapNhat    = 16 // Q
)

// --------------------------------------------------------
// INDEX SHEET: CHI_TIET_PHIEU_NHAP
// --------------------------------------------------------
const (
	IdxCTPN_MaPhieuNhap   = 0 // A
	IdxCTPN_MaSanPham     = 1 // B
	IdxCTPN_TenSanPham    = 2 // C
	IdxCTPN_DonVi         = 3 // D
	IdxCTPN_SoLuong       = 4 // E
	IdxCTPN_DonGiaNhap    = 5 // F
	IdxCTPN_VatPercent    = 6 // G
	IdxCTPN_GiaSauVat     = 7 // H
	IdxCTPN_ChietKhauDong = 8 // I
	IdxCTPN_ThanhTienDong = 9 // J
	IdxCTPN_GiaVonThucTe  = 10 // K
	IdxCTPN_BaoHanhThang  = 11 // L
	IdxCTPN_GhiChuDong    = 12 // M
)

// --------------------------------------------------------
// INDEX SHEET: NHA_CUNG_CAP
// --------------------------------------------------------
const (
	IdxNCC_MaNCC        = 0 // A
	IdxNCC_TenNCC       = 1 // B
	IdxNCC_DienThoai    = 2 // C
	IdxNCC_Email        = 3 // D
	IdxNCC_DiaChi       = 4 // E
	IdxNCC_MaSoThue     = 5 // F
	IdxNCC_NguoiLienHe  = 6 // G
	IdxNCC_NganHang     = 7 // H
	IdxNCC_NoCanTra     = 8 // I
	IdxNCC_TongMua      = 9 // J
	IdxNCC_HanMucCongNo = 10 // K
	IdxNCC_TrangThai    = 11 // L
	IdxNCC_GhiChu       = 12 // M
	IdxNCC_NguoiTao     = 13 // N
	IdxNCC_NgayTao      = 14 // O
	IdxNCC_NgayCapNhat  = 15 // P
)

// --------------------------------------------------------
// INDEX SHEET: SERIAL_SAN_PHAM
// --------------------------------------------------------
const (
	IdxSerial_Imei                 = 0 // A
	IdxSerial_MaSanPham            = 1 // B
	IdxSerial_MaNCC                = 2 // C
	IdxSerial_MaPhieuNhap          = 3 // D
	IdxSerial_MaPhieuXuat          = 4 // E
	IdxSerial_TrangThai            = 5 // F
	IdxSerial_BaoHanhNCC           = 6 // G
	IdxSerial_HanBaoHanhNCC        = 7 // H
	IdxSerial_MaKhachHangHienTai   = 8 // I
	IdxSerial_NgayXuatKho          = 9 // J
	IdxSerial_KichHoatBaoHanhKhach = 10 // K
	IdxSerial_HanBaoHanhKhach      = 11 // L
	IdxSerial_MaKho                = 12 // M
	IdxSerial_GhiChu               = 13 // N
	IdxSerial_NgayCapNhat          = 14 // O
)

// --------------------------------------------------------
// INDEX SHEET: PHIEU_XUAT
// --------------------------------------------------------
const (
	IdxPX_MaPhieuXuat      = 0 // A
	IdxPX_LoaiXuat         = 1 // B
	IdxPX_NgayXuat         = 2 // C
	IdxPX_MaKho            = 3 // D
	IdxPX_MaKhachHang      = 4 // E
	IdxPX_TrangThai        = 5 // F
	IdxPX_MaVoucher        = 6 // G
	IdxPX_TienGiamVoucher  = 7 // H
	IdxPX_TongTienPhieu    = 8 // I
	IdxPX_LinkChungTu      = 9 // J
	IdxPX_DaThu            = 10 // K
	IdxPX_ConNo            = 11 // L
	IdxPX_PhuongThucTT     = 12 // M
	IdxPX_TrangThaiTT      = 13 // N
	IdxPX_PhiVanChuyen     = 14 // O
	IdxPX_NguonDonHang     = 15 // P
	IdxPX_ThongTinGiaoHang = 16 // Q
	IdxPX_GhiChu           = 17 // R
	IdxPX_NguoiTao         = 18 // S
	IdxPX_NgayTao          = 19 // T
	IdxPX_NgayCapNhat      = 20 // U
)

// --------------------------------------------------------
// INDEX SHEET: CHI_TIET_PHIEU_XUAT
// --------------------------------------------------------
const (
	IdxCTPX_MaPhieuXuat   = 0 // A
	IdxCTPX_MaSanPham     = 1 // B
	IdxCTPX_TenSanPham    = 2 // C
	IdxCTPX_DonVi         = 3 // D
	IdxCTPX_SoLuong       = 4 // E
	IdxCTPX_DonGiaBan     = 5 // F
	IdxCTPX_VatPercent    = 6 // G
	IdxCTPX_GiaSauVat     = 7 // H
	IdxCTPX_ChietKhauDong = 8 // I
	IdxCTPX_ThanhTienDong = 9 // J
	IdxCTPX_GiaVon        = 10 // K
	IdxCTPX_BaoHanhThang  = 11 // L
	IdxCTPX_GhiChuDong    = 12 // M
)

// --------------------------------------------------------
// INDEX SHEET: HOA_DON
// --------------------------------------------------------
const (
	IdxHD_MaHoaDon           = 0 // A
	IdxHD_MaTraCuu           = 1 // B
	IdxHD_XmlUrl             = 2 // C
	IdxHD_LoaiHoaDon         = 3 // D
	IdxHD_MaPhieuXuat        = 4 // E
	IdxHD_MaKhachHang        = 5 // F
	IdxHD_NgayHoaDon         = 6 // G
	IdxHD_KyHieu             = 7 // H
	IdxHD_SoHoaDon           = 8 // I
	IdxHD_MauSo              = 9 // J
	IdxHD_LinkChungTu        = 10 // K
	IdxHD_TongTienPhieu      = 11 // L
	IdxHD_TongVat            = 12 // M
	IdxHD_TongTienSauVat     = 13 // N
	IdxHD_TrangThai          = 14 // O
	IdxHD_TrangThaiThanhToan = 15 // P
	IdxHD_GhiChu             = 16 // Q
	IdxHD_NguoiTao           = 17 // R
	IdxHD_NgayTao            = 18 // S
	IdxHD_NgayCapNhat        = 19 // T
)

// --------------------------------------------------------
// INDEX SHEET: HOA_DON_CHI_TIET
// --------------------------------------------------------
const (
	IdxHDCT_MaHoaDon   = 0 // A
	IdxHDCT_MaSanPham  = 1 // B
	IdxHDCT_TenSanPham = 2 // C
	IdxHDCT_DonVi      = 3 // D
	IdxHDCT_SoLuong    = 4 // E
	IdxHDCT_DonGiaBan  = 5 // F
	IdxHDCT_VatPercent = 6 // G
	IdxHDCT_TienVat    = 7 // H
	IdxHDCT_ThanhTien  = 8 // I
)

// --------------------------------------------------------
// INDEX SHEET: KHACH_HANG
// --------------------------------------------------------
const (
	IdxKH_MaKhachHang   = 0 // A
	IdxKH_UserName      = 1 // B
	IdxKH_PasswordHash  = 2 // C
	IdxKH_LoaiKhachHang = 3 // D
	IdxKH_TenKhachHang  = 4 // E
	IdxKH_DienThoai     = 5 // F
	IdxKH_Email         = 6 // G
	IdxKH_UrlFb         = 7 // H
	IdxKH_Zalo          = 8 // I
	IdxKH_UrlTele       = 9 // J
	IdxKH_UrlTiktok     = 10 // K
	IdxKH_DiaChi        = 11 // L
	IdxKH_NgaySinh      = 12 // M
	IdxKH_GioiTinh      = 13 // N
	IdxKH_MaSoThue      = 14 // O
	IdxKH_DangNo        = 15 // P
	IdxKH_TongMua       = 16 // Q
	IdxKH_TrangThai     = 17 // R
	IdxKH_GhiChu        = 18 // S
	IdxKH_NguoiTao      = 19 // T
	IdxKH_NgayTao       = 20 // U
	IdxKH_NgayCapNhat   = 21 // V
)

// --------------------------------------------------------
// INDEX SHEET: PHIEU_THU_CHI
// --------------------------------------------------------
const (
	IdxPTC_MaPhieuThuChi      = 0 // A
	IdxPTC_NgayTaoPhieu       = 1 // B
	IdxPTC_LoaiPhieu          = 2 // C
	IdxPTC_DoiTuongLoai       = 3 // D
	IdxPTC_DoiTuongID         = 4 // E
	IdxPTC_HangMucThuChi      = 5 // F
	IdxPTC_CoHoaDonDo         = 6 // G
	IdxPTC_MaChungTuThamChieu = 7 // H
	IdxPTC_SoTien             = 8 // I
	IdxPTC_PhuongThucTT       = 9 // J
	IdxPTC_TrangThaiDuyet     = 10 // K
	IdxPTC_NguoiDuyet         = 11 // L
	IdxPTC_GhiChu             = 12 // M
	IdxPTC_NguoiTao           = 13 // N
	IdxPTC_NgayTao            = 14 // O
	IdxPTC_NgayCapNhat        = 15 // P
)

// --------------------------------------------------------
// INDEX SHEET: PHIEU_BAO_HANH
// --------------------------------------------------------
const (
	IdxPBH_MaPhieuBaoHanh    = 0 // A
	IdxPBH_LoaiPhieu         = 1 // B
	IdxPBH_SerialImei        = 2 // C
	IdxPBH_MaSanPham         = 3 // D
	IdxPBH_MaKhachHang       = 4 // E
	IdxPBH_TenNguoiGui       = 5 // F
	IdxPBH_SdtNguoiGui       = 6 // G
	IdxPBH_NgayNhan          = 7 // H
	IdxPBH_TinhTrangLoi      = 8 // I
	IdxPBH_HinhThuc          = 9 // J
	IdxPBH_TrangThai         = 10 // K
	IdxPBH_NgayTraDuKien     = 11 // L
	IdxPBH_NgayTraThucTe     = 12 // M
	IdxPBH_ChiPhiSua         = 13 // N
	IdxPBH_PhiThuKhach       = 14 // O
	IdxPBH_KetQuaSuaChua     = 15 // P
	IdxPBH_LinhKienThayThe   = 16 // Q
	IdxPBH_MaNhanVienKyThuat = 17 // R
	IdxPBH_GhiChu            = 18 // S
	IdxPBH_NguoiTao          = 19 // T
	IdxPBH_NgayTao           = 20 // U
	IdxPBH_NgayCapNhat       = 21 // V
)

// --------------------------------------------------------
// INDEX SHEET: DANH_MUC
// --------------------------------------------------------
const (
	IdxDM_MaDanhMuc    = 0 // A
	IdxDM_ThuTuHienThi = 1 // B
	IdxDM_TenDanhMuc   = 2 // C
	IdxDM_Slug         = 3 // D
	IdxDM_MaDanhMucCha = 4 // E
)

// --------------------------------------------------------
// INDEX SHEET: THUONG_HIEU
// --------------------------------------------------------
const (
	IdxTH_MaThuongHieu  = 0 // A
	IdxTH_TenThuongHieu = 1 // B
	IdxTH_LogoUrl       = 2 // C
)

// --------------------------------------------------------
// INDEX SHEET: NHAN_VIEN
// --------------------------------------------------------
const (
	IdxNV_MaNhanVien      = 0 // A
	IdxNV_TenDangNhap     = 1 // B
	IdxNV_Email           = 2 // C
	IdxNV_MatKhauHash     = 3 // D
	IdxNV_HoTen           = 4 // E
	IdxNV_ChucVu          = 5 // F
	IdxNV_MaPin           = 6 // G
	IdxNV_Cookie          = 7 // H
	IdxNV_CookieExpired   = 8 // I
	IdxNV_VaiTroQuyenHan  = 9 // J
	IdxNV_TrangThai       = 10 // K
	IdxNV_LanDangNhapCuoi = 11 // L
)

// --------------------------------------------------------
// INDEX SHEET: KHUYEN_MAI
// --------------------------------------------------------
const (
	IdxKM_MaVoucher      = 0 // A
	IdxKM_TenChuongTrinh = 1 // B
	IdxKM_LoaiGiam       = 2 // C
	IdxKM_GiaTriGiam     = 3 // D
	IdxKM_DonToThieu     = 4 // E
	IdxKM_NgayBatDau     = 5 // F
	IdxKM_NgayKetThuc    = 6 // G
	IdxKM_SoLuongConLai  = 7 // H
	IdxKM_TrangThai      = 8 // I
)

// --------------------------------------------------------
// INDEX SHEET: CAU_HINH_WEB
// --------------------------------------------------------
const (
	IdxCH_MaCauHinh = 0 // A
	IdxCH_GiaTri    = 1 // B
	IdxCH_MoTa      = 2 // C
	IdxCH_TrangThai = 3 // D
)
