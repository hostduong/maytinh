package nghiep_vu

import (
	"log"
	"strings"
	"sync"
	"app/mo_hinh"
)

// BỘ NHỚ ĐỆM PHÂN QUYỀN
// Map cấp 1: Mã Vai Trò (VD: "THU_KHO")
// Map cấp 2: Mã Chức Năng (VD: "product.view") -> true/false
var CachePhanQuyen map[string]map[string]bool
var mtxPhanQuyen sync.RWMutex

// Key Cache để dùng Lock file (tránh tranh chấp khi Admin đang sửa quyền trên Sheet)
var KeyCachePhanQuyen = "CACHE_PHAN_QUYEN"

// 1. HÀM NẠP DỮ LIỆU TỪ SHEET VÀO RAM
func NapDuLieuPhanQuyen() {
	// Gọi hàm load chung từ bo_nho_dem.go (hoặc viết lại logic loadSheetData ở đây nếu cần tách biệt)
	// Ở đây tôi dùng hàm loadSheetData giả định bạn đã có trong bo_nho_dem.go
	// Nếu chưa public hàm loadSheetData, ta sẽ copy logic đọc sheet vào đây.
	
	raw, err := loadSheetData("PHAN_QUYEN", KeyCachePhanQuyen)
	if err != nil {
		log.Println("⚠️ Chưa có Sheet PHAN_QUYEN, hệ thống sẽ chạy chế độ mặc định (Chặn tất cả trừ AdminRoot)")
		return
	}

	mtxPhanQuyen.Lock()
	defer mtxPhanQuyen.Unlock()

	// Khởi tạo lại Map
	CachePhanQuyen = make(map[string]map[string]bool)

	// Các vai trò cần duyệt (Khớp với Header Sheet)
	roles := []string{"ADMIN", "QUAN_LY", "THU_KHO", "KE_TOAN", "SALE", "CONTENT"}
	
	// Khởi tạo map cho từng role
	for _, r := range roles {
		CachePhanQuyen[r] = make(map[string]bool)
	}

	for i, row := range raw {
		if i == 0 { continue } // Bỏ qua dòng tiêu đề
		
		// Lấy Mã Chức Năng (Cột A)
		if len(row) <= mo_hinh.CotPQ_MaChucNang { continue }
		maChucNang := strings.TrimSpace(layString(row, mo_hinh.CotPQ_MaChucNang))
		if maChucNang == "" { continue }

		// Helper: Đọc giá trị 1/0
		checkQuyen := func(colIndex int) bool {
			val := layString(row, colIndex)
			return val == "1" || strings.ToLower(val) == "true"
		}

		// Gán quyền vào RAM
		if checkQuyen(mo_hinh.CotPQ_ADMIN)   { CachePhanQuyen["ADMIN"][maChucNang] = true }
		if checkQuyen(mo_hinh.CotPQ_QUAN_LY) { CachePhanQuyen["QUAN_LY"][maChucNang] = true }
		if checkQuyen(mo_hinh.CotPQ_THU_KHO) { CachePhanQuyen["THU_KHO"][maChucNang] = true }
		if checkQuyen(mo_hinh.CotPQ_KE_TOAN) { CachePhanQuyen["KE_TOAN"][maChucNang] = true }
		if checkQuyen(mo_hinh.CotPQ_SALE)    { CachePhanQuyen["SALE"][maChucNang] = true }
		if checkQuyen(mo_hinh.CotPQ_CONTENT) { CachePhanQuyen["CONTENT"][maChucNang] = true }
	}

	log.Printf("✅ [RBAC] Đã nạp phân quyền cho %d chức năng.", len(raw)-1)
}

// 2. HÀM KIỂM TRA QUYỀN (GATEKEEPER)
// Input: vaiTroHienTai (VD: "thu_kho"), maChucNang (VD: "product.create")
func KiemTraQuyen(vaiTroHienTai string, maChucNang string) bool {
	// 1. Admin Root (Founder) - Quyền lực tối thượng
	// Check cứng trong code để tránh trường hợp lỡ tay xóa quyền trên Sheet
	if vaiTroHienTai == "admin_root" || vaiTroHienTai == "quan_tri_vien_cap_cao" {
		return true
	}

	mtxPhanQuyen.RLock()
	defer mtxPhanQuyen.RUnlock()

	// 2. Chuẩn hóa Key (Về chữ hoa để khớp Map)
	roleKey := strings.ToUpper(vaiTroHienTai)

	// 3. Tra cứu
	if listQuyen, ok := CachePhanQuyen[roleKey]; ok {
		if allowed, exist := listQuyen[maChucNang]; exist {
			return allowed
		}
	}

	// Mặc định: CHẶN (Deny all) nếu không tìm thấy cấu hình
	return false
}
