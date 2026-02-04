package nghiep_vu

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time" // Dùng để sleep tránh rate limit

	"app/cau_hinh"
	"app/kho_du_lieu"
	"app/mo_hinh"
)

// =================================================================================
// 1. ĐỊNH NGHĨA KHO DỮ LIỆU (STORE STRUCTS)
// Mỗi bảng sẽ có một struct Store riêng để quản lý dữ liệu và Key lock
// =================================================================================

type KhoSanPhamStore struct {
	DuLieu   map[string]mo_hinh.SanPham
	DanhSach []mo_hinh.SanPham
	TenKey   string
}
type KhoDanhMucStore struct {
	DuLieu map[string]mo_hinh.DanhMuc
	TenKey string
}
type KhoThuongHieuStore struct {
	DuLieu map[string]mo_hinh.ThuongHieu
	TenKey string
}
type KhoNhaCungCapStore struct {
	DuLieu map[string]mo_hinh.NhaCungCap
	TenKey string
}
type KhoKhachHangStore struct {
	DuLieu map[string]mo_hinh.KhachHang
	TenKey string
}
type KhoNhanVienStore struct {
	DuLieu map[string]mo_hinh.NhanVien
	TenKey string
}
type KhoPhieuNhapStore struct {
	DuLieu   map[string]mo_hinh.PhieuNhap
	DanhSach []mo_hinh.PhieuNhap
	TenKey   string
}
type KhoChiTietPhieuNhapStore struct {
	DanhSach []mo_hinh.ChiTietPhieuNhap
	TenKey   string
}
type KhoPhieuXuatStore struct {
	DuLieu   map[string]mo_hinh.PhieuXuat
	DanhSach []mo_hinh.PhieuXuat
	TenKey   string
}
type KhoChiTietPhieuXuatStore struct {
	DanhSach []mo_hinh.ChiTietPhieuXuat
	TenKey   string
}
type KhoSerialStore struct {
	DuLieu map[string]mo_hinh.SerialSanPham
	TenKey string
}
type KhoHoaDonStore struct {
	DuLieu map[string]mo_hinh.HoaDon
	TenKey string
}
type KhoHoaDonChiTietStore struct {
	DanhSach []mo_hinh.HoaDonChiTiet
	TenKey   string
}
type KhoPhieuThuChiStore struct {
	DuLieu   map[string]mo_hinh.PhieuThuChi
	DanhSach []mo_hinh.PhieuThuChi
	TenKey   string
}
type KhoPhieuBaoHanhStore struct {
	DuLieu   map[string]mo_hinh.PhieuBaoHanh
	DanhSach []mo_hinh.PhieuBaoHanh
	TenKey   string
}
type KhoKhuyenMaiStore struct {
	DuLieu map[string]mo_hinh.KhuyenMai
	TenKey string
}
type KhoCauHinhWebStore struct {
	DuLieu map[string]mo_hinh.CauHinhWeb
	TenKey string
}

// =================================================================================
// 2. BIẾN TOÀN CỤC (GLOBAL CACHE) - ĐỦ 17 BẢNG
// =================================================================================
var (
	CacheSanPham         *KhoSanPhamStore
	CacheDanhMuc         *KhoDanhMucStore
	CacheThuongHieu      *KhoThuongHieuStore
	CacheNhaCungCap      *KhoNhaCungCapStore
	CacheKhachHang       *KhoKhachHangStore
	CacheNhanVien        *KhoNhanVienStore
	CachePhieuNhap       *KhoPhieuNhapStore
	CacheChiTietNhap     *KhoChiTietPhieuNhapStore
	CachePhieuXuat       *KhoPhieuXuatStore
	CacheChiTietXuat     *KhoChiTietPhieuXuatStore
	CacheSerial          *KhoSerialStore
	CacheHoaDon          *KhoHoaDonStore
	CacheHoaDonChiTiet   *KhoHoaDonChiTietStore
	CachePhieuThuChi     *KhoPhieuThuChiStore
	CachePhieuBaoHanh    *KhoPhieuBaoHanhStore
	CacheKhuyenMai       *KhoKhuyenMaiStore
	CacheCauHinhWeb      *KhoCauHinhWebStore
)

// Hàm tạo Key chuẩn: SpreadsheetID__SheetName
func TaoKeyCache(tenSheet string) string {
	idSheet := cau_hinh.BienCauHinh.IdFileSheet
	return fmt.Sprintf("%s__%s", strings.TrimSpace(idSheet), tenSheet)
}

