package chuc_nang

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"app/bo_nho_dem"
	"app/cau_hinh"
	"app/mo_hinh"
	"app/nghiep_vu"

	"github.com/gin-gonic/gin"
)

// TrangQuanLySanPham : Hiển thị danh sách
func TrangQuanLySanPham(c *gin.Context) {
	userID := c.GetString("USER_ID")
	kh, _ := nghiep_vu.LayThongTinKhachHang(userID)

	// 1. Lấy danh sách sản phẩm (Slice)
	listSP := nghiep_vu.LayDanhSachSanPham()
	
	// 2. Lấy Danh Mục (Map -> Slice để sort)
	// Hàm LayDanhSachDanhMuc nằm trong nghiep_vu/truy_xuat.go trả về Map
	mapDM := nghiep_vu.LayDanhSachDanhMuc()
	var listDM []mo_hinh.DanhMuc
	for _, v := range mapDM { listDM = append(listDM, v) }
	
	// Sắp xếp Danh mục A-Z
	sort.Slice(listDM, func(i, j int) bool { return listDM[i].TenDanhMuc < listDM[j].TenDanhMuc })

	// 3. Lấy Thương Hiệu (Map -> Slice để sort)
	mapTH := nghiep_vu.LayDanhSachThuongHieu()
	var listTH []mo_hinh.ThuongHieu
	for _, v := range mapTH { listTH = append(listTH, v) }
	
	// Sắp xếp Thương hiệu A-Z
	sort.Slice(listTH, func(i, j int) bool { return listTH[i].TenThuongHieu < listTH[j].TenThuongHieu })

	c.HTML(http.StatusOK, "quan_tri_san_pham", gin.H{
		"TieuDe":         "Quản lý sản phẩm",
		"NhanVien":       kh,
		"DaDangNhap":     true,
		"TenNguoiDung":   kh.TenKhachHang,
		"QuyenHan":       kh.VaiTroQuyenHan,
		"DanhSach":       listSP,
		"ListDanhMuc":    listDM, 
		"ListThuongHieu": listTH,
	})
}

