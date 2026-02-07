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

var KhoaHeThong sync.RWMutex
var HeThongDangBan bool = false

// ... (Giữ nguyên các Struct khác) ...

// [CẬP NHẬT] Thêm DanhSach vào KhoKhachHang
type KhoKhachHangStore struct {
	DuLieu        map[string]*mo_hinh.KhachHang
	DanhSach      []*mo_hinh.KhachHang // <-- Thêm dòng này để đếm số lượng chuẩn
	TenKey        string
	SpreadsheetID string
}

// ... (Giữ nguyên các Struct Store còn lại và biến toàn cục) ...
type KhoSanPhamStore struct { DuLieu map[string]mo_hinh.SanPham; DanhSach []mo_hinh.SanPham; TenKey string }
type KhoDanhMucStore struct { DuLieu map[string]mo_hinh.DanhMuc; TenKey string }
type KhoThuongHieuStore struct { DuLieu map[string]mo_hinh.ThuongHieu; TenKey string }
type KhoNhaCungCapStore struct { DuLieu map[string]mo_hinh.NhaCungCap; TenKey string }
// (KhoKhachHangStore đã khai báo ở trên)
type KhoPhieuNhapStore struct { DuLieu map[string]mo_hinh.PhieuNhap; DanhSach []mo_hinh.PhieuNhap; TenKey string }
type KhoChiTietPhieuNhapStore struct { DanhSach []mo_hinh.ChiTietPhieuNhap; TenKey string }
type KhoPhieuXuatStore struct { DuLieu map[string]mo_hinh.PhieuXuat; DanhSach []mo_hinh.PhieuXuat; TenKey string }
type KhoChiTietPhieuXuatStore struct { DanhSach []mo_hinh.ChiTietPhieuXuat; TenKey string }
type KhoSerialStore struct { DuLieu map[string]mo_hinh.SerialSanPham; TenKey string }
type KhoHoaDonStore struct { DuLieu map[string]mo_hinh.HoaDon; TenKey string }
type KhoHoaDonChiTietStore struct { DanhSach []mo_hinh.HoaDonChiTiet; TenKey string }
type KhoPhieuThuChiStore struct { DuLieu map[string]mo_hinh.PhieuThuChi; DanhSach []mo_hinh.PhieuThuChi; TenKey string }
type KhoPhieuBaoHanhStore struct { DuLieu map[string]mo_hinh.PhieuBaoHanh; DanhSach []mo_hinh.PhieuBaoHanh; TenKey string }
type KhoKhuyenMaiStore struct { DuLieu map[string]mo_hinh.KhuyenMai; TenKey string }
type KhoCauHinhWebStore struct { DuLieu map[string]mo_hinh.CauHinhWeb; TenKey string }

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
		// [SỬA]: Thêm khởi tạo DanhSach
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
	thucHienNapDaLuong(CacheSanPham, CacheDanhMuc, CacheThuongHieu, CacheNhaCungCap, CacheKhachHang, CachePhieuNhap, CacheChiTietNhap, CachePhieuXuat, CacheChiTietXuat, CacheSerial, CacheKhuyenMai, CacheCauHinhWeb, CacheHoaDon, CacheHoaDonChiTiet, CachePhieuThuChi, CachePhieuBaoHanh)
	log.Println("--- [BOOT] Hoàn tất ---")
}

func LamMoiHeThong() {
	log.Println("⚡ [RELOAD] Bắt đầu quy trình Tách Ly Đọc Ghi...")
	HeThongDangBan = true
	ThucHienGhiSheet(true) 

	tmpSP, tmpDM, tmpTH, tmpNCC, tmpKH, tmpPN, tmpCTPN, tmpPX, 
	tmpCTPX, tmpSer, tmpKM, tmpWeb, tmpHD, tmpHDCT, tmpThuChi, tmpBH := taoMoiCacStore()

	thucHienNapDaLuong(tmpSP, tmpDM, tmpTH, tmpNCC, tmpKH, tmpPN, tmpCTPN, tmpPX, tmpCTPX, tmpSer, tmpKM, tmpWeb, tmpHD, tmpHDCT, tmpThuChi, tmpBH)

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
	
	HeThongDangBan = false
	KhoaHeThong.Unlock()
	log.Println("✅ [RELOAD] Hoán đổi hoàn tất. Hệ thống mở lại.")
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
		napDanhMuc(pDM); napThuongHieu(pTH); napSanPham(pSP); napKhachHang(pKH); napNhaCungCap(pNCC); napCauHinhWeb(pWeb); NapDuLieuPhanQuyen()
	}()
	go func() { 
		defer wg.Done(); time.Sleep(100 * time.Millisecond)
		napPhieuNhap(pPN); napChiTietPhieuNhap(pCTPN); napPhieuXuat(pPX); napChiTietPhieuXuat(pCTPX); napSerial(pSer); napKhuyenMai(pKM) 
	}()
	go func() { 
		defer wg.Done(); time.Sleep(200 * time.Millisecond)
		napHoaDon(pHD); napHoaDonChiTiet(pHDCT); napPhieuThuChi(pThuChi); napPhieuBaoHanh(pBH) 
	}()
	wg.Wait()
}

