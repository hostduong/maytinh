package mo_hinh

// Cấu trúc map với các cột trong Sheet PHAN_QUYEN
// A=MaChucNang, B=Nhom, C=MoTa, D=Admin, E=QuanLy, F=ThuKho, G=KeToan, H=Sale, I=Content
const (
	CotPQ_MaChucNang   = 0
	CotPQ_NhomChucNang = 1
	CotPQ_MoTa         = 2
	CotPQ_ADMIN        = 3
	CotPQ_QUAN_LY      = 4 // Nếu có dùng
	CotPQ_THU_KHO      = 5
	CotPQ_KE_TOAN      = 6
	CotPQ_SALE         = 7
	CotPQ_CONTENT      = 8
)

// Dùng để hiển thị lên giao diện Admin (nếu cần trang quản lý quyền)
type ChucNangHeThong struct {
	MaChucNang   string
	NhomChucNang string
	MoTa         string
}
