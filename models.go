package models

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
}

// =================================================================================
// 2. NHẬP KHO (Import)
// =================================================================================
type PhieuNhap struct {
	MaPhieuNhap        string  `json:"ma_phieu_nhap"`         // Col A
	MaNhaCungCap       string  `json:"ma_nha_cung_cap"`       // Col B
	MaKho              string  `json:"ma_kho"`                // Col C
	NgayNhap           string  `json:"ngay_nhap"`             // Col D
	TrangThai          string  `json:"trang_thai"`            // Col E
	SoHoaDon           string  `json:"so_hoa_don"`            // Col F
	NgayHoaDon         string  `json:"ngay_hoa_don"`          // Col G
	UrlChungTu         string  `json:"url_chung_tu"`          // Col H
	TongTienPhieu      float64 `json:"tong_tien_phieu"`       // Col I
	DaThanhToan        float64 `json:"da_thanh_toan"`         // Col J
	ConNo              float64 `json:"con_no"`                // Col K
	PhuongThucTT       string  `json:"phuong_thuc_thanh_toan"`// Col L
	TrangThaiThanhToan string  `json:"trang_thai_thanh_toan"` // Col M
	GhiChu             string  `json:"ghi_chu"`               // Col N
	NguoiTao           string  `json:"nguoi_tao"`             // Col O
	NgayTao            string  `json:"ngay_tao"`              // Col P
	NgayCapNhat        string  `json:"ngay_cap_nhat"`         // Col Q
}

type ChiTietPhieuNhap struct {
	MaPhieuNhap   string  `json:"ma_phieu_nhap"`   // Col A
	MaSanPham     string  `json:"ma_san_pham"`     // Col B
	TenSanPham    string  `json:"ten_san_pham"`    // Col C
	DonVi         string  `json:"don_vi"`          // Col D
	SoLuong       int     `json:"so_luong"`        // Col E
	DonGiaNhap    float64 `json:"don_gia_nhap"`    // Col F
	VatPercent    float64 `json:"vat_percent"`     // Col G
	GiaSauVat     float64 `json:"gia_sau_vat"`     // Col H
	ChietKhauDong float64 `json:"chiet_khau_dong"` // Col I
	ThanhTienDong float64 `json:"thanh_tien_dong"` // Col J
	GiaVonThucTe  float64 `json:"gia_von_thuc_te"` // Col K
	BaoHanhThang  int     `json:"bao_hanh_thang"`  // Col L
	GhiChuDong    string  `json:"ghi_chu_dong"`    // Col M
}

type NhaCungCap struct {
	MaNhaCungCap string  `json:"ma_nha_cung_cap"` // Col A
	TenNhaCungCap string `json:"ten_nha_cung_cap"`// Col B
	DienThoai    string  `json:"dien_thoai"`      // Col C
	Email        string  `json:"email"`           // Col D
	DiaChi       string  `json:"dia_chi"`         // Col E
	MaSoThue     string  `json:"ma_so_thue"`      // Col F
	NguoiLienHe  string  `json:"nguoi_lien_he"`   // Col G
	NganHang     string  `json:"ngan_hang"`       // Col H
	NoCanTra     float64 `json:"no_can_tra"`      // Col I
	TongMua      float64 `json:"tong_mua"`        // Col J
	HanMucCongNo float64 `json:"han_muc_cong_no"` // Col K
	TrangThai    int     `json:"trang_thai"`      // Col L
	GhiChu       string  `json:"ghi_chu"`         // Col M
	NguoiTao     string  `json:"nguoi_tao"`       // Col N
	NgayTao      string  `json:"ngay_tao"`        // Col O
	NgayCapNhat  string  `json:"ngay_cap_nhat"`   // Col P
}

