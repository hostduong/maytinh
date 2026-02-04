package bao_mat

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// 1. Kiểm tra Tên Đăng Nhập
// - 6-30 ký tự
// - Chỉ gồm chữ thường, số, dấu _, dấu .
// - Không chứa dấu cách, không dấu tiếng Việt
func KiemTraTenDangNhap(user string) bool {
	if len(user) < 6 || len(user) > 30 {
		return false
	}
	// Regex: Bắt đầu và kết thúc chỉ chứa a-z, 0-9, ., _
	match, _ := regexp.MatchString(`^[a-z0-9._]+$`, strings.ToLower(user))
	return match
}

// 2. Kiểm tra Email
// - 6-100 ký tự
// - Đúng định dạng email, không chứa ký tự lạ nguy hiểm
func KiemTraEmail(email string) bool {
	if len(email) < 6 || len(email) > 100 {
		return false
	}
	// Regex email tiêu chuẩn + chặn các ký tự đặc biệt nguy hiểm
	match, _ := regexp.MatchString(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`, email)
	return match
}

// 3. Kiểm tra Mật khẩu
// - 8-30 ký tự
// - Không chứa: ' " < > space
func KiemTraMatKhau(pass string) bool {
	if len(pass) < 8 || len(pass) > 30 {
		return false
	}
	// Kiểm tra các ký tự CẤM
	if strings.ContainsAny(pass, `' "<>`) || strings.Contains(pass, " ") {
		return false
	}
	return true
}

// 4. Kiểm tra Mã PIN
// - Bắt buộc đúng 8 số (Theo yêu cầu mới)
func KiemTraMaPin(pin string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, pin)
	return match
}

// 5. Kiểm tra Họ Tên
// - 6-50 ký tự
// - Chấp nhận tiếng Việt có dấu (Unicode)
// - Không chứa số, không chứa ký tự đặc biệt (trừ khoảng trắng)
func KiemTraHoTen(name string) bool {
	length := utf8.RuneCountInString(name)
	if length < 6 || length > 50 {
		return false
	}
	// Regex: \p{L} là chữ cái Unicode (bao gồm tiếng Việt), \s là khoảng trắng
	match, _ := regexp.MatchString(`^[\p{L}\s]+$`, name)
	return match
}
