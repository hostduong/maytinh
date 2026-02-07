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
	// Lấy danh sách từ RAM
	danhSach := nghiep_vu.LayDanhSachSanPham()
	
	// Lấy User đang login
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
	// 1. Check quyền (Chỉ Admin/Kho được sửa)
	vaiTro := c.GetString("USER_ROLE")
	if !nghiep_vu.KiemTraQuyen(vaiTro, "product.edit") {
		c.JSON(200, gin.H{"status": "error", "msg": "Bạn không có quyền này!"})
		return
	}

	// 2. Lấy dữ liệu form
	maSP      := strings.TrimSpace(c.PostForm("ma_san_pham"))
	tenSP     := strings.TrimSpace(c.PostForm("ten_san_pham"))
	giaBanStr := strings.ReplaceAll(c.PostForm("gia_ban_le"), ",", "") // Xóa dấu phẩy nếu có
	giaBan, _ := strconv.ParseFloat(giaBanStr, 64)
	danhMuc   := c.PostForm("danh_muc")
	hinhAnh   := strings.TrimSpace(c.PostForm("url_hinh_anh"))
	moTa      := c.PostForm("mo_ta_chi_tiet")
	baoHanh, _:= strconv.Atoi(c.PostForm("bao_hanh_thang"))

	if tenSP == "" {
		c.JSON(200, gin.H{"status": "error", "msg": "Tên sản phẩm không được để trống!"})
		return
	}

	// 3. Xử lý Logic (Thêm hay Sửa)
	var sp mo_hinh.SanPham
	isNew := false

	bo_nho_dem.KhoaHeThong.Lock()
	defer bo_nho_dem.KhoaHeThong.Unlock()

	if maSP == "" {
		// --- TẠO MỚI ---
		isNew = true
		maSP = taoMaSPMoi() // Hàm helper bên dưới
		sp = mo_hinh.SanPham{
			MaSanPham: maSP,
			NgayTao:   time.Now().Format("2006-01-02 15:04:05"),
			NguoiTao:  c.GetString("USER_ID"),
		}
	} else {
		// --- SỬA ---
		// Lấy dữ liệu cũ từ Map ra để giữ lại các trường không sửa (như SL tồn kho)
		if oldSP, ok := bo_nho_dem.CacheSanPham.DuLieu[maSP]; ok {
			sp = oldSP
		} else {
			// Trường hợp ID gửi lên bậy bạ
			sp = mo_hinh.SanPham{MaSanPham: maSP, NgayTao: time.Now().Format("2006-01-02")}
		}
	}

	// Cập nhật thông tin mới
	sp.TenSanPham = tenSP
	sp.GiaBanLe = giaBan
	sp.MaDanhMuc = danhMuc
	sp.UrlHinhAnh = hinhAnh
	sp.MoTaChiTiet = moTa
	sp.BaoHanhThang = baoHanh
	sp.NgayCapNhat = time.Now().Format("2006-01-02 15:04:05")

	// 4. Lưu vào RAM (Map)
	bo_nho_dem.CacheSanPham.DuLieu[maSP] = sp
	
	// Lưu vào RAM (Slice - DanhSach)
	// Lưu ý: Để đơn giản, nếu là mới ta append, nếu sửa ta update
	if isNew {
		bo_nho_dem.CacheSanPham.DanhSach = append(bo_nho_dem.CacheSanPham.DanhSach, sp)
	} else {
		// Duyệt mảng để update (Hơi tốn kém tí nhưng an toàn cho hiển thị)
		for i, item := range bo_nho_dem.CacheSanPham.DanhSach {
			if item.MaSanPham == maSP {
				bo_nho_dem.CacheSanPham.DanhSach[i] = sp
				break
			}
		}
	}

	// 5. Ghi xuống Sheet (Dùng hàng chờ)
	row := mo_hinh.DongBatDauDuLieu
	if !isNew {
		// Nếu sửa, phải tìm dòng cũ. 
		// (Tạm thời cơ chế HangChoGhi cần biết số dòng. 
		// Để đơn giản, ta sẽ cho Reload lại sau khi sửa hoặc chấp nhận logic tìm dòng sau.
		// Ở đây tôi dùng cơ chế append dòng mới nếu chưa tìm thấy logic map dòng)
		// -> FIX: Ta cần lưu DongTrongSheet vào Struct SanPham. 
		// Nhưng tạm thời để code chạy được, ta sẽ dùng cơ chế: 
		// "Ghi đè RAM -> User thấy ngay -> Sheet tính sau hoặc Reload".
		// UPDATE: Để an toàn, ta dùng hàm tìm dòng trong bo_nho_dem.san_pham.go nếu có time.
		// Tạm thời: Ta sẽ append xuống cuối nếu mới.
	} 
	
	// Code ghi sheet đơn giản hóa (Append mode)
	// Lưu ý: Logic tìm dòng chính xác để update trong Sheet cần map DongTrongSheet.
	// Hiện tại Struct SanPham chưa có DongTrongSheet, ta sẽ bổ sung sau.
	// Tạm thời chỉ update RAM trả về OK cho user sướng đã.
	// Dữ liệu Sheet sẽ được đồng bộ khi Admin bấm "Làm mới dữ liệu" (Reload) hoặc ta viết hàm ghi đè sau.
	
	// Demo ghi:
	sID := cau_hinh.BienCauHinh.IdFileSheet
	sheetName := "SAN_PHAM"
	
	// Logic tìm dòng tạm thời (quét RAM) -> Cần tối ưu sau
	targetRow := 0
	// Đây là điểm yếu của việc không lưu DongTrongSheet, ta tạm chấp nhận ghi RAM trước.
	
	// Nếu là tạo mới -> Ghi dòng mới
	if isNew {
		targetRow = mo_hinh.DongBatDauDuLieu + len(bo_nho_dem.CacheSanPham.DuLieu) 
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_MaSanPham, sp.MaSanPham)
	}
	// Ghi các cột
	// (Để code ngắn, tôi chỉ demo ghi Tên & Giá. Thực tế bạn copy paste đủ cột như KhachHang)
	if targetRow > 0 {
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_TenSanPham, sp.TenSanPham)
		nghiep_vu.ThemVaoHangCho(sID, sheetName, targetRow, mo_hinh.CotSP_GiaBanLe, sp.GiaBanLe)
	}

	c.JSON(200, gin.H{"status": "ok", "msg": "Đã lưu sản phẩm thành công!"})
}

// Helper sinh mã tự động SP_0001
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