// =================================================================================
// 3. KHỞI TẠO VÀ NẠP DỮ LIỆU (INITIALIZE & LOAD)
// =================================================================================
func KhoiTaoBoNho() {
	log.Println("--- [CACHE] Bắt đầu khởi tạo bộ nhớ cho 17 bảng ---")
	
	// Bước 1: Khởi tạo Struct rỗng và Key
	khoiTaoCacStore()

	// Bước 2: Nạp dữ liệu (Chia làm 3 đợt để tránh Google chặn API)
	var wg sync.WaitGroup

	// --- ĐỢT 1: Master Data (Danh mục, Thương hiệu, SP, KH, NCC, NV, Cấu hình) ---
	log.Println(">> Đợt 1: Nạp Master Data...")
	wg.Add(7)
	go func() { defer wg.Done(); napDanhMuc() }()
	go func() { defer wg.Done(); napThuongHieu() }()
	go func() { defer wg.Done(); napSanPham() }()
	go func() { defer wg.Done(); napKhachHang() }()
	go func() { defer wg.Done(); napNhaCungCap() }()
	go func() { defer wg.Done(); napNhanVien() }()
	go func() { defer wg.Done(); napCauHinhWeb() }()
	wg.Wait()
	
	time.Sleep(1 * time.Second) // Nghỉ 1 xíu

	// --- ĐỢT 2: Transaction Chính (Nhập, Xuất, Serial, Khuyến mãi) ---
	log.Println(">> Đợt 2: Nạp Giao dịch chính...")
	wg.Add(6)
	go func() { defer wg.Done(); napPhieuNhap() }()
	go func() { defer wg.Done(); napChiTietPhieuNhap() }()
	go func() { defer wg.Done(); napPhieuXuat() }()
	go func() { defer wg.Done(); napChiTietPhieuXuat() }()
	go func() { defer wg.Done(); napSerial() }()
	go func() { defer wg.Done(); napKhuyenMai() }()
	wg.Wait()

	time.Sleep(1 * time.Second)

	// --- ĐỢT 3: Tài chính & CSKH (Hóa đơn, Thu chi, Bảo hành) ---
	log.Println(">> Đợt 3: Nạp Tài chính & CSKH...")
	wg.Add(4)
	go func() { defer wg.Done(); napHoaDon() }()
	go func() { defer wg.Done(); napHoaDonChiTiet() }()
	go func() { defer wg.Done(); napPhieuThuChi() }()
	go func() { defer wg.Done(); napPhieuBaoHanh() }()
	wg.Wait()

	log.Println("--- [CACHE] HOÀN TẤT NẠP 100% DỮ LIỆU ---")
}

func khoiTaoCacStore() {
	CacheSanPham = &KhoSanPhamStore{DuLieu: make(map[string]mo_hinh.SanPham), TenKey: TaoKeyCache("SAN_PHAM")}
	CacheDanhMuc = &KhoDanhMucStore{DuLieu: make(map[string]mo_hinh.DanhMuc), TenKey: TaoKeyCache("DANH_MUC")}
	CacheThuongHieu = &KhoThuongHieuStore{DuLieu: make(map[string]mo_hinh.ThuongHieu), TenKey: TaoKeyCache("THUONG_HIEU")}
	CacheNhaCungCap = &KhoNhaCungCapStore{DuLieu: make(map[string]mo_hinh.NhaCungCap), TenKey: TaoKeyCache("NHA_CUNG_CAP")}
	CacheKhachHang = &KhoKhachHangStore{DuLieu: make(map[string]mo_hinh.KhachHang), TenKey: TaoKeyCache("KHACH_HANG")}
	CacheNhanVien = &KhoNhanVienStore{DuLieu: make(map[string]mo_hinh.NhanVien), TenKey: TaoKeyCache("NHAN_VIEN")}
	
	CachePhieuNhap = &KhoPhieuNhapStore{DuLieu: make(map[string]mo_hinh.PhieuNhap), TenKey: TaoKeyCache("PHIEU_NHAP")}
	CacheChiTietNhap = &KhoChiTietPhieuNhapStore{TenKey: TaoKeyCache("CHI_TIET_PHIEU_NHAP")}
	
	CachePhieuXuat = &KhoPhieuXuatStore{DuLieu: make(map[string]mo_hinh.PhieuXuat), TenKey: TaoKeyCache("PHIEU_XUAT")}
	CacheChiTietXuat = &KhoChiTietPhieuXuatStore{TenKey: TaoKeyCache("CHI_TIET_PHIEU_XUAT")}
	
	CacheSerial = &KhoSerialStore{DuLieu: make(map[string]mo_hinh.SerialSanPham), TenKey: TaoKeyCache("SERIAL_SAN_PHAM")}
	CacheKhuyenMai = &KhoKhuyenMaiStore{DuLieu: make(map[string]mo_hinh.KhuyenMai), TenKey: TaoKeyCache("KHUYEN_MAI")}
	CacheCauHinhWeb = &KhoCauHinhWebStore{DuLieu: make(map[string]mo_hinh.CauHinhWeb), TenKey: TaoKeyCache("CAU_HINH_WEB")}
	
	CacheHoaDon = &KhoHoaDonStore{DuLieu: make(map[string]mo_hinh.HoaDon), TenKey: TaoKeyCache("HOA_DON")}
	CacheHoaDonChiTiet = &KhoHoaDonChiTietStore{TenKey: TaoKeyCache("HOA_DON_CHI_TIET")}
	CachePhieuThuChi = &KhoPhieuThuChiStore{DuLieu: make(map[string]mo_hinh.PhieuThuChi), TenKey: TaoKeyCache("PHIEU_THU_CHI")}
	CachePhieuBaoHanh = &KhoPhieuBaoHanhStore{DuLieu: make(map[string]mo_hinh.PhieuBaoHanh), TenKey: TaoKeyCache("PHIEU_BAO_HANH")}
}

