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
	if len(user) < 6 || len(user) > 30 {
		return false
	}
	// Cho phép chữ hoa/thường, số, dấu chấm, gạch dưới
	match, _ := regexp.MatchString(`^[a-zA-Z0-9._]+$`, user)
	return match
}

// 2. Kiểm tra Email
// - Từ 6-100 ký tự
// - Chỉ chứa chữ, số, . _ - + % @
func KiemTraEmail(email string) bool {
	if len(email) < 6 || len(email) > 100 {
		return false
	}
	// Regex email tiêu chuẩn, không cho phép ký tự lạ
	match, _ := regexp.MatchString(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, email)
	return match
}

// 3. Kiểm tra Mật khẩu
// - Từ 8-30 ký tự
// - Không chứa ', ", <, >, dấu cách
func KiemTraDinhDangMatKhau(pass string) bool {
	if len(pass) < 8 || len(pass) > 30 {
		return false
	}
	// Kiểm tra các ký tự cấm: nháy đơn, nháy kép, ngoặc nhọn, dấu cách
	if strings.ContainsAny(pass, `' "<>`) || strings.Contains(pass, " ") {
		return false
	}
	return true
}

// 4. Kiểm tra Mã PIN
// - Đúng 8 số
func KiemTraMaPin(pin string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, pin)
	return match
}

// 5. Kiểm tra Họ Tên
// - Từ 6-50 ký tự
// - Không số, không ký tự đặc biệt (Chỉ chấp nhận chữ cái Unicode và khoảng trắng)
func KiemTraHoTen(name string) bool {
	length := utf8.RuneCountInString(name)
	if length < 6 || length > 50 {
		return false
	}
	// \p{L} là tất cả chữ cái Unicode (bao gồm tiếng Việt), \s là khoảng trắng
	match, _ := regexp.MatchString(`^[\p{L}\s]+$`, name)
	return match
}
