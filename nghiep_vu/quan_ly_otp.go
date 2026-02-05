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

// Cấu hình API gửi thư ngoài (Giữ nguyên yêu cầu)
const (
	URL_API_MAIL = "https://script.google.com/macros/s/AKfycbxd40H4neotKdnL54uQevZgSZpyZKXWfV7kJhNLY0oD9pPPA5Mn75KlFWvFd5WqiokZyA/exec"
	KEY_API_MAIL = "A1qPqCeLaX9oO0ozrMiH1a2IJKFDaj095Dlhmr8STXuS3cCmOe"
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

// Maps lưu trữ trong RAM
var (
	CacheOTP  = make(map[string]ThongTinOTP)
	CacheRate = make(map[string]*BoDemRate)
	mtxOTP    sync.Mutex
)

// --- HÀM GỐC (BẮT BUỘC GIỮ ĐỂ BUILD THÀNH CÔNG) ---

// TaoMaOTP : Sinh mã 8 số ngẫu nhiên (Dùng cho quen_mat_khau.go)
func TaoMaOTP() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(99999999))
	return fmt.Sprintf("%08d", n.Int64())
}

// LuuOTP : Lưu OTP vào RAM (Dùng cho quen_mat_khau.go)
func LuuOTP(userKey string, code string) {
	mtxOTP.Lock()
	defer mtxOTP.Unlock()
	CacheOTP[userKey] = ThongTinOTP{
		MaCode:    code,
		HetHanLuc: time.Now().Add(10 * time.Minute).Unix(),
	}
}

// KiemTraOTP : Kiểm tra OTP (Dùng chung cho cả cũ và mới)
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

// --- HÀM BỔ SUNG CHO TÍNH NĂNG EMAIL OTP ---

// TaoMaOTP6So : Sinh mã 6 số (Theo yêu cầu API gửi mail)
func TaoMaOTP6So() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(999999))
	return fmt.Sprintf("%06d", n.Int64())
}

func KiemTraRateLimit(email string) (bool, string) {
	mtxOTP.Lock()
	defer mtxOTP.Unlock()

	now := time.Now().Unix()
	rd, ok := CacheRate[email]

	if !ok || now > rd.ResetLuc {
		CacheRate[email] = &BoDemRate{LanGuiCuoi: 0, SoLanTrong6h: 0, ResetLuc: now + (6 * 3600)}
		rd = CacheRate[email]
	}

	if now-rd.LanGuiCuoi < 60 {
		return false, fmt.Sprintf("Vui lòng đợi %d giây trước khi gửi lại mã.", 60-(now-rd.LanGuiCuoi))
	}

	if rd.SoLanTrong6h >= 10 {
		return false, "Bạn đã vượt quá 10 lần gửi trong 6 giờ."
	}
	return true, ""
}

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

func GuiMailXacMinhAPI(email, code string) error {
	payload := map[string]string{
		"type": "sender_mail", "api_key": KEY_API_MAIL, "email": email, "code": code,
	}
	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(URL_API_MAIL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil { return err }
	defer resp.Body.Close()

	var resObj struct { Status string `json:"status"`; Messenger string `json:"messenger"` }
	json.NewDecoder(resp.Body).Decode(&resObj)
	if resObj.Status == "true" { return nil }
	return fmt.Errorf("%s", resObj.Messenger)
}
