package nghiep_vu

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"app/cau_hinh"
	"app/kho_du_lieu"
	"app/mo_hinh"
)

// [CƠ CHẾ KHÓA HỆ THỐNG]
var KhoaHeThong sync.RWMutex

// [CỜ TRẠNG THÁI] True = Đang reload, Middleware sẽ chặn request Ghi (POST/PUT/DELETE)
var HeThongDangBan bool = false

// =================================================================================
// 1. ĐỊNH NGHĨA KHO DỮ LIỆU (STORE STRUCTS)
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
	DuLieu        map[string]*mo_hinh.KhachHang
	TenKey        string
	SpreadsheetID string
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
// 2. BIẾN TOÀN CỤC (GLOBAL POINTERS)
// =================================================================================
var (
	CacheSanPham         *KhoSanPhamStore
	CacheDanhMuc         *KhoDanhMucStore
	CacheThuongHieu      *KhoThuongHieuStore
	CacheNhaCungCap      *KhoNhaCungCapStore
	CacheKhachHang       *KhoKhachHangStore
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

func TaoKeyCache(tenSheet string) string {
	idSheet := cau_hinh.BienCauHinh.IdFileSheet
	return fmt.Sprintf("%s__%s", strings.TrimSpace(idSheet), tenSheet)
}

// Hàm khởi tạo struct rỗng (Helper để tạo biến tạm hoặc khởi tạo ban đầu)
func taoMoiCacStore() (
	*KhoSanPhamStore, *KhoDanhMucStore, *KhoThuongHieuStore, *KhoNhaCungCapStore,
	*KhoKhachHangStore, *KhoPhieuNhapStore, *KhoChiTietPhieuNhapStore, *KhoPhieuXuatStore,
	*KhoChiTietPhieuXuatStore, *KhoSerialStore, *KhoKhuyenMaiStore, *KhoCauHinhWebStore,
	*KhoHoaDonStore, *KhoHoaDonChiTietStore, *KhoPhieuThuChiStore, *KhoPhieuBaoHanhStore,
) {
	return &KhoSanPhamStore{DuLieu: make(map[string]mo_hinh.SanPham), TenKey: TaoKeyCache("SAN_PHAM")},
		&KhoDanhMucStore{DuLieu: make(map[string]mo_hinh.DanhMuc), TenKey: TaoKeyCache("DANH_MUC")},
		&KhoThuongHieuStore{DuLieu: make(map[string]mo_hinh.ThuongHieu), TenKey: TaoKeyCache("THUONG_HIEU")},
		&KhoNhaCungCapStore{DuLieu: make(map[string]mo_hinh.NhaCungCap), TenKey: TaoKeyCache("NHA_CUNG_CAP")},
		&KhoKhachHangStore{DuLieu: make(map[string]*mo_hinh.KhachHang), TenKey: TaoKeyCache("KHACH_HANG"), SpreadsheetID: cau_hinh.BienCauHinh.IdFileSheet},
		&KhoPhieuNhapStore{DuLieu: make(map[string]mo_hinh.PhieuNhap), TenKey: TaoKeyCache("PHIEU_NHAP")},
		&KhoChiTietPhieuNhapStore{TenKey: TaoKeyCache("CHI_TIET_PHIEU_NHAP")},
		&KhoPhieuXuatStore{DuLieu: make(map[string]mo_hinh.PhieuXuat), TenKey: TaoKeyCache("PHIEU_XUAT")},
		&KhoChiTietPhieuXuatStore{TenKey: TaoKeyCache("CHI_TIET_PHIEU_XUAT")},
		&KhoSerialStore{DuLieu: make(map[string]mo_hinh.SerialSanPham), TenKey: TaoKeyCache("SERIAL_SAN_PHAM")},
		&KhoKhuyenMaiStore{DuLieu: make(map[string]mo_hinh.KhuyenMai), TenKey: TaoKeyCache("KHUYEN_MAI")},
		&KhoCauHinhWebStore{DuLieu: make(map[string]mo_hinh.CauHinhWeb), TenKey: TaoKeyCache("CAU_HINH_WEB")},
		&KhoHoaDonStore{DuLieu: make(map[string]mo_hinh.HoaDon), TenKey: TaoKeyCache("HOA_DON")},
		&KhoHoaDonChiTietStore{TenKey: TaoKeyCache("HOA_DON_CHI_TIET")},
		&KhoPhieuThuChiStore{DuLieu: make(map[string]mo_hinh.PhieuThuChi), TenKey: TaoKeyCache("PHIEU_THU_CHI")},
		&KhoPhieuBaoHanhStore{DuLieu: make(map[string]mo_hinh.PhieuBaoHanh), TenKey: TaoKeyCache("PHIEU_BAO_HANH")}
}

// =================================================================================
// 3. CÁC HÀM KHỞI TẠO VÀ LÀM MỚI (RELOAD LOGIC)
// =================================================================================

// 1. Khởi tạo lần đầu (Boot) - Nạp thẳng vào biến Global
func KhoiTaoCacStore() {
	CacheSanPham, CacheDanhMuc, CacheThuongHieu, CacheNhaCungCap,
	CacheKhachHang, CachePhieuNhap, CacheChiTietNhap, CachePhieuXuat,
	CacheChiTietXuat, CacheSerial, CacheKhuyenMai, CacheCauHinhWeb,
	CacheHoaDon, CacheHoaDonChiTiet, CachePhieuThuChi, CachePhieuBaoHanh = taoMoiCacStore()
	
	log.Println("✅ [MEMORY] Đã khởi tạo bộ nhớ rỗng.")
}

// 2. Logic Nạp Dữ Liệu (Booting)
func KhoiTaoBoNho() {
	log.Println("--- [BOOT] Bắt đầu nạp dữ liệu ---")
	thucHienNapDaLuong(
		CacheSanPham, CacheDanhMuc, CacheThuongHieu, CacheNhaCungCap,
		CacheKhachHang, CachePhieuNhap, CacheChiTietNhap, CachePhieuXuat,
		CacheChiTietXuat, CacheSerial, CacheKhuyenMai, CacheCauHinhWeb,
		CacheHoaDon, CacheHoaDonChiTiet, CachePhieuThuChi, CachePhieuBaoHanh,
	)
	log.Println("--- [BOOT] Hoàn tất ---")
}

// 3. Logic Làm Mới (Reload - Shadow Load)
// Quy trình: Chặn Ghi -> Tải vào biến tạm -> Hoán đổi -> Mở Ghi
func LamMoiHeThong() {
	log.Println("⚡ [RELOAD] Bắt đầu quy trình Tách Ly Đọc Ghi...")

	// B1: Bật cờ Bận -> Middleware sẽ chặn POST
	HeThongDangBan = true
	
	// B2: Vét sạch dữ liệu tồn đọng trong Hàng Chờ (Flush)
	ThucHienGhiSheet(true) 

	// B3: Tạo biến tạm (Shadow Memory) - RAM Sạch hoàn toàn
	tmpSP, tmpDM, tmpTH, tmpNCC, tmpKH, tmpPN, tmpCTPN, tmpPX, 
	tmpCTPX, tmpSer, tmpKM, tmpWeb, tmpHD, tmpHDCT, tmpThuChi, tmpBH := taoMoiCacStore()

	// B4: Tải dữ liệu vào biến tạm (Mất 3-5s)
	// Trong lúc này User vẫn đọc biến cũ (Cache...) bình thường, không bị lag
	thucHienNapDaLuong(
		tmpSP, tmpDM, tmpTH, tmpNCC, tmpKH, tmpPN, tmpCTPN, tmpPX, 
		tmpCTPX, tmpSer, tmpKM, tmpWeb, tmpHD, tmpHDCT, tmpThuChi, tmpBH,
	)

	// B5: Hoán đổi (Swap) - Chỉ mất vài mili giây
	KhoaHeThong.Lock()
	CacheSanPham = tmpSP
	CacheDanhMuc = tmpDM
	CacheThuongHieu = tmpTH
	CacheNhaCungCap = tmpNCC
	CacheKhachHang = tmpKH
	CachePhieuNhap = tmpPN
	CacheChiTietNhap = tmpCTPN
	CachePhieuXuat = tmpPX
	CacheChiTietXuat = tmpCTPX
	CacheSerial = tmpSer
	CacheKhuyenMai = tmpKM
	CacheCauHinhWeb = tmpWeb
	CacheHoaDon = tmpHD
	CacheHoaDonChiTiet = tmpHDCT
	CachePhieuThuChi = tmpThuChi
	CachePhieuBaoHanh = tmpBH
	
	// Reset lại cờ, mở khóa
	HeThongDangBan = false
	KhoaHeThong.Unlock()

	log.Println("✅ [RELOAD] Hoán đổi hoàn tất. Hệ thống mở lại.")
}

// Hàm chạy Goroutine nạp (Nhận tham số Pointer thay vì dùng Global)
func thucHienNapDaLuong(
	pSP *KhoSanPhamStore, pDM *KhoDanhMucStore, pTH *KhoThuongHieuStore, pNCC *KhoNhaCungCapStore,
	pKH *KhoKhachHangStore, pPN *KhoPhieuNhapStore, pCTPN *KhoChiTietPhieuNhapStore, pPX *KhoPhieuXuatStore,
	pCTPX *KhoChiTietPhieuXuatStore, pSer *KhoSerialStore, pKM *KhoKhuyenMaiStore, pWeb *KhoCauHinhWebStore,
	pHD *KhoHoaDonStore, pHDCT *KhoHoaDonChiTietStore, pThuChi *KhoPhieuThuChiStore, pBH *KhoPhieuBaoHanhStore,
) {
	var wg sync.WaitGroup
	wg.Add(3)

	// Nhóm 1: Master Data
	go func() { 
		defer wg.Done()
		napDanhMuc(pDM)
		napThuongHieu(pTH)
		napSanPham(pSP)
		napKhachHang(pKH)
		napNhaCungCap(pNCC)
		napCauHinhWeb(pWeb)
		NapDuLieuPhanQuyen() // Hàm này tự quản lý biến riêng của nó
	}()
	
	// Nhóm 2: Giao dịch hàng hóa
	go func() { 
		defer wg.Done()
		// Tạo delay giả lập nhỏ nếu cần để tránh Google chặn rate limit quá gắt (Opsional)
		time.Sleep(100 * time.Millisecond) 
		napPhieuNhap(pPN)
		napChiTietPhieuNhap(pCTPN)
		napPhieuXuat(pPX)
		napChiTietPhieuXuat(pCTPX)
		napSerial(pSer)
		napKhuyenMai(pKM) 
	}()
	
	// Nhóm 3: Tài chính
	go func() { 
		defer wg.Done()
		time.Sleep(200 * time.Millisecond)
		napHoaDon(pHD)
		napHoaDonChiTiet(pHDCT)
		napPhieuThuChi(pThuChi)
		napPhieuBaoHanh(pBH) 
	}()

	wg.Wait()
}

// Helper load sheet
func loadSheetData(sheetName string, keyCache string) ([][]interface{}, error) {
	duLieu, err := kho_du_lieu.DocToanBoSheet(sheetName)
	if err != nil {
		log.Printf("LỖI ĐỌC %s: %v", sheetName, err)
		return nil, err
	}
	return duLieu, nil
}

// =================================================================================
// 4. LOGIC NẠP CHI TIẾT CHO TỪNG BẢNG (16 BẢNG)
// =================================================================================

// 1. KHACH_HANG
func napKhachHang(target *KhoKhachHangStore) {
	raw, err := loadSheetData("KHACH_HANG", target.TenKey)
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotKH_MaKhachHang || layString(r, mo_hinh.CotKH_MaKhachHang) == "" { continue }
		item := &mo_hinh.KhachHang{
			DongTrongSheet: i + 1,
			MaKhachHang:    layString(r, mo_hinh.CotKH_MaKhachHang),
			TenDangNhap:    layString(r, mo_hinh.CotKH_TenDangNhap),
			MatKhauHash:    layString(r, mo_hinh.CotKH_MatKhauHash),
			Cookie:         layString(r, mo_hinh.CotKH_Cookie),
			CookieExpired:  int64(layFloat(r, mo_hinh.CotKH_CookieExpired)),
			MaPinHash:      layString(r, mo_hinh.CotKH_MaPinHash),
			TenKhachHang:   layString(r, mo_hinh.CotKH_TenKhachHang),
			Email:          layString(r, mo_hinh.CotKH_Email),
			DienThoai:      layString(r, mo_hinh.CotKH_DienThoai),
			UrlFb:          layString(r, mo_hinh.CotKH_UrlFb),
			Zalo:           layString(r, mo_hinh.CotKH_Zalo),
			UrlTiktok:      layString(r, mo_hinh.CotKH_UrlTiktok),
			DiaChi:         layString(r, mo_hinh.CotKH_DiaChi),
			MaSoThue:       layString(r, mo_hinh.CotKH_MaSoThue),
			DangNo:         layFloat(r, mo_hinh.CotKH_DangNo),
			TongMua:        layFloat(r, mo_hinh.CotKH_TongMua),
			ChucVu:         layString(r, mo_hinh.CotKH_ChucVu),
			VaiTroQuyenHan: layString(r, mo_hinh.CotKH_VaiTroQuyenHan),
			TrangThai:      layInt(r, mo_hinh.CotKH_TrangThai),
			GhiChu:         layString(r, mo_hinh.CotKH_GhiChu),
			NguoiTao:       layString(r, mo_hinh.CotKH_NguoiTao),
			NgayTao:        layString(r, mo_hinh.CotKH_NgayTao),
			NgayCapNhat:    layString(r, mo_hinh.CotKH_NgayCapNhat),
		}
		target.DuLieu[item.MaKhachHang] = item
		if item.TenDangNhap != "" { target.DuLieu[strings.ToLower(item.TenDangNhap)] = item }
		if item.Email != "" { target.DuLieu[strings.ToLower(item.Email)] = item }
	}
}

