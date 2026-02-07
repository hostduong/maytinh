package bo_nho_dem

import (
	"log"
	"strings"
	"app/mo_hinh"
)

// [PHAN QUYEN - Đặt tại đây để tránh Cycle Import]
var CachePhanQuyen map[string]map[string]bool

func NapDuLieuPhanQuyen() {
	raw, err := loadSheetData("PHAN_QUYEN")
	if err != nil { return }

	// Không cần Lock ở đây vì hàm này được gọi trong luồng Boot/Reload đã được Lock từ bên ngoài.
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

func napKhuyenMai(target *KhoKhuyenMaiStore) { raw,_:=loadSheetData("KHUYEN_MAI"); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.KhuyenMai{MaVoucher:layString(r,mo_hinh.CotKM_MaVoucher)}; target.DuLieu[item.MaVoucher]=item } }
func napCauHinhWeb(target *KhoCauHinhWebStore) { raw,_:=loadSheetData("CAU_HINH_WEB"); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.CauHinhWeb{MaCauHinh:layString(r,mo_hinh.CotCH_MaCauHinh), GiaTri:layString(r,mo_hinh.CotCH_GiaTri)}; target.DuLieu[item.MaCauHinh]=item } }
