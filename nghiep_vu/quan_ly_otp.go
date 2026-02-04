package nghiep_vu

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"
)

type ThongTinOTP struct {
	MaCode    string
	HetHanLuc int64
}

// Map lưu OTP: Key là Username -> Value là Code
var CacheOTP = make(map[string]ThongTinOTP)
var mtxOTP sync.Mutex

// Sinh mã 8 số ngẫu nhiên
func TaoMaOTP() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(99999999))
	return fmt.Sprintf("%08d", n.Int64())
}

// Lưu OTP vào RAM (TTL 10 phút)
func LuuOTP(userKey string, code string) {
	mtxOTP.Lock()
	defer mtxOTP.Unlock()
	CacheOTP[userKey] = ThongTinOTP{
		MaCode:    code,
		HetHanLuc: time.Now().Add(10 * time.Minute).Unix(),
	}
}

// Kiểm tra OTP có đúng và còn hạn không
func KiemTraOTP(userKey string, inputCode string) bool {
	mtxOTP.Lock()
	defer mtxOTP.Unlock()
	
	otp, ok := CacheOTP[userKey]
	if !ok { return false }

	// Check hạn
	if time.Now().Unix() > otp.HetHanLuc {
		delete(CacheOTP, userKey) // Xóa rác
		return false
	}

	// Check mã
	if otp.MaCode == inputCode {
		delete(CacheOTP, userKey) // Dùng xong xóa luôn (One-time use)
		return true
	}
	return false
}