// 2. SAN_PHAM
func napSanPham(target *KhoSanPhamStore) {
	raw, err := loadSheetData("SAN_PHAM", target.TenKey)
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
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
		target.DuLieu[item.MaSanPham] = item
		target.DanhSach = append(target.DanhSach, item)
	}
}

// 3. DANH_MUC
func napDanhMuc(target *KhoDanhMucStore) {
	raw, err := loadSheetData("DANH_MUC", target.TenKey)
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotDM_MaDanhMuc || layString(r, mo_hinh.CotDM_MaDanhMuc) == "" { continue }
		item := mo_hinh.DanhMuc{
			MaDanhMuc:    layString(r, mo_hinh.CotDM_MaDanhMuc),
			ThuTuHienThi: layInt(r, mo_hinh.CotDM_ThuTuHienThi),
			TenDanhMuc:   layString(r, mo_hinh.CotDM_TenDanhMuc),
			Slug:         layString(r, mo_hinh.CotDM_Slug),
			MaDanhMucCha: layString(r, mo_hinh.CotDM_MaDanhMucCha),
		}
		target.DuLieu[item.MaDanhMuc] = item
	}
}

// 4. THUONG_HIEU
func napThuongHieu(target *KhoThuongHieuStore) {
	raw, err := loadSheetData("THUONG_HIEU", target.TenKey)
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotTH_MaThuongHieu || layString(r, mo_hinh.CotTH_MaThuongHieu) == "" { continue }
		item := mo_hinh.ThuongHieu{
			MaThuongHieu:  layString(r, mo_hinh.CotTH_MaThuongHieu),
			TenThuongHieu: layString(r, mo_hinh.CotTH_TenThuongHieu),
			LogoUrl:       layString(r, mo_hinh.CotTH_LogoUrl),
		}
		target.DuLieu[item.MaThuongHieu] = item
	}
}

