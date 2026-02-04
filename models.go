package models

// GHI CHÚ CHUNG:
// - PK (Primary Key): Khóa chính, định danh duy nhất.
// - FK (Foreign Key): Khóa ngoại, liên kết sang bảng khác.
// - Enum: Các giá trị cố định (VD: 0, 1, 2).
// - Snapshot: Dữ liệu sao chép tại thời điểm tạo, không đổi theo Master Data.

// =================================================================================
// 1. SẢN PHẨM (Master Data) - Quản lý danh mục hàng hóa
// =================================================================================
type SanPham struct {
	MaSanPham    string  `json:"ma_san_pham"`    // Cột A (PK): Mã duy nhất (VD: SP001).
	TenSanPham   string  `json:"ten_san_pham"`   // Cột B: Tên hiển thị đầy đủ.
	TenRutGon    string  `json:"ten_rut_gon"`    // Cột C: Slug cho SEO (VD: iphone-15-pro-max).
	Sku          string  `json:"sku"`            // Cột D: Mã vạch/Barcode in trên tem.
	MaDanhMuc    string  `json:"ma_danh_muc"`    // Cột E (FK): Link sang sheet DANH_MUC.
	MaThuongHieu string  `json:"ma_thuong_hieu"` // Cột F (FK): Link sang sheet THUONG_HIEU.
	DonVi        string  `json:"don_vi"`         // Cột G: Cái, Hộp, Bộ...
	MauSac       string  `json:"mau_sac"`        // Cột H: Màu sắc (Text).
	UrlHinhAnh   string  `json:"url_hinh_anh"`   // Cột I: Link ảnh đại diện (hoặc chuỗi JSON nhiều ảnh).
	ThongSo      string  `json:"thong_so"`       // Cột J: Tóm tắt cấu hình ngắn.
	MoTaChiTiet  string  `json:"mo_ta_chi_tiet"` // Cột K: Nội dung HTML bài viết.
	BaoHanhThang int     `json:"bao_hanh_thang"` // Cột L: Số tháng bảo hành mặc định.
	TinhTrang    string  `json:"tinh_trang"`     // Cột M: Mới / Cũ (99%) / Like New.
	TrangThai    int     `json:"trang_thai"`     // Cột N: 1=Đang bán, 0=Ngừng kinh doanh.
	GiaBanLe     float64 `json:"gia_ban_le"`     // Cột O: Giá niêm yết hiện tại.
	GhiChu       string  `json:"ghi_chu"`        // Cột P: Ghi chú nội bộ.
	NguoiTao     string  `json:"nguoi_tao"`      // Cột Q: User tạo SP.
	NgayTao      string  `json:"ngay_tao"`       // Cột R: Thời gian tạo.
	NgayCapNhat  string  `json:"ngay_cap_nhat"`  // Cột S: Thời gian sửa cuối.
}

// =================================================================================
// 2. PHIẾU NHẬP (Import Header) - Quản lý hóa đơn nhập mua
// =================================================================================
type PhieuNhap struct {
	MaPhieuNhap        string  `json:"ma_phieu_nhap"`         // Cột A (PK): Mã phiếu (VD: PN2410-001).
	MaNhaCungCap       string  `json:"ma_nha_cung_cap"`       // Cột B (FK): Link sang NHA_CUNG_CAP.
	MaKho              string  `json:"ma_kho"`                // Cột C: Nhập vào kho nào.
	NgayNhap           string  `json:"ngay_nhap"`             // Cột D: Ngày hạch toán kế toán.
	TrangThai          string  `json:"trang_thai"`            // Cột E: NEW (Mới), DONE (Đã nhập kho), CANCEL (Hủy).
	SoHoaDon           string  `json:"so_hoa_don"`            // Cột F: Số hóa đơn đỏ của NCC gửi.
	NgayHoaDon         string  `json:"ngay_hoa_don"`          // Cột G: Ngày trên hóa đơn đỏ.
	UrlChungTu         string  `json:"url_chung_tu"`          // Cột H: Link ảnh chụp phiếu/hóa đơn.
	TongTienPhieu      float64 `json:"tong_tien_phieu"`       // Cột I: = SUM(ThanhTienDong) bên ChiTiet.
	DaThanhToan        float64 `json:"da_thanh_toan"`         // Cột J: Số tiền đã trả cho NCC đợt này.
	ConNo              float64 `json:"con_no"`                // Cột K: = TongTienPhieu - DaThanhToan.
	PhuongThucTT       string  `json:"phuong_thuc_thanh_toan"`// Cột L: TM (Tiền mặt), CK (Chuyển khoản).
	TrangThaiThanhToan string  `json:"trang_thai_thanh_toan"` // Cột M: UNPAID, PARTIAL, PAID.
	GhiChu             string  `json:"ghi_chu"`               // Cột N: Note chung.
	NguoiTao           string  `json:"nguoi_tao"`             // Cột O.
	NgayTao            string  `json:"ngay_tao"`              // Cột P.
	NgayCapNhat        string  `json:"ngay_cap_nhat"`         // Cột Q.
}

