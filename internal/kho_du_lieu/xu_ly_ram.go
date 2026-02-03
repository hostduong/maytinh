package khodulieu

import (
	"sync"
	"fmt"
)

// layKhoaCuaBang: Hàm nội bộ để lấy ổ khóa tương ứng với bảng
func (k *KhoRAM) layKhoaCuaBang(tenBang string) *sync.RWMutex {
	k.KhoaTong.Lock()
	defer k.KhoaTong.Unlock()

	// Nếu chưa có khóa cho bảng này thì tạo mới
	if _, tonTai := k.KhoaBaoVe[tenBang]; !tonTai {
		k.KhoaBaoVe[tenBang] = &sync.RWMutex{}
	}
	return k.KhoaBaoVe[tenBang]
}

// LayDuLieu: Đọc dữ liệu từ RAM (Ưu tiên tốc độ)
func (k *KhoRAM) LayDuLieu(tenBang string) ([][]interface{}, error) {
	mu := k.layKhoaCuaBang(tenBang)
	
	// 1. Khóa đọc (Cho phép nhiều người cùng đọc)
	mu.RLock()
	data, tonTai := k.DuLieu[tenBang]
	mu.RUnlock()

	// 2. Nếu RAM có dữ liệu -> Trả về ngay
	if tonTai {
		return data, nil
	}

	// 3. Nếu RAM chưa có -> Tải từ Google Sheet (Cache Miss)
	// Lưu ý: Lúc này cần khóa Ghi để tải data vào
	mu.Lock()
	defer mu.Unlock()
	
	// Kiểm tra lại lần nữa (Double-check locking)
	if data, tonTai = k.DuLieu[tenBang]; tonTai {
		return data, nil
	}

	// Gọi hàm tải từ Google (sẽ viết ở file ket_noi_google.go)
	dataTuGoogle, err := k.TaiDuLieuGoc(tenBang)
	if err != nil {
		return nil, err
	}

	// Lưu vào RAM
	k.DuLieu[tenBang] = dataTuGoogle
	return dataTuGoogle, nil
}

// ThemMoi: Ghi dữ liệu vào RAM trước, đẩy vào Queue sau
func (k *KhoRAM) ThemMoi(tenBang string, dongMoi []interface{}) error {
	mu := k.layKhoaCuaBang(tenBang)

	// 1. Cập nhật RAM (Khóa Ghi - Chặn người khác đọc lúc đang ghi)
	mu.Lock()
	// Logic append vào slice trong RAM
	if _, tonTai := k.DuLieu[tenBang]; !tonTai {
		k.DuLieu[tenBang] = make([][]interface{}, 0)
	}
	k.DuLieu[tenBang] = append(k.DuLieu[tenBang], dongMoi)
	mu.Unlock()

	// 2. Đẩy vào hàng đợi (Không cần khóa, channel tự xử lý đồng bộ)
	select {
	case k.HangDoi <- YeuCauGhi{
		TenBang:     tenBang,
		LoaiThaoTac: "THEM",
		DongDuLieu:  dongMoi,
	}:
		// Đẩy thành công
		return nil
	default:
		// Queue bị đầy (Trường hợp rất hiếm nếu set buffer lớn)
		return fmt.Errorf("hàng đợi đang bận, vui lòng thử lại sau")
	}
}
