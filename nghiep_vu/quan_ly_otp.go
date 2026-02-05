package nghiep_vu

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"
)

// Cấu hình API gửi thư
const (
	URL_API_MAIL = "https://script.google.com/macros/s/AKfycbxd40H4neotKdnL54uQevZgSZpyZKXWfV7kJhNLY0oD9pPPA5Mn75KlFWvFd5WqiokZyA/exec"
	KEY_API_MAIL = "A1qPqCeLaX9oO0ozrMiH1a2IJKFDaj095Dlhmr8STXuS3cCmOe"
)

type ThongTinOTP struct {
	MaCode    string
	HetHanLuc int64
}

// Cấu trúc đếm để chặn Spam
type BoDemRate struct {
	LanGuiCuoi   int64
	SoLanTrong6h int
	ResetLuc     int64
}

var CacheOTP = make(map[string]ThongTinOTP)
var CacheRate = make(map[string]*BoDemRate)
var mtxOTP sync.Mutex

// --- [NHÓM HÀM CŨ - GIỮ NGUYÊN] ---

// TaoMaOTP : Sinh mã 8 số (Dùng cho logic cũ)
func TaoMaOTP() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(99999999))
	return fmt.Sprintf("%08d", n.Int64())
}

// LuuOTP : Lưu vào RAM
func LuuOTP(userKey string, code string) {
	mtxOTP.Lock()
	defer mtxOTP.Unlock()
	CacheOTP[userKey] = ThongTinOTP{
		MaCode:    code,
		HetHanLuc: time.Now().Add(10 * time.Minute).Unix(),
	}
}

// KiemTraOTP : Kiểm tra mã
func KiemTraOTP(userKey string, inputCode string) bool {
	mtxOTP.Lock()
	defer mtxOTP.Unlock()
	
	otp, ok := CacheOTP[userKey]
	if !ok { return false }

	if time.Now().Unix() > otp.HetHanLuc {
		delete(CacheOTP, userKey)
		return false
	}

	if otp.MaCode == inputCode {
		delete(CacheOTP, userKey)
		return true
	}
	return false
}

// --- [NHÓM HÀM MỚI - BỔ SUNG] ---

// TaoMaOTP6So : Sinh mã 6 số (Dùng cho gửi mail OTP mật khẩu)
func TaoMaOTP6So() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(999999))
	return fmt.Sprintf("%06d", n.Int64())
}

// KiemTraRateLimit : Chặn spam (1 phút/lần, 10 lần/6h)
func KiemTraRateLimit(email string) (bool, string) {
	mtxOTP.Lock()
	defer mtxOTP.Unlock()

	now := time.Now().Unix()
	rd, ok := CacheRate[email]

	if !ok || now > rd.ResetLuc {
		CacheRate[email] = &BoDemRate{ResetLuc: now + (6 * 3600)}
		rd = CacheRate[email]
	}

	if now-rd.LanGuiCuoi < 60 {
		return false, fmt.Sprintf("Vui lòng đợi %d giây.", 60-(now-rd.LanGuiCuoi))
	}
	if rd.SoLanTrong6h >= 10 {
		return false, "Bạn đã vượt quá 10 lần gửi trong 6 giờ."
	}
	rd.LanGuiCuoi = now
	rd.SoLanTrong6h++
	return true, ""
}

// Gọi API bên thứ 3 để gửi mail xác minh (OTP)
func GuiMailXacMinhAPI(email, code string) error {
	payload := map[string]string{
		"type": "sender_mail", "api_key": KEY_API_MAIL, "email": email, "code": code,
	}
	return callEmailAPI(payload)
}

// Gọi API bên thứ 3 để gửi mail thông báo (Nội dung tùy chỉnh - Dùng cho PIN mới)
func GuiMailThongBaoAPI(email, subject, name, body string) error {
	payload := map[string]string{
		"type": "sender", "api_key": KEY_API_MAIL, "email": email,
		"subject": subject, "name": name, "body": body,
	}
	return callEmailAPI(payload)
}

func callEmailAPI(payload interface{}) error {
	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(URL_API_MAIL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil { return err }
	defer resp.Body.Close()

	var res struct { Status string `json:"status"`; Messenger string `json:"messenger"` }
	json.NewDecoder(resp.Body).Decode(&res)

	if res.Status == "true" { return nil }
	return fmt.Errorf("%s", res.Messenger)
}