// =================================================================================
// 3. CHI TIẾT PHIẾU NHẬP (Import Details) - Danh sách hàng hóa
// =================================================================================
type ChiTietPhieuNhap struct {
	MaPhieuNhap    string  `json:"ma_phieu_nhap"`   // Cột A (FK): Link về Header PHIEU_NHAP.
	MaSanPham      string  `json:"ma_san_pham"`     // Cột B (FK): Link về SAN_PHAM.
	TenSanPham     string  `json:"ten_san_pham"`    // Cột C (Snapshot): Lưu cứng tên lúc nhập.
	DonVi          string  `json:"don_vi"`          // Cột D: Đơn vị tính.
	SoLuong        int     `json:"so_luong"`        // Cột E: Số lượng nhập.
	DonGiaNhap     float64 `json:"don_gia_nhap"`    // Cột F: Giá nhập gốc (Chưa VAT).
	VatPercent     float64 `json:"vat_percent"`     // Cột G: % Thuế (0, 5, 8, 10).
	GiaSauVat      float64 `json:"gia_sau_vat"`     // Cột H: = DonGiaNhap * (1 + VAT/100).
	ChietKhauDong  float64 `json:"chiet_khau_dong"` // Cột I: Tiền giảm giá riêng cho dòng này.
	ThanhTienDong  float64 `json:"thanh_tien_dong"` // Cột J: = (SoLuong * GiaSauVat) - ChietKhauDong.
	GiaVonThucTe   float64 `json:"gia_von_thuc_te"` // Cột K: Giá sau khi phân bổ chi phí vận chuyển (nếu có logic tính).
	BaoHanhThang   int     `json:"bao_hanh_thang"`  // Cột L: Bảo hành áp dụng cho lô hàng này.
	GhiChuDong     string  `json:"ghi_chu_dong"`    // Cột M: Ghi chú chi tiết dòng.
}

// =================================================================================
// 4. NHÀ CUNG CẤP (Suppliers) - Đối tác & Công nợ phải trả
// =================================================================================
type NhaCungCap struct {
	MaNhaCungCap string  `json:"ma_nha_cung_cap"` // Cột A (PK).
	TenNhaCungCap string `json:"ten_nha_cung_cap"`// Cột B.
	DienThoai    string  `json:"dien_thoai"`      // Cột C.
	Email        string  `json:"email"`           // Cột D.
	DiaChi       string  `json:"dia_chi"`         // Cột E.
	MaSoThue     string  `json:"ma_so_thue"`      // Cột F.
	NguoiLienHe  string  `json:"nguoi_lien_he"`   // Cột G.
	NganHang     string  `json:"ngan_hang"`       // Cột H: Số tài khoản ngân hàng.
	NoCanTra     float64 `json:"no_can_tra"`      // Cột I: Tổng nợ hiện tại (= Sum PhieuNhap.ConNo).
	TongMua      float64 `json:"tong_mua"`        // Cột J: Tổng doanh số đã mua.
	HanMucCongNo float64 `json:"han_muc_cong_no"` // Cột K: Giới hạn nợ cho phép.
	TrangThai    int     `json:"trang_thai"`      // Cột L: 1=Active, 0=Block.
	GhiChu       string  `json:"ghi_chu"`         // Cột M.
	NguoiTao     string  `json:"nguoi_tao"`       // Cột N.
	NgayTao      string  `json:"ngay_tao"`        // Cột O.
	NgayCapNhat  string  `json:"ngay_cap_nhat"`   // Cột P.
}

