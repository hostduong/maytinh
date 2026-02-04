package mo_hinh

// ========================================================
// 1. BẢNG: SẢN PHẨM (SAN_PHAM)
// ========================================================
type SanPham struct {
	MaSanPham    string  `json:"ma_san_pham"`    // Cột A
	TenSanPham   string  `json:"ten_san_pham"`   // Cột B
	TenRutGon    string  `json:"ten_rut_gon"`    // Cột C
	Sku          string  `json:"sku"`            // Cột D
	MaDanhMuc    string  `json:"ma_danh_muc"`    // Cột E
	MaThuongHieu string  `json:"ma_thuong_hieu"` // Cột F
	DonVi        string  `json:"don_vi"`         // Cột G
	MauSac       string  `json:"mau_sac"`        // Cột H
	UrlHinhAnh   string  `json:"url_hinh_anh"`   // Cột I
	ThongSo      string  `json:"thong_so"`       // Cột J
	MoTaChiTiet  string  `json:"mo_ta_chi_tiet"` // Cột K
	BaoHanhThang int     `json:"bao_hanh_thang"` // Cột L
	TinhTrang    string  `json:"tinh_trang"`     // Cột M
	TrangThai    int     `json:"trang_thai"`     // Cột N
	GiaBanLe     float64 `json:"gia_ban_le"`     // Cột O
	GhiChu       string  `json:"ghi_chu"`        // Cột P
	NguoiTao     string  `json:"nguoi_tao"`      // Cột Q
	NgayTao      string  `json:"ngay_tao"`       // Cột R
	NgayCapNhat  string  `json:"ngay_cap_nhat"`  // Cột S
}

// ========================================================
// 2. BẢNG: PHIẾU NHẬP (PHIEU_NHAP)
// ========================================================
type PhieuNhap struct {
	MaPhieuNhap         string  `json:"ma_phieu_nhap"`          // Cột A
	MaNhaCungCap        string  `json:"ma_nha_cung_cap"`        // Cột B
	MaKho               string  `json:"ma_kho"`                 // Cột C
	NgayNhap            string  `json:"ngay_nhap"`              // Cột D
	TrangThai           string  `json:"trang_thai"`             // Cột E
	SoHoaDon            string  `json:"so_hoa_don"`             // Cột F
	NgayHoaDon          string  `json:"ngay_hoa_don"`           // Cột G
	UrlChungTu          string  `json:"url_chung_tu"`           // Cột H
	TongTienPhieu       float64 `json:"tong_tien_phieu"`        // Cột I
	DaThanhToan         float64 `json:"da_thanh_toan"`          // Cột J
	ConNo               float64 `json:"con_no"`                 // Cột K
	PhuongThucThanhToan string  `json:"phuong_thuc_thanh_toan"` // Cột L
	TrangThaiThanhToan  string  `json:"trang_thai_thanh_toan"`  // Cột M
	GhiChu              string  `json:"ghi_chu"`                // Cột N
	NguoiTao            string  `json:"nguoi_tao"`              // Cột O
	NgayTao             string  `json:"ngay_tao"`               // Cột P
	NgayCapNhat         string  `json:"ngay_cap_nhat"`          // Cột Q
}

// ========================================================
// 3. BẢNG: CHI TIẾT PHIẾU NHẬP (CHI_TIET_PHIEU_NHAP)
// ========================================================
type ChiTietPhieuNhap struct {
	MaPhieuNhap   string  `json:"ma_phieu_nhap"`   // Cột A
	MaSanPham     string  `json:"ma_san_pham"`     // Cột B
	TenSanPham    string  `json:"ten_san_pham"`    // Cột C
	DonVi         string  `json:"don_vi"`          // Cột D
	SoLuong       int     `json:"so_luong"`        // Cột E
	DonGiaNhap    float64 `json:"don_gia_nhap"`    // Cột F
	VatPercent    float64 `json:"vat_percent"`     // Cột G
	GiaSauVat     float64 `json:"gia_sau_vat"`     // Cột H
	ChietKhauDong float64 `json:"chiet_khau_dong"` // Cột I
	ThanhTienDong float64 `json:"thanh_tien_dong"` // Cột J
	GiaVonThucTe  float64 `json:"gia_von_thuc_te"` // Cột K
	BaoHanhThang  int     `json:"bao_hanh_thang"`  // Cột L
	GhiChuDong    string  `json:"ghi_chu_dong"`    // Cột M
}

