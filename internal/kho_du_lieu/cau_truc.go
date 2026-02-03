package khodulieu

import (
	"sync"
	"google.golang.org/api/sheets/v4"
)

// YeuCauGhi: Cấu trúc bản tin gửi vào hàng đợi (Queue)
type YeuCauGhi struct {
	TenBang     string        // Ví dụ: "SanPham", "KhachHang"
	LoaiThaoTac string        // "THEM", "SUA"
	DongDuLieu  []interface{} // Dữ liệu cần ghi
	ViTriDong   int           // Dùng cho việc sửa (nếu cần)
}

// KhoRAM: Database chạy trên RAM
type KhoRAM struct {
	// DuLieu: Key là Tên Bảng (Sheet Name), Value là danh sách các dòng
	DuLieu map[string][][]interface{}

	// KhoaBaoVe: Smart Locks - Mỗi bảng có 1 ổ khóa riêng
	KhoaBaoVe map[string]*sync.RWMutex

	// KhoaTong: Dùng khi tạo bảng mới chưa tồn tại trong map
	KhoaTong sync.Mutex

	// HangDoi: Kênh chứa các yêu cầu chờ ghi xuống Sheet
	HangDoi chan YeuCauGhi

	// DichVuSheet: Client kết nối Google API
	DichVuSheet *sheets.Service
	IDFileSheet string
}

// KhoiTaoKho: Hàm khởi tạo ban đầu
func KhoiTaoKho(srv *sheets.Service, spreadsheetID string) *KhoRAM {
	return &KhoRAM{
		DuLieu:      make(map[string][][]interface{}),
		KhoaBaoVe:   make(map[string]*sync.RWMutex),
		HangDoi:     make(chan YeuCauGhi, 1000), // Buffer 1000 request
		DichVuSheet: srv,
		IDFileSheet: spreadsheetID,
	}
}