// API_LuuSanPham : Xử lý Full trường (Dùng chi_muc.go)
func API_LuuSanPham(c *gin.Context) {
	// 1. Check quyền
	vaiTro := c.GetString("USER_ROLE")
	if !nghiep_vu.KiemTraQuyen(vaiTro, "product.edit") {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền này!"})
		return
	}

	// 2. Lấy dữ liệu form
	maSP        := strings.TrimSpace(c.PostForm("ma_san_pham"))
	tenSP       := strings.TrimSpace(c.PostForm("ten_san_pham"))
	tenRutGon   := strings.TrimSpace(c.PostForm("ten_rut_gon"))
	sku         := strings.TrimSpace(c.PostForm("sku"))
	
	// Xử lý giá tiền (1.500.000 -> 1500000)
	giaBanStr   := strings.ReplaceAll(c.PostForm("gia_ban_le"), ".", "")
	giaBanStr    = strings.ReplaceAll(giaBanStr, ",", "")
	giaBan, _   := strconv.ParseFloat(giaBanStr, 64)

	danhMucRaw  := c.PostForm("ma_danh_muc")
	danhMuc     := xuLyTags(danhMucRaw)

	thuongHieu  := c.PostForm("ma_thuong_hieu")
	donVi       := c.PostForm("don_vi")
	mauSac      := c.PostForm("mau_sac")
	hinhAnh     := strings.TrimSpace(c.PostForm("url_hinh_anh"))
	thongSo     := c.PostForm("thong_so")
	moTa        := c.PostForm("mo_ta_chi_tiet")
	baoHanh, _  := strconv.Atoi(c.PostForm("bao_hanh_thang"))
	tinhTrang   := c.PostForm("tinh_trang")
	ghiChu      := c.PostForm("ghi_chu")
	
	trangThai := 0
	if c.PostForm("trang_thai") == "on" { trangThai = 1 }

	if tenSP == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Tên sản phẩm không được để trống!"})
		return
	}

	// 3. Logic Thêm/Sửa
	var sp mo_hinh.SanPham
	isNew := false
	nowStr := time.Now().Format("2006-01-02 15:04:05")
	userID := c.GetString("USER_ID")

	bo_nho_dem.KhoaHeThong.Lock()
	
	if maSP == "" {
		isNew = true
		maSP = taoMaSPMoi()
		sp = mo_hinh.SanPham{
			MaSanPham: maSP,
			NgayTao:   nowStr,
			NguoiTao:  userID,
		}
	} else {
		if oldSP, ok := bo_nho_dem.CacheSanPham.DuLieu[maSP]; ok {
			sp = oldSP
		} else {
			sp = mo_hinh.SanPham{MaSanPham: maSP, NgayTao: nowStr}
		}
	}

	// Map dữ liệu
	sp.TenSanPham = tenSP
	sp.TenRutGon = tenRutGon
	sp.Sku = sku
	sp.GiaBanLe = giaBan
	sp.MaDanhMuc = danhMuc
	sp.MaThuongHieu = thuongHieu
	sp.DonVi = donVi
	sp.MauSac = mauSac
	sp.UrlHinhAnh = hinhAnh
	sp.ThongSo = thongSo
	sp.MoTaChiTiet = moTa
	sp.BaoHanhThang = baoHanh
	sp.TinhTrang = tinhTrang
	sp.TrangThai = trangThai
	sp.GhiChu = ghiChu
	sp.NgayCapNhat = nowStr

	// 4. Lưu RAM
	bo_nho_dem.CacheSanPham.DuLieu[maSP] = sp
	if isNew {
		bo_nho_dem.CacheSanPham.DanhSach = append(bo_nho_dem.CacheSanPham.DanhSach, sp)
	} else {
		for i, item := range bo_nho_dem.CacheSanPham.DanhSach {
			if item.MaSanPham == maSP {
				bo_nho_dem.CacheSanPham.DanhSach[i] = sp
				break
			}
		}
	}
	bo_nho_dem.KhoaHeThong.Unlock()

	// 5. Ghi Sheet (Sử dụng mo_hinh/chi_muc.go hiện có)
	sID := cau_hinh.BienCauHinh.IdFileSheet
	sheetName := "SAN_PHAM"
	targetRow := 0
	if isNew {
		targetRow = mo_hinh.DongBatDauDuLieu + len(bo_nho_dem.CacheSanPham.DuLieu) - 1
	}

	if targetRow > 0 {
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_MaSanPham, sp.MaSanPham)
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_TenSanPham, sp.TenSanPham)
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_TenRutGon, sp.TenRutGon)
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_Sku, sp.Sku)
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_MaDanhMuc, sp.MaDanhMuc)
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_MaThuongHieu, sp.MaThuongHieu)
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_DonVi, sp.DonVi)
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_MauSac, sp.MauSac)
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_UrlHinhAnh, sp.UrlHinhAnh)
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_ThongSo, sp.ThongSo)
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_MoTaChiTiet, sp.MoTaChiTiet)
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_BaoHanhThang, sp.BaoHanhThang)
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_TinhTrang, sp.TinhTrang)
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_TrangThai, sp.TrangThai)
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_GiaBanLe, sp.GiaBanLe)
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_GhiChu, sp.GhiChu)
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_NguoiTao, sp.NguoiTao)
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_NgayTao, sp.NgayTao)
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_NgayCapNhat, sp.NgayCapNhat)
	}

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã lưu sản phẩm thành công!"})
}

// Helper xử lý Tags
func xuLyTags(raw string) string {
	if !strings.Contains(raw, "[") { return raw }
	res := strings.ReplaceAll(raw, "[", "")
	res = strings.ReplaceAll(res, "]", "")
	res = strings.ReplaceAll(res, "{", "")
	res = strings.ReplaceAll(res, "}", "")
	res = strings.ReplaceAll(res, "\"value\":", "")
	res = strings.ReplaceAll(res, "\"", "")
	return res
}

func taoMaSPMoi() string {
	maxID := 0
	for _, sp := range bo_nho_dem.CacheSanPham.DanhSach {
		parts := strings.Split(sp.MaSanPham, "_")
		if len(parts) == 2 {
			id, _ := strconv.Atoi(parts[1])
			if id > maxID { maxID = id }
		}
	}
	return fmt.Sprintf("SP_%04d", maxID+1)
}
