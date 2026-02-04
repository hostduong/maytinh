package mo_hinh

// GHI CHÚ CHUNG:
// - PK (Primary Key): Khóa chính.
// - FK (Foreign Key): Khóa ngoại.

// =================================================================================
// 1. SẢN PHẨM (Master Data)
// =================================================================================
type SanPham struct {
	MaSanPham    string  `json:"ma_san_pham"`    // Col A
	TenSanPham   string  `json:"ten_san_pham"`   // Col B
	TenRutGon    string  `json:"ten_rut_gon"`    // Col C
	Sku          string  `json:"sku"`            // Col D
	MaDanhMuc    string  `json:"ma_danh_muc"`    // Col E
	MaThuongHieu string  `json:"ma_thuong_hieu"` // Col F
	DonVi        string  `json:"don_vi"`         // Col G
	MauSac       string  `json:"mau_sac"`        // Col H
	UrlHinhAnh   string  `json:"url_hinh_anh"`   // Col I
	ThongSo      string  `json:"thong_so"`       // Col J
	MoTaChiTiet  string  `json:"mo_ta_chi_tiet"` // Col K
	BaoHanhThang int     `json:"bao_hanh_thang"` // Col L
	TinhTrang    string  `json:"tinh_trang"`     // Col M
	TrangThai    int     `json:"trang_thai"`     // Col N
	GiaBanLe     float64 `json:"gia_ban_le"`     // Col O
	GhiChu       string  `json:"ghi_chu"`        // Col P
	NguoiTao     string  `json:"nguoi_tao"`      // Col Q
	NgayTao      string  `json:"ngay_tao"`       // Col R
	NgayCapNhat  string  `json:"ngay_cap_nhat"`  // Col S
	// KhoiLuong int `json:"khoi_luong"` // (Giai đoạn 2)
}

// =================================================================================
// 2. PHIẾU NHẬP
// =================================================================================
type PhieuNhap struct {
	MaPhieuNhap        string  `json:"ma_phieu_nhap"`
	MaNhaCungCap       string  `json:"ma_nha_cung_cap"`
	MaKho              string  `json:"ma_kho"`
	NgayNhap           string  `json:"ngay_nhap"`
	TrangThai          string  `json:"trang_thai"`
	SoHoaDon           string  `json:"so_hoa_don"`
	NgayHoaDon         string  `json:"ngay_hoa_don"`
	UrlChungTu         string  `json:"url_chung_tu"`
	TongTienPhieu      float64 `json:"tong_tien_phieu"`
	DaThanhToan        float64 `json:"da_thanh_toan"`
	ConNo              float64 `json:"con_no"`
	PhuongThucTT       string  `json:"phuong_thuc_thanh_toan"`
	TrangThaiTT        string  `json:"trang_thai_thanh_toan"`
	GhiChu             string  `json:"ghi_chu"`
	NguoiTao           string  `json:"nguoi_tao"`
	NgayTao            string  `json:"ngay_tao"`
	NgayCapNhat        string  `json:"ngay_cap_nhat"`
}

type ChiTietPhieuNhap struct {
	MaPhieuNhap    string  `json:"ma_phieu_nhap"`
	MaSanPham      string  `json:"ma_san_pham"`
	TenSanPham     string  `json:"ten_san_pham"`
	DonVi          string  `json:"don_vi"`
	SoLuong        int     `json:"so_luong"`
	DonGiaNhap     float64 `json:"don_gia_nhap"`
	VatPercent     float64 `json:"vat_percent"`
	GiaSauVat      float64 `json:"gia_sau_vat"`
	ChietKhauDong  float64 `json:"chiet_khau_dong"`
	ThanhTienDong  float64 `json:"thanh_tien_dong"`
	GiaVonThucTe   float64 `json:"gia_von_thuc_te"`
	BaoHanhThang   int     `json:"bao_hanh_thang"`
	GhiChuDong     string  `json:"ghi_chu_dong"`
}