// ========================================================
// 4. BẢNG: NHÀ CUNG CẤP (NHA_CUNG_CAP)
// ========================================================
type NhaCungCap struct {
	MaNhaCungCap  string  `json:"ma_nha_cung_cap"`  // Cột A
	TenNhaCungCap string  `json:"ten_nha_cung_cap"` // Cột B
	DienThoai     string  `json:"dien_thoai"`       // Cột C
	Email         string  `json:"email"`            // Cột D
	DiaChi        string  `json:"dia_chi"`          // Cột E
	MaSoThue      string  `json:"ma_so_thue"`       // Cột F
	NguoiLienHe   string  `json:"nguoi_lien_he"`    // Cột G
	NganHang      string  `json:"ngan_hang"`        // Cột H
	NoCanTra      float64 `json:"no_can_tra"`       // Cột I
	TongMua       float64 `json:"tong_mua"`         // Cột J
	HanMucCongNo  float64 `json:"han_muc_cong_no"`  // Cột K
	TrangThai     int     `json:"trang_thai"`       // Cột L
	GhiChu        string  `json:"ghi_chu"`          // Cột M
	NguoiTao      string  `json:"nguoi_tao"`        // Cột N
	NgayTao       string  `json:"ngay_tao"`         // Cột O
	NgayCapNhat   string  `json:"ngay_cap_nhat"`    // Cột P
}

// ========================================================
// 5. BẢNG: SERIAL SẢN PHẨM (SERIAL_SAN_PHAM)
// ========================================================
type SerialSanPham struct {
	SerialImei           string `json:"serial_imei"`              // Cột A
	MaSanPham            string `json:"ma_san_pham"`              // Cột B
	MaNhaCungCap         string `json:"ma_nha_cung_cap"`          // Cột C
	MaPhieuNhap          string `json:"ma_phieu_nhap"`            // Cột D
	MaPhieuXuat          string `json:"ma_phieu_xuat"`            // Cột E
	TrangThai            int    `json:"trang_thai"`               // Cột F
	BaoHanhNhaCungCap    int    `json:"bao_hanh_nha_cung_cap"`    // Cột G
	HanBaoHanhNhaCungCap string `json:"han_bao_hanh_nha_cung_cap"`// Cột H
	MaKhachHangHienTai   string `json:"ma_khach_hang_hien_tai"`   // Cột I
	NgayXuatKho          string `json:"ngay_xuat_kho"`            // Cột J
	KichHoatBaoHanhKhach string `json:"kich_hoat_bao_hanh_khach"` // Cột K
	HanBaoHanhKhach      string `json:"han_bao_hanh_khach"`       // Cột L
	MaKho                string `json:"ma_kho"`                   // Cột M
	GhiChu               string `json:"ghi_chu"`                  // Cột N
	NgayCapNhat          string `json:"ngay_cap_nhat"`            // Cột O
}

// ========================================================
// 6. BẢNG: PHIẾU XUẤT (PHIEU_XUAT)
// ========================================================
type PhieuXuat struct {
	MaPhieuXuat         string  `json:"ma_phieu_xuat"`          // Cột A
	LoaiXuat            string  `json:"loai_xuat"`              // Cột B
	NgayXuat            string  `json:"ngay_xuat"`              // Cột C
	MaKho               string  `json:"ma_kho"`                 // Cột D
	MaKhachHang         string  `json:"ma_khach_hang"`          // Cột E
	TrangThai           string  `json:"trang_thai"`             // Cột F
	MaVoucher           string  `json:"ma_voucher"`             // Cột G
	TienGiamVoucher     float64 `json:"tien_giam_voucher"`      // Cột H
	TongTienPhieu       float64 `json:"tong_tien_phieu"`        // Cột I
	LinkChungTu         string  `json:"link_chung_tu"`          // Cột J
	DaThu               float64 `json:"da_thu"`                 // Cột K
	ConNo               float64 `json:"con_no"`                 // Cột L
	PhuongThucThanhToan string  `json:"phuong_thuc_thanh_toan"` // Cột M
	TrangThaiThanhToan  string  `json:"trang_thai_thanh_toan"`  // Cột N
	PhiVanChuyen        float64 `json:"phi_van_chuyen"`         // Cột O
	NguonDonHang        string  `json:"nguon_don_hang"`         // Cột P
	ThongTinGiaoHang    string  `json:"thong_tin_giao_hang"`    // Cột Q
	GhiChu              string  `json:"ghi_chu"`                // Cột R
	NguoiTao            string  `json:"nguoi_tao"`              // Cột S
	NgayTao             string  `json:"ngay_tao"`               // Cột T
	NgayCapNhat         string  `json:"ngay_cap_nhat"`          // Cột U
}