// =================================================================================
// 3. XUẤT KHO & BÁN HÀNG (Export & Sales)
// =================================================================================
type PhieuXuat struct {
	MaPhieuXuat        string  `json:"ma_phieu_xuat"`         // Col A
	LoaiXuat           string  `json:"loai_xuat"`             // Col B
	NgayXuat           string  `json:"ngay_xuat"`             // Col C
	MaKho              string  `json:"ma_kho"`                // Col D
	MaKhachHang        string  `json:"ma_khach_hang"`         // Col E
	TrangThai          string  `json:"trang_thai"`            // Col F
	MaVoucher          string  `json:"ma_voucher"`            // Col G
	TienGiamVoucher    float64 `json:"tien_giam_voucher"`     // Col H
	TongTienPhieu      float64 `json:"tong_tien_phieu"`       // Col I
	LinkChungTu        string  `json:"link_chung_tu"`         // Col J
	DaThu              float64 `json:"da_thu"`                // Col K
	ConNo              float64 `json:"con_no"`                // Col L
	PhuongThucTT       string  `json:"phuong_thuc_thanh_toan"`// Col M
	TrangThaiTT        string  `json:"trang_thai_thanh_toan"` // Col N
	PhiVanChuyen       float64 `json:"phi_van_chuyen"`        // Col O
	NguonDonHang       string  `json:"nguon_don_hang"`        // Col P
	ThongTinGiaoHang   string  `json:"thong_tin_giao_hang"`   // Col Q
	GhiChu             string  `json:"ghi_chu"`               // Col R
	NguoiTao           string  `json:"nguoi_tao"`             // Col S
	NgayTao            string  `json:"ngay_tao"`              // Col T
	NgayCapNhat        string  `json:"ngay_cap_nhat"`         // Col U
}

type ChiTietPhieuXuat struct {
	MaPhieuXuat   string  `json:"ma_phieu_xuat"`   // Col A
	MaSanPham     string  `json:"ma_san_pham"`     // Col B
	TenSanPham    string  `json:"ten_san_pham"`    // Col C
	DonVi         string  `json:"don_vi"`          // Col D
	SoLuong       int     `json:"so_luong"`        // Col E
	DonGiaBan     float64 `json:"don_gia_ban"`     // Col F
	VatPercent    float64 `json:"vat_percent"`     // Col G
	GiaSauVat     float64 `json:"gia_sau_vat"`     // Col H
	ChietKhauDong float64 `json:"chiet_khau_dong"` // Col I
	ThanhTienDong float64 `json:"thanh_tien_dong"` // Col J
	GiaVon        float64 `json:"gia_von"`         // Col K
	BaoHanhThang  int     `json:"bao_hanh_thang"`  // Col L
	GhiChuDong    string  `json:"ghi_chu_dong"`    // Col M
}

// =================================================================================
// 4. KHO & BẢO HÀNH (WMS & Warranty)
// =================================================================================
type SerialSanPham struct {
	SerialImei           string `json:"serial_imei"`              // Col A
	MaSanPham            string `json:"ma_san_pham"`              // Col B
	MaNhaCungCap         string `json:"ma_nha_cung_cap"`          // Col C
	MaPhieuNhap          string `json:"ma_phieu_nhap"`            // Col D
	MaPhieuXuat          string `json:"ma_phieu_xuat"`            // Col E
	TrangThai            int    `json:"trang_thai"`               // Col F
	BaoHanhNhaCungCap    int    `json:"bao_hanh_nha_cung_cap"`    // Col G
	HanBaoHanhNhaCungCap string `json:"han_bao_hanh_nha_cung_cap"`// Col H
	MaKhachHangHienTai   string `json:"ma_khach_hang_hien_tai"`   // Col I
	NgayXuatKho          string `json:"ngay_xuat_kho"`            // Col J
	KichHoatBaoHanhKhach string `json:"kich_hoat_bao_hanh_khach"` // Col K
	HanBaoHanhKhach      string `json:"han_bao_hanh_khach"`       // Col L
	MaKho                string `json:"ma_kho"`                   // Col M
	GhiChu               string `json:"ghi_chu"`                  // Col N
	NgayCapNhat          string `json:"ngay_cap_nhat"`            // Col O
}

type PhieuBaoHanh struct {
	MaPhieuBaoHanh    string  `json:"ma_phieu_bao_hanh"`   // Col A
	LoaiPhieu         string  `json:"loai_phieu"`          // Col B
	SerialImei        string  `json:"serial_imei"`         // Col C
	MaSanPham         string  `json:"ma_san_pham"`         // Col D
	MaKhachHang       string  `json:"ma_khach_hang"`       // Col E
	TenNguoiGui       string  `json:"ten_nguoi_gui"`       // Col F
	SdtNguoiGui       string  `json:"sdt_nguoi_gui"`       // Col G
	NgayNhan          string  `json:"ngay_nhan"`           // Col H
	TinhTrangLoi      string  `json:"tinh_trang_loi"`      // Col I
	HinhThuc          string  `json:"hinh_thuc"`           // Col J
	TrangThai         int     `json:"trang_thai"`          // Col K
	NgayTraDuKien     string  `json:"ngay_tra_du_kien"`    // Col L
	NgayTraThucTe     string  `json:"ngay_tra_thuc_te"`    // Col M
	ChiPhiSua         float64 `json:"chi_phi_sua"`         // Col N
	PhiThuKhach       float64 `json:"phi_thu_khach"`       // Col O
	KetQuaSuaChua     string  `json:"ket_qua_sua_chua"`    // Col P
	LinhKienThayThe   string  `json:"linh_kien_thay_the"`  // Col Q
	MaNhanVienKyThuat string  `json:"ma_nhan_vien_ky_thuat"`// Col R
	GhiChu            string  `json:"ghi_chu"`             // Col S
	NguoiTao          string  `json:"nguoi_tao"`           // Col T
	NgayTao           string  `json:"ngay_tao"`            // Col U
	NgayCapNhat       string  `json:"ngay_cap_nhat"`       // Col V
}