// =================================================================================
// 4. CÁC HÀM NẠP DỮ LIỆU CHI TIẾT (IMPLEMENTATION)
// =================================================================================

// --- Helper nạp chung (để giảm lặp code) ---
func loadSheetData(sheetName string, keyCache string) ([][]interface{}, error) {
	// log.Printf("Đang đọc: %s...", sheetName)
	duLieu, err := kho_du_lieu.DocToanBoSheet(sheetName)
	if err != nil {
		log.Printf("LỖI ĐỌC %s: %v", sheetName, err)
		return nil, err
	}
	
	// Khóa GHI trước khi nạp
	khoa := BoQuanLyKhoa.LayKhoa(keyCache)
	khoa.Lock()
	// Lưu ý: Unlock sẽ được gọi ở hàm cha sau khi xử lý xong
	return duLieu, nil
}

// 1. SAN_PHAM
func napSanPham() {
	raw, err := loadSheetData("SAN_PHAM", CacheSanPham.TenKey)
	if err != nil { return }
	defer BoQuanLyKhoa.LayKhoa(CacheSanPham.TenKey).Unlock()

	for i, r := range raw {
		if i < mo_hinh.DongBatDauDuLieu { continue }
		if len(r) <= mo_hinh.CotSP_MaSanPham || layString(r, mo_hinh.CotSP_MaSanPham) == "" { continue }
		
		item := mo_hinh.SanPham{
			MaSanPham:    layString(r, mo_hinh.CotSP_MaSanPham),
			TenSanPham:   layString(r, mo_hinh.CotSP_TenSanPham),
			TenRutGon:    layString(r, mo_hinh.CotSP_TenRutGon),
			Sku:          layString(r, mo_hinh.CotSP_Sku),
			MaDanhMuc:    layString(r, mo_hinh.CotSP_MaDanhMuc),
			MaThuongHieu: layString(r, mo_hinh.CotSP_MaThuongHieu),
			DonVi:        layString(r, mo_hinh.CotSP_DonVi),
			MauSac:       layString(r, mo_hinh.CotSP_MauSac),
			UrlHinhAnh:   layString(r, mo_hinh.CotSP_UrlHinhAnh),
			ThongSo:      layString(r, mo_hinh.CotSP_ThongSo),
			MoTaChiTiet:  layString(r, mo_hinh.CotSP_MoTaChiTiet),
			BaoHanhThang: layInt(r, mo_hinh.CotSP_BaoHanhThang),
			TinhTrang:    layString(r, mo_hinh.CotSP_TinhTrang),
			TrangThai:    layInt(r, mo_hinh.CotSP_TrangThai),
			GiaBanLe:     layFloat(r, mo_hinh.CotSP_GiaBanLe),
			GhiChu:       layString(r, mo_hinh.CotSP_GhiChu),
			NguoiTao:     layString(r, mo_hinh.CotSP_NguoiTao),
			NgayTao:      layString(r, mo_hinh.CotSP_NgayTao),
			NgayCapNhat:  layString(r, mo_hinh.CotSP_NgayCapNhat),
		}
		CacheSanPham.DuLieu[item.MaSanPham] = item
		CacheSanPham.DanhSach = append(CacheSanPham.DanhSach, item)
	}
}

// 2. DANH_MUC
func napDanhMuc() {
	raw, err := loadSheetData("DANH_MUC", CacheDanhMuc.TenKey)
	if err != nil { return }
	defer BoQuanLyKhoa.LayKhoa(CacheDanhMuc.TenKey).Unlock()

	for i, r := range raw {
		if i < mo_hinh.DongBatDauDuLieu { continue }
		if len(r) <= mo_hinh.CotDM_MaDanhMuc || layString(r, mo_hinh.CotDM_MaDanhMuc) == "" { continue }
		
		item := mo_hinh.DanhMuc{
			MaDanhMuc:    layString(r, mo_hinh.CotDM_MaDanhMuc),
			ThuTuHienThi: layInt(r, mo_hinh.CotDM_ThuTuHienThi),
			TenDanhMuc:   layString(r, mo_hinh.CotDM_TenDanhMuc),
			Slug:         layString(r, mo_hinh.CotDM_Slug),
			MaDanhMucCha: layString(r, mo_hinh.CotDM_MaDanhMucCha),
		}
		CacheDanhMuc.DuLieu[item.MaDanhMuc] = item
	}
}