// =================================================================================
// 5. SERIAL / IMEI (Tracking) - Quản lý từng sản phẩm & Bảo hành
// =================================================================================
type SerialSanPham struct {
	SerialImei           string `json:"serial_imei"`              // Cột A (PK): Số Serial/IMEI duy nhất.
	MaSanPham            string `json:"ma_san_pham"`              // Cột B (FK).
	MaNhaCungCap         string `json:"ma_nha_cung_cap"`          // Cột C (FK): Để biết bảo hành đầu vào ở đâu.
	MaPhieuNhap          string `json:"ma_phieu_nhap"`            // Cột D (FK).
	MaPhieuXuat          string `json:"ma_phieu_xuat"`            // Cột E (FK): Rỗng nếu chưa bán.
	TrangThai            int    `json:"trang_thai"`               // Cột F: 0=Kho, 1=Đã bán, 2=Lỗi, 3=Chuyển kho.
	BaoHanhNhaCungCap    int    `json:"bao_hanh_nha_cung_cap"`    // Cột G: Số tháng NCC bảo hành.
	HanBaoHanhNhaCungCap string `json:"han_bao_hanh_nha_cung_cap"`// Cột H: Ngày hết hạn bảo hành gốc.
	MaKhachHangHienTai   string `json:"ma_khach_hang_hien_tai"`   // Cột I: Người đang sở hữu máy.
	NgayXuatKho          string `json:"ngay_xuat_kho"`            // Cột J: Ngày bán thực tế.
	KichHoatBaoHanhKhach string `json:"kich_hoat_bao_hanh_khach"` // Cột K: Ngày bắt đầu tính BH cho khách.
	HanBaoHanhKhach      string `json:"han_bao_hanh_khach"`       // Cột L: Ngày hết trách nhiệm với khách.
	MaKho                string `json:"ma_kho"`                   // Cột M: Vị trí hiện tại.
	GhiChu               string `json:"ghi_chu"`                  // Cột N.
	NgayCapNhat          string `json:"ngay_cap_nhat"`            // Cột O.
}

// =================================================================================
// 6. PHIẾU XUẤT (Sales Order) - Đơn bán hàng / Web Order
// =================================================================================
type PhieuXuat struct {
	MaPhieuXuat        string  `json:"ma_phieu_xuat"`         // Cột A (PK): VD: PX2410-001.
	LoaiXuat           string  `json:"loai_xuat"`             // Cột B: BAN_LE, BAN_BUON, XUAT_HUY.
	NgayXuat           string  `json:"ngay_xuat"`             // Cột C.
	MaKho              string  `json:"ma_kho"`                // Cột D.
	MaKhachHang        string  `json:"ma_khach_hang"`         // Cột E (FK).
	TrangThai          string  `json:"trang_thai"`            // Cột F: MOI, DANG_GIAO, HOAN_THANH, HUY.
	MaVoucher          string  `json:"ma_voucher"`            // Cột G (FK): Mã giảm giá khách dùng.
	TienGiamVoucher    float64 `json:"tien_giam_voucher"`     // Cột H: Số tiền được giảm.
	TongTienPhieu      float64 `json:"tong_tien_phieu"`       // Cột I: = Tổng hàng + Ship - Voucher.
	LinkChungTu        string  `json:"link_chung_tu"`         // Cột J.
	DaThu              float64 `json:"da_thu"`                // Cột K: Khách đã trả.
	ConNo              float64 `json:"con_no"`                // Cột L: = TongTien - DaThu.
	PhuongThucTT       string  `json:"phuong_thuc_thanh_toan"`// Cột M: COD, BANK, MOMO.
	TrangThaiTT        string  `json:"trang_thai_thanh_toan"` // Cột N.
	PhiVanChuyen       float64 `json:"phi_van_chuyen"`        // Cột O: Phí ship thu của khách.
	NguonDonHang       string  `json:"nguon_don_hang"`        // Cột P: WEB, POS, SOCIAL.
	ThongTinGiaoHang   string  `json:"thong_tin_giao_hang"`   // Cột Q (JSON): {Ten, SDT, DiaChi} người nhận.
	GhiChu             string  `json:"ghi_chu"`               // Cột R.
	NguoiTao           string  `json:"nguoi_tao"`             // Cột S.
	NgayTao            string  `json:"ngay_tao"`              // Cột T.
	NgayCapNhat        string  `json:"ngay_cap_nhat"`         // Cột U.
}

