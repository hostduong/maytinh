package bo_nho_dem

import (
	"log"
	"strings"
	"app/mo_hinh"
)

// [PHAN QUYEN - Đặt tại đây để tránh Cycle Import]
var CachePhanQuyen map[string]map[string]bool

func napPhanQuyen() {
	raw, err := loadSheetData("PHAN_QUYEN")
	if err != nil { return }

	CachePhanQuyen = make(map[string]map[string]bool)
	roles := []string{"ADMIN", "QUAN_LY", "THU_KHO", "KE_TOAN", "SALE", "CONTENT"}
	for _, r := range roles { CachePhanQuyen[r] = make(map[string]bool) }

	for i, row := range raw {
		if i == 0 || len(row) <= mo_hinh.CotPQ_MaChucNang { continue }
		maChucNang := strings.TrimSpace(layString(row, mo_hinh.CotPQ_MaChucNang))
		
		check := func(idx int) bool { return layString(row, idx) == "1" || strings.ToLower(layString(row, idx)) == "true" }
		
		if check(mo_hinh.CotPQ_ADMIN)   { CachePhanQuyen["ADMIN"][maChucNang] = true }
		if check(mo_hinh.CotPQ_QUAN_LY) { CachePhanQuyen["QUAN_LY"][maChucNang] = true }
		if check(mo_hinh.CotPQ_THU_KHO) { CachePhanQuyen["THU_KHO"][maChucNang] = true }
		if check(mo_hinh.CotPQ_KE_TOAN) { CachePhanQuyen["KE_TOAN"][maChucNang] = true }
		if check(mo_hinh.CotPQ_SALE)    { CachePhanQuyen["SALE"][maChucNang] = true }
		if check(mo_hinh.CotPQ_CONTENT) { CachePhanQuyen["CONTENT"][maChucNang] = true }
	}
	log.Println("✅ [RBAC] Đã nạp phân quyền.")
}
