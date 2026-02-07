package bo_nho_dem

import (
	"log"
	"sync"
	"time"
	"app/cau_hinh"
	"app/mo_hinh"
)

// BIẾN TOÀN CỤC (GLOBAL POINTERS)
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
	log.Println("✅ [MEMORY] Đã khởi tạo bộ nhớ rỗng.")
}

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

func LamMoiHeThong() {
	log.Println("⚡ [RELOAD] Bắt đầu quy trình Tách Ly Đọc Ghi...")
	HeThongDangBan = true
	if CallbackGhiSheet != nil { CallbackGhiSheet(true) }

	tmpSP, tmpDM, tmpTH, tmpNCC, tmpKH, tmpPN, tmpCTPN, tmpPX, 
	tmpCTPX, tmpSer, tmpKM, tmpWeb, tmpHD, tmpHDCT, tmpThuChi, tmpBH := taoMoiCacStore()

	thucHienNapDaLuong(tmpSP, tmpDM, tmpTH, tmpNCC, tmpKH, tmpPN, tmpCTPN, tmpPX, tmpCTPX, tmpSer, tmpKM, tmpWeb, tmpHD, tmpHDCT, tmpThuChi, tmpBH)

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

func thucHienNapDaLuong(
	pSP *KhoSanPhamStore, pDM *KhoDanhMucStore, pTH *KhoThuongHieuStore, pNCC *KhoNhaCungCapStore,
	pKH *KhoKhachHangStore, pPN *KhoPhieuNhapStore, pCTPN *KhoChiTietPhieuNhapStore, pPX *KhoPhieuXuatStore,
	pCTPX *KhoChiTietPhieuXuatStore, pSer *KhoSerialStore, pKM *KhoKhuyenMaiStore, pWeb *KhoCauHinhWebStore,
	pHD *KhoHoaDonStore, pHDCT *KhoHoaDonChiTietStore, pThuChi *KhoPhieuThuChiStore, pBH *KhoPhieuBaoHanhStore,
) {
	var wg sync.WaitGroup
	wg.Add(3)

	go func() { 
		defer wg.Done()
		napDanhMuc(pDM); napThuongHieu(pTH); napSanPham(pSP)
		napKhachHang(pKH); napNhaCungCap(pNCC); napCauHinhWeb(pWeb)
		napPhanQuyen() // Trong file phan_quyen.go (cùng package bo_nho_dem)
	}()
	go func() { 
		defer wg.Done(); time.Sleep(100 * time.Millisecond)
		napPhieuNhap(pPN); napChiTietPhieuNhap(pCTPN); napPhieuXuat(pPX); napChiTietPhieuXuat(pCTPX)
		napSerial(pSer); napKhuyenMai(pKM) 
	}()
	go func() { 
		defer wg.Done(); time.Sleep(200 * time.Millisecond)
		napHoaDon(pHD); napHoaDonChiTiet(pHDCT); napPhieuThuChi(pThuChi); napPhieuBaoHanh(pBH) 
	}()
	wg.Wait()
}