// =================================================================================
// 5. TÀI CHÍNH & HÓA ĐƠN (Accounting)
// =================================================================================
type HoaDon struct {
	MaHoaDon           string  `json:"ma_hoa_don"`            // Col A
	MaTraCuu           string  `json:"ma_tra_cuu"`            // Col B
	XmlUrl             string  `json:"xml_url"`               // Col C
	LoaiHoaDon         string  `json:"loai_hoa_don"`          // Col D
	MaPhieuXuat        string  `json:"ma_phieu_xuat"`         // Col E
	MaKhachHang        string  `json:"ma_khach_hang"`         // Col F
	NgayHoaDon         string  `json:"ngay_hoa_don"`          // Col G
	KyHieu             string  `json:"ky_hieu"`               // Col H
	SoHoaDon           string  `json:"so_hoa_don"`            // Col I
	MauSo              string  `json:"mau_so"`                // Col J
	LinkChungTu        string  `json:"link_chung_tu"`         // Col K
	TongTienPhieu      float64 `json:"tong_tien_phieu"`       // Col L
	TongVat            float64 `json:"tong_vat"`              // Col M
	TongTienSauVat     float64 `json:"tong_tien_sau_vat"`     // Col N
	TrangThai          string  `json:"trang_thai"`            // Col O
	TrangThaiThanhToan string  `json:"trang_thai_thanh_toan"` // Col P
	GhiChu             string  `json:"ghi_chu"`               // Col Q
	NguoiTao           string  `json:"nguoi_tao"`             // Col R
	NgayTao            string  `json:"ngay_tao"`              // Col S
	NgayCapNhat        string  `json:"ngay_cap_nhat"`         // Col T
}

type HoaDonChiTiet struct {
	MaHoaDon   string  `json:"ma_hoa_don"`   // Col A
	MaSanPham  string  `json:"ma_san_pham"`  // Col B
	TenSanPham string  `json:"ten_san_pham"` // Col C
	DonVi      string  `json:"don_vi"`       // Col D
	SoLuong    int     `json:"so_luong"`     // Col E
	DonGiaBan  float64 `json:"don_gia_ban"`  // Col F
	VatPercent float64 `json:"vat_percent"`  // Col G
	TienVat    float64 `json:"tien_vat"`     // Col H
	ThanhTien  float64 `json:"thanh_tien"`   // Col I
}

type PhieuThuChi struct {
	MaPhieuThuChi      string  `json:"ma_phieu_thu_chi"`      // Col A
	NgayTaoPhieu       string  `json:"ngay_tao_phieu"`        // Col B
	LoaiPhieu          string  `json:"loai_phieu"`            // Col C
	DoiTuongLoai       string  `json:"doi_tuong_loai"`        // Col D
	DoiTuongID         string  `json:"doi_tuong_id"`          // Col E
	HangMucThuChi      string  `json:"hang_muc_thu_chi"`      // Col F
	CoHoaDonDo         bool    `json:"co_hoa_don_do"`         // Col G
	MaChungTuThamChieu string  `json:"ma_chung_tu_tham_chieu"`// Col H
	SoTien             float64 `json:"so_tien"`               // Col I
	PhuongThucTT       string  `json:"phuong_thuc_thanh_toan"`// Col J
	TrangThaiDuyet     int     `json:"trang_thai_duyet"`      // Col K
	NguoiDuyet         string  `json:"nguoi_duyet"`           // Col L
	GhiChu             string  `json:"ghi_chu"`               // Col M
	NguoiTao           string  `json:"nguoi_tao"`             // Col N
	NgayTao            string  `json:"ngay_tao"`              // Col O
	NgayCapNhat        string  `json:"ngay_cap_nhat"`         // Col P
}