// ========================================================
// 7. BẢNG: CHI TIẾT PHIẾU XUẤT (CHI_TIET_PHIEU_XUAT)
// ========================================================
type ChiTietPhieuXuat struct {
	MaPhieuXuat   string  `json:"ma_phieu_xuat"`   // Cột A
	MaSanPham     string  `json:"ma_san_pham"`     // Cột B
	TenSanPham    string  `json:"ten_san_pham"`    // Cột C
	DonVi         string  `json:"don_vi"`          // Cột D
	SoLuong       int     `json:"so_luong"`        // Cột E
	DonGiaBan     float64 `json:"don_gia_ban"`     // Cột F
	VatPercent    float64 `json:"vat_percent"`     // Cột G
	GiaSauVat     float64 `json:"gia_sau_vat"`     // Cột H
	ChietKhauDong float64 `json:"chiet_khau_dong"` // Cột I
	ThanhTienDong float64 `json:"thanh_tien_dong"` // Cột J
	GiaVon        float64 `json:"gia_von"`         // Cột K
	BaoHanhThang  int     `json:"bao_hanh_thang"`  // Cột L
	GhiChuDong    string  `json:"ghi_chu_dong"`    // Cột M
}

// ========================================================
// 8. BẢNG: HÓA ĐƠN (HOA_DON)
// ========================================================
type HoaDon struct {
	MaHoaDon            string  `json:"ma_hoa_don"`             // Cột A
	MaTraCuu            string  `json:"ma_tra_cuu"`             // Cột B
	XmlUrl              string  `json:"xml_url"`                // Cột C
	LoaiHoaDon          string  `json:"loai_hoa_don"`           // Cột D
	MaPhieuXuat         string  `json:"ma_phieu_xuat"`          // Cột E
	MaKhachHang         string  `json:"ma_khach_hang"`          // Cột F
	NgayHoaDon          string  `json:"ngay_hoa_don"`           // Cột G
	KyHieu              string  `json:"ky_hieu"`                // Cột H
	SoHoaDon            string  `json:"so_hoa_don"`             // Cột I
	MauSo               string  `json:"mau_so"`                 // Cột J
	LinkChungTu         string  `json:"link_chung_tu"`          // Cột K
	TongTienPhieu       float64 `json:"tong_tien_phieu"`        // Cột L
	TongVat             float64 `json:"tong_vat"`               // Cột M
	TongTienSauVat      float64 `json:"tong_tien_sau_vat"`      // Cột N
	TrangThai           string  `json:"trang_thai"`             // Cột O
	TrangThaiThanhToan  string  `json:"trang_thai_thanh_toan"`  // Cột P
	GhiChu              string  `json:"ghi_chu"`                // Cột Q
	NguoiTao            string  `json:"nguoi_tao"`              // Cột R
	NgayTao             string  `json:"ngay_tao"`               // Cột S
	NgayCapNhat         string  `json:"ngay_cap_nhat"`          // Cột T
}

// ========================================================
// 9. BẢNG: HÓA ĐƠN CHI TIẾT (HOA_DON_CHI_TIET)
// ========================================================
type HoaDonChiTiet struct {
	MaHoaDon   string  `json:"ma_hoa_don"`   // Cột A
	MaSanPham  string  `json:"ma_san_pham"`  // Cột B
	TenSanPham string  `json:"ten_san_pham"` // Cột C
	DonVi      string  `json:"don_vi"`       // Cột D
	SoLuong    int     `json:"so_luong"`     // Cột E
	DonGiaBan  float64 `json:"don_gia_ban"`  // Cột F
	VatPercent float64 `json:"vat_percent"`  // Cột G
	TienVat    float64 `json:"tien_vat"`     // Cột H
	ThanhTien  float64 `json:"thanh_tien"`   // Cột I
}