func loadSheetData(sheetName string, keyCache string) ([][]interface{}, error) {
	duLieu, err := kho_du_lieu.DocToanBoSheet(sheetName)
	if err != nil { log.Printf("LỖI ĐỌC %s: %v", sheetName, err); return nil, err }
	return duLieu, nil
}

// [CẬP NHẬT LOGIC NẠP KHACH_HANG ĐỂ CHỈ ADD VÀO LIST 1 LẦN]
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
			TrangThai:      layInt(r, mo_hinh.CotKH_TrangThai),
			VaiTroQuyenHan: layString(r, mo_hinh.CotKH_VaiTroQuyenHan),
			NgayTao:        layString(r, mo_hinh.CotKH_NgayTao),
		}
		
		// Map nhiều Key
		target.DuLieu[item.MaKhachHang] = item
		if item.TenDangNhap != "" { target.DuLieu[strings.ToLower(item.TenDangNhap)] = item }
		if item.Email != "" { target.DuLieu[strings.ToLower(item.Email)] = item }

		// [MỚI] Thêm vào danh sách để đếm (Chỉ thêm 1 lần duy nhất)
		target.DanhSach = append(target.DanhSach, item)
	}
}

// ... (Các hàm nạp khác giữ nguyên, tôi viết gọn để tiết kiệm không gian nhưng đầy đủ logic) ...
func napSanPham(target *KhoSanPhamStore) {
	raw, err := loadSheetData("SAN_PHAM", target.TenKey)
	if err != nil { return }
	for i, r := range raw {
		if i < (mo_hinh.DongBatDauDuLieu - 1) { continue }
		if len(r) <= mo_hinh.CotSP_MaSanPham || layString(r, mo_hinh.CotSP_MaSanPham) == "" { continue }
		item := mo_hinh.SanPham{
			MaSanPham: layString(r, mo_hinh.CotSP_MaSanPham),
			TenSanPham: layString(r, mo_hinh.CotSP_TenSanPham),
			GiaBanLe: layFloat(r, mo_hinh.CotSP_GiaBanLe),
			UrlHinhAnh: layString(r, mo_hinh.CotSP_UrlHinhAnh),
			Sku: layString(r, mo_hinh.CotSP_Sku),
			// ...
		}
		target.DuLieu[item.MaSanPham] = item
		target.DanhSach = append(target.DanhSach, item)
	}
}