// =================================================================================
// 7. CHI TIẾT PHIẾU XUẤT (Sales Details)
// =================================================================================
type ChiTietPhieuXuat struct {
	MaPhieuXuat   string  `json:"ma_phieu_xuat"`   // Cột A (FK).
	MaSanPham     string  `json:"ma_san_pham"`     // Cột B (FK).
	TenSanPham    string  `json:"ten_san_pham"`    // Cột C (Snapshot).
	DonVi         string  `json:"don_vi"`          // Cột D.
	SoLuong       int     `json:"so_luong"`        // Cột E.
	DonGiaBan     float64 `json:"don_gia_ban"`     // Cột F: Giá bán ra (Chưa VAT).
	VatPercent    float64 `json:"vat_percent"`     // Cột G.
	GiaSauVat     float64 `json:"gia_sau_vat"`     // Cột H.
	ChietKhauDong float64 `json:"chiet_khau_dong"` // Cột I.
	ThanhTienDong float64 `json:"thanh_tien_dong"` // Cột J: = (SL * GiaSauVat) - CK.
	GiaVon        float64 `json:"gia_von"`         // Cột K: Giá vốn tại thời điểm bán (để tính lãi).
	BaoHanhThang  int     `json:"bao_hanh_thang"`  // Cột L: Cam kết BH cho đơn này.
	GhiChuDong    string  `json:"ghi_chu_dong"`    // Cột M.
}

// =================================================================================
// 8. HÓA ĐƠN ĐIỆN TỬ (VAT Invoice) - Dữ liệu xuất thuế
// =================================================================================
type HoaDon struct {
	MaHoaDon           string  `json:"ma_hoa_don"`            // Cột A (PK).
	MaTraCuu           string  `json:"ma_tra_cuu"`            // Cột B: Mã bí mật để khách tra cứu.
	XmlUrl             string  `json:"xml_url"`               // Cột C: Link file XML pháp lý.
	LoaiHoaDon         string  `json:"loai_hoa_don"`          // Cột D: GTGT, TRUC_TIEP.
	MaPhieuXuat        string  `json:"ma_phieu_xuat"`         // Cột E (FK): Link về đơn hàng gốc.
	MaKhachHang        string  `json:"ma_khach_hang"`         // Cột F.
	NgayHoaDon         string  `json:"ngay_hoa_don"`          // Cột G.
	KyHieu             string  `json:"ky_hieu"`               // Cột H: VD: 1C24TYY.
	SoHoaDon           string  `json:"so_hoa_don"`            // Cột I: VD: 0001234.
	MauSo              string  `json:"mau_so"`                // Cột J: VD: 1/001.
	LinkChungTu        string  `json:"link_chung_tu"`         // Cột K: Bản thể hiện PDF.
	TongTienPhieu      float64 `json:"tong_tien_phieu"`       // Cột L: Tổng chưa thuế.
	TongVat            float64 `json:"tong_vat"`              // Cột M: Tổng tiền thuế.
	TongTienSauVat     float64 `json:"tong_tien_sau_vat"`     // Cột N: Tổng thanh toán.
	TrangThai          string  `json:"trang_thai"`            // Cột O: PENDING, SIGNED (Đã ký), CANCEL.
	TrangThaiThanhToan string  `json:"trang_thai_thanh_toan"` // Cột P.
	GhiChu             string  `json:"ghi_chu"`               // Cột Q.
	NguoiTao           string  `json:"nguoi_tao"`             // Cột R.
	NgayTao            string  `json:"ngay_tao"`              // Cột S.
	NgayCapNhat        string  `json:"ngay_cap_nhat"`         // Cột T.
}