// ========================================================
// 10. BẢNG: KHÁCH HÀNG (KHACH_HANG)
// ========================================================
type KhachHang struct {
	MaKhachHang   string  `json:"ma_khach_hang"`   // Cột A
	UserName      string  `json:"user_name"`       // Cột B
	PasswordHash  string  `json:"-"`               // Cột C (Ẩn JSON)
	LoaiKhachHang string  `json:"loai_khach_hang"` // Cột D
	TenKhachHang  string  `json:"ten_khach_hang"`  // Cột E
	DienThoai     string  `json:"dien_thoai"`      // Cột F
	Email         string  `json:"email"`           // Cột G
	UrlFb         string  `json:"url_fb"`          // Cột H
	Zalo          string  `json:"zalo"`            // Cột I
	UrlTele       string  `json:"url_tele"`        // Cột J
	UrlTiktok     string  `json:"url_tiktok"`      // Cột K
	DiaChi        string  `json:"dia_chi"`         // Cột L
	NgaySinh      string  `json:"ngay_sinh"`       // Cột M
	GioiTinh      string  `json:"gioi_tinh"`       // Cột N
	MaSoThue      string  `json:"ma_so_thue"`      // Cột O
	DangNo        float64 `json:"dang_no"`         // Cột P
	TongMua       float64 `json:"tong_mua"`        // Cột Q
	TrangThai     int     `json:"trang_thai"`      // Cột R
	GhiChu        string  `json:"ghi_chu"`         // Cột S
	NguoiTao      string  `json:"nguoi_tao"`       // Cột T
	NgayTao       string  `json:"ngay_tao"`        // Cột U
	NgayCapNhat   string  `json:"ngay_cap_nhat"`   // Cột V
}

// ========================================================
// 11. BẢNG: PHIẾU THU CHI (PHIEU_THU_CHI)
// ========================================================
type PhieuThuChi struct {
	MaPhieuThuChi       string  `json:"ma_phieu_thu_chi"`       // Cột A
	NgayTaoPhieu        string  `json:"ngay_tao_phieu"`         // Cột B
	LoaiPhieu           string  `json:"loai_phieu"`             // Cột C
	DoiTuongLoai        string  `json:"doi_tuong_loai"`         // Cột D
	DoiTuongId          string  `json:"doi_tuong_id"`           // Cột E
	HangMucThuChi       string  `json:"hang_muc_thu_chi"`       // Cột F
	CoHoaDonDo          bool    `json:"co_hoa_don_do"`          // Cột G
	MaChungTuThamChieu  string  `json:"ma_chung_tu_tham_chieu"` // Cột H
	SoTien              float64 `json:"so_tien"`                // Cột I
	PhuongThucThanhToan string  `json:"phuong_thuc_thanh_toan"` // Cột J
	TrangThaiDuyet      int     `json:"trang_thai_duyet"`       // Cột K
	NguoiDuyet          string  `json:"nguoi_duyet"`            // Cột L
	GhiChu              string  `json:"ghi_chu"`                // Cột M
	NguoiTao            string  `json:"nguoi_tao"`              // Cột N
	NgayTao             string  `json:"ngay_tao"`               // Cột O
	NgayCapNhat         string  `json:"ngay_cap_nhat"`          // Cột P
}

// ========================================================
// 12. BẢNG: PHIẾU BẢO HÀNH (PHIEU_BAO_HANH)
// ========================================================
type PhieuBaoHanh struct {
	MaPhieuBaoHanh    string  `json:"ma_phieu_bao_hanh"`    // Cột A
	LoaiPhieu         string  `json:"loai_phieu"`           // Cột B
	SerialImei        string  `json:"serial_imei"`          // Cột C
	MaSanPham         string  `json:"ma_san_pham"`          // Cột D
	MaKhachHang       string  `json:"ma_khach_hang"`        // Cột E
	TenNguoiGui       string  `json:"ten_nguoi_gui"`        // Cột F
	SdtNguoiGui       string  `json:"sdt_nguoi_gui"`        // Cột G
	NgayNhan          string  `json:"ngay_nhan"`            // Cột H
	TinhTrangLoi      string  `json:"tinh_trang_loi"`       // Cột I
	HinhThuc          string  `json:"hinh_thuc"`            // Cột J
	TrangThai         int     `json:"trang_thai"`           // Cột K
	NgayTraDuKien     string  `json:"ngay_tra_du_kien"`     // Cột L
	NgayTraThucTe     string  `json:"ngay_tra_thuc_te"`     // Cột M
	ChiPhiSua         float64 `json:"chi_phi_sua"`          // Cột N
	PhiThuKhach       float64 `json:"phi_thu_khach"`        // Cột O
	KetQuaSuaChua     string  `json:"ket_qua_sua_chua"`     // Cột P
	LinhKienThayThe   string  `json:"linh_kien_thay_the"`   // Cột Q
	MaNhanVienKyThuat string  `json:"ma_nhan_vien_ky_thuat"`// Cột R
	GhiChu            string  `json:"ghi_chu"`              // Cột S
	NguoiTao          string  `json:"nguoi_tao"`            // Cột T
	NgayTao           string  `json:"ngay_tao"`             // Cột U
	NgayCapNhat       string  `json:"ngay_cap_nhat"`        // Cột V
}