// 5. NHA_CUNG_CAP
func napNhaCungCap(target *KhoNhaCungCapStore) {
	raw, err := loadSheetData("NHA_CUNG_CAP", target.TenKey)
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
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
		target.DuLieu[item.MaNhaCungCap] = item
	}
}

// 6. PHIEU_XUAT
func napPhieuXuat(target *KhoPhieuXuatStore) {
	raw, err := loadSheetData("PHIEU_XUAT", target.TenKey)
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
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
		target.DuLieu[item.MaPhieuXuat] = item
		target.DanhSach = append(target.DanhSach, item)
	}
}

// 7. CHI_TIET_PHIEU_XUAT
func napChiTietPhieuXuat(target *KhoChiTietPhieuXuatStore) {
	raw, err := loadSheetData("CHI_TIET_PHIEU_XUAT", target.TenKey)
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
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
		target.DanhSach = append(target.DanhSach, item)
	}
}

// 8. PHIEU_NHAP
func napPhieuNhap(target *KhoPhieuNhapStore) {
	raw, err := loadSheetData("PHIEU_NHAP", target.TenKey)
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
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
		target.DuLieu[item.MaPhieuNhap] = item
		target.DanhSach = append(target.DanhSach, item)
	}
}

// 9. CHI_TIET_PHIEU_NHAP
func napChiTietPhieuNhap(target *KhoChiTietPhieuNhapStore) {
	raw, err := loadSheetData("CHI_TIET_PHIEU_NHAP", target.TenKey)
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotCTPN_MaPhieuNhap || layString(r, mo_hinh.CotCTPN_MaPhieuNhap) == "" { continue }
		item := mo_hinh.ChiTietPhieuNhap{
			MaPhieuNhap:   layString(r, mo_hinh.CotCTPN_MaPhieuNhap),
			MaSanPham:     layString(r, mo_hinh.CotCTPN_MaSanPham),
			SoLuong:       layInt(r, mo_hinh.CotCTPN_SoLuong),
			DonGiaNhap:    layFloat(r, mo_hinh.CotCTPN_DonGiaNhap),
			ThanhTienDong: layFloat(r, mo_hinh.CotCTPN_ThanhTienDong),
		}
		target.DanhSach = append(target.DanhSach, item)
	}
}