type HoaDonChiTiet struct {
	MaHoaDon   string  `json:"ma_hoa_don"`   // Cột A.
	MaSanPham  string  `json:"ma_san_pham"`  // Cột B.
	TenSanPham string  `json:"ten_san_pham"` // Cột C.
	DonVi      string  `json:"don_vi"`       // Cột D.
	SoLuong    int     `json:"so_luong"`     // Cột E.
	DonGiaBan  float64 `json:"don_gia_ban"`  // Cột F.
	VatPercent float64 `json:"vat_percent"`  // Cột G.
	TienVat    float64 `json:"tien_vat"`     // Cột H.
	ThanhTien  float64 `json:"thanh_tien"`   // Cột I.
}

// =================================================================================
// 10. KHÁCH HÀNG (CRM) - Thông tin người mua
// =================================================================================
type KhachHang struct {
	MaKhachHang   string  `json:"ma_khach_hang"`   // Cột A (PK).
	UserName      string  `json:"user_name"`       // Cột B: Tài khoản đăng nhập Web.
	PasswordHash  string  `json:"-"`               // Cột C: Mật khẩu mã hóa (Không trả về JSON).
	LoaiKhachHang string  `json:"loai_khach_hang"` // Cột D: LE, DAI_LY, VIP.
	TenKhachHang  string  `json:"ten_khach_hang"`  // Cột E.
	DienThoai     string  `json:"dien_thoai"`      // Cột F.
	Email         string  `json:"email"`           // Cột G.
	UrlFb         string  `json:"url_fb"`          // Cột H.
	Zalo          string  `json:"zalo"`            // Cột I.
	UrlTele       string  `json:"url_tele"`        // Cột J.
	UrlTiktok     string  `json:"url_tiktok"`      // Cột K.
	DiaChi        string  `json:"dia_chi"`         // Cột L.
	NgaySinh      string  `json:"ngay_sinh"`       // Cột M: YYYY-MM-DD.
	GioiTinh      string  `json:"gioi_tinh"`       // Cột N: NAM, NU, KHAC.
	MaSoThue      string  `json:"ma_so_thue"`      // Cột O: Nếu là khách doanh nghiệp.
	DangNo        float64 `json:"dang_no"`         // Cột P: Tiền khách đang nợ mình.
	TongMua       float64 `json:"tong_mua"`        // Cột Q: Tổng tiền đã mua (để xếp hạng VIP).
	TrangThai     int     `json:"trang_thai"`      // Cột R.
	GhiChu        string  `json:"ghi_chu"`         // Cột S.
	NguoiTao      string  `json:"nguoi_tao"`       // Cột T.
	NgayTao       string  `json:"ngay_tao"`        // Cột U.
	NgayCapNhat   string  `json:"ngay_cap_nhat"`   // Cột V.
}

