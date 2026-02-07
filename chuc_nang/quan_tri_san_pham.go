package chuc_nang

import (
	"fmt"
	"net/http"
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
	danhSach := nghiep_vu.LayDanhSachSanPham()
	userID := c.GetString("USER_ID")
	kh, _ := nghiep_vu.LayThongTinKhachHang(userID)

	c.HTML(http.StatusOK, "quan_tri_san_pham", gin.H{
		"TieuDe":       "Quản lý sản phẩm",
		"NhanVien":     kh,
		"DaDangNhap":   true,
		"TenNguoiDung": kh.TenKhachHang,
		"QuyenHan":     kh.VaiTroQuyenHan,
		"DanhSach":     danhSach,
	})
}

// API_LuuSanPham : Xử lý Thêm mới hoặc Cập nhật
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
	giaBanStr   := strings.ReplaceAll(c.PostForm("gia_ban_le"), ",", "")
	giaBan, _   := strconv.ParseFloat(giaBanStr, 64)
	danhMuc     := c.PostForm("danh_muc")
	hinhAnh     := strings.TrimSpace(c.PostForm("url_hinh_anh"))
	moTa        := c.PostForm("mo_ta_chi_tiet")
	baoHanh, _  := strconv.Atoi(c.PostForm("bao_hanh_thang"))
	
	// Tạo tên rút gọn (Logic đơn giản: lấy 2 từ cuối hoặc giữ nguyên nếu ngắn)
	// Bạn có thể tùy chỉnh logic này sau. Tạm thời copy tên full.
	tenRutGon   := tenSP 

	if tenSP == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Tên sản phẩm không được để trống!"})
		return
	}

	// 3. Xử lý Logic (Thêm hay Sửa)
	var sp mo_hinh.SanPham
	isNew := false
	nowStr := time.Now().Format("2006-01-02 15:04:05")
	userID := c.GetString("USER_ID")

	bo_nho_dem.KhoaHeThong.Lock()
	
	if maSP == "" {
		// --- TẠO MỚI ---
		isNew = true
		maSP = taoMaSPMoi()
		sp = mo_hinh.SanPham{
			MaSanPham: maSP,
			NgayTao:   nowStr,
			NguoiTao:  userID,
			TrangThai: 1, // Mặc định đang bán
		}
	} else {
		// --- SỬA ---
		if oldSP, ok := bo_nho_dem.CacheSanPham.DuLieu[maSP]; ok {
			sp = oldSP
		} else {
			sp = mo_hinh.SanPham{MaSanPham: maSP, NgayTao: nowStr}
		}
	}

	// Cập nhật thông tin
	sp.TenSanPham = tenSP
	sp.TenRutGon = tenRutGon // Cập nhật tên rút gọn
	sp.GiaBanLe = giaBan
	sp.MaDanhMuc = danhMuc
	sp.UrlHinhAnh = hinhAnh
	sp.MoTaChiTiet = moTa
	sp.BaoHanhThang = baoHanh
	sp.NgayCapNhat = nowStr

	// 4. Lưu vào RAM
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
	
	// Mở khóa sớm để hệ thống khác chạy
	bo_nho_dem.KhoaHeThong.Unlock()

	// 5. Ghi xuống Sheet (Async)
	sID := cau_hinh.BienCauHinh.IdFileSheet
	sheetName := "SAN_PHAM"
	
	// [QUAN TRỌNG] Logic tìm dòng để ghi
	// Vì ta chưa lưu DongTrongSheet trong Struct SanPham,
	// nên tạm thời chỉ hỗ trợ GHI MỚI (Append).
	// Nếu SỬA -> Dữ liệu Sheet cũ sẽ không đổi ngay, 
	// nhưng RAM đã đổi -> Web hiển thị đúng.
	// Admin cần bấm "Reload" hoặc chờ cơ chế đồng bộ định kỳ (sẽ làm sau).
	
	targetRow := 0
	
	if isNew {
		// Tính dòng cuối dựa trên số lượng hiện có (tương đối)
		// Để chính xác tuyệt đối cần query sheet, nhưng để nhanh ta dùng RAM.
		// + DongBatDauDuLieu + len(RAM)
		targetRow = mo_hinh.DongBatDauDuLieu + len(bo_nho_dem.CacheSanPham.DanhSach) - 1
	}

	// Ghi đầy đủ các cột (A -> S)
	if targetRow > 0 {
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_MaSanPham, sp.MaSanPham)       // A
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_TenSanPham, sp.TenSanPham)     // B
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_TenRutGon, sp.TenRutGon)       // C
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_Sku, sp.Sku)                   // D
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_MaDanhMuc, sp.MaDanhMuc)       // E
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_MaThuongHieu, sp.MaThuongHieu) // F
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_DonVi, sp.DonVi)               // G
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_MauSac, sp.MauSac)             // H
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_UrlHinhAnh, sp.UrlHinhAnh)     // I
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_ThongSo, sp.ThongSo)           // J
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_MoTaChiTiet, sp.MoTaChiTiet)   // K
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_BaoHanhThang, sp.BaoHanhThang) // L
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_TinhTrang, sp.TinhTrang)       // M
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_TrangThai, sp.TrangThai)       // N
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_GiaBanLe, sp.GiaBanLe)         // O
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_GhiChu, sp.GhiChu)             // P
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_NguoiTao, sp.NguoiTao)         // Q
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_NgayTao, sp.NgayTao)           // R
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_NgayCapNhat, sp.NgayCapNhat)   // S
	}

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã lưu sản phẩm thành công!"})
}

// Helper sinh mã tự động
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