// 3. THUONG_HIEU
func napThuongHieu() {
	raw, err := loadSheetData("THUONG_HIEU", CacheThuongHieu.TenKey)
	if err != nil { return }
	defer BoQuanLyKhoa.LayKhoa(CacheThuongHieu.TenKey).Unlock()

	for i, r := range raw {
		if i < mo_hinh.DongBatDauDuLieu { continue }
		if len(r) <= mo_hinh.CotTH_MaThuongHieu || layString(r, mo_hinh.CotTH_MaThuongHieu) == "" { continue }
		
		item := mo_hinh.ThuongHieu{
			MaThuongHieu:  layString(r, mo_hinh.CotTH_MaThuongHieu),
			TenThuongHieu: layString(r, mo_hinh.CotTH_TenThuongHieu),
			LogoUrl:       layString(r, mo_hinh.CotTH_LogoUrl),
		}
		CacheThuongHieu.DuLieu[item.MaThuongHieu] = item
	}
}

// 4. KHACH_HANG
func napKhachHang() {
	raw, err := loadSheetData("KHACH_HANG", CacheKhachHang.TenKey)
	if err != nil { return }
	defer BoQuanLyKhoa.LayKhoa(CacheKhachHang.TenKey).Unlock()

	for i, r := range raw {
		if i < mo_hinh.DongBatDauDuLieu { continue }
		if len(r) <= mo_hinh.CotKH_MaKhachHang || layString(r, mo_hinh.CotKH_MaKhachHang) == "" { continue }

		item := mo_hinh.KhachHang{
			MaKhachHang:   layString(r, mo_hinh.CotKH_MaKhachHang),
			UserName:      layString(r, mo_hinh.CotKH_UserName),
			PasswordHash:  layString(r, mo_hinh.CotKH_PasswordHash),
			LoaiKhachHang: layString(r, mo_hinh.CotKH_LoaiKhachHang),
			TenKhachHang:  layString(r, mo_hinh.CotKH_TenKhachHang),
			DienThoai:     layString(r, mo_hinh.CotKH_DienThoai),
			Email:         layString(r, mo_hinh.CotKH_Email),
			UrlFb:         layString(r, mo_hinh.CotKH_UrlFb),
			DiaChi:        layString(r, mo_hinh.CotKH_DiaChi),
			NgaySinh:      layString(r, mo_hinh.CotKH_NgaySinh),
			GioiTinh:      layString(r, mo_hinh.CotKH_GioiTinh),
			MaSoThue:      layString(r, mo_hinh.CotKH_MaSoThue),
			DangNo:        layFloat(r, mo_hinh.CotKH_DangNo),
			TongMua:       layFloat(r, mo_hinh.CotKH_TongMua),
			TrangThai:     layInt(r, mo_hinh.CotKH_TrangThai),
		}
		CacheKhachHang.DuLieu[item.MaKhachHang] = item
	}
}

// 5. NHA_CUNG_CAP
func napNhaCungCap() {
	raw, err := loadSheetData("NHA_CUNG_CAP", CacheNhaCungCap.TenKey)
	if err != nil { return }
	defer BoQuanLyKhoa.LayKhoa(CacheNhaCungCap.TenKey).Unlock()

	for i, r := range raw {
		if i < mo_hinh.DongBatDauDuLieu { continue }
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
		CacheNhaCungCap.DuLieu[item.MaNhaCungCap] = item
	}
}

// 6. NHAN_VIEN
func napNhanVien() {
	raw, err := loadSheetData("NHAN_VIEN", CacheNhanVien.TenKey)
	if err != nil { return }
	defer BoQuanLyKhoa.LayKhoa(CacheNhanVien.TenKey).Unlock()

	for i, r := range raw {
		if i < mo_hinh.DongBatDauDuLieu { continue }
		if len(r) <= mo_hinh.CotNV_MaNhanVien || layString(r, mo_hinh.CotNV_MaNhanVien) == "" { continue }
		
		item := mo_hinh.NhanVien{
			MaNhanVien:     layString(r, mo_hinh.CotNV_MaNhanVien),
			TenDangNhap:    layString(r, mo_hinh.CotNV_TenDangNhap),
			MatKhauHash:    layString(r, mo_hinh.CotNV_MatKhauHash),
			HoTen:          layString(r, mo_hinh.CotNV_HoTen),
			ChucVu:         layString(r, mo_hinh.CotNV_ChucVu),
			MaPin:          layString(r, mo_hinh.CotNV_MaPin),
			Cookie:         layString(r, mo_hinh.CotNV_Cookie),
			VaiTroQuyenHan: layString(r, mo_hinh.CotNV_VaiTroQuyenHan),
			TrangThai:      layInt(r, mo_hinh.CotNV_TrangThai),
		}
		CacheNhanVien.DuLieu[item.MaNhanVien] = item
	}
}