// =================================================================================
// 6. KHÁCH HÀNG & NHÂN VIÊN (CRM & HR)
// =================================================================================
type KhachHang struct {
	MaKhachHang   string  `json:"ma_khach_hang"`   // Col A
	UserName      string  `json:"user_name"`       // Col B
	PasswordHash  string  `json:"-"`               // Col C
	LoaiKhachHang string  `json:"loai_khach_hang"` // Col D
	TenKhachHang  string  `json:"ten_khach_hang"`  // Col E
	DienThoai     string  `json:"dien_thoai"`      // Col F
	Email         string  `json:"email"`           // Col G
	UrlFb         string  `json:"url_fb"`          // Col H
	Zalo          string  `json:"zalo"`            // Col I
	UrlTele       string  `json:"url_tele"`        // Col J
	UrlTiktok     string  `json:"url_tiktok"`      // Col K
	DiaChi        string  `json:"dia_chi"`         // Col L
	NgaySinh      string  `json:"ngay_sinh"`       // Col M
	GioiTinh      string  `json:"gioi_tinh"`       // Col N
	MaSoThue      string  `json:"ma_so_thue"`      // Col O
	DangNo        float64 `json:"dang_no"`         // Col P
	TongMua       float64 `json:"tong_mua"`        // Col Q
	TrangThai     int     `json:"trang_thai"`      // Col R
	GhiChu        string  `json:"ghi_chu"`         // Col S
	NguoiTao      string  `json:"nguoi_tao"`       // Col T
	NgayTao       string  `json:"ngay_tao"`        // Col U
	NgayCapNhat   string  `json:"ngay_cap_nhat"`   // Col V
}

type NhanVien struct {
	MaNhanVien      string `json:"ma_nhan_vien"`       // Col A
	TenDangNhap     string `json:"ten_dang_nhap"`      // Col B
	Email           string `json:"email"`              // Col C
	MatKhauHash     string `json:"-"`                  // Col D
	HoTen           string `json:"ho_ten"`             // Col E
	ChucVu          string `json:"chuc_vu"`            // Col F
	MaPin           string `json:"-"`                  // Col G
	Cookie          string `json:"-"`                  // Col H
	CookieExpired   string `json:"cookie_expired"`     // Col I
	VaiTroQuyenHan  string `json:"vai_tro_quyen_han"`  // Col J
	TrangThai       int    `json:"trang_thai"`         // Col K
	LanDangNhapCuoi string `json:"lan_dang_nhap_cuoi"` // Col L
}

// =================================================================================
// 7. CẤU HÌNH & DANH MỤC (Settings)
// =================================================================================
type DanhMuc struct {
	MaDanhMuc    string `json:"ma_danh_muc"`     // Col A
	ThuTuHienThi int    `json:"thu_tu_hien_thi"` // Col B
	TenDanhMuc   string `json:"ten_danh_muc"`    // Col C
	Slug         string `json:"slug"`            // Col D
	MaDanhMucCha string `json:"ma_danh_muc_cha"` // Col E
}

type ThuongHieu struct {
	MaThuongHieu  string `json:"ma_thuong_hieu"`  // Col A
	TenThuongHieu string `json:"ten_thuong_hieu"` // Col B
	LogoUrl       string `json:"logo_url"`        // Col C
}

type KhuyenMai struct {
	MaVoucher      string  `json:"ma_voucher"`       // Col A
	TenChuongTrinh string  `json:"ten_chuong_trinh"` // Col B
	LoaiGiam       string  `json:"loai_giam"`        // Col C
	GiaTriGiam     float64 `json:"gia_tri_giam"`     // Col D
	DonToThieu     float64 `json:"don_to_thieu"`     // Col E
	NgayBatDau     string  `json:"ngay_bat_dau"`     // Col F
	NgayKetThuc    string  `json:"ngay_ket_thuc"`    // Col G
	SoLuongConLai  int     `json:"so_luong_con_lai"` // Col H
	TrangThai      int     `json:"trang_thai"`       // Col I
}

type CauHinhWeb struct {
	MaCauHinh string `json:"ma_cau_hinh"` // Col A
	GiaTri    string `json:"gia_tri"`     // Col B
	MoTa      string `json:"mo_ta"`       // Col C
	TrangThai int    `json:"trang_thai"`  // Col D
}
