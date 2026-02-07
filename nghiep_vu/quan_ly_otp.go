package nghiep_vu

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
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

var CacheOTP = make(map[string]ThongTinOTP)
var CacheRate = make(map[string]*BoDemRate)
var mtxOTP sync.Mutex

func TaoMaOTP() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(99999999))
	return fmt.Sprintf("%08d", n.Int64())
}

func LuuOTP(userKey string, code string) {
	mtxOTP.Lock(); defer mtxOTP.Unlock()
	CacheOTP[userKey] = ThongTinOTP{MaCode: code, HetHanLuc: time.Now().Add(10 * time.Minute).Unix()}
}

func KiemTraOTP(userKey string, inputCode string) bool {
	mtxOTP.Lock(); defer mtxOTP.Unlock()
	otp, ok := CacheOTP[userKey]
	if !ok || time.Now().Unix() > otp.HetHanLuc { return false }
	if otp.MaCode == inputCode { delete(CacheOTP, userKey); return true }
	return false
}

func TaoMaOTP6So() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(999999))
	return fmt.Sprintf("%06d", n.Int64())
}

func KiemTraRateLimit(email string) (bool, string) {
	mtxOTP.Lock(); defer mtxOTP.Unlock()
	now := time.Now().Unix()
	rd, ok := CacheRate[email]
	if !ok || now > rd.ResetLuc { CacheRate[email] = &BoDemRate{ResetLuc: now + 21600}; rd = CacheRate[email] }
	if now-rd.LanGuiCuoi < 60 { return false, fmt.Sprintf("Vui lòng đợi %d giây.", 60-(now-rd.LanGuiCuoi)) }
	if rd.SoLanTrong6h >= 10 { return false, "Vượt quá 10 lần gửi/6h." }
	rd.LanGuiCuoi = now; rd.SoLanTrong6h++; return true, ""
}

func GuiMailXacMinhAPI(email, code string) error {
	return callApi(map[string]string{"type": "sender_mail", "api_key": KEY_API_MAIL, "email": email, "code": code})
}

func GuiMailThongBaoAPI(email, subject, name, body string) error {
	return callApi(map[string]string{"type": "sender", "api_key": KEY_API_MAIL, "email": email, "subject": subject, "name": name, "body": body})
}

// [CẬP NHẬT QUAN TRỌNG]: Lấy đúng field 'messenger' để trả về lỗi
func callApi(payload interface{}) error {
	b, _ := json.Marshal(payload)
	resp, err := http.Post(URL_API_MAIL, "application/json", bytes.NewBuffer(b))
	if err != nil { return fmt.Errorf("Lỗi kết nối đến Google: %v", err) }
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	// log.Printf("📧 [MAIL DEBUG]: %s", string(bodyBytes)) // Bật lên nếu cần debug

	// Cấu trúc hứng phản hồi từ Google Apps Script
	var r struct { 
		Status string `json:"status"` 
		Messenger string `json:"messenger"` 
	}
	
	if err := json.Unmarshal(bodyBytes, &r); err != nil {
		return fmt.Errorf("Lỗi định dạng phản hồi từ Mail Service")
	}

	if r.Status == "true" { 
		return nil 
	}
	
	// Trả về nguyên văn thông báo lỗi từ Google Script (VD: "Email không đúng định dạng", "Sai API Key"...)
	if r.Messenger != "" {
		return fmt.Errorf("%s", r.Messenger)
	}
	return fmt.Errorf("Gửi mail thất bại (Lỗi không xác định)")
}
