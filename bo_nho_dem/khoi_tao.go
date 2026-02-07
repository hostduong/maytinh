package bo_nho_dem

import (
	"log"
	"sync"
	"app/cau_hinh"
	"app/mo_hinh"
)

// BIẾN TOÀN CỤC (Giữ nguyên)
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

// Helper tạo struct rỗng (Giữ nguyên)
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
		&KhoKhachHangStore{DuLieu: make(map[string]*mo_hinh.KhachHang), DanhSach: []*mo_hinh.KhachHang{}, TenKey: TaoKeyCache("KHACH_HANG"), SpreadsheetID: cau_hinh.BienCauHinh.IdFileSheet},
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

func KhoiTaoCacStore() {
	CacheSanPham, CacheDanhMuc, CacheThuongHieu, CacheNhaCungCap,
	CacheKhachHang, CachePhieuNhap, CacheChiTietNhap, CachePhieuXuat,
	CacheChiTietXuat, CacheSerial, CacheKhuyenMai, CacheCauHinhWeb,
	CacheHoaDon, CacheHoaDonChiTiet, CachePhieuThuChi, CachePhieuBaoHanh = taoMoiCacStore()
}

// =============================================================================
// [MỚI] CHIẾN LƯỢC TẢI DỮ LIỆU THÔNG MINH
// =============================================================================

// 1. Chỉ tải những gì cần thiết để Đăng nhập & Check Quyền (Chạy tuần tự)
func NapDuLieuCotLoi() {
	var wg sync.WaitGroup
	wg.Add(3) // 3 bảng quan trọng nhất

	go func() { defer wg.Done(); napKhachHang(CacheKhachHang) }()
	go func() { defer wg.Done(); napCauHinhWeb(CacheCauHinhWeb) }()
	go func() { defer wg.Done(); napPhanQuyen() }() // Hàm trong phan_quyen.go

	wg.Wait()
	log.Println("✅ [BOOT 1/2] Đã nạp xong Dữ Liệu Cốt Lõi (User, Config).")
}

// 2. Tải phần còn lại (Chạy ngầm)
func NapDuLieuNen() {
	log.Println("⏳ [BOOT 2/2] Bắt đầu nạp dữ liệu nền (Sản phẩm, Đơn hàng)...")
	
	// Gom các biến toàn cục vào tham số để gọi hàm helper
	thucHienNapPhanConLai(
		CacheSanPham, CacheDanhMuc, CacheThuongHieu, CacheNhaCungCap,
		CachePhieuNhap, CacheChiTietNhap, CachePhieuXuat, CacheChiTietXuat,
		CacheSerial, CacheKhuyenMai, CacheHoaDon, CacheHoaDonChiTiet, 
		CachePhieuThuChi, CachePhieuBaoHanh,
	)
	
	log.Println("✅ [BOOT 2/2] Hoàn tất nạp 100% dữ liệu.")
}

// =============================================================================
// LOGIC RELOAD (VẪN PHẢI TẢI FULL RỒI MỚI SWAP ĐỂ AN TOÀN)
// =============================================================================

func LamMoiHeThong() {
	log.Println("⚡ [RELOAD] Bắt đầu quy trình Tách Ly Đọc Ghi...")
	HeThongDangBan = true
	if CallbackGhiSheet != nil { CallbackGhiSheet(true) }

	// Tạo bộ nhớ tạm mới tinh
	tmpSP, tmpDM, tmpTH, tmpNCC, tmpKH, tmpPN, tmpCTPN, tmpPX, 
	tmpCTPX, tmpSer, tmpKM, tmpWeb, tmpHD, tmpHDCT, tmpThuChi, tmpBH := taoMoiCacStore()

	// Tải Song Song cả 2 nhóm (Cốt lõi + Nền) vào biến tạm
	var wg sync.WaitGroup
	wg.Add(2)

	// Nhóm 1: Cốt lõi
	go func() {
		defer wg.Done()
		napKhachHang(tmpKH)
		napCauHinhWeb(tmpWeb)
		napPhanQuyen() // Lưu ý: Biến PhanQuyen là map riêng, reload thẳng
	}()

	// Nhóm 2: Nền
	go func() {
		defer wg.Done()
		thucHienNapPhanConLai(tmpSP, tmpDM, tmpTH, tmpNCC, tmpPN, tmpCTPN, tmpPX, tmpCTPX, tmpSer, tmpKM, tmpHD, tmpHDCT, tmpThuChi, tmpBH)
	}()

	wg.Wait()

	// Swap
	KhoaHeThong.Lock()
	CacheSanPham = tmpSP; CacheDanhMuc = tmpDM; CacheThuongHieu = tmpTH; CacheNhaCungCap = tmpNCC
	CacheKhachHang = tmpKH
	CachePhieuNhap = tmpPN; CacheChiTietNhap = tmpCTPN
	CachePhieuXuat = tmpPX; CacheChiTietXuat = tmpCTPX
	CacheSerial = tmpSer
	CacheKhuyenMai = tmpKM; CacheCauHinhWeb = tmpWeb
	CacheHoaDon = tmpHD; CacheHoaDonChiTiet = tmpHDCT
	CachePhieuThuChi = tmpThuChi; CachePhieuBaoHanh = tmpBH
	
	HeThongDangBan = false
	KhoaHeThong.Unlock()
	log.Println("✅ [RELOAD] Hoán đổi hoàn tất.")
}

// Helper nạp các bảng còn lại (Trừ KH, Config, PhanQuyen)
func thucHienNapPhanConLai(
	pSP *KhoSanPhamStore, pDM *KhoDanhMucStore, pTH *KhoThuongHieuStore, pNCC *KhoNhaCungCapStore,
	pPN *KhoPhieuNhapStore, pCTPN *KhoChiTietPhieuNhapStore, pPX *KhoPhieuXuatStore,
	pCTPX *KhoChiTietPhieuXuatStore, pSer *KhoSerialStore, pKM *KhoKhuyenMaiStore,
	pHD *KhoHoaDonStore, pHDCT *KhoHoaDonChiTietStore, pThuChi *KhoPhieuThuChiStore, pBH *KhoPhieuBaoHanhStore,
) {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() { 
		defer wg.Done()
		napDanhMuc(pDM); napThuongHieu(pTH); napSanPham(pSP); napNhaCungCap(pNCC)
	}()
	go func() { 
		defer wg.Done()
		napPhieuNhap(pPN); napChiTietPhieuNhap(pCTPN); napPhieuXuat(pPX); napChiTietPhieuXuat(pCTPX)
		napSerial(pSer); napKhuyenMai(pKM) 
	}()
	go func() { 
		defer wg.Done()
		napHoaDon(pHD); napHoaDonChiTiet(pHDCT); napPhieuThuChi(pThuChi); napPhieuBaoHanh(pBH) 
	}()
	wg.Wait()
}