// 7. PHIEU_XUAT (ĐƠN HÀNG)
func napPhieuXuat() {
	raw, err := loadSheetData("PHIEU_XUAT", CachePhieuXuat.TenKey)
	if err != nil { return }
	defer BoQuanLyKhoa.LayKhoa(CachePhieuXuat.TenKey).Unlock()

	for i, r := range raw {
		if i < mo_hinh.DongBatDauDuLieu { continue }
		if len(r) <= mo_hinh.CotPX_MaPhieuXuat || layString(r, mo_hinh.CotPX_MaPhieuXuat) == "" { continue }

		item := mo_hinh.PhieuXuat{
			MaPhieuXuat:      layString(r, mo_hinh.CotPX_MaPhieuXuat),
			LoaiXuat:         layString(r, mo_hinh.CotPX_LoaiXuat),
			NgayXuat:         layString(r, mo_hinh.CotPX_NgayXuat),
			MaKho:            layString(r, mo_hinh.CotPX_MaKho),
			MaKhachHang:      layString(r, mo_hinh.CotPX_MaKhachHang),
			TrangThai:        layString(r, mo_hinh.CotPX_TrangThai),
			MaVoucher:        layString(r, mo_hinh.CotPX_MaVoucher),
			TienGiamVoucher:  layFloat(r, mo_hinh.CotPX_TienGiamVoucher),
			TongTienPhieu:    layFloat(r, mo_hinh.CotPX_TongTienPhieu),
			DaThu:            layFloat(r, mo_hinh.CotPX_DaThu),
			ConNo:            layFloat(r, mo_hinh.CotPX_ConNo),
			PhuongThucThanhToan:     layString(r, mo_hinh.CotPX_PhuongThucThanhToan),
			PhiVanChuyen:     layFloat(r, mo_hinh.CotPX_PhiVanChuyen),
			NguonDonHang:     layString(r, mo_hinh.CotPX_NguonDonHang),
			ThongTinGiaoHang: layString(r, mo_hinh.CotPX_ThongTinGiaoHang),
			NguoiTao:         layString(r, mo_hinh.CotPX_NguoiTao),
		}
		CachePhieuXuat.DuLieu[item.MaPhieuXuat] = item
		CachePhieuXuat.DanhSach = append(CachePhieuXuat.DanhSach, item)
	}
}

// 8. CHI_TIET_PHIEU_XUAT
func napChiTietPhieuXuat() {
	raw, err := loadSheetData("CHI_TIET_PHIEU_XUAT", CacheChiTietXuat.TenKey)
	if err != nil { return }
	defer BoQuanLyKhoa.LayKhoa(CacheChiTietXuat.TenKey).Unlock()

	for i, r := range raw {
		if i < mo_hinh.DongBatDauDuLieu { continue }
		if len(r) <= mo_hinh.CotCTPX_MaPhieuXuat || layString(r, mo_hinh.CotCTPX_MaPhieuXuat) == "" { continue }

		item := mo_hinh.ChiTietPhieuXuat{
			MaPhieuXuat:   layString(r, mo_hinh.CotCTPX_MaPhieuXuat),
			MaSanPham:     layString(r, mo_hinh.CotCTPX_MaSanPham),
			TenSanPham:    layString(r, mo_hinh.CotCTPX_TenSanPham),
			SoLuong:       layInt(r, mo_hinh.CotCTPX_SoLuong),
			DonGiaBan:     layFloat(r, mo_hinh.CotCTPX_DonGiaBan),
			ThanhTienDong: layFloat(r, mo_hinh.CotCTPX_ThanhTienDong),
			GiaVon:        layFloat(r, mo_hinh.CotCTPX_GiaVon),
		}
		CacheChiTietXuat.DanhSach = append(CacheChiTietXuat.DanhSach, item)
	}
}

// 9. PHIEU_NHAP
func napPhieuNhap() {
	raw, err := loadSheetData("PHIEU_NHAP", CachePhieuNhap.TenKey)
	if err != nil { return }
	defer BoQuanLyKhoa.LayKhoa(CachePhieuNhap.TenKey).Unlock()

	for i, r := range raw {
		if i < mo_hinh.DongBatDauDuLieu { continue }
		if len(r) <= mo_hinh.CotPN_MaPhieuNhap || layString(r, mo_hinh.CotPN_MaPhieuNhap) == "" { continue }

		item := mo_hinh.PhieuNhap{
			MaPhieuNhap:   layString(r, mo_hinh.CotPN_MaPhieuNhap),
			MaNhaCungCap:  layString(r, mo_hinh.CotPN_MaNhaCungCap),
			MaKho:         layString(r, mo_hinh.CotPN_MaKho),
			NgayNhap:      layString(r, mo_hinh.CotPN_NgayNhap),
			TrangThai:     layString(r, mo_hinh.CotPN_TrangThai),
			TongTienPhieu: layFloat(r, mo_hinh.CotPN_TongTienPhieu),
			DaThanhToan:   layFloat(r, mo_hinh.CotPN_DaThanhToan),
			ConNo:         layFloat(r, mo_hinh.CotPN_ConNo),
		}
		CachePhieuNhap.DuLieu[item.MaPhieuNhap] = item
		CachePhieuNhap.DanhSach = append(CachePhieuNhap.DanhSach, item)
	}
}