// =================================================================================
// 3. NHÀ CUNG CẤP
// =================================================================================
type NhaCungCap struct {
	MaNhaCungCap string  `json:"ma_nha_cung_cap"`
	TenNhaCungCap string `json:"ten_nha_cung_cap"`
	DienThoai    string  `json:"dien_thoai"`
	Email        string  `json:"email"`
	DiaChi       string  `json:"dia_chi"`
	MaSoThue     string  `json:"ma_so_thue"`
	NguoiLienHe  string  `json:"nguoi_lien_he"`
	NganHang     string  `json:"ngan_hang"`
	NoCanTra     float64 `json:"no_can_tra"`
	TongMua      float64 `json:"tong_mua"`
	HanMucCongNo float64 `json:"han_muc_cong_no"`
	TrangThai    int     `json:"trang_thai"`
	GhiChu       string  `json:"ghi_chu"`
	NguoiTao     string  `json:"nguoi_tao"`
	NgayTao      string  `json:"ngay_tao"`
	NgayCapNhat  string  `json:"ngay_cap_nhat"`
}

// =================================================================================
// 4. SERIAL (Theo dõi bảo hành)
// =================================================================================
type SerialSanPham struct {
	SerialImei           string `json:"serial_imei"`
	MaSanPham            string `json:"ma_san_pham"`
	MaNhaCungCap         string `json:"ma_nha_cung_cap"`
	MaPhieuNhap          string `json:"ma_phieu_nhap"`
	MaPhieuXuat          string `json:"ma_phieu_xuat"`
	TrangThai            int    `json:"trang_thai"`
	BaoHanhNhaCungCap    int    `json:"bao_hanh_nha_cung_cap"`
	HanBaoHanhNhaCungCap string `json:"han_bao_hanh_nha_cung_cap"`
	MaKhachHangHienTai   string `json:"ma_khach_hang_hien_tai"`
	NgayXuatKho          string `json:"ngay_xuat_kho"`
	KichHoatBaoHanhKhach string `json:"kich_hoat_bao_hanh_khach"`
	HanBaoHanhKhach      string `json:"han_bao_hanh_khach"`
	MaKho                string `json:"ma_kho"`
	GhiChu               string `json:"ghi_chu"`
	NgayCapNhat          string `json:"ngay_cap_nhat"`
}

// =================================================================================
// 5. PHIẾU XUẤT (Đơn hàng)
// =================================================================================
type PhieuXuat struct {
	MaPhieuXuat        string  `json:"ma_phieu_xuat"`
	LoaiXuat           string  `json:"loai_xuat"`
	NgayXuat           string  `json:"ngay_xuat"`
	MaKho              string  `json:"ma_kho"`
	MaKhachHang        string  `json:"ma_khach_hang"`
	TrangThai          string  `json:"trang_thai"`
	MaVoucher          string  `json:"ma_voucher"`
	TienGiamVoucher    float64 `json:"tien_giam_voucher"`
	TongTienPhieu      float64 `json:"tong_tien_phieu"`
	LinkChungTu        string  `json:"link_chung_tu"`
	DaThu              float64 `json:"da_thu"`
	ConNo              float64 `json:"con_no"`
	PhuongThucTT       string  `json:"phuong_thuc_thanh_toan"`
	TrangThaiTT        string  `json:"trang_thai_thanh_toan"`
	PhiVanChuyen       float64 `json:"phi_van_chuyen"`
	NguonDonHang       string  `json:"nguon_don_hang"`
	ThongTinGiaoHang   string  `json:"thong_tin_giao_hang"`
	GhiChu             string  `json:"ghi_chu"`
	NguoiTao           string  `json:"nguoi_tao"`
	NgayTao            string  `json:"ngay_tao"`
	NgayCapNhat        string  `json:"ngay_cap_nhat"`
}

type ChiTietPhieuXuat struct {
	MaPhieuXuat   string  `json:"ma_phieu_xuat"`
	MaSanPham     string  `json:"ma_san_pham"`
	TenSanPham    string  `json:"ten_san_pham"`
	DonVi         string  `json:"don_vi"`
	SoLuong       int     `json:"so_luong"`
	DonGiaBan     float64 `json:"don_gia_ban"`
	VatPercent    float64 `json:"vat_percent"`
	GiaSauVat     float64 `json:"gia_sau_vat"`
	ChietKhauDong float64 `json:"chiet_khau_dong"`
	ThanhTienDong float64 `json:"thanh_tien_dong"`
	GiaVon        float64 `json:"gia_von"`
	BaoHanhThang  int     `json:"bao_hanh_thang"`
	GhiChuDong    string  `json:"ghi_chu_dong"`
}