// 10. SERIAL_SAN_PHAM
func napSerial(target *KhoSerialStore) {
	raw, err := loadSheetData("SERIAL_SAN_PHAM", target.TenKey)
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

// 11. KHUYEN_MAI
func napKhuyenMai(target *KhoKhuyenMaiStore) {
	raw, err := loadSheetData("KHUYEN_MAI", target.TenKey)
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotKM_MaVoucher || layString(r, mo_hinh.CotKM_MaVoucher) == "" { continue }
		item := mo_hinh.KhuyenMai{
			MaVoucher:      layString(r, mo_hinh.CotKM_MaVoucher),
			LoaiGiam:       layString(r, mo_hinh.CotKM_LoaiGiam),
			GiaTriGiam:     layFloat(r, mo_hinh.CotKM_GiaTriGiam),
			DonToThieu:     layFloat(r, mo_hinh.CotKM_DonToThieu),
			SoLuongConLai:  layInt(r, mo_hinh.CotKM_SoLuongConLai),
			TrangThai:      layInt(r, mo_hinh.CotKM_TrangThai),
		}
		target.DuLieu[item.MaVoucher] = item
	}
}

// 12. CAU_HINH_WEB
func napCauHinhWeb(target *KhoCauHinhWebStore) {
	raw, err := loadSheetData("CAU_HINH_WEB", target.TenKey)
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotCH_MaCauHinh || layString(r, mo_hinh.CotCH_MaCauHinh) == "" { continue }
		item := mo_hinh.CauHinhWeb{
			MaCauHinh: layString(r, mo_hinh.CotCH_MaCauHinh),
			GiaTri:    layString(r, mo_hinh.CotCH_GiaTri),
			TrangThai: layInt(r, mo_hinh.CotCH_TrangThai),
		}
		target.DuLieu[item.MaCauHinh] = item
	}
}

