package nghiep_vu

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"app/mo_hinh"
)

var mtxNV sync.Mutex // Khóa an toàn để tránh xung đột khi nhiều người đăng ký cùng lúc

// =================================================================================
// 1. CÁC HÀM TÌM KIẾM (READ)
// =================================================================================

// Tìm nhân viên theo Cookie (Dùng cho Middleware kiểm tra đăng nhập)
func TimNhanVienTheoCookie(cookie string) (*mo_hinh.NhanVien, bool) {
	for _, nv := range CacheNhanVien.DuLieu {
		if nv.Cookie == cookie {
			return nv, true
		}
	}
	return nil, false
}

// Tìm nhân viên theo Tên Đăng Nhập (Dùng kiểm tra tồn tại)
func TimNhanVienTheoUsername(username string) (*mo_hinh.NhanVien, bool) {
	for _, nv := range CacheNhanVien.DuLieu {
		if nv.TenDangNhap == username {
			return nv, true
		}
	}
	return nil, false
}

// Tìm nhân viên theo User HOẶC Email (Dùng cho trang Đăng Nhập)
func TimNhanVienTheoUserHoacEmail(input string) (*mo_hinh.NhanVien, bool) {
	for _, nv := range CacheNhanVien.DuLieu {
		// So sánh input với Tên đăng nhập OR Email
		if nv.TenDangNhap == input || nv.Email == input {
			return nv, true
		}
	}
	return nil, false
}

// Kiểm tra xem Username hoặc Email đã tồn tại chưa (Dùng cho trang Đăng Ký)
func KiemTraTonTaiUserOrEmail(username string, email string) bool {
	for _, nv := range CacheNhanVien.DuLieu {
		if nv.TenDangNhap == username || (email != "" && nv.Email == email) {
			return true
		}
	}
	return false
}

// =================================================================================
// 2. CÁC HÀM LOGIC NGHIỆP VỤ (GENERATE ID & COUNT)
// =================================================================================

// Sinh Mã Nhân Viên Mới (NV_0001 -> NV_0002...)
func TaoMaNhanVienMoi() string {
	maxID := 0
	for _, nv := range CacheNhanVien.DuLieu {
		// Chỉ đếm các mã bắt đầu bằng "NV_"
		if strings.HasPrefix(nv.MaNhanVien, "NV_") {
			parts := strings.Split(nv.MaNhanVien, "_")
			if len(parts) == 2 {
				id, err := strconv.Atoi(parts[1])
				if err == nil && id > maxID {
					maxID = id
				}
			}
		}
	}
	return fmt.Sprintf("NV_%04d", maxID+1)
}

// Sinh Mã Khách Hàng Mới (KH_0001 -> KH_0002...)
func TaoMaKhachHangMoi() string {
	maxID := 0
	for _, nv := range CacheNhanVien.DuLieu {
		// Chỉ đếm các mã bắt đầu bằng "KH_"
		if strings.HasPrefix(nv.MaNhanVien, "KH_") {
			parts := strings.Split(nv.MaNhanVien, "_")
			if len(parts) == 2 {
				id, err := strconv.Atoi(parts[1])
				if err == nil && id > maxID {
					maxID = id
				}
			}
		}
	}
	return fmt.Sprintf("KH_%04d", maxID+1)
}

// Đếm tổng số lượng (Để xác định ai là Admin đầu tiên)
func DemSoLuongNhanVien() int {
	return len(CacheNhanVien.DuLieu)
}

// Lấy dòng trong Sheet của nhân viên
func LayDongNhanVien(maNV string) int {
	if nv, ok := CacheNhanVien.DuLieu[maNV]; ok {
		return nv.DongTrongSheet
	}
	return 0
}

// =================================================================================
// 3. CÁC HÀM GHI DỮ LIỆU (WRITE)
// =================================================================================

// Cập nhật Phiên làm việc (Khi đăng nhập thành công)
func CapNhatPhienDangNhap(maNV string, newCookie string, newExpired int64) {
	mtxNV.Lock()
	defer mtxNV.Unlock()

	nv, ok := CacheNhanVien.DuLieu[maNV]
	if !ok { return }

	// 1. Cập nhật RAM
	nv.Cookie = newCookie
	nv.CookieExpired = newExpired

	// 2. Đẩy vào Hàng Chờ Ghi Sheet
	ThemVaoHangCho(CacheNhanVien.SpreadsheetID, "NHAN_VIEN", nv.DongTrongSheet, mo_hinh.CotNV_Cookie, newCookie)
	ThemVaoHangCho(CacheNhanVien.SpreadsheetID, "NHAN_VIEN", nv.DongTrongSheet, mo_hinh.CotNV_CookieExpired, newExpired)
}

// Gia hạn Cookie (Auto Renew khi sắp hết hạn)
func CapNhatHanCookieRAM(maNV string, newExpired int64) {
	mtxNV.Lock()
	defer mtxNV.Unlock()
	
	nv, ok := CacheNhanVien.DuLieu[maNV]
	if !ok { return }
	
	nv.CookieExpired = newExpired
	ThemVaoHangCho(CacheNhanVien.SpreadsheetID, "NHAN_VIEN", nv.DongTrongSheet, mo_hinh.CotNV_CookieExpired, newExpired)
}

// Thêm Nhân Viên/Khách Hàng Mới (Đăng Ký)
func ThemNhanVienMoi(nv *mo_hinh.NhanVien) {
	mtxNV.Lock()
	defer mtxNV.Unlock()

	// 1. Tìm dòng trống tiếp theo trong Sheet
	maxRow := mo_hinh.DongBatDauDuLieu - 1
	for _, item := range CacheNhanVien.DuLieu {
		if item.DongTrongSheet > maxRow {
			maxRow = item.DongTrongSheet
		}
	}
	newRow := maxRow + 1
	nv.DongTrongSheet = newRow // Gán dòng mới cho user

	// 2. Lưu vào RAM ngay lập tức
	CacheNhanVien.DuLieu[nv.MaNhanVien] = nv

	// 3. Đẩy vào Hàng Chờ Ghi (Ghi đầy đủ 12 cột từ A -> L)
	sID := CacheNhanVien.SpreadsheetID
	sName := "NHAN_VIEN"

	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_MaNhanVien, nv.MaNhanVien)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_TenDangNhap, nv.TenDangNhap)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_Email, nv.Email)           // Cột C
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_MatKhauHash, nv.MatKhauHash)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_HoTen, nv.HoTen)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_ChucVu, nv.ChucVu)         // Cột F
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_MaPin, nv.MaPin)           // Cột G
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_Cookie, nv.Cookie)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_CookieExpired, nv.CookieExpired)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_VaiTroQuyenHan, nv.VaiTroQuyenHan)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_TrangThai, nv.TrangThai)
	ThemVaoHangCho(sID, sName, newRow, mo_hinh.CotNV_LanDangNhapCuoi, nv.LanDangNhapCuoi)
}
