package bao_mat

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// 1. Kiểm tra Tên Đăng Nhập
// - Từ 6-30 ký tự
// - Chỉ gồm chữ, số, dấu _ và dấu .
func KiemTraTenDangNhap(user string) bool {
	if len(user) < 6 || len(user) > 30 { return false }
	// Regex mới: Chữ, số, chấm, gạch dưới
	match, _ := regexp.MatchString(`^[a-zA-Z0-9._]{6,30}$`, user)
	return match
}

// 2. Kiểm tra Email
// - Từ 6-100 ký tự
// - Chuẩn định dạng email, không dấu cách
func KiemTraEmail(email string) bool {
	if len(email) < 6 || len(email) > 100 { return false }
	match, _ := regexp.MatchString(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, email)
	return match
}

// 3. Kiểm tra Mật khẩu
// - Từ 8-30 ký tự
// - Không chứa ', ", <, >, dấu cách, Emoji
func KiemTraDinhDangMatKhau(pass string) bool {
	if len(pass) < 8 || len(pass) > 30 { return false }
	// Regex: Không chứa ký tự cấm và không chứa Symbol (Emoji)
	match, _ := regexp.MatchString(`^[^'"<>\s\p{So}]{8,30}$`, pass)
	return match
}

// 4. Kiểm tra Mã PIN
// - Đúng 8 số
func KiemTraMaPin(pin string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, pin)
	return match
}

// 5. Kiểm tra Họ Tên
// - Từ 6-50 ký tự
// - Chỉ chấp nhận chữ cái Unicode và khoảng trắng
func KiemTraHoTen(name string) bool {
	// Trim space trước khi đếm
	name = strings.TrimSpace(name)
	length := utf8.RuneCountInString(name)
	if length < 6 || length > 50 { return false }
	match, _ := regexp.MatchString(`^[\p{L}\s]+$`, name)
	return match
}