// 10. CHI_TIET_PHIEU_NHAP
func napChiTietPhieuNhap() {
	raw, err := loadSheetData("CHI_TIET_PHIEU_NHAP", CacheChiTietNhap.TenKey)
	if err != nil { return }
	defer BoQuanLyKhoa.LayKhoa(CacheChiTietNhap.TenKey).Unlock()

	for i, r := range raw {
		if i < mo_hinh.DongBatDauDuLieu { continue }
		if len(r) <= mo_hinh.CotCTPN_MaPhieuNhap || layString(r, mo_hinh.CotCTPN_MaPhieuNhap) == "" { continue }

		item := mo_hinh.ChiTietPhieuNhap{
			MaPhieuNhap:   layString(r, mo_hinh.CotCTPN_MaPhieuNhap),
			MaSanPham:     layString(r, mo_hinh.CotCTPN_MaSanPham),
			SoLuong:       layInt(r, mo_hinh.CotCTPN_SoLuong),
			DonGiaNhap:    layFloat(r, mo_hinh.CotCTPN_DonGiaNhap),
			ThanhTienDong: layFloat(r, mo_hinh.CotCTPN_ThanhTienDong),
		}
		CacheChiTietNhap.DanhSach = append(CacheChiTietNhap.DanhSach, item)
	}
}

// 11. SERIAL_SAN_PHAM
func napSerial() {
	raw, err := loadSheetData("SERIAL_SAN_PHAM", CacheSerial.TenKey)
	if err != nil { return }
	defer BoQuanLyKhoa.LayKhoa(CacheSerial.TenKey).Unlock()

	for i, r := range raw {
		if i < mo_hinh.DongBatDauDuLieu { continue }
		if len(r) <= mo_hinh.CotSerial_SerialImei || layString(r, mo_hinh.CotSerial_SerialImei) == "" { continue }

		item := mo_hinh.SerialSanPham{
			SerialImei:         layString(r, mo_hinh.CotSerial_SerialImei),
			MaSanPham:          layString(r, mo_hinh.CotSerial_MaSanPham),
			MaPhieuNhap:        layString(r, mo_hinh.CotSerial_MaPhieuNhap),
			MaPhieuXuat:        layString(r, mo_hinh.CotSerial_MaPhieuXuat),
			TrangThai:          layInt(r, mo_hinh.CotSerial_TrangThai),
			MaKhachHangHienTai: layString(r, mo_hinh.CotSerial_MaKhachHangHienTai),
		}
		CacheSerial.DuLieu[item.SerialImei] = item
	}
}

// 12. KHUYEN_MAI
func napKhuyenMai() {
	raw, err := loadSheetData("KHUYEN_MAI", CacheKhuyenMai.TenKey)
	if err != nil { return }
	defer BoQuanLyKhoa.LayKhoa(CacheKhuyenMai.TenKey).Unlock()

	for i, r := range raw {
		if i < mo_hinh.DongBatDauDuLieu { continue }
		if len(r) <= mo_hinh.CotKM_MaVoucher || layString(r, mo_hinh.CotKM_MaVoucher) == "" { continue }

		item := mo_hinh.KhuyenMai{
			MaVoucher:      layString(r, mo_hinh.CotKM_MaVoucher),
			LoaiGiam:       layString(r, mo_hinh.CotKM_LoaiGiam),
			GiaTriGiam:     layFloat(r, mo_hinh.CotKM_GiaTriGiam),
			DonToThieu:     layFloat(r, mo_hinh.CotKM_DonToThieu),
			SoLuongConLai:  layInt(r, mo_hinh.CotKM_SoLuongConLai),
			TrangThai:      layInt(r, mo_hinh.CotKM_TrangThai),
		}
		CacheKhuyenMai.DuLieu[item.MaVoucher] = item
	}
}

// 13. CAU_HINH_WEB
func napCauHinhWeb() {
	raw, err := loadSheetData("CAU_HINH_WEB", CacheCauHinhWeb.TenKey)
	if err != nil { return }
	defer BoQuanLyKhoa.LayKhoa(CacheCauHinhWeb.TenKey).Unlock()

	for i, r := range raw {
		if i < mo_hinh.DongBatDauDuLieu { continue }
		if len(r) <= mo_hinh.CotCH_MaCauHinh || layString(r, mo_hinh.CotCH_MaCauHinh) == "" { continue }

		item := mo_hinh.CauHinhWeb{
			MaCauHinh: layString(r, mo_hinh.CotCH_MaCauHinh),
			GiaTri:    layString(r, mo_hinh.CotCH_GiaTri),
			TrangThai: layInt(r, mo_hinh.CotCH_TrangThai),
		}
		CacheCauHinhWeb.DuLieu[item.MaCauHinh] = item
	}
}