// =================================================================================
// 11. PHIẾU THU CHI (Cashbook) - Sổ quỹ & Hạch toán lãi lỗ
// =================================================================================
type PhieuThuChi struct {
	MaPhieuThuChi      string  `json:"ma_phieu_thu_chi"`      // Cột A (PK).
	NgayTaoPhieu       string  `json:"ngay_tao_phieu"`        // Cột B.
	LoaiPhieu          string  `json:"loai_phieu"`            // Cột C: THU (In), CHI (Out).
	DoiTuongLoai       string  `json:"doi_tuong_loai"`        // Cột D: KHACH_HANG, NCC, NHAN_VIEN, KHAC.
	DoiTuongID         string  `json:"doi_tuong_id"`          // Cột E: ID đối tượng tương ứng.
	HangMucThuChi      string  `json:"hang_muc_thu_chi"`      // Cột F: THANH_TOAN_HANG, DIEN_NUOC, LUONG, THUE...
	CoHoaDonDo         bool    `json:"co_hoa_don_do"`         // Cột G: True (Có VAT đầu vào), False.
	MaChungTuThamChieu string  `json:"ma_chung_tu_tham_chieu"`// Cột H: Link tới Mã Phiếu Nhập/Xuất.
	SoTien             float64 `json:"so_tien"`               // Cột I.
	PhuongThucTT       string  `json:"phuong_thuc_thanh_toan"`// Cột J.
	TrangThaiDuyet     int     `json:"trang_thai_duyet"`      // Cột K: 0 (Chờ), 1 (Đã duyệt).
	NguoiDuyet         string  `json:"nguoi_duyet"`           // Cột L: Admin nào duyệt.
	GhiChu             string  `json:"ghi_chu"`               // Cột M.
	NguoiTao           string  `json:"nguoi_tao"`             // Cột N.
	NgayTao            string  `json:"ngay_tao"`              // Cột O.
	NgayCapNhat        string  `json:"ngay_cap_nhat"`         // Cột P.
}

// =================================================================================
// 12. PHIẾU BẢO HÀNH (Warranty) - Dịch vụ sửa chữa
// =================================================================================
type PhieuBaoHanh struct {
	MaPhieuBaoHanh    string  `json:"ma_phieu_bao_hanh"`   // Cột A (PK).
	LoaiPhieu         string  `json:"loai_phieu"`          // Cột B: BAO_HANH (Free), SUA_CHUA (Tính phí).
	SerialImei        string  `json:"serial_imei"`         // Cột C (FK).
	MaSanPham         string  `json:"ma_san_pham"`         // Cột D (FK).
	MaKhachHang       string  `json:"ma_khach_hang"`       // Cột E (FK).
	TenNguoiGui       string  `json:"ten_nguoi_gui"`       // Cột F: Người mang máy đến.
	SdtNguoiGui       string  `json:"sdt_nguoi_gui"`       // Cột G.
	NgayNhan          string  `json:"ngay_nhan"`           // Cột H.
	TinhTrangLoi      string  `json:"tinh_trang_loi"`      // Cột I: Mô tả lỗi (Vỡ màn, không lên nguồn...).
	HinhThuc          string  `json:"hinh_thuc"`           // Cột J: Tai cua hang / Gui buu dien.
	TrangThai         int     `json:"trang_thai"`          // Cột K: 0:Mới, 1:Đang sửa, 2:Xong, 3:Đã trả.
	NgayTraDuKien     string  `json:"ngay_tra_du_kien"`    // Cột L.
	NgayTraThucTe     string  `json:"ngay_tra_thuc_te"`    // Cột M.
	ChiPhiSua         float64 `json:"chi_phi_sua"`         // Cột N: Chi phí nội bộ (linh kiện + công).
	PhiThuKhach       float64 `json:"phi_thu_khach"`       // Cột O: Tiền thu của khách.
	KetQuaSuaChua     string  `json:"ket_qua_sua_chua"`    // Cột P: Kết luận kỹ thuật.
	LinhKienThayThe   string  `json:"linh_kien_thay_the"`  // Cột Q: Danh sách SKU linh kiện đã dùng (JSON).
	MaNhanVienKyThuat string  `json:"ma_nhan_vien_ky_thuat"`// Cột R: Ai sửa.
	GhiChu            string  `json:"ghi_chu"`             // Cột S.
	NguoiTao          string  `json:"nguoi_tao"`           // Cột T.
	NgayTao           string  `json:"ngay_tao"`            // Cột U.
	NgayCapNhat       string  `json:"ngay_cap_nhat"`       // Cột V.
}