// =================================================================================
// 6. HÓA ĐƠN ĐIỆN TỬ
// =================================================================================
type HoaDon struct {
	MaHoaDon           string  `json:"ma_hoa_don"`
	MaTraCuu           string  `json:"ma_tra_cuu"`
	XmlUrl             string  `json:"xml_url"`
	LoaiHoaDon         string  `json:"loai_hoa_don"`
	MaPhieuXuat        string  `json:"ma_phieu_xuat"`
	MaKhachHang        string  `json:"ma_khach_hang"`
	NgayHoaDon         string  `json:"ngay_hoa_don"`
	KyHieu             string  `json:"ky_hieu"`
	SoHoaDon           string  `json:"so_hoa_don"`
	MauSo              string  `json:"mau_so"`
	LinkChungTu        string  `json:"link_chung_tu"`
	TongTienPhieu      float64 `json:"tong_tien_phieu"`
	TongVat            float64 `json:"tong_vat"`
	TongTienSauVat     float64 `json:"tong_tien_sau_vat"`
	TrangThai          string  `json:"trang_thai"`
	TrangThaiThanhToan string  `json:"trang_thai_thanh_toan"`
	GhiChu             string  `json:"ghi_chu"`
	NguoiTao           string  `json:"nguoi_tao"`
	NgayTao            string  `json:"ngay_tao"`
	NgayCapNhat        string  `json:"ngay_cap_nhat"`
}

type HoaDonChiTiet struct {
	MaHoaDon   string  `json:"ma_hoa_don"`
	MaSanPham  string  `json:"ma_san_pham"`
	TenSanPham string  `json:"ten_san_pham"`
	DonVi      string  `json:"don_vi"`
	SoLuong    int     `json:"so_luong"`
	DonGiaBan  float64 `json:"don_gia_ban"`
	VatPercent float64 `json:"vat_percent"`
	TienVat    float64 `json:"tien_vat"`
	ThanhTien  float64 `json:"thanh_tien"`
}

// =================================================================================
// 7. KHÁCH HÀNG & NHÂN VIÊN
// =================================================================================
type KhachHang struct {
	MaKhachHang   string  `json:"ma_khach_hang"`
	UserName      string  `json:"user_name"`
	PasswordHash  string  `json:"-"`
	LoaiKhachHang string  `json:"loai_khach_hang"`
	TenKhachHang  string  `json:"ten_khach_hang"`
	DienThoai     string  `json:"dien_thoai"`
	Email         string  `json:"email"`
	UrlFb         string  `json:"url_fb"`
	Zalo          string  `json:"zalo"`
	UrlTele       string  `json:"url_tele"`
	UrlTiktok     string  `json:"url_tiktok"`
	DiaChi        string  `json:"dia_chi"`
	NgaySinh      string  `json:"ngay_sinh"`
	GioiTinh      string  `json:"gioi_tinh"`
	MaSoThue      string  `json:"ma_so_thue"`
	DangNo        float64 `json:"dang_no"`
	TongMua       float64 `json:"tong_mua"`
	TrangThai     int     `json:"trang_thai"`
	GhiChu        string  `json:"ghi_chu"`
	NguoiTao      string  `json:"nguoi_tao"`
	NgayTao       string  `json:"ngay_tao"`
	NgayCapNhat   string  `json:"ngay_cap_nhat"`
}

type NhanVien struct {
	MaNhanVien      string `json:"ma_nhan_vien"`
	TenDangNhap     string `json:"ten_dang_nhap"`
	Email           string `json:"email"`
	MatKhauHash     string `json:"-"`
	HoTen           string `json:"ho_ten"`
	ChucVu          string `json:"chuc_vu"`
	MaPin           string `json:"-"`
	Cookie          string `json:"-"`
	CookieExpired   string `json:"cookie_expired"`
	VaiTroQuyenHan  string `json:"vai_tro_quyen_han"`
	TrangThai       int    `json:"trang_thai"`
	LanDangNhapCuoi string `json:"lan_dang_nhap_cuoi"`
}