// 13. HOA_DON
func napHoaDon(target *KhoHoaDonStore) {
	raw, err := loadSheetData("HOA_DON", target.TenKey)
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
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
		target.DuLieu[item.MaHoaDon] = item
	}
}

// 14. HOA_DON_CHI_TIET
func napHoaDonChiTiet(target *KhoHoaDonChiTietStore) {
	raw, err := loadSheetData("HOA_DON_CHI_TIET", target.TenKey)
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotHDCT_MaHoaDon || layString(r, mo_hinh.CotHDCT_MaHoaDon) == "" { continue }
		item := mo_hinh.HoaDonChiTiet{
			MaHoaDon:   layString(r, mo_hinh.CotHDCT_MaHoaDon),
			MaSanPham:  layString(r, mo_hinh.CotHDCT_MaSanPham),
			SoLuong:    layInt(r, mo_hinh.CotHDCT_SoLuong),
			ThanhTien:  layFloat(r, mo_hinh.CotHDCT_ThanhTien),
		}
		target.DanhSach = append(target.DanhSach, item)
	}
}

// 15. PHIEU_THU_CHI
func napPhieuThuChi(target *KhoPhieuThuChiStore) {
	raw, err := loadSheetData("PHIEU_THU_CHI", target.TenKey)
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotPTC_MaPhieuThuChi || layString(r, mo_hinh.CotPTC_MaPhieuThuChi) == "" { continue }
		item := mo_hinh.PhieuThuChi{
			MaPhieuThuChi: layString(r, mo_hinh.CotPTC_MaPhieuThuChi),
			LoaiPhieu:     layString(r, mo_hinh.CotPTC_LoaiPhieu),
			SoTien:        layFloat(r, mo_hinh.CotPTC_SoTien),
			HangMucThuChi: layString(r, mo_hinh.CotPTC_HangMucThuChi),
			TrangThaiDuyet: layInt(r, mo_hinh.CotPTC_TrangThaiDuyet),
		}
		target.DuLieu[item.MaPhieuThuChi] = item
		target.DanhSach = append(target.DanhSach, item)
	}
}

