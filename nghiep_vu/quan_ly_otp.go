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

// Cấu hình API ngoài
const (
	URL_API_GUIMAIL = "https://script.google.com/macros/s/AKfycbxd40H4neotKdnL54uQevZgSZpyZKXWfV7kJhNLY0oD9pPPA5Mn75KlFWvFd5WqiokZyA/exec"
	KEY_API_GUIMAIL = "A1qPqCeLaX9oO0ozrMiH1a2IJKFDaj095Dlhmr8STXuS3cCmOe"
)

type ThongTinOTP struct {
	MaCode    string
	HetHanLuc int64
}

type BoDemRate struct {
	LanGuiCuoi   int64
	SoLanTrong6h int
	ResetLuc     int64
}

var (
	CacheOTP  = make(map[string]ThongTinOTP)
	CacheRate = make(map[string]*BoDemRate)
	mtxOTP    sync.Mutex
)

// Sinh mã 6 số ngẫu nhiên
func TaoMaOTP6So() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(999999))
	return fmt.Sprintf("%06d", n.Int64())
}

// Kiểm tra Rate Limit
func KiemTraRateLimit(email string) (bool, string) {
	mtxOTP.Lock()
	defer mtxOTP.Unlock()

	now := time.Now().Unix()
	rd, tonTai := CacheRate[email]

	// Nếu chưa có hoặc đã qua chu kỳ 6h thì reset bộ đếm
	if !tonTai || now > rd.ResetLuc {
		CacheRate[email] = &BoDemRate{
			LanGuiCuoi:   0,
			SoLanTrong6h: 0,
			ResetLuc:     now + (6 * 3600),
		}
		rd = CacheRate[email]
	}

	// 1. Kiểm tra giãn cách 60 giây
	if now-rd.LanGuiCuoi < 60 {
		giayConLai := 60 - (now - rd.LanGuiCuoi)
		return false, fmt.Sprintf("Vui lòng đợi %d giây trước khi gửi lại mã.", giayConLai)
	}

	// 2. Kiểm tra giới hạn 10 lần
	if rd.SoLanTrong6h >= 10 {
		return false, "Bạn đã vượt quá 10 lần gửi trong 6 giờ."
	}

	return true, ""
}

// Lưu OTP và cập nhật bộ đếm
func LuuOTPVaUpdateRate(email string, code string) {
	mtxOTP.Lock()
	defer mtxOTP.Unlock()

	CacheOTP[email] = ThongTinOTP{
		MaCode:    code,
		HetHanLuc: time.Now().Add(10 * time.Minute).Unix(),
	}

	rd := CacheRate[email]
	rd.LanGuiCuoi = time.Now().Unix()
	rd.SoLanTrong6h++
}

// Gọi API ngoài để gửi thư
func GuiMailXacMinhAPI(email, code string) error {
	payload := map[string]string{
		"type":    "sender_mail",
		"api_key": KEY_API_GUIMAIL,
		"email":   email,
		"code":    code,
	}
	jsonData, _ := json.Marshal(payload)

	resp, err := http.Post(URL_API_GUIMAIL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		Status    string `json:"status"`
		Messenger string `json:"messenger"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	if result.Status == "true" {
		return nil
	}
	return fmt.Errorf("%s", result.Messenger)
}

// Verifier (One-time use)
func KiemTraOTP(email string, inputCode string) bool {
	mtxOTP.Lock()
	defer mtxOTP.Unlock()

	otp, ok := CacheOTP[email]
	if !ok { return false }

	if time.Now().Unix() > otp.HetHanLuc {
		delete(CacheOTP, email)
		return false
	}

	if otp.MaCode == inputCode {
		delete(CacheOTP, email)
		return true
	}
	return false
}