// =================================================================================
// 8. TÀI CHÍNH & BẢO HÀNH
// =================================================================================
type PhieuThuChi struct {
	MaPhieuThuChi      string  `json:"ma_phieu_thu_chi"`
	NgayTaoPhieu       string  `json:"ngay_tao_phieu"`
	LoaiPhieu          string  `json:"loai_phieu"`
	DoiTuongLoai       string  `json:"doi_tuong_loai"`
	DoiTuongID         string  `json:"doi_tuong_id"`
	HangMucThuChi      string  `json:"hang_muc_thu_chi"`
	CoHoaDonDo         bool    `json:"co_hoa_don_do"`
	MaChungTuThamChieu string  `json:"ma_chung_tu_tham_chieu"`
	SoTien             float64 `json:"so_tien"`
	PhuongThucTT       string  `json:"phuong_thuc_thanh_toan"`
	TrangThaiDuyet     int     `json:"trang_thai_duyet"`
	NguoiDuyet         string  `json:"nguoi_duyet"`
	GhiChu             string  `json:"ghi_chu"`
	NguoiTao           string  `json:"nguoi_tao"`
	NgayTao            string  `json:"ngay_tao"`
	NgayCapNhat        string  `json:"ngay_cap_nhat"`
}

type PhieuBaoHanh struct {
	MaPhieuBaoHanh    string  `json:"ma_phieu_bao_hanh"`
	LoaiPhieu         string  `json:"loai_phieu"`
	SerialImei        string  `json:"serial_imei"`
	MaSanPham         string  `json:"ma_san_pham"`
	MaKhachHang       string  `json:"ma_khach_hang"`
	TenNguoiGui       string  `json:"ten_nguoi_gui"`
	SdtNguoiGui       string  `json:"sdt_nguoi_gui"`
	NgayNhan          string  `json:"ngay_nhan"`
	TinhTrangLoi      string  `json:"tinh_trang_loi"`
	HinhThuc          string  `json:"hinh_thuc"`
	TrangThai         int     `json:"trang_thai"`
	NgayTraDuKien     string  `json:"ngay_tra_du_kien"`
	NgayTraThucTe     string  `json:"ngay_tra_thuc_te"`
	ChiPhiSua         float64 `json:"chi_phi_sua"`
	PhiThuKhach       float64 `json:"phi_thu_khach"`
	KetQuaSuaChua     string  `json:"ket_qua_sua_chua"`
	LinhKienThayThe   string  `json:"linh_kien_thay_the"`
	MaNhanVienKyThuat string  `json:"ma_nhan_vien_ky_thuat"`
	GhiChu            string  `json:"ghi_chu"`
	NguoiTao          string  `json:"nguoi_tao"`
	NgayTao           string  `json:"ngay_tao"`
	NgayCapNhat       string  `json:"ngay_cap_nhat"`
}

// =================================================================================
// 9. DANH MỤC, THƯƠNG HIỆU, CẤU HÌNH
// =================================================================================
type DanhMuc struct {
	MaDanhMuc    string `json:"ma_danh_muc"`
	ThuTuHienThi int    `json:"thu_tu_hien_thi"`
	TenDanhMuc   string `json:"ten_danh_muc"`
	Slug         string `json:"slug"`
	MaDanhMucCha string `json:"ma_danh_muc_cha"`
}

type ThuongHieu struct {
	MaThuongHieu  string `json:"ma_thuong_hieu"`
	TenThuongHieu string `json:"ten_thuong_hieu"`
	LogoUrl       string `json:"logo_url"`
}

type KhuyenMai struct {
	MaVoucher      string  `json:"ma_voucher"`
	TenChuongTrinh string  `json:"ten_chuong_trinh"`
	LoaiGiam       string  `json:"loai_giam"`
	GiaTriGiam     float64 `json:"gia_tri_giam"`
	DonToThieu     float64 `json:"don_to_thieu"`
	NgayBatDau     string  `json:"ngay_bat_dau"`
	NgayKetThuc    string  `json:"ngay_ket_thuc"`
	SoLuongConLai  int     `json:"so_luong_con_lai"`
	TrangThai      int     `json:"trang_thai"`
}

type CauHinhWeb struct {
	MaCauHinh string `json:"ma_cau_hinh"`
	GiaTri    string `json:"gia_tri"`
	MoTa      string `json:"mo_ta"`
	TrangThai int    `json:"trang_thai"`
}
