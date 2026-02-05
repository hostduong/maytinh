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

// Cấu hình API Mail
const (
	URL_API_MAIL = "https://script.google.com/macros/s/AKfycbxd40H4neotKdnL54uQevZgSZpyZKXWfV7kJhNLY0oD9pPPA5Mn75KlFWvFd5WqiokZyA/exec"
	KEY_API_MAIL = "A1qPqCeLaX9oO0ozrMiH1a2IJKFDaj095Dlhmr8STXuS3cCmOe"
)

type ThongTinOTP struct {
	MaCode    string
	HetHanLuc int64
}

// Struct mới cho Rate Limit
type BoDemRate struct {
	LanGuiCuoi   int64
	SoLanTrong6h int
	ResetLuc     int64
}

var CacheOTP = make(map[string]ThongTinOTP)
var CacheRate = make(map[string]*BoDemRate) // Map mới
var mtxOTP sync.Mutex

// --- [HÀM CŨ - GIỮ NGUYÊN] ---
func TaoMaOTP() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(99999999))
	return fmt.Sprintf("%08d", n.Int64())
}

func LuuOTP(userKey string, code string) {
	mtxOTP.Lock()
	defer mtxOTP.Unlock()
	CacheOTP[userKey] = ThongTinOTP{
		MaCode:    code,
		HetHanLuc: time.Now().Add(10 * time.Minute).Unix(),
	}
}

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

// --- [HÀM MỚI BỔ SUNG - KHÔNG XÓA CÁI CŨ] ---

// 1. Sinh mã 6 số (dùng cho Quên mật khẩu OTP)
func TaoMaOTP6So() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(999999))
	return fmt.Sprintf("%06d", n.Int64())
}

// 2. Kiểm tra Rate Limit gửi mail
func KiemTraRateLimit(email string) (bool, string) {
	mtxOTP.Lock()
	defer mtxOTP.Unlock()

	now := time.Now().Unix()
	rd, ok := CacheRate[email]

	if !ok || now > rd.ResetLuc {
		CacheRate[email] = &BoDemRate{ResetLuc: now + 21600} // 6 tiếng
		rd = CacheRate[email]
	}

	if now-rd.LanGuiCuoi < 60 {
		return false, fmt.Sprintf("Vui lòng đợi %d giây.", 60-(now-rd.LanGuiCuoi))
	}
	if rd.SoLanTrong6h >= 10 {
		return false, "Vượt quá giới hạn 10 lần/6 giờ."
	}
	
	rd.LanGuiCuoi = now
	rd.SoLanTrong6h++
	return true, ""
}

// 3. Gọi API gửi thư (OTP)
func GuiMailXacMinhAPI(email, code string) error {
	p := map[string]string{"type": "sender_mail", "api_key": KEY_API_MAIL, "email": email, "code": code}
	return callApi(p)
}

// 4. Gọi API gửi thư (Thông báo PIN mới)
func GuiMailThongBaoAPI(email, subject, name, body string) error {
	p := map[string]string{
		"type": "sender", "api_key": KEY_API_MAIL, "email": email,
		"subject": subject, "name": name, "body": body,
	}
	return callApi(p)
}

func callApi(payload interface{}) error {
	b, _ := json.Marshal(payload)
	resp, err := http.Post(URL_API_MAIL, "application/json", bytes.NewBuffer(b))
	if err != nil { return err }
	defer resp.Body.Close()
	
	var r struct { Status string `json:"status"`; Messenger string `json:"messenger"` }
	json.NewDecoder(resp.Body).Decode(&r)
	
	if r.Status == "true" { return nil }
	return fmt.Errorf("%s", r.Messenger)
}
