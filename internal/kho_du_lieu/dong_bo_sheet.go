package khodulieu

import (
	"fmt"
	"time"
)

// KhoiDongTacVuNgam: Bắt đầu chạy Worker
// done: channel báo hiệu khi worker đã tắt hẳn
func (k *KhoRAM) KhoiDongTacVuNgam(tinHieuTatServer chan bool) {
	fmt.Println("🤖 Tác vụ ngầm đã khởi động...")
	
	// Bộ đếm thời gian: 5 giây gõ 1 lần
	ticker := time.NewTicker(5 * time.Second)
	
	// Bộ đệm tạm để gom request (Batching)
	var danhSachChoGhi []YeuCauGhi

	for {
		select {
		// Trường hợp 1: Có yêu cầu mới vào hàng đợi
		case yeuCau := <-k.HangDoi:
			danhSachChoGhi = append(danhSachChoGhi, yeuCau)
			
			// Nếu gom đủ 50 yêu cầu thì ghi luôn, không chờ 5s nữa
			if len(danhSachChoGhi) >= 50 {
				k.GhiMeXuongSheet(danhSachChoGhi)
				danhSachChoGhi = nil // Reset bộ đệm
			}

		// Trường hợp 2: Đã hết 5 giây
		case <-ticker.C:
			if len(danhSachChoGhi) > 0 {
				fmt.Printf("⏳ Đang đồng bộ %d dòng xuống Sheet...\n", len(danhSachChoGhi))
				k.GhiMeXuongSheet(danhSachChoGhi)
				danhSachChoGhi = nil
			}

		// Trường hợp 3: Server nhận lệnh tắt (Graceful Shutdown)
		case <-tinHieuTatServer:
			fmt.Println("⚠️ Đang tắt server! Ghi nốt dữ liệu còn lại...")
			ticker.Stop()
			
			// Ghi nốt những gì còn trong bộ đệm
			if len(danhSachChoGhi) > 0 {
				k.GhiMeXuongSheet(danhSachChoGhi)
			}
			
			// Ghi nốt những gì còn sót trong channel (HangDoi)
			close(k.HangDoi)
			for yeuCau := range k.HangDoi {
				danhSachChoGhi = append(danhSachChoGhi, yeuCau)
			}
			if len(danhSachChoGhi) > 0 {
				k.GhiMeXuongSheet(danhSachChoGhi)
			}
			
			fmt.Println("✅ Đã lưu toàn bộ dữ liệu an toàn.")
			return // Thoát vòng lặp, kết thúc goroutine
		}
	}
}

// GhiMeXuongSheet: Hàm thực hiện gọi API Google (Bulk Update)
func (k *KhoRAM) GhiMeXuongSheet(danhSach []YeuCauGhi) {
	// Logic gom nhóm dữ liệu theo từng Bảng (Sheet)
	// Để tối ưu số lượng request gửi lên Google
	
	// (Phần này sẽ implement logic gọi sheets.values.batchUpdate 
	// hoặc values.append trong bước tiếp theo)
	
	// Giả lập log
	fmt.Println("--> Đã ghi thành công xuống Google Sheets")
}
