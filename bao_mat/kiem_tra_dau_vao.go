package bao_mat

import (
	"regexp"
	"strings"
	"unicode/utf8"
)

// 1. Tên Đăng Nhập
// - 6-30 ký tự, a-z, 0-9, . _
// - Không bắt đầu/kết thúc bằng . _
// - Không chứa .. hoặc __ hoặc ._ hoặc _.
func KiemTraTenDangNhap(user string) bool {
	// Check độ dài
	if len(user) < 6 || len(user) > 30 { return false }
	
	// Check ký tự hợp lệ (Chỉ a-z, 0-9, ., _)
	// Lưu ý: Đã loại bỏ @, +, % theo yêu cầu
	match, _ := regexp.MatchString(`^[a-z0-9._]+$`, user)
	if !match { return false }

	// Check ký tự đầu và cuối (Phải là chữ hoặc số)
	firstChar := user[0]
	lastChar := user[len(user)-1]
	if !isAlphaNumeric(firstChar) || !isAlphaNumeric(lastChar) {
		return false
	}

	// Check liên tiếp (Thay thế cho Regex Lookahead)
	// Chặn .. và __ và ._ và _.
	if strings.Contains(user, "..") || strings.Contains(user, "__") || 
	   strings.Contains(user, "._") || strings.Contains(user, "_.") {
		return false
	}

	return true
}

// 2. Email
func KiemTraEmail(email string) bool {
	if len(email) < 6 || len(email) > 100 { return false }
	
	// Regex chặn @mail..com
	// Go hỗ trợ Non-capturing group (?:...) nhưng ko hỗ trợ Lookahead
	// Logic check domain: (?:[a-z0-9-]+\.)+ nghĩa là (Cụm-từ + Chấm) lặp lại
	match, _ := regexp.MatchString(`^[a-z0-9._%+\-]+@(?:[a-z0-9-]+\.)+[a-z]{2,}$`, email)
	return match
}

// 3. Mật khẩu
func KiemTraDinhDangMatKhau(pass string) bool {
	if len(pass) < 8 || len(pass) > 30 { return false }
	// Whitelist: a-z, A-Z, 0-9 và các ký tự đặc biệt cho phép
	match, _ := regexp.MatchString(`^[a-zA-Z0-9!@#$%^&*()\-+_.,?]+$`, pass)
	return match
}

// 4. Mã PIN
func KiemTraMaPin(pin string) bool {
	match, _ := regexp.MatchString(`^\d{8}$`, pin)
	return match
}

// 5. Họ Tên
func KiemTraHoTen(name string) bool {
	name = strings.TrimSpace(name)
	length := utf8.RuneCountInString(name)
	if length < 6 || length > 50 { return false }
	// Unicode và khoảng trắng
	match, _ := regexp.MatchString(`^[\p{L}\s]+$`, name)
	return match
}

// Helper kiểm tra chữ/số (byte)
func isAlphaNumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
}
