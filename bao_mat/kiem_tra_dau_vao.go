package bao_mat

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// 1. Kiểm tra Tên Đăng Nhập
func KiemTraTenDangNhap(user string) bool {
	if len(user) < 6 || len(user) > 30 {
		return false
	}
	match, _ := regexp.MatchString(`^[a-z0-9._]+$`, strings.ToLower(user))
	return match
}

// 2. Kiểm tra Email
func KiemTraEmail(email string) bool {
	if len(email) < 6 || len(email) > 100 {
		return false
	}
	match, _ := regexp.MatchString(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, email)
	return match
}

// 3. Kiểm tra Mật khẩu (Validation)
// [ĐÃ SỬA TÊN HÀM ĐỂ TRÁNH TRÙNG]
func KiemTraDinhDangMatKhau(pass string) bool {
	if len(pass) < 8 || len(pass) > 30 {
		return false
	}
	if strings.ContainsAny(pass, `' "<>`) || strings.Contains(pass, " ") {
		return false
	}
	return true
}

// 4. Kiểm tra Mã PIN
func KiemTraMaPin(pin string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, pin)
	return match
}

// 5. Kiểm tra Họ Tên
func KiemTraHoTen(name string) bool {
	length := utf8.RuneCountInString(name)
	if length < 6 || length > 50 {
		return false
	}
	match, _ := regexp.MatchString(`^[\p{L}\s]+$`, name)
	return match
}