// =================================================================================
// 13, 14. DANH MỤC & THƯƠNG HIỆU (Categories & Brands)
// =================================================================================
type DanhMuc struct {
	MaDanhMuc    string `json:"ma_danh_muc"`     // Cột A (PK).
	ThuTuHienThi int    `json:"thu_tu_hien_thi"` // Cột B: Số nhỏ hiện trước (0, 1, 2...).
	TenDanhMuc   string `json:"ten_danh_muc"`    // Cột C.
	Slug         string `json:"slug"`            // Cột D: URL thân thiện.
	MaDanhMucCha string `json:"ma_danh_muc_cha"` // Cột E: Parent ID (cho menu đa cấp).
}

type ThuongHieu struct {
	MaThuongHieu  string `json:"ma_thuong_hieu"`  // Cột A (PK).
	TenThuongHieu string `json:"ten_thuong_hieu"` // Cột B.
	LogoUrl       string `json:"logo_url"`        // Cột C.
}

// =================================================================================
// 15. NHÂN VIÊN (Staff) - Quản trị nội bộ
// =================================================================================
type NhanVien struct {
	MaNhanVien      string `json:"ma_nhan_vien"`       // Cột A (PK).
	TenDangNhap     string `json:"ten_dang_nhap"`      // Cột B.
	Email           string `json:"email"`              // Cột C.
	MatKhauHash     string `json:"-"`                  // Cột D: Mật khẩu Admin/POS.
	HoTen           string `json:"ho_ten"`             // Cột E.
	ChucVu          string `json:"chuc_vu"`            // Cột F: ADMIN, SALE, KHO, KETOAN.
	MaPin           string `json:"-"`                  // Cột G: Mã số nhanh cho POS.
	Cookie          string `json:"-"`                  // Cột H: Token phiên đăng nhập.
	CookieExpired   string `json:"cookie_expired"`     // Cột I.
	VaiTroQuyenHan  string `json:"vai_tro_quyen_han"`  // Cột J: JSON chi tiết quyền (VD: {"can_view_report": true}).
	TrangThai       int    `json:"trang_thai"`         // Cột K.
	LanDangNhapCuoi string `json:"lan_dang_nhap_cuoi"` // Cột L: Log bảo mật.
}

// =================================================================================
// 16. KHUYẾN MÃI (Voucher)
// =================================================================================
type KhuyenMai struct {
	MaVoucher      string  `json:"ma_voucher"`       // Cột A (PK): Mã nhập (VD: SALE50).
	TenChuongTrinh string  `json:"ten_chuong_trinh"` // Cột B.
	LoaiGiam       string  `json:"loai_giam"`        // Cột C: PERCENT (%), AMOUNT (Số tiền).
	GiaTriGiam     float64 `json:"gia_tri_giam"`     // Cột D: Giá trị (VD: 10 (%), 50000 (VND)).
	DonToThieu     float64 `json:"don_to_thieu"`     // Cột E: Giá trị đơn tối thiểu để áp dụng.
	NgayBatDau     string  `json:"ngay_bat_dau"`     // Cột F.
	NgayKetThuc    string  `json:"ngay_ket_thuc"`    // Cột G.
	SoLuongConLai  int     `json:"so_luong_con_lai"` // Cột H: Giảm dần khi có người dùng.
	TrangThai      int     `json:"trang_thai"`       // Cột I.
}

// =================================================================================
// 17. CẤU HÌNH WEB (Settings) - Banner, Popup...
// =================================================================================
type CauHinhWeb struct {
	MaCauHinh string `json:"ma_cau_hinh"` // Cột A (PK): VD: BANNER_HOME.
	GiaTri    string `json:"gia_tri"`     // Cột B: URL ảnh hoặc JSON Config.
	MoTa      string `json:"mo_ta"`       // Cột C.
	TrangThai int    `json:"trang_thai"`  // Cột D.
}
