package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	
	"app/cau_hinh"
	"app/chuc_nang"
	"app/kho_du_lieu"
	"app/nghiep_vu" // Import gói nghiệp vụ

	"github.com/gin-gonic/gin"
)

func main() {
	// ... (Các bước khởi tạo cũ: 1, 2, 3 giữ nguyên) ...
	cau_hinh.KhoiTaoCauHinh()
	kho_du_lieu.KhoiTaoKetNoiGoogle()
	nghiep_vu.KhoiTaoBoNho()

	// --- [MỚI] KHỞI ĐỘNG WORKER GHI SHEET ---
	nghiep_vu.KhoiTaoWorkerGhiSheet()
    chuc_nang.KhoiTaoBoDemRateLimit() // Worker reset bộ đếm request
	// 4. Khởi tạo Web Server
	router := gin.Default()
	// ... (Các config router giữ nguyên) ...
	
	// ... (Đoạn router.Run cũ XÓA ĐI hoặc comment lại) ...
	
	// --- [MỚI] CHẠY SERVER VỚI CƠ CHẾ GRACEFUL SHUTDOWN ---
	port := cau_hinh.BienCauHinh.CongChayWeb
	if port == "" { port = "8080" }

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Chạy server trong 1 Goroutine riêng
	go func() {
		log.Println(">> Server đang chạy tại cổng: " + port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Lỗi khởi động: %s\n", err)
		}
	}()

	// --- LẮNG NGHE TÍN HIỆU TẮT (SIGTERM) ---
	quit := make(chan os.Signal, 1)
	// Cloud Run gửi SIGTERM khi tắt container
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	
	<-quit // Code sẽ dừng ở đây chờ tín hiệu...
	log.Println("⚠️  Đang tắt Server... Đang xả hàng đợi lần cuối!")

	// GỌI HÀM XẢ HÀNG KHẨN CẤP
	nghiep_vu.ThucHienGhiSheet(true)

	log.Println("✅ Server đã tắt an toàn.")
}
