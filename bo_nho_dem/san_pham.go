package bo_nho_dem

import "app/mo_hinh"

func napSanPham(target *KhoSanPhamStore) {
	raw, err := loadSheetData("SAN_PHAM")
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
			MaDanhMuc: layString(r, mo_hinh.CotSP_MaDanhMuc),
			MaThuongHieu: layString(r, mo_hinh.CotSP_MaThuongHieu),
			BaoHanhThang: layInt(r, mo_hinh.CotSP_BaoHanhThang),
			ThongSo: layString(r, mo_hinh.CotSP_ThongSo),
		}
		target.DuLieu[item.MaSanPham] = item
		target.DanhSach = append(target.DanhSach, item)
	}
}

func napDanhMuc(target *KhoDanhMucStore) { 
	raw,_:=loadSheetData("DANH_MUC")
	for i,r:=range raw{ 
		if i<10{continue}
		item:=mo_hinh.DanhMuc{
			MaDanhMuc:layString(r,mo_hinh.CotDM_MaDanhMuc),
			TenDanhMuc:layString(r,mo_hinh.CotDM_TenDanhMuc),
		}
		target.DuLieu[item.MaDanhMuc]=item 
	} 
}

func napThuongHieu(target *KhoThuongHieuStore) { 
	raw,_:=loadSheetData("THUONG_HIEU")
	for i,r:=range raw{ 
		if i<10{continue}
		item:=mo_hinh.ThuongHieu{
			MaThuongHieu:layString(r,mo_hinh.CotTH_MaThuongHieu),
			LogoUrl:layString(r,mo_hinh.CotTH_LogoUrl),
		}
		target.DuLieu[item.MaThuongHieu]=item 
	} 
}