// 14. HOA_DON
func napHoaDon() {
	raw, err := loadSheetData("HOA_DON", CacheHoaDon.TenKey)
	if err != nil { return }
	defer BoQuanLyKhoa.LayKhoa(CacheHoaDon.TenKey).Unlock()

	for i, r := range raw {
		if i < mo_hinh.DongBatDauDuLieu { continue }
		if len(r) <= mo_hinh.CotHD_MaHoaDon || layString(r, mo_hinh.CotHD_MaHoaDon) == "" { continue }

		item := mo_hinh.HoaDon{
			MaHoaDon:      layString(r, mo_hinh.CotHD_MaHoaDon),
			MaTraCuu:      layString(r, mo_hinh.CotHD_MaTraCuu),
			XmlUrl:        layString(r, mo_hinh.CotHD_XmlUrl),
			MaPhieuXuat:   layString(r, mo_hinh.CotHD_MaPhieuXuat),
			TongTienPhieu: layFloat(r, mo_hinh.CotHD_TongTienPhieu),
			TongVat:       layFloat(r, mo_hinh.CotHD_TongVat),
			TrangThai:     layString(r, mo_hinh.CotHD_TrangThai),
		}
		CacheHoaDon.DuLieu[item.MaHoaDon] = item
	}
}

// 15. HOA_DON_CHI_TIET
func napHoaDonChiTiet() {
	raw, err := loadSheetData("HOA_DON_CHI_TIET", CacheHoaDonChiTiet.TenKey)
	if err != nil { return }
	defer BoQuanLyKhoa.LayKhoa(CacheHoaDonChiTiet.TenKey).Unlock()

	for i, r := range raw {
		if i < mo_hinh.DongBatDauDuLieu { continue }
		if len(r) <= mo_hinh.CotHDCT_MaHoaDon || layString(r, mo_hinh.CotHDCT_MaHoaDon) == "" { continue }

		item := mo_hinh.HoaDonChiTiet{
			MaHoaDon:   layString(r, mo_hinh.CotHDCT_MaHoaDon),
			MaSanPham:  layString(r, mo_hinh.CotHDCT_MaSanPham),
			SoLuong:    layInt(r, mo_hinh.CotHDCT_SoLuong),
			ThanhTien:  layFloat(r, mo_hinh.CotHDCT_ThanhTien),
		}
		CacheHoaDonChiTiet.DanhSach = append(CacheHoaDonChiTiet.DanhSach, item)
	}
}

// 16. PHIEU_THU_CHI
func napPhieuThuChi() {
	raw, err := loadSheetData("PHIEU_THU_CHI", CachePhieuThuChi.TenKey)
	if err != nil { return }
	defer BoQuanLyKhoa.LayKhoa(CachePhieuThuChi.TenKey).Unlock()

	for i, r := range raw {
		if i < mo_hinh.DongBatDauDuLieu { continue }
		if len(r) <= mo_hinh.CotPTC_MaPhieuThuChi || layString(r, mo_hinh.CotPTC_MaPhieuThuChi) == "" { continue }

		item := mo_hinh.PhieuThuChi{
			MaPhieuThuChi: layString(r, mo_hinh.CotPTC_MaPhieuThuChi),
			LoaiPhieu:     layString(r, mo_hinh.CotPTC_LoaiPhieu),
			SoTien:        layFloat(r, mo_hinh.CotPTC_SoTien),
			HangMucThuChi: layString(r, mo_hinh.CotPTC_HangMucThuChi),
			TrangThaiDuyet: layInt(r, mo_hinh.CotPTC_TrangThaiDuyet),
		}
		CachePhieuThuChi.DuLieu[item.MaPhieuThuChi] = item
		CachePhieuThuChi.DanhSach = append(CachePhieuThuChi.DanhSach, item)
	}
}

// 17. PHIEU_BAO_HANH
func napPhieuBaoHanh() {
	raw, err := loadSheetData("PHIEU_BAO_HANH", CachePhieuBaoHanh.TenKey)
	if err != nil { return }
	defer BoQuanLyKhoa.LayKhoa(CachePhieuBaoHanh.TenKey).Unlock()

	for i, r := range raw {
		if i < mo_hinh.DongBatDauDuLieu { continue }
		if len(r) <= mo_hinh.CotPBH_MaPhieuBaoHanh || layString(r, mo_hinh.CotPBH_MaPhieuBaoHanh) == "" { continue }

		item := mo_hinh.PhieuBaoHanh{
			MaPhieuBaoHanh: layString(r, mo_hinh.CotPBH_MaPhieuBaoHanh),
			SerialImei:     layString(r, mo_hinh.CotPBH_SerialImei),
			TrangThai:      layInt(r, mo_hinh.CotPBH_TrangThai),
			TinhTrangLoi:   layString(r, mo_hinh.CotPBH_TinhTrangLoi),
		}
		CachePhieuBaoHanh.DuLieu[item.MaPhieuBaoHanh] = item
		CachePhieuBaoHanh.DanhSach = append(CachePhieuBaoHanh.DanhSach, item)
	}
}

