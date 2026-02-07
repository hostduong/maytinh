package bo_nho_dem

import "app/mo_hinh"

func napPhieuNhap(target *KhoPhieuNhapStore) { raw,_:=loadSheetData("PHIEU_NHAP"); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.PhieuNhap{MaPhieuNhap:layString(r,mo_hinh.CotPN_MaPhieuNhap)}; target.DuLieu[item.MaPhieuNhap]=item; target.DanhSach=append(target.DanhSach,item) } }
func napChiTietPhieuNhap(target *KhoChiTietPhieuNhapStore) { raw,_:=loadSheetData("CHI_TIET_PHIEU_NHAP"); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.ChiTietPhieuNhap{MaPhieuNhap:layString(r,mo_hinh.CotCTPN_MaPhieuNhap)}; target.DanhSach=append(target.DanhSach,item) } }

func napPhieuXuat(target *KhoPhieuXuatStore) { 
	raw,_:=loadSheetData("PHIEU_XUAT")
	for i,r:=range raw{ 
		if i<10{continue}
		item:=mo_hinh.PhieuXuat{
			MaPhieuXuat:layString(r,mo_hinh.CotPX_MaPhieuXuat),
			TongTienPhieu:layFloat(r,mo_hinh.CotPX_TongTienPhieu),
			TrangThai:layString(r,mo_hinh.CotPX_TrangThai),
			NgayTao:layString(r,mo_hinh.CotPX_NgayTao),
		}
		target.DuLieu[item.MaPhieuXuat]=item
		target.DanhSach=append(target.DanhSach,item) 
	} 
}
func napChiTietPhieuXuat(target *KhoChiTietPhieuXuatStore) { raw,_:=loadSheetData("CHI_TIET_PHIEU_XUAT"); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.ChiTietPhieuXuat{MaPhieuXuat:layString(r,mo_hinh.CotCTPX_MaPhieuXuat), MaSanPham:layString(r,mo_hinh.CotCTPX_MaSanPham), DonGiaBan:layFloat(r,mo_hinh.CotCTPX_DonGiaBan)}; target.DanhSach=append(target.DanhSach,item) } }

func napHoaDon(target *KhoHoaDonStore) { raw,_:=loadSheetData("HOA_DON"); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.HoaDon{MaHoaDon:layString(r,mo_hinh.CotHD_MaHoaDon)}; target.DuLieu[item.MaHoaDon]=item } }
func napHoaDonChiTiet(target *KhoHoaDonChiTietStore) { raw,_:=loadSheetData("HOA_DON_CHI_TIET"); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.HoaDonChiTiet{MaHoaDon:layString(r,mo_hinh.CotHDCT_MaHoaDon)}; target.DanhSach=append(target.DanhSach,item) } }

func napPhieuThuChi(target *KhoPhieuThuChiStore) { raw,_:=loadSheetData("PHIEU_THU_CHI"); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.PhieuThuChi{MaPhieuThuChi:layString(r,mo_hinh.CotPTC_MaPhieuThuChi)}; target.DuLieu[item.MaPhieuThuChi]=item; target.DanhSach=append(target.DanhSach,item) } }
func napPhieuBaoHanh(target *KhoPhieuBaoHanhStore) { raw,_:=loadSheetData("PHIEU_BAO_HANH"); for i,r:=range raw{ if i<10{continue}; item:=mo_hinh.PhieuBaoHanh{MaPhieuBaoHanh:layString(r,mo_hinh.CotPBH_MaPhieuBaoHanh)}; target.DuLieu[item.MaPhieuBaoHanh]=item; target.DanhSach=append(target.DanhSach,item) } }
