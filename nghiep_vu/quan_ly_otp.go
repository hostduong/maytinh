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

const (
	URL_API_MAIL = "https://script.google.com/macros/s/AKfycbxd40H4neotKdnL54uQevZgSZpyZKXWfV7kJhNLY0oD9pPPA5Mn75KlFWvFd5WqiokZyA/exec"
	KEY_API_MAIL = "A1qPqCeLaX9oO0ozrMiH1a2IJKFDaj095Dlhmr8STXuS3cCmOe"
)

type ThongTinOTP struct { MaCode string; HetHanLuc int64 }
type BoDemRate struct { LanGuiCuoi int64; SoLanTrong6h int; ResetLuc int64 }

var (
	CacheOTP  = make(map[string]ThongTinOTP)
	CacheRate = make(map[string]*BoDemRate)
	mtxOTP    sync.Mutex
)

// HÀM GỐC 8 SỐ
func TaoMaOTP() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(99999999))
	return fmt.Sprintf("%08d", n.Int64())
}

// HÀM MỚI 6 SỐ
func TaoMaOTP6So() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(999999))
	return fmt.Sprintf("%06d", n.Int64())
}

func LuuOTP(key string, code string) {
	mtxOTP.Lock(); defer mtxOTP.Unlock()
	CacheOTP[key] = ThongTinOTP{MaCode: code, HetHanLuc: time.Now().Add(10 * time.Minute).Unix()}
}

func KiemTraOTP(key string, code string) bool {
	mtxOTP.Lock(); defer mtxOTP.Unlock()
	otp, ok := CacheOTP[key]
	if !ok || time.Now().Unix() > otp.HetHanLuc || otp.MaCode != code { return false }
	delete(CacheOTP, key); return true
}

func KiemTraRateLimit(email string) (bool, string) {
	mtxOTP.Lock(); defer mtxOTP.Unlock()
	now := time.Now().Unix()
	rd, ok := CacheRate[email]
	if !ok || now > rd.ResetLuc { CacheRate[email] = &BoDemRate{ResetLuc: now + 21600}; rd = CacheRate[email] }
	if now-rd.LanGuiCuoi < 60 { return false, fmt.Sprintf("Đợi %d giây.", 60-(now-rd.LanGuiCuoi)) }
	if rd.SoLanTrong6h >= 10 { return false, "Vượt giới hạn 10 lần/6h." }
	rd.LanGuiCuoi = now; rd.SoLanTrong6h++; return true, ""
}

func GuiMailXacMinhAPI(email, code string) error {
	p := map[string]string{"type": "sender_mail", "api_key": KEY_API_MAIL, "email": email, "code": code}
	return callEmailAPI(p)
}

func GuiMailThongBaoAPI(email, subject, name, body string) error {
	p := map[string]string{"type": "sender", "api_key": KEY_API_MAIL, "email": email, "subject": subject, "name": name, "body": body}
	return callEmailAPI(p)
}

func callEmailAPI(payload interface{}) error {
	b, _ := json.Marshal(payload)
	resp, err := http.Post(URL_API_MAIL, "application/json", bytes.NewBuffer(b))
	if err != nil { return err }
	defer resp.Body.Close()
	var r struct{ Status string `json:"status"` }; json.NewDecoder(resp.Body).Decode(&r)
	if r.Status == "true" { return nil }
	return fmt.Errorf("API error")
}
