package bao_mat

import (
	"crypto/rand"
	"encoding/hex"
	"golang.org/x/crypto/bcrypt"
)

// HashMatKhau : (Giữ nguyên)
func HashMatKhau(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// KiemTraMatKhau : (Giữ nguyên)
func KiemTraMatKhau(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// --- [MỚI] Tạo Session ID siêu dài (128 ký tự) ---
// Dùng thư viện crypto/rand để đảm bảo tính ngẫu nhiên bảo mật cao nhất
func TaoSessionIDAnToan() string {
	b := make([]byte, 64) // 64 byte
	_, err := rand.Read(b)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(b) // Chuyển sang Hex -> Thành 128 ký tự
}