// 16. PHIEU_BAO_HANH
func napPhieuBaoHanh(target *KhoPhieuBaoHanhStore) {
	raw, err := loadSheetData("PHIEU_BAO_HANH", target.TenKey)
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotPBH_MaPhieuBaoHanh || layString(r, mo_hinh.CotPBH_MaPhieuBaoHanh) == "" { continue }
		item := mo_hinh.PhieuBaoHanh{
			MaPhieuBaoHanh: layString(r, mo_hinh.CotPBH_MaPhieuBaoHanh),
			SerialImei:     layString(r, mo_hinh.CotPBH_SerialImei),
			TrangThai:      layInt(r, mo_hinh.CotPBH_TrangThai),
			TinhTrangLoi:   layString(r, mo_hinh.CotPBH_TinhTrangLoi),
		}
		target.DuLieu[item.MaPhieuBaoHanh] = item
		target.DanhSach = append(target.DanhSach, item)
	}
}

// =================================================================================
// 5. HELPER FUNCTIONS
// =================================================================================

func layString(dong []interface{}, index int) string {
	if index >= len(dong) || dong[index] == nil { return "" }
	return fmt.Sprintf("%v", dong[index])
}

func layInt(dong []interface{}, index int) int {
	str := layString(dong, index)
	if str == "" { return 0 }
	
	str = strings.ReplaceAll(str, ".", "")
	str = strings.ReplaceAll(str, ",", "")
	str = strings.ReplaceAll(str, " ", "")
	
	val, _ := strconv.Atoi(str)
	return val
}

func layFloat(dong []interface{}, index int) float64 {
	str := layString(dong, index)
	if str == "" { return 0 }

	str = strings.ReplaceAll(str, "đ", "")
	str = strings.ReplaceAll(str, "USD", "")
	str = strings.ReplaceAll(str, " ", "")
	
	str = strings.ReplaceAll(str, ".", "")
	str = strings.ReplaceAll(str, ",", "")

	val, _ := strconv.ParseFloat(str, 64)
	return val
}