// --------------------------------------------------------
// HELPER FUNCTIONS (TIỆN ÍCH - Dán vào cuối file)
// --------------------------------------------------------

func layString(dong []interface{}, index int) string {
	if index >= len(dong) || dong[index] == nil { return "" }
	return fmt.Sprintf("%v", dong[index])
}

// Hàm xử lý số nguyên thông minh (Hiểu được cả 1.000)
func layInt(dong []interface{}, index int) int {
	str := layString(dong, index)
	if str == "" { return 0 }
	
	// Xóa các ký tự phân cách (VD: 1.000 -> 1000)
	str = strings.ReplaceAll(str, ".", "")
	str = strings.ReplaceAll(str, ",", "")
	str = strings.ReplaceAll(str, " ", "")
	
	val, _ := strconv.Atoi(str)
	return val
}

// Hàm xử lý tiền tệ/số thực thông minh (Hiểu được cả 1.500.000 đ)
func layFloat(dong []interface{}, index int) float64 {
	str := layString(dong, index)
	if str == "" { return 0 }

	// 1. Xóa ký tự tiền tệ và khoảng trắng (VD: "1.200.000 đ" -> "1.200.000")
	str = strings.ReplaceAll(str, "đ", "")
	str = strings.ReplaceAll(str, "USD", "")
	str = strings.ReplaceAll(str, " ", "")
	
	// 2. Xóa dấu chấm phân cách ngàn (Kiểu Việt Nam: 1.000.000 -> 1000000)
	str = strings.ReplaceAll(str, ".", "")

	// 3. Xóa dấu phẩy (Kiểu Mỹ: 1,000,000 -> 1000000)
	str = strings.ReplaceAll(str, ",", "")

	// 4. Parse sang số
	val, _ := strconv.ParseFloat(str, 64)
	return val
}

// -------------------------------------------------------------
// BIẾN TOÀN CỤC CHO NHÂN VIÊN
// -------------------------------------------------------------
var CacheNhanVien = struct {
	DuLieu        map[string]*mo_hinh.NhanVien
	SpreadsheetID string
	TenKey        string
}{
	DuLieu: make(map[string]*mo_hinh.NhanVien),
	TenKey: "NHAN_VIEN",
}

// -------------------------------------------------------------
// HÀM NẠP NHÂN VIÊN (Dùng chuẩn mo_hinh.CotNV_...)
// -------------------------------------------------------------
func napNhanVien() {
	// 1. Lấy ID Sheet (Giả sử chung với cấu hình chính)
	CacheNhanVien.SpreadsheetID = cau_hinh.BienCauHinh.IdGoogleSheet

	// 2. Tải dữ liệu thô
	raw, err := loadSheetData("NHAN_VIEN", CacheNhanVien.TenKey)
	if err != nil { return }
	
	// 3. Reset map cũ
	CacheNhanVien.DuLieu = make(map[string]*mo_hinh.NhanVien)

	for i, r := range raw {
		// Bỏ qua dòng tiêu đề (Dòng < 11 theo quy ước DongBatDauDuLieu = 11)
		// Nhưng lưu ý: Với Sheet Nhân viên thường ít dòng, ta hay để tiêu đề dòng 1.
		// Tạm thời tôi theo logic chung là check DongBatDauDuLieu
		if i < mo_hinh.DongBatDauDuLieu { continue }
		
		// Cần có Mã NV (Cột 0) để làm Key
		maNV := layString(r, mo_hinh.CotNV_MaNhanVien)
		if maNV == "" { continue }

		// Map dữ liệu vào Struct
		nv := &mo_hinh.NhanVien{
			DongTrongSheet:  i + 1, // Lưu lại vị trí dòng (Excel index start 1)
			
			MaNhanVien:      maNV,
			TenDangNhap:     layString(r, mo_hinh.CotNV_TenDangNhap),
			Email:           layString(r, mo_hinh.CotNV_Email),
			MatKhauHash:     layString(r, mo_hinh.CotNV_MatKhauHash),
			HoTen:           layString(r, mo_hinh.CotNV_HoTen),
			ChucVu:          layString(r, mo_hinh.CotNV_ChucVu),
			MaPin:           layString(r, mo_hinh.CotNV_MaPin),
			Cookie:          layString(r, mo_hinh.CotNV_Cookie),
			
			// Convert số từ Sheet sang int64
			CookieExpired:   int64(layFloat(r, mo_hinh.CotNV_CookieExpired)),
			
			VaiTroQuyenHan:  layString(r, mo_hinh.CotNV_VaiTroQuyenHan),
			TrangThai:       layInt(r, mo_hinh.CotNV_TrangThai),
			LanDangNhapCuoi: layString(r, mo_hinh.CotNV_LanDangNhapCuoi),
		}

		CacheNhanVien.DuLieu[maNV] = nv
	}
	log.Printf("   -> Đã nạp %d nhân viên.", len(CacheNhanVien.DuLieu))
}
