package bao_mat

import (
	"golang.org/x/crypto/bcrypt"
)

// HashMatKhau : Biến mật khẩu thô thành chuỗi mã hóa
func HashMatKhau(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14) // Cost 14 là rất an toàn
	return string(bytes), err
}

// KiemTraMatKhau : So sánh mật khẩu nhập vào với mật khẩu đã mã hóa
func KiemTraMatKhau(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