// ========================================================
// 13. BẢNG: DANH MỤC (DANH_MUC)
// ========================================================
type DanhMuc struct {
	MaDanhMuc    string `json:"ma_danh_muc"`     // Cột A
	ThuTuHienThi int    `json:"thu_tu_hien_thi"` // Cột B
	TenDanhMuc   string `json:"ten_danh_muc"`    // Cột C
	Slug         string `json:"slug"`            // Cột D
	MaDanhMucCha string `json:"ma_danh_muc_cha"` // Cột E
}

// ========================================================
// 14. BẢNG: THƯƠNG HIỆU (THUONG_HIEU)
// ========================================================
type ThuongHieu struct {
	MaThuongHieu  string `json:"ma_thuong_hieu"`  // Cột A
	TenThuongHieu string `json:"ten_thuong_hieu"` // Cột B
	LogoUrl       string `json:"logo_url"`        // Cột C
}

// ========================================================
// 15. BẢNG: NHÂN VIÊN (NHAN_VIEN)
// ========================================================
type NhanVien struct {
	MaNhanVien      string `json:"ma_nhan_vien"`       // Cột A
	TenDangNhap     string `json:"ten_dang_nhap"`      // Cột B
	Email           string `json:"email"`              // Cột C
	MatKhauHash     string `json:"-"`                  // Cột D (Ẩn JSON)
	HoTen           string `json:"ho_ten"`             // Cột E
	ChucVu          string `json:"chuc_vu"`            // Cột F
	MaPin           string `json:"-"`                  // Cột G (Ẩn JSON)
	Cookie          string `json:"-"`                  // Cột H (Ẩn JSON)
	CookieExpired   string `json:"cookie_expired"`     // Cột I
	VaiTroQuyenHan  string `json:"vai_tro_quyen_han"`  // Cột J
	TrangThai       int    `json:"trang_thai"`         // Cột K
	LanDangNhapCuoi string `json:"lan_dang_nhap_cuoi"` // Cột L
}

// ========================================================
// 16. BẢNG: KHUYẾN MÃI (KHUYEN_MAI)
// ========================================================
type KhuyenMai struct {
	MaVoucher      string  `json:"ma_voucher"`       // Cột A
	TenChuongTrinh string  `json:"ten_chuong_trinh"` // Cột B
	LoaiGiam       string  `json:"loai_giam"`        // Cột C
	GiaTriGiam     float64 `json:"gia_tri_giam"`     // Cột D
	DonToThieu     float64 `json:"don_to_thieu"`     // Cột E
	NgayBatDau     string  `json:"ngay_bat_dau"`     // Cột F
	NgayKetThuc    string  `json:"ngay_ket_thuc"`    // Cột G
	SoLuongConLai  int     `json:"so_luong_con_lai"` // Cột H
	TrangThai      int     `json:"trang_thai"`       // Cột I
}

// ========================================================
// 17. BẢNG: CẤU HÌNH WEB (CAU_HINH_WEB)
// ========================================================
type CauHinhWeb struct {
	MaCauHinh string `json:"ma_cau_hinh"` // Cột A
	GiaTri    string `json:"gia_tri"`     // Cột B
	MoTa      string `json:"mo_ta"`       // Cột C
	TrangThai int    `json:"trang_thai"`  // Cột D
}