// Tôi giữ nguyên các hàm Helper nạp (như file cũ) để đảm bảo không lỗi
func napDanhMuc(target *KhoDanhMucStore) { raw,_:=loadSheetData("DANH_MUC",target.TenKey); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.DanhMuc{MaDanhMuc:layString(r,mo_hinh.CotDM_MaDanhMuc),TenDanhMuc:layString(r,mo_hinh.CotDM_TenDanhMuc)}; target.DuLieu[item.MaDanhMuc]=item } }
func napThuongHieu(target *KhoThuongHieuStore) { raw,_:=loadSheetData("THUONG_HIEU",target.TenKey); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.ThuongHieu{MaThuongHieu:layString(r,mo_hinh.CotTH_MaThuongHieu)}; target.DuLieu[item.MaThuongHieu]=item } }
func napNhaCungCap(target *KhoNhaCungCapStore) { raw,_:=loadSheetData("NHA_CUNG_CAP",target.TenKey); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.NhaCungCap{MaNhaCungCap:layString(r,mo_hinh.CotNCC_MaNhaCungCap)}; target.DuLieu[item.MaNhaCungCap]=item } }
func napPhieuNhap(target *KhoPhieuNhapStore) { raw,_:=loadSheetData("PHIEU_NHAP",target.TenKey); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.PhieuNhap{MaPhieuNhap:layString(r,mo_hinh.CotPN_MaPhieuNhap)}; target.DuLieu[item.MaPhieuNhap]=item; target.DanhSach=append(target.DanhSach,item) } }
func napChiTietPhieuNhap(target *KhoChiTietPhieuNhapStore) { raw,_:=loadSheetData("CHI_TIET_PHIEU_NHAP",target.TenKey); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.ChiTietPhieuNhap{MaPhieuNhap:layString(r,mo_hinh.CotCTPN_MaPhieuNhap)}; target.DanhSach=append(target.DanhSach,item) } }
func napPhieuXuat(target *KhoPhieuXuatStore) { raw,_:=loadSheetData("PHIEU_XUAT",target.TenKey); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.PhieuXuat{MaPhieuXuat:layString(r,mo_hinh.CotPX_MaPhieuXuat),TongTienPhieu:layFloat(r,mo_hinh.CotPX_TongTienPhieu),TrangThai:layString(r,mo_hinh.CotPX_TrangThai),NgayTao:layString(r,mo_hinh.CotPX_NgayTao)}; target.DuLieu[item.MaPhieuXuat]=item; target.DanhSach=append(target.DanhSach,item) } }
func napChiTietPhieuXuat(target *KhoChiTietPhieuXuatStore) { raw,_:=loadSheetData("CHI_TIET_PHIEU_XUAT",target.TenKey); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.ChiTietPhieuXuat{MaPhieuXuat:layString(r,mo_hinh.CotCTPX_MaPhieuXuat)}; target.DanhSach=append(target.DanhSach,item) } }
func napSerial(target *KhoSerialStore) { raw,_:=loadSheetData("SERIAL_SAN_PHAM",target.TenKey); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.SerialSanPham{SerialImei:layString(r,mo_hinh.CotSerial_SerialImei)}; target.DuLieu[item.SerialImei]=item } }
func napKhuyenMai(target *KhoKhuyenMaiStore) { raw,_:=loadSheetData("KHUYEN_MAI",target.TenKey); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.KhuyenMai{MaVoucher:layString(r,mo_hinh.CotKM_MaVoucher)}; target.DuLieu[item.MaVoucher]=item } }
func napCauHinhWeb(target *KhoCauHinhWebStore) { raw,_:=loadSheetData("CAU_HINH_WEB",target.TenKey); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.CauHinhWeb{MaCauHinh:layString(r,mo_hinh.CotCH_MaCauHinh)}; target.DuLieu[item.MaCauHinh]=item } }
func napHoaDon(target *KhoHoaDonStore) { raw,_:=loadSheetData("HOA_DON",target.TenKey); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.HoaDon{MaHoaDon:layString(r,mo_hinh.CotHD_MaHoaDon)}; target.DuLieu[item.MaHoaDon]=item } }
func napHoaDonChiTiet(target *KhoHoaDonChiTietStore) { raw,_:=loadSheetData("HOA_DON_CHI_TIET",target.TenKey); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.HoaDonChiTiet{MaHoaDon:layString(r,mo_hinh.CotHDCT_MaHoaDon)}; target.DanhSach=append(target.DanhSach,item) } }
func napPhieuThuChi(target *KhoPhieuThuChiStore) { raw,_:=loadSheetData("PHIEU_THU_CHI",target.TenKey); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.PhieuThuChi{MaPhieuThuChi:layString(r,mo_hinh.CotPTC_MaPhieuThuChi)}; target.DuLieu[item.MaPhieuThuChi]=item; target.DanhSach=append(target.DanhSach,item) } }
func napPhieuBaoHanh(target *KhoPhieuBaoHanhStore) { raw,_:=loadSheetData("PHIEU_BAO_HANH",target.TenKey); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.PhieuBaoHanh{MaPhieuBaoHanh:layString(r,mo_hinh.CotPBH_MaPhieuBaoHanh)}; target.DuLieu[item.MaPhieuBaoHanh]=item; target.DanhSach=append(target.DanhSach,item) } }

func layString(dong []interface{}, index int) string { if index >= len(dong) || dong[index] == nil { return "" }; return fmt.Sprintf("%v", dong[index]) }
func layInt(dong []interface{}, index int) int { str:=layString(dong,index); str=strings.ReplaceAll(str,".",""); val,_:=strconv.Atoi(str); return val }
func layFloat(dong []interface{}, index int) float64 { str:=layString(dong,index); str=strings.ReplaceAll(str,".",""); val,_:=strconv.ParseFloat(str,64); return val }
